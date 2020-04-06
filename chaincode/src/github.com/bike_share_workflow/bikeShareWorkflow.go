package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	// "strings"
	// "time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type BikeShareWorkflowChaincode struct {
	devMode bool
}

func (t *BikeShareWorkflowChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("BikeShareWorkflow Initialization")
	return shim.Success(nil)
}

func (t *BikeShareWorkflowChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var creatorOrg, creatorCertIssuer string
	var err error

	fmt.Println("BikeShareWorkflow Invoke")

	if !t.devMode {
		creatorOrg, creatorCertIssuer, err = getTxCreatorInfo(stub)
		if err != nil {
			fmt.Errorf("Error extracting creator identity info: %s", err.Error())
			return shim.Error(err.Error())
		}
		fmt.Printf("BikeShareWorkflow invoked by '%s', '%s'.\n", creatorOrg, creatorCertIssuer)
	}

	function, args := stub.GetFunctionAndParameters()
	if function == "registerBike" {
		// Provider registers a bike
		return t.registerBike(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "discardBike" {
		// Provider discards a bike
		return t.discardBike(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "updateBikeLocation" {
		// Provider updates the location of a bike
		return t.updateBikeLocation(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getBikesByStatus" {
		// Get all bikes with specified status
		return t.getBikesByStatus(stub, creatorOrg, creatorCertIssuer, args)
	}

	return shim.Error("Invalid invoke function name.")
}

// Register a bike
func (t *BikeShareWorkflowChaincode) registerBike(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get the state from the ledger
	bikeKey, err := getBikeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) != 0 {
		err = errors.New(fmt.Sprintf("Bike %s already registered.", args[0]))
		return shim.Error(err.Error())
	}

	// Create bike object
	bike := &Bike{BIKE, args[0], []float32{}, BIKE_AVAILABLE}
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Bike %s registered.\n", args[0])

	return shim.Success(nil)
}

// Discard a bike
func (t *BikeShareWorkflowChaincode) discardBike(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var bike *Bike

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get the state from the ledger
	bikeKey, err := getBikeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) == 0 {
		err = errors.New(fmt.Sprintf("Bike %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bikeBytes, &bike)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify that the bike can be discarded
	if bike.Status == BIKE_AVAILABLE {
		bike.Status = BIKE_DISCARDED
		bikeBytes, err = json.Marshal(bike)
		if err != nil {
			return shim.Error("Error marshaling bike structure.")
		}
		// Write the state to the ledger
		err = stub.PutState(bikeKey, bikeBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("Bike %s discarded.\n", args[0])
	} else if bike.Status == BIKE_IN_USE {
		// fmt.Printf("Bike %s is in use.\n", args[0])
		// return shim.Error("Bike in use")
		err = errors.New(fmt.Sprintf("Bike %s in use.", args[0]))
		return shim.Error(err.Error())
	} else if bike.Status == BIKE_REPAIRING {
		// fmt.Printf("Bike %s is being repaired.\n", args[0])
		// eturn shim.Error("Bike repairing")
		err = errors.New(fmt.Sprintf("Bike %s repairing.", args[0]))
		return shim.Error(err.Error())
	} else if bike.Status == BIKE_DISCARDED {
		fmt.Printf("Bike %s already discarded.\n", args[0])
	}

	return shim.Success(nil)
}

// Update the location of a bike
func (t *BikeShareWorkflowChaincode) updateBikeLocation(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var bike *Bike

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 3 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {ID, Longitude, Latitude}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get the state from the ledger
	bikeKey, err := getBikeKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) == 0 {
		err = errors.New(fmt.Sprintf("Bike %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bikeBytes, &bike)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if bike is not discarded
	if bike.Status == BIKE_DISCARDED {
		// fmt.Printf("Bike %s is discarded.\n", args[0])
		// return shim.Error("Bike discarded.")
		err = errors.New(fmt.Sprintf("Bike %s already discarded.", args[0]))
		return shim.Error(err.Error())
	}

	// Parse longitude and latitude
	longitude, err := strconv.ParseFloat(string(args[1]), 8)
	if err != nil {
		return shim.Error(err.Error())
	}
	latitude, err := strconv.ParseFloat(string(args[2]), 8)
	if err != nil {
		return shim.Error(err.Error())
	}

	bike.Location = []float32{float32(longitude), float32(latitude)}
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("The location of bike %s updated.\n", args[0])

	return shim.Success(nil)
}

// Construct JSON array from a given query results iterator
func constructQueryResponseFromIterator(iterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	var queryResponseArray bytes.Buffer

	queryResponseArray.WriteString("[")
	isArrayMemberAlreadyWritten := false
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if isArrayMemberAlreadyWritten == true {
			queryResponseArray.WriteString(",")
		}
		queryResponseArray.WriteString("{\"Key\":")
		queryResponseArray.WriteString("\"")
		queryResponseArray.WriteString(queryResponse.Key)
		queryResponseArray.WriteString("\"")
		queryResponseArray.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		queryResponseArray.WriteString(string(queryResponse.Value))
		queryResponseArray.WriteString("}")
		isArrayMemberAlreadyWritten = true
	}
	queryResponseArray.WriteString("]")

	return &queryResponseArray, nil
}

func getQueryResponse(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("Query String:\n%s\n", queryString)

	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	queryResponse, err := constructQueryResponseFromIterator(iterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Query Result:\n%s\n", queryResponse.String())

	return queryResponse.Bytes(), nil
}

func (t *BikeShareWorkflowChaincode) getBikesByStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Provider/User/Repairer Org member can invoke this transaction
	if !t.devMode && !(authenticateProviderOrg(creatorOrg, creatorCertIssuer) || authenticateUserOrg(creatorOrg, creatorCertIssuer) || authenticateRepairerOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Status}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"status\":\"%s\"}}", BIKE, args[0])
	queryResponse, err := getQueryResponse(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResponse)
}

func main() {
	bswc := new(BikeShareWorkflowChaincode)
	bswc.devMode = true
	err := shim.Start(bswc)
	if err != nil {
		fmt.Printf("Error starting Bike Share Workflow chaincode: %s", err)
	}
}
