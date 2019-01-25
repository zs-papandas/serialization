#!/bin/sh

set -x

go build -o pscc

sleep 1

CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./pscc