package main

import (
	"encoding/json"
    "fmt"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("main")

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

type PersonalInfo struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	DOB             string `json:"DOB"`
	Email           string `json:"email"`
	Mobile          string `json:"mobile"`
	Company         string `json:"company"`
}

/*type ProductInfo struct {
	SerialId       	string 			`json:"serialid"`
	Created        	string 			`json:"created"`
	Creator    		PersonalInfo	`json:"creator"`
	Name           	string 			`json:"name"`
	Expire          string 			`json:"expire"`
	GTIN         	string 			`json:"gtin"`
	LotNumber       string 			`json:"lotnumber"`
	Status         	string 			`json:"status"`
	TotalQty        int 			`json:"totalqty"`
	AvailQty        int 			`json:"availqty"`
}*/


// Init is called during chaincode instantiation to initialize any data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("Instantiated chaincode!!!")
	return shim.Success(nil)

}


// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The 'set'
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("Invoked chaincode!!!")

	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	
	var result string
	var err error
	if fn == "set" {
			result, err = set(stub, args)
	} else {
			result, err = get(stub, args)
	}
	if err != nil {
			return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}


// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 7 {
			return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	
	//personalInfo := PersonalInfo()
	//jsonBytes, err := json.Marshal(personalInfo)

	var personalInfo PersonalInfo
 	personalInfo = PersonalInfo{args[1], args[2], args[3], args[4], args[5], args[6]}
 	jsonBytes, err := json.Marshal (&personalInfo)
	

	err = stub.PutState(args[0], jsonBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
			return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
			return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
			return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}


// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
			fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
