package contracts

import (
	"encoding/json"
	"time"
	"fmt"
	"strconv"
	//"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

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

	if toProduct.AvailQty == 0 && toProduct.ProductType != types.ItemProduct {
		productLogger.Error("Inventory/stock over")
		return shim.Error("Inventory/stock over")
	}

	if toProduct.ProductType != types.PalletProduct {

		if toProduct.ProductType == types.ItemProduct && toProduct.Status != "CREATED" {
			productLogger.Error("Item not available. Sold Out.")
			return shim.Error("Item not available. Sold Out.")
		}

		fromProduct, err := utils.GetProduct(APIstub, toProduct.ParentProduct)
		if err != nil {
			switch e := err.(type) {
			case *utils.WarningResult:
				productLogger.Warning(err.Error())
				productLogger.Warning(e.JSONBytes())
			default:
				productLogger.Error(err.Error())
			}
		}

		/* 
		If current product inventory is zero, update the parent about it
		Packet LEVEL 
		*/

		if fromProduct.AvailQty > 0 {

			fromProduct.AvailQty = fromProduct.AvailQty - 1
			fromProduct.Status = "ITEM_SOLD"
	
			fromProductBytes, err := json.Marshal(fromProduct)
			if err != nil {
				productLogger.Error(err.Error())
			}
			if err := APIstub.PutState(fromProduct.SerialId, fromProductBytes); err != nil {
				productLogger.Error(err.Error())
			}

			/*
			if stock available is 0
			product type is Packet
			*/

			if fromProduct.AvailQty == 0 && fromProduct.ProductType == types.PacketProduct {

				fmt.Println("Packet Product has Zero Inventory Available.")

				fromfromProduct, err := utils.GetProduct(APIstub, fromProduct.ParentProduct)
				if err != nil {
					switch e := err.(type) {
					case *utils.WarningResult:
						productLogger.Warning(err.Error())
						productLogger.Warning(e.JSONBytes())
					default:
						productLogger.Error(err.Error())
					}
				}

				fromfromProduct.AvailQty = fromfromProduct.AvailQty - 1


				fromfromProductBytes, err := json.Marshal(fromfromProduct)
				if err != nil {
					productLogger.Error(err.Error())
				}
				if err := APIstub.PutState(fromfromProduct.SerialId, fromfromProductBytes); err != nil {
					productLogger.Error(err.Error())
				}

				

				if(fromfromProduct.AvailQty == 0){
					fmt.Println("Update Box Product about Packet Producted inventory running zero")

					fromfromfromProduct, err := utils.GetProduct(APIstub, fromfromProduct.ParentProduct)
					if err != nil {
						switch e := err.(type) {
						case *utils.WarningResult:
							productLogger.Warning(err.Error())
							productLogger.Warning(e.JSONBytes())
						default:
							productLogger.Error(err.Error())
						}
					}
	
					fromfromfromProduct.AvailQty = fromfromfromProduct.AvailQty - 1
	
	
					fromfromfromProductBytes, err := json.Marshal(fromfromfromProduct)
					if err != nil {
						productLogger.Error(err.Error())
					}
					if err := APIstub.PutState(fromfromfromProduct.SerialId, fromfromfromProductBytes); err != nil {
						productLogger.Error(err.Error())
					}

					if(fromfromfromProduct.AvailQty == 0){
						fmt.Println("Update Pallet Proudct when Box Product Inventory found zero")

						fromfromfromfromProduct, err := utils.GetProduct(APIstub, fromfromfromProduct.ParentProduct)
						if err != nil {
							switch e := err.(type) {
							case *utils.WarningResult:
								productLogger.Warning(err.Error())
								productLogger.Warning(e.JSONBytes())
							default:
								productLogger.Error(err.Error())
							}
						}
		
						fromfromfromfromProduct.AvailQty = fromfromfromfromProduct.AvailQty - 1
		
		
						fromfromfromfromProductBytes, err := json.Marshal(fromfromfromfromProduct)
						if err != nil {
							productLogger.Error(err.Error())
						}
						if err := APIstub.PutState(fromfromfromfromProduct.SerialId, fromfromfromfromProductBytes); err != nil {
							productLogger.Error(err.Error())
						}

					}
				}
			}  

			/*
			if stock available is 0
			product type is Box
			*/

			if fromProduct.AvailQty == 0 && fromProduct.ProductType == types.BoxProduct {

				fmt.Println("Box Product has Zero Inventory Available.")

				fromfromProduct, err := utils.GetProduct(APIstub, fromProduct.ParentProduct)
				if err != nil {
					switch e := err.(type) {
					case *utils.WarningResult:
						productLogger.Warning(err.Error())
						productLogger.Warning(e.JSONBytes())
					default:
						productLogger.Error(err.Error())
					}
				}

				fromfromProduct.AvailQty = fromfromProduct.AvailQty - 1


				fromfromProductBytes, err := json.Marshal(fromfromProduct)
				if err != nil {
					productLogger.Error(err.Error())
				}
				if err := APIstub.PutState(fromfromProduct.SerialId, fromfromProductBytes); err != nil {
					productLogger.Error(err.Error())
				}

				

				if(fromfromProduct.AvailQty == 0){

					fmt.Println("Update Pallet Product about Box Producted inventory running zero")

					fromfromfromProduct, err := utils.GetProduct(APIstub, fromfromProduct.ParentProduct)
					if err != nil {
						switch e := err.(type) {
						case *utils.WarningResult:
							productLogger.Warning(err.Error())
							productLogger.Warning(e.JSONBytes())
						default:
							productLogger.Error(err.Error())
						}
					}
	
					fromfromfromProduct.AvailQty = fromfromfromProduct.AvailQty - 1
	
	
					fromfromfromProductBytes, err := json.Marshal(fromfromfromProduct)
					if err != nil {
						productLogger.Error(err.Error())
					}
					if err := APIstub.PutState(fromfromfromProduct.SerialId, fromfromfromProductBytes); err != nil {
						productLogger.Error(err.Error())
					}

					
				}
			}

			


		}else{
			productLogger.Error("Inventory is Empty. Sold Out.")
			return shim.Error("Inventory is Empty. Sold Out.")
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

	productLogger.Infof("toProduct.ProductType, args=%s\n", toProduct.ProductType)

	if toProduct.ProductType == types.PalletProduct {
		
		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"product_type": 2,
			},
		}
	
		queryBytes, err := json.Marshal(query)
		if err != nil {
			productLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		productLogger.Infof("Query string = '%s'", string(queryBytes))
		resultsIterator, err := APIstub.GetQueryResult(string(queryBytes))
		
		if err != nil {
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		defer resultsIterator.Close()
	
		results := make([]*models.Product, 0)
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				accountLogger.Error(err.Error())
				return shim.Error(err.Error())
			}
			productLogger.Infof("=================================")
			productLogger.Infof("Query Response = '%s'", string(queryResponse.Value))
			fmt.Println(queryResponse.Value)
		}


	}else if toProduct.ProductType == types.BoxProduct {

	}else if  toProduct.ProductType == types.PacketProduct {

	}else{

	}

	//=

	return shim.Success(toProductBytes)
		
	

	
	
}


/*func mockChannelProvider(channelID string) context.ChannelProvider {

	channelProvider := func() (context.Channel, error) {
		return mocks.NewMockChannel(channelID)
	}

	return channelProvider
}*/


// '{"selector":{"product_type":1}}'
// ChangeOwner : Change owner of a product.
func (ac *ProductContract) TestQueryInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productLogger.Infof("invoke TestQueryInfo, args=%s\n", args)

	//queryString := fmt.Sprintf("{\"selector\":{\"lastname\":\"Harry\",\"owner\":\"%s\"}}", owner)
	/*queryString := fmt.Sprintf("{\"selector\":{\"lastname\":\"Harry\"}}")

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)*/

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"product_type": 2,
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		productLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	productLogger.Infof("Query string = '%s'", string(queryBytes))
	resultsIterator, err := APIstub.GetQueryResult(string(queryBytes))
	
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := make([]*models.Product, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		account := new(models.Product)
		if err := json.Unmarshal(queryResponse.Value, account); err != nil {
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		results = append(results, account)
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)

	return shim.Success([]byte("Reply from TestQueryInfo"))
}


/*func getQueryResultForQueryString(APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}*/