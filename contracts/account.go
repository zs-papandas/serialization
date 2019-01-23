package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/zs-papandas/serialization/models"
	"github.com/zs-papandas/serialization/types"
)

var accountLogger = shim.NewLogger("contracts/account")

// AccountContract : a struct to handle Account.
type AccountContract struct {
}

// ListAccount : return a list of all accounts.
func (ac *AccountContract) ListAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke ListAccount, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = no argument, Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	var UserTypeInt types.UserType
	switch args[0] {
	case "manufacturer":
		UserTypeInt = types.ManufacturerUser
	case "wholesaler":
		UserTypeInt = types.WholesalerUser
	case "retailer":
		UserTypeInt = types.RetailerUser
	case "patient":
		UserTypeInt = types.PatientUser
	default:
		UserTypeInt = types.UnKnownUser
	}

	accountLogger.Infof("UserTypeInt string = '%s'", string(UserTypeInt))

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"email": "papan.das@zs.com",
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	accountLogger.Infof("Query string = '%s'", string(queryBytes))
	resultsIterator, err := APIstub.GetQueryResult(string(queryBytes))
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := make([]*models.Account, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		account := new(models.Account)
		accountLogger.Infof("Query Response Key = '%s'", string(queryResponse.Key))
		accountLogger.Infof("Query Response Value = '%s'", string(queryResponse.Value))
		//accountLogger.Infof("Query Response account = '%s'", account)
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
	
	/*return shim.Success([]byte("Reply from ListAccount"))*/
}


// ListAccount : return a list of all accounts.
func (ac *AccountContract) CreateAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke CreateAccount, args=%s\n", args)

	if len(args) != 8 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['name'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	var UserTypeInt types.UserType
	switch args[7] {
	case "manufacturer":
		UserTypeInt = types.ManufacturerUser
	case "wholesaler":
		UserTypeInt = types.WholesalerUser
	case "retailer":
		UserTypeInt = types.RetailerUser
	case "patient":
		UserTypeInt = types.PatientUser
	default:
		UserTypeInt = types.UnKnownUser
	}

	var personalInfo models.Account
 	personalInfo = models.Account{args[1], args[2], args[3], args[4], args[5], args[6], UserTypeInt}
 	jsonBytes, err := json.Marshal(&personalInfo)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	if err := APIstub.PutState(args[0], jsonBytes); err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}


// ListAccount : return a list of all accounts.
func (ac *AccountContract) RetrieveAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	accountLogger.Infof("invoke RetrieveAccount, args=%s\n", args)

	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	account, err := APIstub.GetState(args[0])
	if err != nil {
		errMsg1 := fmt.Sprintf("Failed to get asset: %s with error: %s", args[0], err)
		accountLogger.Error(errMsg1)
		return shim.Error(errMsg1)
	}

	return shim.Success(account)
	//return shim.Success([]byte("Reply from RetrieveAccount"))
}