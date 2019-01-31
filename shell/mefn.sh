#!/bin/bash

PACKAGE_VERSION="15"

 YELLOW='\033[1;33m'
 RED='\033[1;31m'
 GREEN='\033[1;32m'
 RESET='\033[0m'
 GREY='\033[2m'

# get the command line options
NETWORK_NAME="product-serialization-network"
HLF_INSTALL_PATH="/home/ubuntu/fabric-tools"
BNA_VERSION="0.0."${PACKAGE_VERSION}


# indent text on echo
function indent() {
  c='s/^/       /'
  case $(uname) in
    Darwin) sed -l "$c";;
    *)      sed -u "$c";;
  esac
}

# displays where we are, uses the indent function (above) to indent each line
function showStep ()
{
    echo -e "${GREY}=====================================================" | indent
    echo -e "${RESET}-----> $*" | indent
    echo -e "${GREY}=====================================================${RESET}" | indent
}

# Grab the current directory
function getCurrent() 
{
    showStep "getting current directory"
    DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
    echo "DIR in startup.sh is $DIR"
    THIS_SCRIPT=`basename "$0"`
    showStep "Running '${THIS_SCRIPT}'"
}


echo  "Parameters:"
echo -e "Network Name is: ${GREEN} $NETWORK_NAME ${RESET}" | indent

showStep "update package.json version"
sed -i '/version/s/[^.]*$/'"${PACKAGE_VERSION}\",/" package.json

showStep "creating archive"
composer archive create --sourceType dir --sourceName . -a ./dist/$NETWORK_NAME.bna

showStep "running getCurrent"
getCurrent

showStep "using execs from previous installation, stored in ${HLF_INSTALL_PATH}"
cd "${HLF_INSTALL_PATH}"
showStep "starting fabric"
./startFabric.sh

showStep "creating new PeerAdmin card"
./createPeerAdminCard.sh 

showStep "copying admin card to ~/.hfc-key-store"
CA_PEM_SOURCE="$DIR/controller/restapi/features/composer/creds"
PEER_SOURCE="$HOME/.composer/client-data/PeerAdmin@hlfv1/*"
HFC_KEY_STORE="$HOME/.hfc-key-store"
echo "CA_PEM_SOURCE is: $CA_PEM_SOURCE"
echo "PEER_SOURCE is: $PEER_SOURCE"
echo "HFC_KEY_STORE is: $HFC_KEY_STORE"
rm -R $HFC_KEY_STORE/
mkdir $HFC_KEY_STORE
cp  -Rv ${CA_PEM_SOURCE}/ca.pem ${HFC_KEY_STORE}/
cp  -Rv ${PEER_SOURCE} ${HFC_KEY_STORE}/
cp  -Rv ${PEER_SOURCE} ${CA_PEM_SOURCE}/


showStep 'Listing current cards'
composer card list -c PeerAdmin@hlfv1
showStep "start up complete"


showStep "deploying network"
cd $DIR/dist

showStep "creating connection JSON file"
cat > ./connection.json <<EOF
{
    "name": "hlfv1",
    "x-type": "hlfv1",
    "version": "1.0.0",
    "peers": {
        "peer0.org1.example.com": {
            "url": "grpc://localhost:7051"
        }
    },
    "certificateAuthorities": {
        "ca.org1.example.com": {
            "url": "http://localhost:7054",
            "caName": "ca.org1.example.com"
        }
    },
    "orderers": {
        "orderer.example.com": {
            "url": "grpc://localhost:7050"
        }
    },
    "organizations": {
        "Org1": {
            "mspid": "Org1MSP",
            "peers": [
                "peer0.org1.example.com"
            ],
            "certificateAuthorities": [
                "ca.org1.example.com"
            ]
        }
    },
    "channels": {
        "composerchannel": {
            "orderers": [
                "orderer.example.com"
            ],
            "peers": {
                "peer0.org1.example.com": {
                    "endorsingPeer": true,
                    "chaincodeQuery": true,
                    "eventSource": true
                }
            }
        }
    },
    "client": {
        "organization": "Org1",
        "connection": {
            "timeout": {
                "peer": {
                    "endorser": "300",
                    "eventHub": "300",
                    "eventReg": "300"
                },
                "orderer": "300"
            }
        }
    }
}
EOF


showStep "Creating a business network card for the Hyperledger Fabric administrator"
composer card create -p connection.json -u PeerAdmin -c ${HLF_INSTALL_PATH}/fabric-scripts/hlfv12/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem -k ${HLF_INSTALL_PATH}/fabric-scripts/hlfv12/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/*_sk -r PeerAdmin -r ChannelAdmin




showStep "Importing the business network card for the Hyperledger Fabric administrator"
composer card delete -c PeerAdmin@hlfv1
sleep 1
composer card import -f PeerAdmin@hlfv1.card



showStep "Installing the Hyperledger Composer business network onto the Hyperledger Fabric peer nodes"
#composer network install --card PeerAdmin@hlfv1 --businessNetworkName $NETWORK_NAME
composer network install -c PeerAdmin@hlfv1 -a $NETWORK_NAME.bna

showStep "starting network"
#composer runtime start -c PeerAdmin@hlfv1 -A admin -S adminpw -a $NETWORK_NAME.bna --file networkadmin.card
composer network start --networkName $NETWORK_NAME --networkVersion $BNA_VERSION -c PeerAdmin@hlfv1 -A admin -S adminpw --file PeerAdmin@hlfv1.card
#composer network start --networkName product-serialization-network --networkVersion 0.0.2 -c PeerAdmin@hlfv1 -A admin -S adminpw --file ./dist/PeerAdmin@hlfv1.card


showStep "importing networkadmin card"
#if composer card list -n admin@$NETWORK_NAME > /dev/null; then
#    composer card delete -n admin@$NETWORK_NAME
#fi
composer card delete -c admin@$NETWORK_NAME
composer card import --file $DIR/dist/PeerAdmin@hlfv1.card



showStep "pinging admin@$NETWORK_NAME card"
composer network ping --card admin@$NETWORK_NAME

rm connection.json

showStep "Erase mefn.sh"
cd ..
pwd
> mefn.sh && cat mefn.sh

showStep "Start Composer Playground on port 8085"
