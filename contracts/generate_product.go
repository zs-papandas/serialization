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

var generateProductLogger = shim.NewLogger("contracts/generate_product")

var totalPallet int = 2
var totalBox int = 2
var totalPacket int = 2
var totalItem int = 2

var countPallet int = 0
var countBox int = 0
var countPacket int = 0
var countItem int = 0

var currCat int = 0;

var PalletArr []string
var BoxArr []string
var PacketArr []string
var ItemArr []string

var identity string
var pname string
var expired string
var gtin string
var lotnum string
var status string
var amt string
var myStr string
var productType string

// GenerateProductContract : a struct to handle auto generate Product.
type GenerateProductContract struct {
}


//CreateProduct : save a product
func (ac *GenerateProductContract) CreateProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	generateProductLogger.Infof("invoke Generate Product -> CreateProduct, args=%s\n", args)
	generateProductLogger.Infof("invoke Generate Product length%s\n", len(args))

	/*if len(args) != 9 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. %s\n", args)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
	}*/

	
	identity = args[0]
	pname = args[1]
	expired = args[2]
	gtin = args[3]
	lotnum = args[4]
	status = "CREATED"
	amt = args[5]
	myStr = args[6]
	productType = args[7]

	totalPallet, _ = strconv.Atoi(args[8])
    totalBox, _ = strconv.Atoi(args[9])
	totalPacket, _ = strconv.Atoi(args[10])
	totalItem, _ = strconv.Atoi(args[11])

	


	myInt, err := strconv.Atoi(myStr)
	if err != nil {
        errMsg := fmt.Sprintf("Failed: string to int. %s\n", myStr)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
    }
	totqty := myInt
	avaiqty := myInt

	i := 0

	for i < 10 {

		no, err := utils.GetSerialNo(APIstub)
		if err != nil {
			generateProductLogger.Error(err.Error())
			return shim.Error(err.Error())
			break
		}
		today := time.Now().Format(time.RFC3339)
		
		// GET USER ACCOUNT DETAIL
		owner, err := utils.GetAccount(APIstub, identity)
		if err != nil {
			switch e := err.(type) {
			case *utils.WarningResult:
				generateProductLogger.Warning(err.Error())
				return shim.Success(e.JSONBytes())
				break
			default:
				generateProductLogger.Error(err.Error())
				return shim.Error(err.Error())
				break
			}
		}

		//generateProductLogger.Infof("User Account %s\n", owner.Firstname)
		

		// Get the Product Type
		var ProductTypeInt types.ProductType
		switch productType {
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

		parentProduct := ""
		//"WIXnkuHMYZL5fGaE"
		

		var productInfo models.Product
		productInfo = models.Product{no, today, *owner, pname, expired, gtin, lotnum, status, amt, totqty, avaiqty, ProductTypeInt, parentProduct}
		jsonBytes, err := json.Marshal(&productInfo)
		if err != nil {
			generateProductLogger.Error(err.Error())
			return shim.Error(err.Error())
			break
		}

		if err := APIstub.PutState(no, jsonBytes); err != nil {
			generateProductLogger.Error(err.Error())
			return shim.Error(err.Error())
			break
		}

		generateProductLogger.Infof("productType %s\n", productType, " - ", no)
		//return shim.Success(jsonBytes)

		i++
	}

	
	return shim.Success([]byte("Auto generating process over."))
}




