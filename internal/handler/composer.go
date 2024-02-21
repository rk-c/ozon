package handler

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/diphantxm/ozon-api-client/ozon"
	"google.golang.org/api/drive/v3"
	"log"
	"ozon-broker/internal/dal/db"
	"ozon-broker/internal/dal/google"
	o "ozon-broker/internal/dal/ozon-broker"
	"sync"
)

func ProductComposer(ctx context.Context, client *ozon.Client, psql *sql.DB, wg *sync.WaitGroup, srv *drive.Service) {
	buf := &bytes.Buffer{}
	defer wg.Done()

	products, err := o.GetListOfProducts(ctx, client)

	if err != nil {
		log.Fatalf("Cannot get list of products: %s", err)
	}

	conn, err := psql.Conn(ctx)
	if err != nil {
		log.Fatalf("Cannot get database connection: %s", err)
	}
	defer conn.Close()

	var wg1 sync.WaitGroup

	for _, product := range products {
		wg1.Add(1)

		go func(product ozon.GetListOfProductsResultItem) {
			defer wg1.Done()
			detailedProduct, err := o.GetProductDetails(ctx, product, client)
			if err != nil {
				log.Fatalf("Cannot get details for product: %s", err)
			}

			err = db.UpdateProduct(ctx, conn, detailedProduct)
			if err != nil {
				log.Fatalf("Cannot update product: %s", err)
			}

		}(product)

	}

	wg1.Wait()

	dbProducts, err := db.GetProducts(psql)

	ProductToCsvFile(dbProducts, buf)
	google.DeleteFile(ProductsFile, srv)
	google.UploadFile(ProductsFile, buf, srv)

}

func OrderComposer(ctx context.Context, client *ozon.Client, psql *sql.DB, wg *sync.WaitGroup, srv *drive.Service) {
	buf := &bytes.Buffer{}
	defer wg.Done()
	orders, err := o.GetListOfOrders(ctx, client)

	if err != nil {
		log.Fatalf("Cannot get list of orders: %s", err)
	}

	conn, err := psql.Conn(ctx)
	if err != nil {
		log.Fatalf("Cannot get database connection: %s", err)
	}
	defer conn.Close()

	var wg1 sync.WaitGroup

	for _, order := range orders {
		wg1.Add(1)

		go func(order ozon.GetFBOShipmentsListResult) {
			defer wg1.Done()
			detailedOrder := o.GetOrderDetails(order)

			err = db.UpdateOrder(ctx, conn, detailedOrder)

			if err != nil {
				log.Fatalf("Cannot update order: %s", err)
			}
		}(order)
	}

	wg1.Wait()

	dbOrders, err := db.GetOrders(psql)

	OrdersToCsvFile(dbOrders, buf)
	google.DeleteFile(OrdersFile, srv)
	google.UploadFile(OrdersFile, buf, srv)
}
