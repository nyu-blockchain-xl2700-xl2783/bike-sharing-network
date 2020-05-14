
var user_fcn = ['getUsers', 'getBikes','updateBikeLocation', 'startRide', 'endRide', 'reportIssue'];
var user_args = ['','','BIKE_ID LONGITUDE LATITUDE' ,'USER_ID RIDE_ID BIKE_ID TIMESTAMP LONGITUDE LATITUDE','USER_ID RIDE_ID TIMESTAMP LONGITUDE LATITUDE', 'USER_ID ISSUE_ID RIDE_ID'];
var provider_fcn = ['getUsers','getRepairers','getBikes','getBikeById','getBikesByStatus','getRides','getRideById','getRidesByUser','getRidesByBike','getRidesByStatus','getIssues','getIssueById','getIssuesByUser','getIssuesByBike','getIssueByRide','getIssuesByStatus','getRepairs','getRepairById','getRepairsByBike','getRepairsByRepairer','getRepairsByStatus','registerBike','reactivateBike','discardBike','updateBikeLocation','acceptIssue','rejectIssue','requestRepair'];
var provider_args = ['','','','BIKE_ID','BIKE_STATUS','','RIDE_ID','USER_ID','BIKE_ID','RIDE_STATUS','','ISSUE_ID','USER_ID','BIKE_ID','RIDE_ID','ISSUE_STATUS','','REPAIR_ID','BIKE_ID','REPAIRER_ID','REPAIR_STATUS','BIKE_ID','BIKE_ID','BIKE_ID','BIKE_ID LONGITUDE LATITUDE','ISSUE_ID','ISSUE_ID','REPAIR_ID BIKE_ID REPAIRER_ID'];
var repairer_fcn = ['getRepairers','getIssues','getIssueById','getIssuesByUser','getIssuesByBike','getIssueByRide','getIssuesByStatus','getRepairs','getRepairById','getRepairsByBike','getRepairsByRepairer','getRepairsByStatus','updateBikeLocation','acceptRepair','rejectRepair','completeRepair'];
var repairer_args = ['','','ISSUE_ID','USER_ID','BIKE_ID','RIDE_ID','ISSUE_STATUS','','REPAIR_ID','BIKE_ID','REPAIRER_ID','REPAIR_STATUS','BIKE_ID LONGITUDE LATITUDE','REPAIRER_ID REPAIR_ID','REPAIRER_ID REPAIR_ID','REPAIRER_ID REPAIR_ID'];
var ccversion = "v0";

function execute(org){
    var args = [];
    var fcn = [];
    switch(org){
        case('userorg'):
            fcn = user_fcn;
            break;
        case('repairerorg'):
            fcn = repairer_fcn;
            break;
        case('providerorg'):
            fcn = provider_fcn;
            break;
        default:
    }
    var httptype = "POST";
    var selectvalue = document.getElementById('fcn').value;
    var textbox = document.getElementById('text1');
    if (fcn[selectvalue].substring(0,3) == "get"){
        httptype = "GET";
    } else {
        httptype = "POST";
    }
    var arg1 = document.getElementById('arg1');
    var arg2 = document.getElementById('arg2');
    var arg3 = document.getElementById('arg3');
    var arg4 = document.getElementById('arg4');
    var arg5 = document.getElementById('arg5');
    var arg6 = document.getElementById('arg6');
    if (arg1.value.length>0) args.push(arg1.value);
    if (arg2.value.length>0) args.push(arg2.value);
    if (arg3.value.length>0) args.push(arg3.value);
    if (arg4.value.length>0) args.push(arg4.value);
    if (arg5.value.length>0) args.push(arg5.value);
    if (arg6.value.length>0) args.push(arg6.value);
    var token = getCookie("my_token");
    
    if (httptype == "POST"){
        $.ajax({
            type: "POST",
            url: "/chaincode/"+fcn[selectvalue] ,
            headers:{
                "content-type": "application/json",
                "authorization": "Bearer "+token
            },
            data: JSON.stringify({
                "ccversion": ccversion,
                "args": args,
            }),
            success: function (result) {
                textbox.value = result.message;
            },
            error : function() {
                textbox.value = "Error";
            }
        });
    } else {
        var arg = "";
        if (args.length>0){
            arg = args[0];
        }
        $.ajax({
            type: "GET",
            url: "/chaincode/"+fcn[selectvalue] ,
            headers:{
                "content-type": "application/json",
                "authorization": "Bearer "+token
            },
            processData: false,
            data: $.param({'ccversion' : ccversion, 'args':arg}),
            success: function (result) {
                textbox.value = result.message;
            },
            error : function() {
                textbox.value = "Error";
            }
        });
    }
    arg1.value = "";
    arg2.value = "";
    arg3.value = "";
    arg4.value = "";
    arg5.value = "";
    arg6.value = "";
}


function fcn(org){
    var args = [];
    var textbox = document.getElementById('text1');
    var selectvalue = document.getElementById('fcn').value;
    switch(org){
        case('userorg'):
            args = user_args;
            break;
        case('repairerorg'):
            args = repairer_args;
            break;
        case('providerorg'):
            args = provider_args;
            break;
        default:
    }
    textbox.value = args[selectvalue];

}


function getCookie(name) {
    var cookieValue = "";
    if (document.cookie && document.cookie !== '') {
        var cookies = document.cookie.split(';');
        for (var i = 0; i < cookies.length; i++) {
            var cookie = $.trim(cookies[i]);
            if (cookie.substring(0, name.length + 1) === (name + '=')) {
                cookieValue = decodeURIComponent(cookie.substring(name.length + 1));
                break;
            }
        }
    }
    return cookieValue;}