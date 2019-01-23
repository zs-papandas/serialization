package contracts

import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

)

var productLogger = shim.NewLogger("contracts/product")

// ProductContract : a struct to handle Product.
type ProductContract struct {
}

//CreateProduct : save a product
func (ac *ProductContract) CreateProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke CreateProduct, args=%s\n", args)

	if len(args) != 7 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['name'], Actual = %s\n", args)
		productLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	no := "PRO" + time.Now()
	today := time.Now()
	creator := args[0]
	name := args[1]
	expired := args[2]
	gtin := args[3]
	lotnum := args[4]
	status := "CREATED"
	amt := args[5]
	totqty := args[6]
	avaiqty := args[6]

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

	return shim.Success([]byte("Reply from RetrieveProduct"))
	
}