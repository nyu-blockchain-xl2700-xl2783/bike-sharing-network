package main

type Bike struct {
	ObjectType 		string 		`json:"docType"`
	Id				string		`json:"id"`
	Location		[]float32	`json:"location"`
	Status			string		`json:"status"`
}

type Ride struct {
	ObjectType 		string 		`json:"docType"`
	Id				string		`json:"id"`
	UserId			string		`json:"userId"`
	BikeId			string		`json:"bikeId"`
	StartTime		string		`json:"startTime"`
	StartLocation	[]float32	`json:"startLocation"`
	EndTime			string		`json:"endTime"`
	EndLocation		[]float32	`json:"endLocation"`
	Status			string		`json:"status"`
}

type Issue struct {
	ObjectType 		string 		`json:"docType"`
	Id				string		`json:"id"`
	UserId			string		`json:"userId"`
	BikeId			string		`json:"bikeId"`
	RideId			string		`json:"rideId"`
	Status			string		`json:"status"`
}

type Repair struct {
	ObjectType 		string 		`json:"docType"`
	Id				string		`json:"id"`
	BikeId			string		`json:"bikeId"`
	RepairerId		string		`json:"repairerId"`
	Status			string		`json:"status"`
}
