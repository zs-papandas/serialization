/*
 Package contracts provides the smart contracts for Hyperledger/fabric 1.1.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/zs-papandas/serialization/models"
)

var accountLogger = shim.NewLogger("contracts/account")

// AccountContract : a struct to handle Account.
type AccountContract struct {
}

// ListAccount : return a list of all accounts.
func (ac *AccountContract) ListAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke ListAccount, args=%s\n", args)
	if len(args) != 0 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = no argument, Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	
	return shim.Success([]byte("Reply from ListAccount"))
}


// ListAccount : return a list of all accounts.
func (ac *AccountContract) CreateAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke CreateAccount, args=%s\n", args)

	if len(args) != 7 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['name'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	var personalInfo models.Account
 	personalInfo = models.Account{args[1], args[2], args[3], args[4], args[5], args[6]}
 	jsonBytes, err := json.Marshal(&personalInfo)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	if err := APIstub.PutState(args[0], jsonBytes); err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}


// ListAccount : return a list of all accounts.
func (ac *AccountContract) RetrieveAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke RetrieveAccount")
	return shim.Success([]byte("Reply from RetrieveAccount"))
}