#!/bin/bash

echo
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo
echo "end-to-end test"
echo
CHANNEL_NAME="$1"
DELAY="$2"
LANGUAGE="$3"
TIMEOUT="$4"
VERBOSE="$5"
NO_CHAINCODE="$6"
: ${CHANNEL_NAME:="tlchannel"}
: ${DELAY:="3"}
: ${LANGUAGE:="golang"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
: ${NO_CHAINCODE:="false"}
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=10

CC_SRC_PATH="local/chaincode"

echo "Channel name : "$CHANNEL_NAME

# import utils
. scripts/utils.sh

createChannel() {
	setGlobals 0 1

	if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
                set -x
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx >&log.txt
		res=$?
                set +x
	else
				set -x
		peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
		res=$?
				set +x
	fi
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel '$CHANNEL_NAME' created ===================== "
	echo
}

joinChannel () {
	for org in 1 2; do
	    for peer in 0 1; do
		joinChannelWithRetry $peer $org
		echo "===================== peer${peer}.org${org} joined channel '$CHANNEL_NAME' ===================== "
		sleep $DELAY
		echo
	    done
	done
}

## Create channel
echo "Creating channel..."
createChannel

## Join all the peers to the channel
echo "Having all peers join the channel..."
joinChannel

## Set the anchor peers for each org in the channel
echo "Updating anchor peers for org1..."
updateAnchorPeers 0 1
echo "Updating anchor peers for org2..."
updateAnchorPeers 0 2

if [ "${NO_CHAINCODE}" != "true" ]; then

	## Install chaincode on peer0.org1 and peer0.org2
	echo "Install chaincode on peer0.org1..."
	installChaincode 0 1
	installChaincode 1 1
	echo "Install chaincode on peer0.org2..."
	installChaincode 0 2

	# Instantiate chaincode on peer0.org2
	# DATA: ["2fafd2b3973957950192593d6e72e6fad8c0a9f1","{\"user\":\"TestUser\",\"time\":\"1234567890\",\"md5\":\"8b70272ef79d09d5\"}"]
	echo "Instantiating chaincode on org1..."
	instantiateChaincode 0 1 '{"Args":["init","2fafd2b3973957950192593d6e72e6fad8c0a9f1","{\"user\":\"TestUser\",\"time\":\"1234567890\",\"md5\":\"8b70272ef79d09d5\"}"]}'
	instantiateChaincode 1 1 '{"Args":["init","2fafd2b3973957950192593d6e72e6fad8c0a9f1","{\"user\":\"TestUser\",\"time\":\"1234567890\",\"md5\":\"8b70272ef79d09d5\"}"]}'

	sleep $DELAY
#Use golang for test

	# Query chaincode on peer0.org1
	#echo "Querying chaincode on peer0.org2..."
	#chaincodeQuery 0 2 '{"Args":["query","2fafd2b3973957950192593d6e72e6fad8c0a9f1"]}'
	#Expect "{\"user\":\"TestUser\",\"time\":\"1234567890\",\"md5\":\"8b70272ef79d09d5\"}"
fi

echo
echo "========= All GOOD, execution completed =========== "
echo

echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

exit 0
