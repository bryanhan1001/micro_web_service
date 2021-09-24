package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	dsn          string
	projectPath  string
	goMod        string
	gormPath     string
	daoPath      string
	tableMatcher string
)

const (
	gormerHeader = `// Code generated by gormer. DO NOT EDIT.
package %s

import "time"
`
	daoHeader = `// Code generated by gormer. DO NOT EDIT.
package %s

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	
	"%s/%s"
)
`
)

// go run main.go -dsn root:123456@tcp(127.0.0.1:3306)/demo
func main() {
	flag.StringVar(&dsn, "dsn", "root:123456@tcp(127.0.0.1:3306)/demo", "mysql 连接串")
	flag.StringVar(&projectPath, "projectPath", "/Users/didi/Study/micro_web_service/", "项目在本机的目录")
	flag.StringVar(&gormPath, "gormPath", "internal/gormer/", "生成的GORM结构体相对路径")
	flag.StringVar(&daoPath, "daoPath", "internal/dao/", "生成的Dao层代码相对路径")
	flag.StringVar(&goMod, "goMod", "github.com/Junedayday/micro_web_service", "包名")
	flag.StringVar(&tableMatcher, "tableMatcher", "orders:order", "将table名称做一次映射，一般用于去掉复数")
	flag.Parse()

	// 创建文件夹（如果已存在会报错，不影响）
	for _, path := range []string{gormPath, daoPath} {
		os.MkdirAll(path, os.ModePerm)
	}

	if len(dsn) == 0 {
		os.Exit(1)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("start to generate gorm structs")

	tables, err := getAllTables(db)
	if err != nil {
		fmt.Printf("getAllTables error %+v", err)
		os.Exit(1)
	}

	tMatcher := getTableMatcher(tableMatcher, tables)

	for _, table := range tables {
		// 1.生成结构
		structResult, err := Generate(db, table, tMatcher[table])
		if err != nil {
			fmt.Printf("Generate table %s error %+v", table, err)
			os.Exit(1)
		}

		// 生成gormer file
		if gormPath[len(gormPath)-1] == '/' {
			gormPath = gormPath[:len(gormPath)-1]
		}
		dirs := strings.Split(gormPath, "/")
		header := fmt.Sprintf(gormerHeader, dirs[len(dirs)-1])
		err = parseToFile(gormPath, tMatcher[table], header, structResult, parseToGormerTmpl)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}

		// 生成dao file
		if daoPath[len(daoPath)-1] == '/' {
			daoPath = daoPath[:len(daoPath)-1]
		}
		dirs = strings.Split(daoPath, "/")
		header = fmt.Sprintf(daoHeader, dirs[len(dirs)-1], goMod, gormPath)
		err = parseToFile(daoPath, tMatcher[table], header, structResult, parseToDaoTmpl)
		if err != nil {
			fmt.Printf("parseToFile error %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Generate Table %s Finished\n", table)
	}

	// go fmt files
	exec.Command("go", "fmt", gormPath+"...").Run()
}

func parseToFile(filePath, name, fileHeader string, structResult StructLevel, parseFunc func(StructLevel) (string, error)) error {
	result, err := parseFunc(structResult)
	if err != nil {
		return errors.Wrapf(err, "parseToDaoTmpl structResult %s", structResult)
	}
	path := fmt.Sprintf("%s%s/%s.go", projectPath, filePath, name)
	file, err := os.OpenFile(path, os.O_WRONLY+os.O_CREATE+os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "OpenFile path %s", path)
	}
	defer file.Close()

	_, err = file.WriteString(fileHeader + result)
	if err != nil {
		return errors.Wrap(err, "WriteString to file")
	}
	return nil
}
