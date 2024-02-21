package handler

import (
	"bytes"
	"encoding/csv"
	"log"
	"strconv"

	"ozon-broker/internal/dal/ozon-broker"
)

const (
	ProductsFile = "products.tsv"
	OrdersFile   = "orders.tsv"
)

func ProductToCsvFile(products []ozon_broker.Product, buf *bytes.Buffer) {

	writer := csv.NewWriter(buf)
	writer.Comma = '\t'
	defer writer.Flush()
	header := []string{
		"ID",
		"Offer_ID",
		"SKU",
		"Навзвание товара",
		"Баркод",
		"Все баркоды",
		"Изображение",
		"Минимальная цена",
		"Цена",
		"Цена для премиум",
	}

	err := writer.Write(header)
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range products {
		barcodes := ""
		for _, barcode := range product.Barcodes {
			barcodes += barcode + " | "
		}
		strArr := []string{
			strconv.Itoa(int(product.Id)),
			product.OfferID,
			strconv.Itoa(int(product.SKU)),
			product.Name,
			product.Barcode,
			barcodes,
			product.Image,
			product.MinimumPrice,
			product.Price,
			product.PremiumPrice,
		}
		err := writer.Write(strArr)
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	}

}

func OrdersToCsvFile(orders []ozon_broker.Order, buf *bytes.Buffer) {

	writer := csv.NewWriter(buf)
	writer.Comma = '\t'
	defer writer.Flush()

	header := []string{
		"PostingNumber",
		"OrderNumber",
		"OrderId",
		"Status",
		"CancelReasonId",
		"CreatedAt",
		"InProcessAt",
		"SKU",
		"Qty",
		"Price",
		"CurrencyCode",
		"OfferId",
		"Name",
		"DeliveryType",
		"Region",
		"City",
		"WarehouseName",
		"PostingServicesTotal",
		"ItemServicesTotal",
		"CommissionAmount",
		"Payout",
	}

	err := writer.Write(header)
	if err != nil {
		log.Fatal(err)
	}

	for _, order := range orders {

		strArr := []string{
			order.PostingNumber,
			order.OrderNumber,
			strconv.FormatInt(order.OrderId, 10),
			order.Status,
			strconv.FormatInt(order.CancelReasonId, 10),
			order.CreatedAt.String(),
			order.InProcessAt.String(),
			strconv.FormatInt(order.SKU, 10),
			strconv.FormatInt(order.Qty, 10),
			order.Price,
			order.CurrencyCode,
			order.OfferId,
			order.Name,
			order.DeliveryType,
			order.Region,
			order.City,
			order.WarehouseName,
			strconv.FormatFloat(order.PostingServicesTotal, 'f', -1, 64),
			strconv.FormatFloat(order.ItemServicesTotal, 'f', -1, 64),
			strconv.FormatFloat(order.CommissionAmount, 'f', -1, 64),
			strconv.FormatFloat(order.Payout, 'f', -1, 64),
		}

		err := writer.Write(strArr)
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()

	}

}
