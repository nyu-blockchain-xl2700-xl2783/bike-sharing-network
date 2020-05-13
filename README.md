# Bike-sharing Network

## Functions and Constants

### Registration

* `registerUser USER_ID BALANCE`
* `registerRepairer REPAIRER_ID`

### Transaction

* `registerBike BIKE_ID`
* `reactivateBike BIKE_ID`
* `discardBike BIKE_ID`
* `updateBikeLocation BIKE_ID LONGITUDE LATITUDE`
* `startRide USER_ID RIDE_ID BIKE_ID TIMESTAMP LONGITUDE LATITUDE`
* `endRide USER_ID RIDE_ID TIMESTAMP LONGITUDE LATITUDE`
* `reportIssue USER_ID ISSUE_ID RIDE_ID`
* `acceptIssue ISSUE_ID`
* `rejectIssue ISSUE_ID`
* `requestRepair REPAIR_ID BIKE_ID REPAIRER_ID`
* `acceptRepair REPAIRER_ID REPAIR_ID`
* `rejectRepair REPAIRER_ID REPAIR_ID`
* `completeRepair REPAIRER_ID REPAIR_ID`

### Query

* `getUsers`
* `getRepairers`
* `getBikes`
* `getBikeById BIKE_ID`
* `getBikesByStatus BIKE_STATUS`
* `getRides`
* `getRideById RIDE_ID`
* `getRidesByUser USER_ID`
* `getRidesByBike BIKE_ID`
* `getRidesByStatus RIDE_STATUS`
* `getIssues`
* `getIssueById ISSUE_ID`
* `getIssuesByUser USER_ID`
* `getIssuesByBike BIKE_ID`
* `getIssueByRide RIDE_ID`
* `getIssuesByStatus ISSUE_STATUS`
* `getRepairs`
* `getRepairById REPAIR_ID`
* `getRepairsByBike BIKE_ID`
* `getRepairsByRepairer REPAIRER_ID`
* `getRepairsByStatus REPAIR_STATUS`

### Status

* User
    - `USER_FREE`
    - `USER_IN_RIDE`
* Bike
    - `BIKE_AVAILABLE`
    - `BIKE_IN_USE`
    - `BIKE_TO_REPAIR`
    - `BIKE_REPAIRING`
    - `BIKE_REPAIRED`
    - `BIKE_DISCARDED`
* Ride
    - `RIDE_ONGOING`
    - `RIDE_COMPLETED`
    - `RIDE_ISSUE_OPEN`
    - `RIDE_ISSUE_CLOSED`
* Issue
    - `ISSUE_OPEN`
    - `ISSUE_CLOSED`
* Repair
    - `REPAIR_REQUESTED`
    - `REPAIR_ACCEPTED`
    - `REPAIR_REJECTED`
    - `REPAIR_COMPLETED`

