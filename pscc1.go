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

// Init is called during chaincode instantiation to initialize any data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("Instantiated chaincode!!!")
        
        // Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
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
	if len(args) != 2 {
			return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	
	//personalInfo := PersonalInfo()
	//jsonBytes, err := json.Marshal(personalInfo)

	var personalInfo PersonalInfo
 	personalInfo = PersonalInfo{"Papan","Das","1987-11-14","papan.das@zs.com","9641443962","ZS Associates India Pvt Ltd"}
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
