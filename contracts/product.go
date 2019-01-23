package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/zs-papandas/serialization/models"
	"github.com/zs-papandas/serialization/types"
)

var productLogger = shim.NewLogger("contracts/product")

// ProductContract : a struct to handle Product.
type ProductContract struct {
}

//CreateProduct : save a product
func (ac *AccountContract) CreateProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke CreateProduct, args=%s\n", args)

	
	return shim.Success([]byte("Reply from CreateProduct"))
	
}

// RetrieveProduct : return a product.
func (ac *AccountContract) RetrieveProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke CreateProduct, args=%s\n", args)

	return shim.Success([]byte("Reply from RetrieveProduct"))
	
}