package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func getUserKey(stub shim.ChaincodeStubInterface, userID string) (string, error) {
	userKey, err := stub.CreateCompositeKey("User-", []string{userID})
	if err != nil {
		return "", err
	} else {
		return userKey, nil
	}
}

func getRepairerKey(stub shim.ChaincodeStubInterface, repairerID string) (string, error) {
	repairerKey, err := stub.CreateCompositeKey("Repairer-", []string{repairerID})
	if err != nil {
		return "", err
	} else {
		return repairerKey, nil
	}
}

func getBikeKey(stub shim.ChaincodeStubInterface, bikeID string) (string, error) {
	bikeKey, err := stub.CreateCompositeKey("Bike-", []string{bikeID})
	if err != nil {
		return "", err
	} else {
		return bikeKey, nil
	}
}

func getRideKey(stub shim.ChaincodeStubInterface, rideID string) (string, error) {
	rideKey, err := stub.CreateCompositeKey("Ride-", []string{rideID})
	if err != nil {
		return "", err
	} else {
		return rideKey, nil
	}
}

func getIssueKey(stub shim.ChaincodeStubInterface, issueID string) (string, error) {
	issueKey, err := stub.CreateCompositeKey("Issue-", []string{issueID})
	if err != nil {
		return "", err
	} else {
		return issueKey, nil
	}
}

func getRepairKey(stub shim.ChaincodeStubInterface, repairID string) (string, error) {
	repairKey, err := stub.CreateCompositeKey("Repair-", []string{repairID})
	if err != nil {
		return "", err
	} else {
		return repairKey, nil
	}
}
