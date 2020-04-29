# Bike-sharing Network

## Transactions

* `registerUser USER_ID`
* `registerRepairer REPAIRER_ID`
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

