/*
 Package models provides the model of state objects.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package models

// Account: Account model
type Account struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	DOB             string `json:"DOB"`
	Email           string `json:"email"`
	Mobile          string `json:"mobile"`
	Company         string `json:"company"`
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
