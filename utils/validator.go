/*
 Package utils provides some utility functions.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/zs-papandas/serialization/models"
)

// GetAccount : get an account from state db using account no.
func GetAccount(APIstub shim.ChaincodeStubInterface, no string) (*models.Account, error) {
	var account = new(models.Account)
	accountBytes, err := APIstub.GetState(no)
	if err != nil {
		return account, err
	} else if accountBytes == nil {
		msg := fmt.Sprintf("Account does not exist, no = %s", no)
		warning := &WarningResult{StatusCode: 404, Message: msg}
		return account, warning
	}
	if err := json.Unmarshal(accountBytes, account); err != nil {
		return account, err
	}
	return account, nil
}

// GetAccount : get an account from state db using account no.
func GetProduct(APIstub shim.ChaincodeStubInterface, no string) (*models.Product, error) {
	var product = new(models.Product)
	productBytes, err := APIstub.GetState(no)
	if err != nil {
		return product, err
	} else if productBytes == nil {
		msg := fmt.Sprintf("Product does not exist, no = %s", no)
		warning := &WarningResult{StatusCode: 404, Message: msg}
		return product, warning
	}
	if err := json.Unmarshal(productBytes, product); err != nil {
		return product, err
	}
	return product, nil
}

// GetAmount : convert amount to int and validate it
func GetAmount(amountStr string) (int, error) {
	var amount int
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		msg := fmt.Sprintf("amount is not integer, amount = %s", amountStr)
		warning := &WarningResult{StatusCode: 400, Message: msg}
		return amount, warning
	}
	if amount < 0 {
		msg := fmt.Sprintf("amount is less than zero, amount = %d", amount)
		warning := &WarningResult{StatusCode: 400, Message: msg}
		return amount, warning
	}
	return amount, nil
}
