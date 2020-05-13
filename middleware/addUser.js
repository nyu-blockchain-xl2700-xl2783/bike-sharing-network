'use stricts'



var ClientUtils = require('./clientUtils.js');


// Enroll a new user of userorg and sign the channel configuration as that user (signing identity)
function addUser(org, username) {
	client._userContext = null;

	return ClientUtils.getClientUser(org, username)
	.then((user) => {
		console.log('Successfully enrolled user \'user',username,'\' for', org);

	});
}