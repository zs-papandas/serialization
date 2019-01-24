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

