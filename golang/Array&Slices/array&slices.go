package main

import (
	"fmt"
	"time"
	"strconv"
)

var totalPallet int = 2
var totalBox int = 2
var totalPacket int = 2
var totalItem int = 2

var countPallet int = 0
var countBox int = 0
var countPacket int = 0
var countItem int = 0

var currCat int = 0;
//var ArrColl = make(map[productCont])
var PalletArr []string
var BoxArr []string
var PacketArr []string
var ItemArr []string

type productCont struct{
	ProductId string
	ProductType string
	ProductObj []string
}

func AddProduct() {
	fmt.Println("======================================")
	fmt.Println("Category:", currCat)
}

func ForceSleep(){
	time.Sleep(1000 * time.Millisecond)
}

func shortCount(param int){
	
	if param == 0 {
		uid := "P"+strconv.Itoa(countPallet)
		PalletArr = append(PalletArr, uid) 
	}else if param == 1 {
		uid := "P"+strconv.Itoa(countPallet)+"-B"+strconv.Itoa(countBox)
		BoxArr = append(BoxArr, uid)
	}else if param == 2 {
		uid := "P"+strconv.Itoa(countPallet)+"-B"+strconv.Itoa(countBox)+"-P"+strconv.Itoa(countPacket)
		PacketArr = append(PacketArr, uid)
	}else{
		uid := "P"+strconv.Itoa(countPallet)+"-B"+strconv.Itoa(countBox)+"-P"+strconv.Itoa(countPacket)+"-I"+strconv.Itoa(countItem)
		ItemArr = append(ItemArr, uid)
	}
}

func initProductChain(){

	AddProduct()
	shortCount(currCat)
	//fmt.Println("Category:", currCat)
	//fmt.Println("Pallet:", countPallet,", Box:", countBox,", Packet:", countPacket,", Item:", countItem)
	

	if countPallet < totalPallet {
		fmt.Println("Total Pallet", len(PalletArr))
		if len(PalletArr) == 0 {
			currCat++
		}else{
			fmt.Println("Total Box", len(BoxArr))
			if len(BoxArr) == 0 {
				currCat++	
			}else{
				fmt.Println("Total Packet", len(PacketArr))
				if len(PacketArr) == 0 {
					currCat++
				}else{
					fmt.Println("Total Item", len(ItemArr))
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
		fmt.Printf("Pallet LOADED\n")
		

	}

	ForceSleep()

	if countPallet < totalPallet {
		initProductChain()
	}else{
		fmt.Println(PalletArr)
		fmt.Println(BoxArr)
		fmt.Println(PacketArr)
		fmt.Println(ItemArr)
	}
}

func main(){

	initProductChain();
	
	//fmt.Printf("PalletItems %v\n", palletItem)

	//fmt.Printf("BoxItem %v\n", boxItem)

	//fmt.Printf("AllItems %v\n", AllItems)
}