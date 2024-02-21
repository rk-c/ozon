package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"

	o "ozon-broker/internal/dal/ozon-broker"
)

func InitDB(host string, port int, user string, psw string, dbname string) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, psw, dbname)

	db, err := sql.Open("postgres", psqlconn)

	CheckError(err)

	return db, err
}

func GetProducts(db *sql.DB) ([]o.Product, error) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []o.Product
	for rows.Next() {
		var product o.Product
		err = rows.Scan(
			&product.Id,
			&product.OfferID,
			&product.SKU,
			&product.Name,
			&product.Barcode,
			pq.Array(&product.Barcodes),
			&product.Image,
			&product.MinimumPrice,
			&product.Price,
			&product.PremiumPrice,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func UpdateProduct(ctx context.Context, conn *sql.Conn, product o.Product) error {
	query := `
		INSERT INTO products (
			id, offerid, sku, name, barcode, barcodes, image, minimumprice, price, premiumprice
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (Id) DO UPDATE SET
			Id = $1, OfferID = $2, SKU = $3, Name = $4, Barcode = $5, Barcodes = $6, Image = $7, 
			MinimumPrice = $8, Price = $9, PremiumPrice = $10
	`
	_, err := conn.ExecContext(ctx, query, product.Id, product.OfferID, product.SKU, product.Name, product.Barcode,
		pq.Array(product.Barcodes), product.Image, product.MinimumPrice, product.Price, product.PremiumPrice)

	return err
}

func GetOrders(db *sql.DB) ([]o.Order, error) {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []o.Order
	for rows.Next() {
		var order o.Order
		err = rows.Scan(
			&order.PostingNumber,
			&order.OrderNumber,
			&order.OrderId,
			&order.Status,
			&order.CancelReasonId,
			&order.CreatedAt,
			&order.InProcessAt,
			&order.SKU,
			&order.Qty,
			&order.Price,
			&order.CurrencyCode,
			&order.OfferId,
			&order.Name,
			&order.DeliveryType,
			&order.Region,
			&order.City,
			&order.WarehouseName,
			&order.PostingServicesTotal,
			&order.ItemServicesTotal,
			&order.CommissionAmount,
			&order.Payout,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func UpdateOrder(ctx context.Context, conn *sql.Conn, order o.Order) error {
	query := `
		INSERT INTO orders (
			PostingNumber, OrderNumber, OrderId, Status, CancelReasonId, CreatedAt, InProcessAt, SKU, Qty, Price, 
			CurrencyCode, OfferId, Name, DeliveryType, Region, City, WarehouseName, PostingServicesTotal, 
			ItemServicesTotal, CommissionAmount, Payout
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		ON CONFLICT (OrderId) DO UPDATE SET
			PostingNumber = $1, OrderNumber = $2, OrderId = $3, Status = $4, CancelReasonId = $5, CreatedAt = $6, InProcessAt = $7, 
			SKU = $8, Qty = $9, Price = $10, CurrencyCode = $11, OfferId = $12, Name = $13, DeliveryType = $14, 
			Region = $15, City = $16, WarehouseName = $17, PostingServicesTotal = $18, ItemServicesTotal = $19, 
			CommissionAmount = $20, Payout = $21
	`
	_, err := conn.ExecContext(ctx, query, order.PostingNumber, order.OrderNumber, order.OrderId, order.Status, order.CancelReasonId,
		order.CreatedAt, order.InProcessAt, order.SKU, order.Qty, order.Price, order.CurrencyCode, order.OfferId,
		order.Name, order.DeliveryType, order.Region, order.City, order.WarehouseName, order.PostingServicesTotal,
		order.ItemServicesTotal, order.CommissionAmount, order.Payout)

	return err
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
