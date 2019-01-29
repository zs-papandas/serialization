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

/*var totalPallet int = 2
var totalBox int = 2
var totalPacket int = 2
var totalItem int = 2*/



// GenerateProductContract : a struct to handle auto generate Product.
type GenerateProductContract struct {
}

func ForceSleep(){
	time.Sleep(3000 * time.Millisecond)
}

//CreateProduct : save a product
func (ac *GenerateProductContract) CreateProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	generateProductLogger.Infof("invoke Generate Product -> CreateProduct, args=%s\n", args)
	//generateProductLogger.Infof("invoke Generate Product length%s\n", len(args))

	if len(args) != 10 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. %s\n", args)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	var countPallet int = 0
	var countBox int = 0
	var countPacket int = 0
	var countItem int = 0



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
	var productType string

	
	identity = args[0]
	pname = args[1]
	expired = args[2]
	gtin = args[3]
	lotnum = args[4]
	status = "CREATED"
	amt = args[5]
	

	/*totalPallet, _ = strconv.Atoi(args[7])
    totalBox, _ = strconv.Atoi(args[8])
	totalPacket, _ = strconv.Atoi(args[9])
	totalItem, _ = strconv.Atoi(args[10])

	identity = "a"
	pname ="Product Name"
	expired = "Expired"
	gtin = "GTIN-23432"
	lotnum = "LOTNUM"
	status = "CREATED"
	amt = "10000"
	myStr = "10"
	productType = "pallet"*/
	
	str1 := args[6]
	str2 := args[7]
	str3 := args[8]
	str4 := args[9]
	totalPallet, err := strconv.Atoi(str1)
	if err != nil {
        errMsg := fmt.Sprintf("Failed: string to int. %s\n", str1)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	
    totalBox, err := strconv.Atoi(str2)
	if err != nil {
        errMsg := fmt.Sprintf("Failed: string to int. %s\n", str2)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	totalPacket, err := strconv.Atoi(str3)
	if err != nil {
        errMsg := fmt.Sprintf("Failed: string to int. %s\n", str3)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	totalItem, err := strconv.Atoi(str4)
	if err != nil {
        errMsg := fmt.Sprintf("Failed: string to int. %s\n", str4)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	currCat := 0
	
	/*myInt, err := strconv.Atoi(myStr)
	if err != nil {
        errMsg := fmt.Sprintf("Failed: string to int. %s\n", myStr)
		generateProductLogger.Error(errMsg)
		return shim.Error(errMsg)
    }*/
	totqty := totalBox
	avaiqty := totalBox


	//=================[ xGENESIS ]===================

	xno, err := utils.GetSerialNo(APIstub)
	if err != nil {
		generateProductLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	xtoday := time.Now().Format(time.RFC3339)
	xowner, err := utils.GetAccount(APIstub, identity)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			generateProductLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			generateProductLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}
	PalletArr = append(PalletArr, xno) 
	var xproductInfo models.Product
	xproductInfo = models.Product{xno, xtoday, *xowner, pname, expired, gtin, lotnum, status, amt, totqty, avaiqty, types.PalletProduct, ""}
	xjsonBytes, err := json.Marshal(&xproductInfo)
	if err != nil {
		generateProductLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	if err := APIstub.PutState(xno, xjsonBytes); err != nil {
		generateProductLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	//==============================================

	

	for {

		//===========================================

		if countPallet < totalPallet {
			//fmt.Println("Total Pallet", len(PalletArr))
			if len(PalletArr) == 0 {
				currCat++
			}else{
				//fmt.Println(" - Total Box", len(BoxArr))
				if len(BoxArr) == 0 {
					currCat++	
				}else{
					//fmt.Println("  - Total Packet", len(PacketArr))
					if len(PacketArr) == 0 {
						currCat++
					}else{
						//fmt.Println("    - Total Item", len(ItemArr))
						if len(ItemArr) == 0{
							currCat++
						}else{
							if len(ItemArr) == totalItem {
								//fmt.Println("Items LOADED")
								if len(PacketArr) == totalPacket {
									//fmt.Println("Packets LOADED")
									if len(BoxArr) == totalBox {
										//fmt.Println("Box LOADED")
	
										// RESET
										countBox=0
										countPacket=0
										countItem=0
	
										countPallet++
	
										currCat=0
	
										//PalletArr=nil
										BoxArr=nil
										PacketArr=nil
										ItemArr=nil
									}else{
										countItem=0
										countPacket=0
	
	
										countBox++
	
										currCat=1
										
										PacketArr=nil
										ItemArr=nil
									}
								}else{
									countItem=0
									countPacket++
									currCat=2
									ItemArr=nil
	
								}
	
							}else{
								
								countItem++
							}
						}
					}
				}
			
			}
			
		}else{
			//fmt.Printf("Pallet LOADED\n")
		}
	
		
	
		if countPallet < totalPallet {
			
			no, err := "SR"+time.Now().Format(time.RFC3339) 
			//utils.GetSerialNo(APIstub)
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
			
			parentProduct := ""
			
			// Get the Product Type
			var ProductTypeInt types.ProductType
			if currCat == 0 {
				PalletArr = append(PalletArr, no) 
				productType = "pallet"
				parentProduct = ""
				ProductTypeInt = types.PalletProduct
				totqty = totalBox
				avaiqty = totalBox

			}else if currCat == 1 {
				BoxArr = append(BoxArr, no)
				productType = "box"
				parentProduct = PalletArr[countPallet]
				ProductTypeInt = types.BoxProduct
				totqty = totalPacket
				avaiqty = totalPacket
				fmt.Printf("  ")
			}else if currCat == 2 {
				PacketArr = append(PacketArr, no)
				productType = "packet"
				parentProduct = BoxArr[countBox]
				ProductTypeInt = types.PacketProduct
				totqty = totalItem
				avaiqty = totalItem
				fmt.Printf("    ")
			}else if currCat == 3 {
				ItemArr = append(ItemArr, no)
				productType = "item"
				parentProduct = PacketArr[countPacket]
				ProductTypeInt = types.ItemProduct
				totqty = 0
				avaiqty = 0
				fmt.Printf("      ")
			}else {
				productType = "unknown"
				parentProduct = ""
				ProductTypeInt = types.UnKnownProduct
			}
			
			fmt.Println("->", productType, " SERIALNO: ", no, "PARENT", parentProduct)
			
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

			ForceSleep()
			//return shim.Success(jsonBytes)
			
		}else{
			PalletArr=nil
			BoxArr=nil
			PacketArr=nil
			ItemArr=nil

			countBox=0
			countPacket=0
			countItem=0
			countPallet=0

			totalPallet=0;
			totalBox=0;
			totalPacket=0;
			totalItem=0;

			currCat=0
			break
		}

		//============================================

		

		
		
	}
	fmt.Println("AUTO-PRODUCT GENERATED")
	return shim.Success([]byte("Auto generating process over."))
}




