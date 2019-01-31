#!/bin/bash


# Exit on first error, print all commands.



# =============================================
# vim mefn.sh && chmod +x mefn.sh && ./mefn.sh
# =============================================
# > mefn.sh && vim mefn.sh && ./mefn.sh
# =============================================
# Check the Composer Card list
# =============================================

export FABRIC_VERSION=hlfv12

composer card list -c PeerAdmin@hlfv1 > /dev/null
PEERADMINEXIST=$?
if [ "${PEERADMINEXIST}" -eq "0" ] ; then
    echo "=================="
    echo "0. Removing PeerAdmin Card"
    echo "=================="
    composer card delete -c PeerAdmin@hlfv1
    sleep 1
else
    echo "=================="
    echo "0. PeerAdmin Card do not exist. Creating a new PeerAdmin Card."
    echo "=================="
fi


composer card list -c admin@product-serialization-network > /dev/null
ADMINEXIST=$?
if [ "${ADMINEXIST}" -eq "0" ] ; then
    echo "=================="
    echo "1. Removing a Product serialization network admin Card"
    echo "=================="
    composer card delete -c admin@product-serialization-network
    sleep 1
else
    echo "=================="
    echo "1. Product serialization network admin Card do not exist. Creating a new card."
    echo "=================="
fi

sleep 1

echo "=================="
echo "2. Create BNA Files"
echo "=================="
composer archive create -t dir -n .

echo "=================="
echo "3. Stop Fabric"
echo "=================="
~/fabric-tools/stopFabric.sh

sleep 2

echo "=================="
echo "5. Start Fabric"
echo "=================="
~/fabric-tools/startFabric.sh


echo "=================="
echo "6. Building a connection profile"
echo "=================="

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


echo "=================="
echo "7. Creating a business network card for the Hyperledger Fabric administrator"
echo "=================="
composer card create -p connection.json -u PeerAdmin -c /home/ubuntu/fabric-tools/fabric-scripts/hlfv12/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem -k /home/ubuntu/fabric-tools/fabric-scripts/hlfv12/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/114aab0e76bf0c78308f89efc4b8c9423e31568da0c340ca187a9b17aa9a4457_sk -r PeerAdmin -r ChannelAdmin



echo "=================="
echo "8. Importing the business network card for the Hyperledger Fabric administrator"
echo "=================="
composer card import -f PeerAdmin@hlfv1.card

echo "=================="
echo "9. Installing the Hyperledger Composer business network onto the Hyperledger Fabric peer nodes"
echo "=================="
composer network install -c PeerAdmin@hlfv1 -a product-serialization-network@0.0.1.bna


echo "=================="
echo "10. Starting the blockchain business network"
echo "=================="
composer network start --networkName product-serialization-network --networkVersion 0.0.1 -c PeerAdmin@hlfv1 -A admin -S adminpw --file /home/ubuntu/pd19385/product-serialization-network/PeerAdmin@hlfv1.card

echo "=================="
echo "11. Importing the business network card for the business network administrator"
echo "=================="
composer card import -f /home/ubuntu/pd19385/product-serialization-network/PeerAdmin@hlfv1.card


echo "=================="
echo "12. Testing the connection to the blockchain business network"
echo "=================="
composer network ping -c admin@product-serialization-network


echo "=================="
echo "13. PROCESS OVER"
echo "=================="
rm ./connection.json
rm ./product-serialization-network@0.0.1.bna
#rm ./*.card

cat > ./mefn.sh <<EOF
#!/bin/bash
npm start
EOF

./mefn.sh