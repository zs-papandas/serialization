package models


import (
	"github.com/zs-papandas/serialization/types"
)

// Account: Account model
type Account struct {
	Firstname	string 			`json:"firstname"`
	Lastname    string 			`json:"lastname"`
	DOB         string 			`json:"DOB"`
	Email       string 			`json:"email"`
	Mobile      string 			`json:"mobile"`
	Company     string 			`json:"company"`
	UserType 	types.UserType 	`json:"user_type"`
}

/*type ProductInfo struct {
	SerialId       	string 			`json:"serialid"`
	Created        	string 			`json:"created"`
	Creator    		PersonalInfo	`json:"creator"`
	Name           	string 			`json:"name"`
	Expire          string 			`json:"expire"`
	GTIN         	string 			`json:"gtin"`
	LotNumber       string 			`json:"lotnumber"`
	Status         	string 			`json:"status"`
	TotalQty        int 			`json:"totalqty"`
	AvailQty        int 			`json:"availqty"`
}*/
