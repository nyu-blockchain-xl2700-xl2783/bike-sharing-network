var ccversion = "v0";
var balance = "10000";

function login(){
    var username = document.getElementById('inputEmail').value;
    var password = document.getElementById('inputPassword').value;
    var org = document.getElementById('org').value;
    $.ajax({
            type: "POST",
            dataType: "json",
            url: "/login" ,
            headers:{
                'content-type': 'application/x-www-form-urlencoded'
            },
            data: {
                'username':username,
                'password':password,
                'orgName':org,
            },
            success: function (result) {
                console.log(result);
                setTokenToCookie(result.token);

                if (username != "admin"){
                    if (org == "userorg"){
                        addUserToDB(username, result.token, ccversion);
                    } else if (org == "repairerorg") {
                        addRepairerToDB(username, result.token, ccversion);
                    }
                } else {
                    if (org == 'providerorg'){
                        switchPage("providers");
                    }
                }
                

            },
            error : function() {
                alert("Error");
            }
        });
}

function addUserToDB(username, token, ccversion){
    $.ajax({
        type: "POST",
        dataType: "json",
        url: "/chaincode/registerUser" ,
        headers:{
            "content-type": "application/json",
            "authorization": "Bearer "+token
        },
        data: JSON.stringify({
            "ccversion": ccversion,
            "args": [username,balance]
        }),
        success: function (result) {
            console.log("Add user to db succeed");
            switchPage("users");

        },
        error : function() {
            alert("Error");
        }
    });
}

function addRepairerToDB(username, token, ccversion){
    $.ajax({
        type: "POST",
        dataType: "json",
        url: "/chaincode/registerRepairer" ,
        headers:{
            "content-type": "application/json",
            "authorization": "Bearer "+token
        },
        data: JSON.stringify({
            "ccversion":ccversion,
            "args": [username]
        }),
        success: function (result) {
            console.log("Add repairer to db succeed");
            switchPage("repairers");

        },
        error : function() {
            alert("Error");
        }
    });
}


function setTokenToCookie(value) {
    var Days = 1; 
    var exp = new Date();
    exp.setTime(exp.getTime() + Days * 24 * 60 * 60 * 1000);
    document.cookie = "my_token =" + escape(value) + ";expires=" + exp.toGMTString();
}


function switchPage(location){
    window.location.replace("http://localhost:8080/"+location);
}