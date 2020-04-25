package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	// "strings"
	"time"

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
	if function == "registerUser" {
		// User registers a user
		return t.registerUser(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "registerRepairer" {
		// Repairer registers a repairer
		return t.registerRepairer(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "registerBike" {
		// Provider registers a bike
		return t.registerBike(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "reactivateBike" {
		// Provider reactivates a bike
		return t.reactivateBike(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "discardBike" {
		// Provider discards a bike
		return t.discardBike(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "updateBikeLocation" {
		// Provider updates the location of a bike
		return t.updateBikeLocation(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "startRide" {
		// User starts a ride
		return t.startRide(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "endRide" {
		// User ends a ride
		return t.endRide(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "reportIssue" {
		// User reports an issue
		return t.reportIssue(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "acceptIssue" {
		// Provider accepts an issue
		return t.acceptIssue(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "rejectIssue" {
		// Provider rejects an issue
		return t.rejectIssue(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "requestRepair" {
		// Provider requests a repair
		return t.requestRepair(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "acceptRepair" {
		// Repairer accepts a repair
		return t.acceptRepair(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "rejectRepair" {
		// Repairer rejects a repair
		return t.rejectRepair(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "completeRepair" {
		// Repairer completes a repair
		return t.completeRepair(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getBikesByStatus" {
		// Provider/User/Repairer gets all bikes with specified status
		return t.getBikesByStatus(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getRidesByStatus" {
		// Provider/User gets all rides with specified status
		return t.getRidesByStatus(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getIssuesByStatus" {
		// Provider/User gets all issues with specified status
		return t.getIssuesByStatus(stub, creatorOrg, creatorCertIssuer, args)
	} else if function == "getRepairsByStatus" {
		// Provider/Repairer gets all repairs with specified status
		return t.getRepairsByStatus(stub, creatorOrg, creatorCertIssuer, args)
	} 

	return shim.Error("Invalid invoke function name.")
}

// Register a user
func (t *BikeShareWorkflowChaincode) registerUser(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a User Org member can invoke this transaction
	if !t.devMode && !authenticateUserOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of User Org. Access denied.")
	}

	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {User ID, Balance}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get user state from the ledger
	userKey, err := getUserKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	userBytes, err := stub.GetState(userKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(userBytes) != 0 {
		err = errors.New(fmt.Sprintf("User %s already registered.", args[0]))
		return shim.Error(err.Error())
	}

	// Parse balance
	balance, err := strconv.ParseFloat(string(args[1]), 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create user object
	user := &User{USER, args[0], float32(balance), "", USER_FREE}
	userBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error("Error marshaling user structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("User %s registered.\n", args[0])

	return shim.Success(nil)
}

// Register a repairer
func (t *BikeShareWorkflowChaincode) registerRepairer(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Repairer Org member can invoke this transaction
	if !t.devMode && !authenticateRepairerOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Repairer Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Repairer ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get repairer state from the ledger
	repairerKey, err := getRepairerKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairerBytes, err := stub.GetState(repairerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairerBytes) != 0 {
		err = errors.New(fmt.Sprintf("Repairer %s already registered.", args[0]))
		return shim.Error(err.Error())
	}

	// Create user object
	repairer := &Repairer{REPAIRER, args[0]}
	repairerBytes, err = json.Marshal(repairer)
	if err != nil {
		return shim.Error("Error marshaling repairer structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(repairerKey, repairerBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Repairer %s registered.\n", args[0])

	return shim.Success(nil)
}

// Register a bike
func (t *BikeShareWorkflowChaincode) registerBike(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Bike ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
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

// Reactivate a bike
func (t *BikeShareWorkflowChaincode) reactivateBike(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var bike *Bike

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Bike ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
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

	// Verify if bike can be reactivated
	if bike.Status == BIKE_DISCARDED {
		err = errors.New(fmt.Sprintf("Bike %s discarded.", args[0]))
		return shim.Error(err.Error())
	} else if bike.Status == BIKE_REPAIRING {
		err = errors.New(fmt.Sprintf("Bike %s repairing.", args[0]))
		return shim.Error(err.Error())
	} else if bike.Status != BIKE_TO_REPAIR && bike.Status != BIKE_REPAIRED {
		err = errors.New(fmt.Sprintf("Bike %s active.", args[0]))
		return shim.Error(err.Error())
	}

	bike.Status = BIKE_AVAILABLE
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}
	
	// Write the state to the ledger
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Bike %s reactivated.\n", args[0])

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
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Bike ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
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

	// Verify if bike is available
	if bike.Status != BIKE_AVAILABLE {
		err = errors.New(fmt.Sprintf("Bike %s not available.", args[0]))
		return shim.Error(err.Error())
	}

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
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Bike ID, Longitude, Latitude}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
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

// Start a ride
func (t *BikeShareWorkflowChaincode) startRide(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var user *User
	var bike *Bike

	// Access control: Only a User Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of User Org. Access denied.")
	}

	if len(args) != 4 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 4: {User ID, Bike ID, Longitude, Latitude}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get user state from the ledger
	userKey, err := getUserKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	userBytes, err := stub.GetState(userKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(userBytes) == 0 {
		err = errors.New(fmt.Sprintf("User %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if user is free
	if user.Status != USER_FREE {
		err = errors.New(fmt.Sprintf("User %s has another ongoing ride.", args[0]))
		return shim.Error(err.Error())
	}

	// Verify if user has positive balance
	if user.Balance <= 0 {
		err = errors.New(fmt.Sprintf("User %s has negative balance.", args[0]))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
	bikeKey, err := getBikeKey(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) == 0 {
		err = errors.New(fmt.Sprintf("Bike %s not found.", args[1]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bikeBytes, &bike)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if bike is available
	if bike.Status != BIKE_AVAILABLE {
		err = errors.New(fmt.Sprintf("Bike %s not available.", args[1]))
		return shim.Error(err.Error())
	}

	// Parse longitude and latitude
	longitude, err := strconv.ParseFloat(string(args[2]), 8)
	if err != nil {
		return shim.Error(err.Error())
	}
	latitude, err := strconv.ParseFloat(string(args[3]), 8)
	if err != nil {
		return shim.Error(err.Error())
	}

	startTime := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	rideId := fmt.Sprintf("%s-%s-%s", args[0], args[1], startTime)

	// Get ride state from the ledger
	rideKey, err := getRideKey(stub, rideId)
	if err != nil {
		return shim.Error(err.Error())
	}
	rideBytes, err := stub.GetState(rideKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(rideBytes) != 0 {
		err = errors.New(fmt.Sprintf("Ride %s already started.", rideId))
		return shim.Error(err.Error())
	}

	// Create ride object
	ride := &Ride{RIDE, rideId, args[0], args[1], startTime, []float32{float32(longitude), float32(latitude)}, "", []float32{}, 0, RIDE_ONGOING}
	rideBytes, err = json.Marshal(ride)
	if err != nil {
		return shim.Error("Error marshaling ride structure.")
	}

	user.RideId = rideId
	user.Status = USER_IN_RIDE
	userBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error("Error marshaling user structure.")
	}

	bike.Location = []float32{float32(longitude), float32(latitude)}
	bike.Status = BIKE_IN_USE
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(rideKey, rideBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Ride %s started.\n", rideId)

	return shim.Success(nil)
}

// End a ride
func (t *BikeShareWorkflowChaincode) endRide(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var user *User
	var bike *Bike
	var ride *Ride

	// Access control: Only a User Org member can invoke this transaction
	if !t.devMode && !authenticateUserOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of User Org. Access denied.")
	}

	if len(args) != 4 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 4: {User ID, Bike ID, Longitude, Latitude}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get user state from the ledger
	userKey, err := getUserKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	userBytes, err := stub.GetState(userKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(userBytes) == 0 {
		err = errors.New(fmt.Sprintf("User %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if user is in a ride
	if user.Status != USER_IN_RIDE {
		err = errors.New(fmt.Sprintf("User %s doesn't have an ongoing ride.", args[0]))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
	bikeKey, err := getBikeKey(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) == 0 {
		err = errors.New(fmt.Sprintf("Bike %s not found.", args[1]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bikeBytes, &bike)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if bike is in use
	if bike.Status != BIKE_IN_USE {
		err = errors.New(fmt.Sprintf("Bike %s not in use.", args[1]))
		return shim.Error(err.Error())
	}

	// Get ride state from the ledger
	rideKey, err := getRideKey(stub, user.RideId)
	if err != nil {
		return shim.Error(err.Error())
	}
	rideBytes, err := stub.GetState(rideKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(rideBytes) == 0 {
		err = errors.New(fmt.Sprintf("Ride %s not found.", user.RideId))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(rideBytes, &ride)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if bike ID matches
	if ride.BikeId != args[1] {
		err = errors.New(fmt.Sprintf("Actual bike %s and requested bike %s not match.", ride.BikeId, args[1]))
		return shim.Error(err.Error())
	}

	// Verify if ride is ongoing
	if ride.Status != RIDE_ONGOING {
		err = errors.New(fmt.Sprintf("Ride %s not ongoing.", user.RideId))
		return shim.Error(err.Error())
	}

	// Parse longitude and latitude
	longitude, err := strconv.ParseFloat(string(args[2]), 8)
	if err != nil {
		return shim.Error(err.Error())
	}
	latitude, err := strconv.ParseFloat(string(args[3]), 8)
	if err != nil {
		return shim.Error(err.Error())
	}

	startTimeInt, err := strconv.ParseInt(ride.StartTime, 10, 64)
    if err != nil {
        return shim.Error(err.Error())
    }
    startTime := time.Unix(startTimeInt, 0)
	endTime := time.Now()
	duration := endTime.Sub(startTime).Minutes()
	cost := float32(duration) * 0.1
	ride.EndTime = strconv.FormatInt(endTime.UTC().Unix(), 10)
	ride.EndLocation = []float32{float32(longitude), float32(latitude)}
	ride.Cost = cost
	ride.Status = RIDE_COMPLETED
	rideBytes, err = json.Marshal(ride)
	if err != nil {
		return shim.Error("Error marshaling ride structure.")
	}
	
	bike.Location = []float32{float32(longitude), float32(latitude)}
	bike.Status = BIKE_AVAILABLE
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}
	
	user.Balance -= cost
	user.Status = USER_FREE
	userBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error("Error marshaling user structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(rideKey, rideBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Ride %s ended.\n", user.RideId)

	return shim.Success(nil)
}

// Report an issue
func (t *BikeShareWorkflowChaincode) reportIssue(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var ride *Ride

	// Access control: Only a User Org member can invoke this transaction
	if !t.devMode && !authenticateUserOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of User Org. Access denied.")
	}

	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {User ID, Ride ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get ride state from the ledger
	rideKey, err := getRideKey(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	rideBytes, err := stub.GetState(rideKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(rideBytes) == 0 {
		err = errors.New(fmt.Sprintf("Ride %s not found.", args[1]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(rideBytes, &ride)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if user matches
	if ride.UserId != args[0] {
		err = errors.New(fmt.Sprintf("Actual user %s and requested user %s not match.", ride.UserId, args[0]))
		return shim.Error(err.Error())
	}

	// Verify if ride is completed
	if ride.Status != RIDE_COMPLETED {
		err = errors.New(fmt.Sprintf("Ride %s not completed.", args[1]))
		return shim.Error(err.Error())
	}

	issueId := fmt.Sprintf("%s-%s", args[1], strconv.FormatInt(time.Now().UTC().Unix(), 10))

	// Get issue state from the ledger
	issueKey, err := getIssueKey(stub, issueId)
	if err != nil {
		return shim.Error(err.Error())
	}
	issueBytes, err := stub.GetState(issueKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(issueBytes) != 0 {
		err = errors.New(fmt.Sprintf("Issue %s already opened.", issueId))
		return shim.Error(err.Error())
	}

	// Create issue object
	issue := &Issue{ISSUE, issueId, args[0], ride.BikeId, args[1], ISSUE_OPEN}
	issueBytes, err = json.Marshal(issue)
	if err != nil {
		return shim.Error("Error marshaling issue structure.")
	}

	ride.Status = RIDE_ISSUE_OPEN
	rideBytes, err = json.Marshal(ride)
	if err != nil {
		return shim.Error("Error marshaling ride structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(issueKey, issueBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(rideKey, rideBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Issue %s opened.\n", issueId)

	return shim.Success(nil)
}

// Accept an issue
func (t *BikeShareWorkflowChaincode) acceptIssue(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var user *User
	var ride *Ride
	var issue *Issue

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Issue ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get issue state from the ledger
	issueKey, err := getIssueKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	issueBytes, err := stub.GetState(issueKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(issueBytes) == 0 {
		err = errors.New(fmt.Sprintf("Issue %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(issueBytes, &issue)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if issue is open
	if issue.Status != ISSUE_OPEN {
		err = errors.New(fmt.Sprintf("Issue %s not open.", args[0]))
		return shim.Error(err.Error())
	}

	// Get user state from the ledger
	userKey, err := getUserKey(stub, issue.UserId)
	if err != nil {
		return shim.Error(err.Error())
	}
	userBytes, err := stub.GetState(userKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(userBytes) == 0 {
		err = errors.New(fmt.Sprintf("User %s not found.", issue.UserId))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Get ride state from the ledger
	rideKey, err := getRideKey(stub, issue.RideId)
	if err != nil {
		return shim.Error(err.Error())
	}
	rideBytes, err := stub.GetState(rideKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(rideBytes) == 0 {
		err = errors.New(fmt.Sprintf("Ride %s not found.", issue.RideId))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(rideBytes, &ride)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if issue is open
	if ride.Status != RIDE_ISSUE_OPEN {
		err = errors.New(fmt.Sprintf("Ride %s not associated with an issue.", issue.RideId))
		return shim.Error(err.Error())
	}

	issue.Status = ISSUE_CLOSED
	issueBytes, err = json.Marshal(issue)
	if err != nil {
		return shim.Error("Error marshaling issue structure.")
	}

	user.Balance += ride.Cost
	userBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error("Error marshaling user structure.")
	}

	ride.Cost = 0
	ride.Status = RIDE_ISSUE_CLOSED
	rideBytes, err = json.Marshal(ride)
	if err != nil {
		return shim.Error("Error marshaling ride structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(issueKey, issueBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(rideKey, rideBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Issue %s accepted.\n", args[0])

	return shim.Success(nil)
}

// Reject an issue
func (t *BikeShareWorkflowChaincode) rejectIssue(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var ride *Ride
	var issue *Issue

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Issue ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get issue state from the ledger
	issueKey, err := getIssueKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	issueBytes, err := stub.GetState(issueKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(issueBytes) == 0 {
		err = errors.New(fmt.Sprintf("Issue %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(issueBytes, &issue)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if issue is open
	if issue.Status != ISSUE_OPEN {
		err = errors.New(fmt.Sprintf("Issue %s not open.", args[0]))
		return shim.Error(err.Error())
	}

	// Get ride state from the ledger
	rideKey, err := getRideKey(stub, issue.RideId)
	if err != nil {
		return shim.Error(err.Error())
	}
	rideBytes, err := stub.GetState(rideKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(rideBytes) == 0 {
		err = errors.New(fmt.Sprintf("Ride %s not found.", issue.RideId))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(rideBytes, &ride)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if issue is open
	if ride.Status != RIDE_ISSUE_OPEN {
		err = errors.New(fmt.Sprintf("Ride %s not associated with an issue.", issue.RideId))
		return shim.Error(err.Error())
	}

	issue.Status = ISSUE_CLOSED
	issueBytes, err = json.Marshal(issue)
	if err != nil {
		return shim.Error("Error marshaling issue structure.")
	}

	ride.Status = RIDE_ISSUE_CLOSED
	rideBytes, err = json.Marshal(ride)
	if err != nil {
		return shim.Error("Error marshaling ride structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(issueKey, issueBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(rideKey, rideBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Issue %s rejected.\n", args[0])

	return shim.Success(nil)
}

// Request to repair a bike
func (t *BikeShareWorkflowChaincode) requestRepair(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var repairer *Repairer
	var bike *Bike

	// Access control: Only a Provider Org member can invoke this transaction
	if !t.devMode && !authenticateProviderOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Provider Org. Access denied.")
	}

	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {Bike ID, Repairer ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
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

	// Verify if bike is available
	if bike.Status != BIKE_AVAILABLE {
		err = errors.New(fmt.Sprintf("Bike %s not available.", args[0]))
		return shim.Error(err.Error())
	}

	// Get repairer state from the ledger
	repairerKey, err := getRepairerKey(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairerBytes, err := stub.GetState(repairerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairerBytes) == 0 {
		err = errors.New(fmt.Sprintf("Repairer %s not found.", args[1]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(repairerBytes, &repairer)
	if err != nil {
		return shim.Error(err.Error())
	}

	repairId := fmt.Sprintf("%s-%s-%s", args[0], args[1], strconv.FormatInt(time.Now().UTC().Unix(), 10))

	// Get repair state from the ledger
	repairKey, err := getRepairKey(stub, repairId)
	if err != nil {
		return shim.Error(err.Error())
	}
	repairBytes, err := stub.GetState(repairKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairBytes) != 0 {
		err = errors.New(fmt.Sprintf("Repair %s already requested.", repairId))
		return shim.Error(err.Error())
	}

	// Create repair object
	repair := &Repair{REPAIR, repairId, args[0], args[1], REPAIR_REQUESTED}
	repairBytes, err = json.Marshal(repair)
	if err != nil {
		return shim.Error("Error marshaling repair structure.")
	}

	bike.Status = BIKE_TO_REPAIR
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(repairKey, repairBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Repair %s requested.\n", repairId)

	return shim.Success(nil)
}

// Accept the request to repair a bike
func (t *BikeShareWorkflowChaincode) acceptRepair(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var repairer *Repairer
	var bike *Bike
	var repair *Repair

	// Access control: Only a Repairer Org member can invoke this transaction
	if !t.devMode && !authenticateRepairerOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Repairer Org. Access denied.")
	}

	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {Repairer ID, Repair ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get repairer state from the ledger
	repairerKey, err := getRepairerKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairerBytes, err := stub.GetState(repairerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairerBytes) == 0 {
		err = errors.New(fmt.Sprintf("Repairer %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(repairerBytes, &repairer)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Get repair state from the ledger
	repairKey, err := getRepairKey(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairBytes, err := stub.GetState(repairKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairBytes) == 0 {
		err = errors.New(fmt.Sprintf("Repair %s not found.", args[1]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(repairBytes, &repair)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if repairer matches
	if repair.RepairerId != args[0] {
		err = errors.New(fmt.Sprintf("Actual repairer %s and requested repairer %s not match.", repair.RepairerId, args[0]))
		return shim.Error(err.Error())
	}

	// Verify if repair is requested
	if repair.Status != REPAIR_REQUESTED {
		err = errors.New(fmt.Sprintf("Repair %s already processed.", args[1]))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
	bikeKey, err := getBikeKey(stub, repair.BikeId)
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) == 0 {
		err = errors.New(fmt.Sprintf("Bike %s not found.", repair.BikeId))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bikeBytes, &bike)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if bike is ready to repair
	if bike.Status != BIKE_TO_REPAIR {
		err = errors.New(fmt.Sprintf("Bike %s not ready to repair.", repair.BikeId))
		return shim.Error(err.Error())
	}

	repair.Status = REPAIR_ACCEPTED
	repairBytes, err = json.Marshal(repair)
	if err != nil {
		return shim.Error("Error marshaling repair structure.")
	}

	bike.Status = BIKE_REPAIRING
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(repairKey, repairBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Repair %s accepted.\n", args[1])

	return shim.Success(nil)
}

// Reject the request to repair a bike
func (t *BikeShareWorkflowChaincode) rejectRepair(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var repairer *Repairer
	var bike *Bike
	var repair *Repair

	// Access control: Only a Repairer Org member can invoke this transaction
	if !t.devMode && !authenticateRepairerOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Repairer Org. Access denied.")
	}

	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {Repairer ID, Repair ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get repairer state from the ledger
	repairerKey, err := getRepairerKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairerBytes, err := stub.GetState(repairerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairerBytes) == 0 {
		err = errors.New(fmt.Sprintf("Repairer %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(repairerBytes, &repairer)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Get repair state from the ledger
	repairKey, err := getRepairKey(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairBytes, err := stub.GetState(repairKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairBytes) == 0 {
		err = errors.New(fmt.Sprintf("Repair %s not found.", args[1]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(repairBytes, &repair)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if repairer matches
	if repair.RepairerId != args[0] {
		err = errors.New(fmt.Sprintf("Actual repairer %s and requested repairer %s not match.", repair.RepairerId, args[0]))
		return shim.Error(err.Error())
	}

	// Verify if repair is requested
	if repair.Status != REPAIR_REQUESTED {
		err = errors.New(fmt.Sprintf("Repair %s already processed.", args[1]))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
	bikeKey, err := getBikeKey(stub, repair.BikeId)
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) == 0 {
		err = errors.New(fmt.Sprintf("Bike %s not found.", repair.BikeId))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bikeBytes, &bike)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if bike is ready to repair
	if bike.Status != BIKE_TO_REPAIR {
		err = errors.New(fmt.Sprintf("Bike %s not ready to repair.", repair.BikeId))
		return shim.Error(err.Error())
	}

	repair.Status = REPAIR_REJECTED
	repairBytes, err = json.Marshal(repair)
	if err != nil {
		return shim.Error("Error marshaling repair structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(repairKey, repairBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Repair %s rejected.\n", args[1])

	return shim.Success(nil)
}

// Complete the repair of a bike
func (t *BikeShareWorkflowChaincode) completeRepair(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error
	var repairer *Repairer
	var bike *Bike
	var repair *Repair

	// Access control: Only a Repairer Org member can invoke this transaction
	if !t.devMode && !authenticateRepairerOrg(creatorOrg, creatorCertIssuer) {
		return shim.Error("Caller not a member of Repairer Org. Access denied.")
	}

	if len(args) != 2 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 2: {Repairer ID, Repair ID}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	// Get repairer state from the ledger
	repairerKey, err := getRepairerKey(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairerBytes, err := stub.GetState(repairerKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairerBytes) == 0 {
		err = errors.New(fmt.Sprintf("Repairer %s not found.", args[0]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(repairerBytes, &repairer)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Get repair state from the ledger
	repairKey, err := getRepairKey(stub, args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	repairBytes, err := stub.GetState(repairKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(repairBytes) == 0 {
		err = errors.New(fmt.Sprintf("Repair %s not found.", args[1]))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(repairBytes, &repair)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if repairer matches
	if repair.RepairerId != args[0] {
		err = errors.New(fmt.Sprintf("Actual repairer %s and requested repairer %s not match.", repair.RepairerId, args[0]))
		return shim.Error(err.Error())
	}

	// Verify if repair is accepted
	if repair.Status != REPAIR_ACCEPTED {
		err = errors.New(fmt.Sprintf("Repair %s not accepted.", args[1]))
		return shim.Error(err.Error())
	}

	// Get bike state from the ledger
	bikeKey, err := getBikeKey(stub, repair.BikeId)
	if err != nil {
		return shim.Error(err.Error())
	}
	bikeBytes, err := stub.GetState(bikeKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(bikeBytes) == 0 {
		err = errors.New(fmt.Sprintf("Bike %s not found.", repair.BikeId))
		return shim.Error(err.Error())
	}

	// Unmarshal the JSON
	err = json.Unmarshal(bikeBytes, &bike)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify if bike is repairing
	if bike.Status != BIKE_REPAIRING {
		err = errors.New(fmt.Sprintf("Bike %s not repairing.", repair.BikeId))
		return shim.Error(err.Error())
	}

	repair.Status = REPAIR_COMPLETED
	repairBytes, err = json.Marshal(repair)
	if err != nil {
		return shim.Error("Error marshaling repair structure.")
	}
	
	bike.Status = BIKE_REPAIRED
	bikeBytes, err = json.Marshal(bike)
	if err != nil {
		return shim.Error("Error marshaling bike structure.")
	}

	// Write the state to the ledger
	err = stub.PutState(repairKey, repairBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(bikeKey, bikeBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Repair %s completed.\n", args[1])

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
		queryResponseArray.WriteString(",\"Value\":")
		// Value is a JSON object, so we write as-is
		queryResponseArray.WriteString(string(queryResponse.Value))
		queryResponseArray.WriteString("}")
		isArrayMemberAlreadyWritten = true
	}
	queryResponseArray.WriteString("]")

	return &queryResponseArray, nil
}

func getQueryResponse(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("Query String: %s\n", queryString)

	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	queryResponse, err := constructQueryResponseFromIterator(iterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Query Result: %s\n", queryResponse.String())

	return queryResponse.Bytes(), nil
}

// Get all bikes with specified status
func (t *BikeShareWorkflowChaincode) getBikesByStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Provider/User/Repairer Org member can invoke this transaction
	if !t.devMode && !(authenticateProviderOrg(creatorOrg, creatorCertIssuer) || authenticateUserOrg(creatorOrg, creatorCertIssuer) || authenticateRepairerOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Provider/User/Repairer Org. Access denied.")
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

// Get all rides with specified status
func (t *BikeShareWorkflowChaincode) getRidesByStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Provider/User Org member can invoke this transaction
	if !t.devMode && !(authenticateProviderOrg(creatorOrg, creatorCertIssuer) || authenticateUserOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Provider/User Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Status}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"status\":\"%s\"}}", RIDE, args[0])
	queryResponse, err := getQueryResponse(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResponse)
}

// Get all issues with specified status
func (t *BikeShareWorkflowChaincode) getIssuesByStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Provider/User Org member can invoke this transaction
	if !t.devMode && !(authenticateProviderOrg(creatorOrg, creatorCertIssuer) || authenticateUserOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Provider/User Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Status}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"status\":\"%s\"}}", ISSUE, args[0])
	queryResponse, err := getQueryResponse(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResponse)
}

// Get all repairs with specified status
func (t *BikeShareWorkflowChaincode) getRepairsByStatus(stub shim.ChaincodeStubInterface, creatorOrg string, creatorCertIssuer string, args []string) pb.Response {
	var err error

	// Access control: Only a Provider/Repairer Org member can invoke this transaction
	if !t.devMode && !(authenticateProviderOrg(creatorOrg, creatorCertIssuer) || authenticateRepairerOrg(creatorOrg, creatorCertIssuer)) {
		return shim.Error("Caller not a member of Provider/Repairer Org. Access denied.")
	}

	if len(args) != 1 {
		err = errors.New(fmt.Sprintf("Incorrect number of arguments. Expecting 1: {Status}. Found %d.", len(args)))
		return shim.Error(err.Error())
	}

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"status\":\"%s\"}}", REPAIR, args[0])
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
