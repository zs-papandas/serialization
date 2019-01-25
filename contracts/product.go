package contracts

import (
	"encoding/json"
	"time"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric/core/scc/qscc"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	sdkCtx "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"

	"github.com/zs-papandas/serialization/models"
	"github.com/zs-papandas/serialization/types"
	"github.com/zs-papandas/serialization/utils"

)

var productLogger = shim.NewLogger("contracts/product")

// ProductContract : a struct to handle Product.
type ProductContract struct {
}

//CreateProduct : save a product
func (ac *ProductContract) CreateProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke CreateProduct, args=%s\n", args)

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

	//SerialId Created Creator  Name Expire GTIN LotNumber Status Amount  TotalQty  AvailQty
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

	if len(args) != 2 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}


	toProduct, err := utils.GetProduct(APIstub, args[0])
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

	toOwner, err := utils.GetAccount(APIstub, args[1])
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

	productLogger.Infof("Product LotNumber")
	//productLogger.Infof(*product.LotNumber)
	productLogger.Infof(toProduct.SerialId)
	productLogger.Infof(toOwner.Firstname)

	toProduct.Owner = *toOwner
	toProduct.Status = "OWNERSHIP_CHANGED"

	toProductBytes, err := json.Marshal(toProduct)
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(toProduct.SerialId, toProductBytes); err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success(toProductBytes)
	
}


// ChangeOwner : Change owner of a product.
func (ac *ProductContract) TestQueryInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke TestQueryInfo, args=%s\n", args)

	c, err := New(mockChannelProvider("myc"))
	if err != nil {
		fmt.Println("failed to create client")
	}

	bci, err := c.QueryInfo()
	if err != nil {
		fmt.Printf("failed to query for blockchain info: %s\n", err)
	}

	if bci != nil {
		fmt.Println("Retrieved ledger info")
	}

	/*client, err := ledger.New(channelContext)
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	block, err := client.QueryBlockByHash(blockHash)
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	block, err = client.QueryBlock(blockNumber)
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}


	//res, err := qscc.Invoke([][]byte([]byte(GetChainInfo), []byte(myc)))
	res, err := qscc.invoke("GetChainInfo", "myc")
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	

	response, err := chClient.Query(chClient.Request{ChaincodeID: "qscc", Fcn: "invoke", Args: integration.ExampleCCQueryArgs("GetChainInfo")})

	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	productLogger.Infof("PASS THE TSTs")
	productLogger.Infof(response)*/

	return shim.Success([]byte("Reply from TestQueryInfo"))
}