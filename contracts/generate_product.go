package contracts

import (
	"encoding/json"
	"time"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/zs-papandas/serialization/models"
	"github.com/zs-papandas/serialization/types"
	"github.com/zs-papandas/serialization/utils"

)

var productLogger = shim.NewLogger("contracts/generate_product")

// GenerateProductContract : a struct to handle auto generate Product.
type GenerateProductContract struct {
}

//CreateProduct : save a product
func (ac *GenerateProductContract) CreateProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke Generate Product -> CreateProduct, args=%s\n", args)

	if len(args) != 9 {
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
	owner, err := utils.GetAccount(APIstub, args[0])
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

	productLogger.Infof("User Account %s\n", owner.Firstname)
	
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

	// Get the Product Type
	var ProductTypeInt types.ProductType
	switch args[7] {
	case "pallet":
		ProductTypeInt = types.PalletProduct
	case "box":
		ProductTypeInt = types.BoxProduct
	case "packet":
		ProductTypeInt = types.PacketProduct
	case "item":
		ProductTypeInt = types.ItemProduct
	default:
		ProductTypeInt = types.UnKnownProduct
	}

	parentProduct := args[8]
	

	var productInfo models.Product
	productInfo = models.Product{no, today, *owner, name, expired, gtin, lotnum, status, amt, totqty, avaiqty, ProductTypeInt, parentProduct}
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
	
}
