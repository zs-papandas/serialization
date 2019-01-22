package main

import (
    "fmt"

    "github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	
	"github.com/zs-papandas/serialization/contracts"
)

var logger = shim.NewLogger("main")

var accountContract = new(contracts.AccountContract)

// EntryPoint implements a simple chaincode to manage an asset
type EntryPoint struct {
}


// Init is called during chaincode instantiation to initialize any data.
func (t *EntryPoint) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Instantiated chaincode!!!")
	return shim.Success(nil)
}


// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The 'set'
// method may create a new asset by specifying a new key-value pair.
func (t *EntryPoint) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Invoked chaincode!!!")

	function, args := APIstub.GetFunctionAndParameters()

	logger.Info("Query Function: %s\n", function)

	var result string
	var err error
	switch function {
		case "listAccount":
			return accountContract.ListAccount(APIstub, args)
		case "createAccount":
			return accountContract.CreateAccount(APIstub, args)
		case "retrieveAccount":
			return accountContract.RetrieveAccount(APIstub, args)
	}

	msg := fmt.Sprintf("No such function. function = %s, args = %s", function, args)
	logger.Error(msg)
	return shim.Error(msg)
}

// Get returns the value of the specified asset key
/*func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
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
}*/


// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(EntryPoint)); err != nil {
		logger.Errorf("Error creating new Chaincode. Error = %s\n", err)
	}
}
