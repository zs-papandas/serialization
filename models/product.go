package models

import (
	"github.com/zs-papandas/serialization/types"
)

// Product: Product model
type Product struct {
	SerialId       	string 				`json:"serialid"`
	Created        	string 				`json:"created"`
	Owner    		Account				`json:"owner"`
	Name           	string 				`json:"name"`
	Expire          string 				`json:"expire"`
	GTIN         	string 				`json:"gtin"`
	LotNumber       string 				`json:"lotnumber"`
	Status         	string 				`json:"status"`
	Amount 			string				`json:"amount"`
	TotalQty        int 				`json:"totalqty"`
	AvailQty        int 				`json:"availqty"`
	ProductType 	types.ProductType	`json:"product_type"`
	ParentProduct	string				`json:"parent_product"`
}
