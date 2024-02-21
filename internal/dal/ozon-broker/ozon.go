package ozon_broker

import (
	"context"
	"github.com/diphantxm/ozon-api-client/ozon"
	"log"
	"net/http"
	"time"
)

type Product struct {
	Id           int64
	OfferID      string
	SKU          int64
	Name         string
	Barcode      string
	Barcodes     []string
	Image        string
	MinimumPrice string
	Price        string
	PremiumPrice string
}

type Order struct {
	PostingNumber        string
	OrderNumber          string
	OrderId              int64
	Status               string
	CancelReasonId       int64
	CreatedAt            time.Time
	InProcessAt          time.Time
	SKU                  int64
	Qty                  int64
	Price                string
	CurrencyCode         string
	OfferId              string
	Name                 string
	DeliveryType         string
	Region               string
	City                 string
	WarehouseName        string
	PostingServicesTotal float64
	ItemServicesTotal    float64
	CommissionAmount     float64
	Payout               float64
}

func GetListOfProducts(ctx context.Context, client *ozon.Client) ([]ozon.GetListOfProductsResultItem, error) {
	resp, err := client.Products().GetListOfProducts(ctx, &ozon.GetListOfProductsParams{
		Limit: 1000,
		Filter: ozon.GetListOfProductsFilter{
			Visibility: "ALL",
		},
	})

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return resp.Result.Items, err
}

func GetProductDetails(ctx context.Context, item ozon.GetListOfProductsResultItem, client *ozon.Client) (Product, error) {
	id := item.ProductId
	offerId := item.OfferId
	r, err := client.Products().GetProductDetails(ctx, &ozon.GetProductDetailsParams{
		ProductId: id,
	})

	if err != nil || r.StatusCode != http.StatusOK {
		return Product{}, err
	}
	name := r.Result.Name
	barcode := r.Result.Barcode
	barcodes := r.Result.Barcodes
	sku := r.Result.SKU
	image := r.Result.PrimaryImage
	minimumPrice := r.Result.MinPrice
	price := r.Result.Price
	premiumPrice := r.Result.PremiumPrice
	product := Product{
		id,
		offerId,
		sku,
		name,
		barcode,
		barcodes,
		image,
		minimumPrice,
		price,
		premiumPrice}
	return product, err
}

func GetListOfOrders(ctx context.Context, client *ozon.Client) ([]ozon.GetFBOShipmentsListResult, error) {

	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)
	resp, err := client.FBO().GetShipmentsList(ctx, &ozon.GetFBOShipmentsListParams{
		Limit:     1000,
		Direction: "ASC",
		Filter: ozon.GetFBOShipmentsListFilter{
			Status: "",
			Since:  oneYearAgo,
			To:     now,
		},
		With: ozon.GetFBOShipmentsListWith{
			AnalyticsData: true,
			FinancialData: true,
		}})
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Print(resp.Message)
		return nil, err
	}
	//log.Print(resp.Message)
	return resp.Result, err

}

func GetOrderDetails(orderResult ozon.GetFBOShipmentsListResult) Order {

	postingNumber := orderResult.PostingNumber
	orderNumber := orderResult.OrderNumber
	orderId := orderResult.OrderId
	status := orderResult.Status
	cancelReasonId := orderResult.CancelReasonId
	createdAt := orderResult.CreatedAt
	inProcessAt := orderResult.InProccessAt
	deliveryType := orderResult.AnalyticsData.DeliveryType
	region := orderResult.AnalyticsData.Region
	city := orderResult.AnalyticsData.City
	warehouseName := orderResult.AnalyticsData.WarehouseName
	s := orderResult.FinancialData.PostingServices
	postingServicesTotal := s.DropoffFF + s.DropoffSC + s.DropoffPVZ + s.DeliveryToCustomer + s.DirectFlowTrans + s.ReturnFlowTrans + s.ReturnAfterDeliveryToCustomer + s.Fulfillment + s.Pickup + s.ReturnNotDeliveryToCustomer + s.ReturnPartGoodsCustomer

	products := orderResult.Products
	for _, product := range products {
		sku := product.SKU
		qty := product.Quantity
		price := product.Price
		currencyCode := product.CurrencyCode
		offerId := product.OfferId
		name := product.Name
		//id := product
		for _, item := range orderResult.FinancialData.Products {
			if item.ProductId == sku {
				payout := item.Payout
				commissionAmount := item.CommissionAmount
				i := item.ItemServices
				itemServicesTotal := i.DropoffFF + i.DropoffSC + i.DropoffPVZ + i.DeliveryToCustomer + i.DirectFlowTrans + i.ReturnFlowTrans + i.ReturnAfterDeliveryToCustomer + i.Fulfillment + i.Pickup + i.ReturnNotDeliveryToCustomer + i.ReturnPartGoodsCustomer

				order := Order{
					PostingNumber:        postingNumber,
					OrderNumber:          orderNumber,
					OrderId:              orderId,
					Status:               status,
					CancelReasonId:       cancelReasonId,
					CreatedAt:            createdAt,
					InProcessAt:          inProcessAt,
					SKU:                  sku,
					Qty:                  qty,
					Price:                price,
					CurrencyCode:         currencyCode,
					OfferId:              offerId,
					Name:                 name,
					DeliveryType:         deliveryType,
					Region:               region,
					City:                 city,
					WarehouseName:        warehouseName,
					PostingServicesTotal: postingServicesTotal,
					ItemServicesTotal:    itemServicesTotal,
					CommissionAmount:     commissionAmount,
					Payout:               payout,
				}

				return order
			}
		}

	}
	order := Order{}
	return order
}
