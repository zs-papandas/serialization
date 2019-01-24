package contracts

import (
	"encoding/json"
	"time"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/zs-papandas/serialization/models"
	"github.com/zs-papandas/serialization/utils"

)

var productLogger = shim.NewLogger("contracts/product")

// ProductContract : a struct to handle Product.
type ProductContract struct {
}

//CreateProduct : save a product
func (ac *ProductContract) CreateProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke CreateProduct, args=%s\n", args)

	if len(args) != 7 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. %s\n", args)
		productLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	no, err := utils.GetSerialNo(APIstub)
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	today := time.Now().Format(time.RFC3339)
	
	// GET USER ACCOUNT DETAIL
	toAccount, err := utils.GetAccount(APIstub, args[0])
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			productLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			productLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	productLogger.Infof("User Account %s\n", toAccount)
	creator := toAccount
	name := args[1]
	expired := args[2]
	gtin := args[3]
	lotnum := args[4]
	status := "CREATED"
	amt := args[5]
	myStr := args[6]
	myInt, err := strconv.Atoi(myStr)
	if err != nil {
        errMsg := fmt.Sprintf("Failed: string to int. %s\n", myStr)
		productLogger.Error(errMsg)
		return shim.Error(errMsg)
    }
	totqty := myInt
	avaiqty := myInt

	//SerialId Created Creator  Name Expire GTIN LotNumber Status Amount  TotalQty  AvailQty
	var productInfo models.Product
	productInfo = models.Product{no, today, creator, name, expired, gtin, lotnum, status, amt, totqty, avaiqty}
 	jsonBytes, err := json.Marshal(&productInfo)
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	if err := APIstub.PutState(no, jsonBytes); err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)

	
	//return shim.Success([]byte("Reply from CreateProduct"))
	
}

// RetrieveProduct : return a product.
func (ac *ProductContract) RetrieveProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke CreateProduct, args=%s\n", args)

	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	product, err := APIstub.GetState(args[0])
	if err != nil {
		errMsg1 := fmt.Sprintf("Failed to get asset: %s with error: %s", args[0], err)
		accountLogger.Error(errMsg1)
		return shim.Error(errMsg1)
	}

	return shim.Success(product)
	//return shim.Success([]byte("Reply from RetrieveProduct"))
	
}


// ChangeOwner : Change owner of a product.
func (ac *ProductContract) ChangeOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke ChangeOwner, args=%s\n", args)

	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	
	/*serialId := args[0]
	toOwner := args[1]

	product, err := APIstub.GetState(args[0])
	if err != nil {
		errMsg1 := fmt.Sprintf("Failed to get asset: %s with error: %s", args[0], err)
		accountLogger.Error(errMsg1)
		return shim.Error(errMsg1)
	}*/

	//return shim.Success(product)
	return shim.Success([]byte("Reply from ChangeOwner"))
	
}