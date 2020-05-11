/*
 * Copyright 2018 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an 'AS IS' BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

var os = require('os');
var path = require('path');

var tempdir = "../network/client-certs";
//path.join(os.tmpdir(), 'hfc');

// Frame the endorsement policy
var THREE_ORG_MEMBERS_AND_ADMIN = [{
	role: {
		name: 'member',
		mspId: 'ProviderOrgMSP'
	}
}, {
	role: {
		name: 'member',
		mspId: 'UserOrgMSP'
	}
}, {
	role: {
		name: 'member',
		mspId: 'RepairerOrgMSP'
	}
}, {
	role: {
		name: 'admin',
		mspId: 'OrdererMSP'
	}
}];


var ONE_OF_THREE_ORG_MEMBER = {
	identities: THREE_ORG_MEMBERS_AND_ADMIN,
	policy: {
		'1-of': [{ 'signed-by': 0 }, { 'signed-by': 1 }, { 'signed-by': 2 }]
	}
};

var ALL_THREE_ORG_MEMBERS = {
	identities: THREE_ORG_MEMBERS_AND_ADMIN,
	policy: {
		'3-of': [{ 'signed-by': 0 }, { 'signed-by': 1 }, { 'signed-by': 2 }]
	}
};


var ACCEPT_ALL = {
	identities: [],
	policy: {
		'0-of': []
	}
};

var chaincodeLocation = '../chaincode';

var networkId = 'bike-sharing-network';

var networkConfig = './config.json';

var networkLocation = '../network';

var channelConfig = 'channel-artifacts/channel.tx';

var PROVIDER_ORG = 'providerorg';
var USER_ORG = 'userorg';
var REPAIRER_ORG = 'repairerorg';


var CHANNEL_NAME = 'bsnchannel';
var CHAINCODE_PATH = 'github.com/bike_share_workflow';
var CHAINCODE_ID = 'bsncc';
var CHAINCODE_VERSION = 'v0';

var TRANSACTION_ENDORSEMENT_POLICY = ALL_THREE_ORG_MEMBERS;

module.exports = {
	tempdir: tempdir,
	chaincodeLocation: chaincodeLocation,
	networkId: networkId,
	networkConfig: networkConfig,
	networkLocation: networkLocation,
	channelConfig: channelConfig,
	PROVIDER_ORG: PROVIDER_ORG,
	USER_ORG: USER_ORG,
	REPAIRER_ORG: REPAIRER_ORG,
	CHANNEL_NAME: CHANNEL_NAME,
	CHAINCODE_PATH: CHAINCODE_PATH,
	CHAINCODE_ID: CHAINCODE_ID,
	CHAINCODE_VERSION: CHAINCODE_VERSION,
	ALL_THREE_ORG_MEMBERS: ALL_THREE_ORG_MEMBERS,
	TRANSACTION_ENDORSEMENT_POLICY: TRANSACTION_ENDORSEMENT_POLICY
};
