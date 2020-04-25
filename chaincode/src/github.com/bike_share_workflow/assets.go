package main

type User struct {
	ObjectType 		string 		`json:"docType"`
	Id				string		`json:"id"`
	Balance			float32		`json:"balance"`
	RideId			string		`json:"rideId"`			// Most receent ride ID
	Status			string		`json:"status"`
}

type Repairer struct {
	ObjectType 		string 		`json:"docType"`
	Id				string		`json:"id"`
}

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
	Cost			float32		`json:"cost"`
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
