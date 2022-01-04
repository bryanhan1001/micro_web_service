// Code generated by gormer. DO NOT EDIT.
package model

import (
	"context"
	"time"

	"github.com/Junedayday/micro_web_service/internal/gormer"
)

type OrderModel interface {
	AddOrder(ctx context.Context, order *gormer.Order) (err error)
	AddOrders(ctx context.Context, orders []*gormer.Order) (err error)
	QueryOrders(ctx context.Context, pageNumber, pageSize int, condition *gormer.OrderOptions) (orders []gormer.Order, err error)
	CountOrders(ctx context.Context, condition *gormer.OrderOptions) (count int64, err error)
	UpdateOrder(ctx context.Context, updated, condition *gormer.OrderOptions) (err error)
	DeleteOrder(ctx context.Context, condition *gormer.OrderOptions) (err error)

	// defined in genQueries in gormer.yaml

	// QueryOrdersDesc 根据id逆序查询
	QueryOrdersDesc(ctx context.Context, pageNumber, pageSize int, condition *gormer.OrderOptions) (orders []gormer.Order, err error)
	// CountOrdersDesc 根据id逆序查询
	CountOrdersDesc(ctx context.Context, condition *gormer.OrderOptions) (count int64, err error)
	// QueryOrdersByNamesAndCreateTime 根据名称和创建时间查询
	QueryOrdersByNamesAndCreateTime(ctx context.Context, names []string, createTime time.Time, pageNumber, pageSize int, condition *gormer.OrderOptions) (orders []gormer.Order, err error)
	// CountOrdersByNamesAndCreateTime 根据名称和创建时间查询
	CountOrdersByNamesAndCreateTime(ctx context.Context, names []string, createTime time.Time, condition *gormer.OrderOptions) (count int64, err error)
	// QueryOrdersByTimeRange 根据创建时间的范围查询
	QueryOrdersByTimeRange(ctx context.Context, startTime time.Time, endTime time.Time, pageNumber, pageSize int, condition *gormer.OrderOptions) (orders []gormer.Order, err error)
	// CountOrdersByTimeRange 根据创建时间的范围查询
	CountOrdersByTimeRange(ctx context.Context, startTime time.Time, endTime time.Time, condition *gormer.OrderOptions) (count int64, err error)

	// Implement Your Method in ext model
	OrderExtModel
}
