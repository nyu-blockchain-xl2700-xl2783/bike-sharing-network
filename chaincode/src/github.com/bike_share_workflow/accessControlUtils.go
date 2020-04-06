package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"crypto/x509"
)

func getCustomAttribute(stub shim.ChaincodeStubInterface, attr string) (string, bool, error) {
	var value string
	var found bool
	var err error

	value, found, err = cid.GetAttributeValue(stub, attr)
	if err != nil {
		fmt.Printf("Error getting MSP identity: %s\n", err.Error())
		return "", found, err
	}

	return value, found, nil
}

func getTxCreatorInfo(stub shim.ChaincodeStubInterface) (string, string, error) {
	var mspid string
	var err error
	var cert *x509.Certificate

	mspid, err = cid.GetMSPID(stub)
	if err != nil {
		fmt.Printf("Error getting MSP identity: %s\n", err.Error())
		return "", "", err
	}

	cert, err = cid.GetX509Certificate(stub)
	if err != nil {
		fmt.Printf("Error getting client certificate: %s\n", err.Error())
		return "", "", err
	}

	return mspid, cert.Issuer.CommonName, nil
}

// For simplicity, just hardcode an ACL

func authenticateProviderOrg(mspID string, certCN string) bool {
	return (mspID == "ProviderOrgMSP") && (certCN == "ca.providerorg.bikeshare.com")
}

func authenticateUserOrg(mspID string, certCN string) bool {
	return (mspID == "UserOrgMSP") && (certCN == "ca.userorg.bikeshare.com")
}

func authenticateRepairerOrg(mspID string, certCN string) bool {
	return (mspID == "RepairerOrgMSP") && (certCN == "ca.repairerorg.bikeshare.com")
}
