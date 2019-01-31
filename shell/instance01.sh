#!/bin/bash


# Exit on first error, print all commands.
set -ev

# =============================================
# vim mefn.sh && chmod +x mefn.sh && ./mefn.sh
# =============================================
# > mefn.sh && vim mefn.sh && ./mefn.sh
# =============================================
# Check the Composer Card list
# =============================================



#composer card delete -c PeerAdmin@hlfv1
#composer card delete -c PeerAdmin@zs-network
#composer card delete -c admin@patient-concent-network
#composer card list

echo "=================="
echo "1. Export Fabric Version to hlfv12"
echo "=================="
export FABRIC_VERSION=hlfv12

echo "=================="
echo "2. Create BNA Files"
echo "=================="
composer archive create -t dir -n .

echo "=================="
echo "3. Stop Fabric"
echo "=================="
~/fabric-tools/stopFabric.sh

#echo "=================="
##echo "Tear Down Fabric"
#echo "=================="
#~/fabric-tools/teardownFabric.sh

echo "=================="
echo "Tear Down Docker"
echo "=================="
#~/fabric-tools/teardownAllDocker.sh

echo "=================="
echo "Download Fabric"
echo "=================="
#~/fabric-tools/downloadFabric.sh

#echo "=================="
#echo 4. "PRE-Composer Card List"
#echo "=================="
#composer card list

echo "SLEEP 2"
sleep 2

### IMPORTANT NOTE
### composer card delete -c admin@tutorial-network
### IMPORTANT NOTE
### composer card delete -c admin@patient-concent-network  

### echo "=================="
### echo "POST-Composer Card List"
### echo "=================="
### composer card list

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
composer card create -p /home/ubuntu/pd19385/patient-concent-network/connection.json -u PeerAdmin -c /home/ubuntu/fabric-tools/fabric-scripts/hlfv12/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem -k /home/ubuntu/fabric-tools/fabric-scripts/hlfv12/composer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/*_sk -r PeerAdmin -r ChannelAdmin



echo "=================="
echo "8. Importing the business network card for the Hyperledger Fabric administrator"
echo "=================="
composer card delete -c PeerAdmin@hlfv1
sleep 1
composer card import -f PeerAdmin@hlfv1.card

echo "=================="
echo "9. Installing the Hyperledger Composer business network onto the Hyperledger Fabric peer nodes"
echo "=================="
composer network install -c PeerAdmin@hlfv1 -a /home/ubuntu/pd19385/patient-concent-network/dist/patient-concent-network@0.0.1.bna


echo "=================="
echo "10. Starting the blockchain business network"
echo "=================="
composer network start --networkName patient-concent-network --networkVersion 0.0.1 -c PeerAdmin@hlfv1 -A admin -S adminpw --file /home/ubuntu/pd19385/patient-concent-network/PeerAdmin@hlfv1.card

echo "=================="
echo "11. Importing the business network card for the business network administrator"
echo "=================="
composer card delete -c admin@patient-concent-network
sleep 1
composer card import -f /home/ubuntu/pd19385/patient-concent-network/PeerAdmin@hlfv1.card


echo "=================="
echo "12. Testing the connection to the blockchain business network"
echo "=================="
composer network ping -c admin@patient-concent-network


echo "=================="
echo "13. PROCESS OVER"
echo "=================="
rm ./connection.json
rm ./patient-concent-network@0.0.1.bna
rm ./*.card
#rm ./admin@patient-concent-network.card

echo "=================="
echo "14. POPULATE DATA with PYTHON"
echo "=================="
echo " "
echo "next: CHECK RESTful API SERVICE STATUS"
echo " "

#forever stop /home/ubuntu/pd19385/patient-concent-network/index.js
sleep 2
#forever start /home/ubuntu/pd19385/patient-concent-network/index.js
sleep 2

#forever list | grep running > /dev/null
#result=$?
#if [ "${result}" -eq "0" ] ; then
#    echo "=================="
#    echo "Forever service running"
#    echo "=================="
#else
#    echo "=================="
#    echo "Forever service is not running"
#    echo "=================="
#    forever start /home/ubuntu/pd19385/patient-concent-network/index.js
#    sleep 5
#fi


echo "INJECT DATA NOW"
echo " "

cat > ./inject.py <<EOF

#!/usr/bin/python
import requests 
import sys, os
import json
import unicodedata
import time

global BASE_URL, TOKEN, WHATTODO
BASE_URL = "http://localhost:6002"

#FIRST_INSTALL


def apiresponce( method, URL, PARAMS=None, TOKEN=None ):
    _url = BASE_URL + URL
    headers = {
        'Content-Type': 'application/json',
        'authorization': TOKEN
    }
    params = {}
    #payload = PARAMS
    if method=="GET":
        r = requests.get(url = _url, headers=headers)
    elif method=="POST":
        r = requests.post(url = _url, headers=headers, params=params, data=PARAMS)
    else:
        r = requests.put(url = _url, params = PARAMS)
    
    data = r.json()
    return data

def unicodeToString(param):
    try:
        return unicodedata.normalize('NFKD', param).encode('ascii','ignore')
    except:
        return param

def fakeaddr():
    os.system('''curl -H "Content-Type: application/json" https://randomuser.me/api/  >> random.txt''')
    file = open("random.txt","r")
    os.system('''rm ./random.txt''')
    randomuser = json.loads(file.read())
    
    first_name = unicodeToString(randomuser['results'][0]['name']['first'])
    last_name = unicodeToString(randomuser['results'][0]['name']['last'])
    gender = unicodeToString(randomuser['results'][0]['gender'])
    street = unicodeToString(randomuser['results'][0]['location']['street'])
    city = unicodeToString(randomuser['results'][0]['location']['city'])
    state = unicodeToString(randomuser['results'][0]['location']['state'])
    postcode = randomuser['results'][0]['location']['postcode']
    email = unicodeToString(randomuser['results'][0]['email'])
    username = unicodeToString(randomuser['results'][0]['login']['username'])
    dob = unicodeToString(randomuser['results'][0]['dob']['date'])
    phone = unicodeToString(randomuser['results'][0]['phone'])
    cell = unicodeToString(randomuser['results'][0]['cell'])
    someid = unicodeToString(randomuser['results'][0]['id']['name']) + unicodeToString(randomuser['results'][0]['id']['value'])
    return first_name, last_name, gender, street, city, state, postcode, email, username, dob, phone, cell, someid


test_res = apiresponce("GET", "/api/test")

if test_res['result'] == 'Hello World!.':
    print("Express service is working.")
    ping_res = apiresponce("GET", "/api/ping")

    if(len(ping_res['ping']['participant']) > 1):
        print("Hyperledger is running.")
        
        

        ############### MANUFACTURER ###################
        param = {"identity":"manufacturer@zs.com","password":"password","name":"Manufacturer 01","drug":"Propellium","services":"MedicalKit","activitylogs":[],"paramslogs":[]}
        signup_res = apiresponce("POST", "/api/participant/Manufacturer", json.dumps(param))
        if(signup_res['result']=='success'):
            print("Manufacturer added successfully.")
            time.sleep(3)

        ############### VENDOR COPAY ###################
        param = {"identity":"vencopay@zs.com","password":"password","name":"Vendor 03 CoPay","drug":"Propellium","services":"CoPay","activitylogs":[],"paramslogs":[]}
        signup_res = apiresponce("POST", "/api/participant/Vendor", json.dumps(param))
        if(signup_res['result']=='success'):
            print("Copay vendor added successfully.")
            time.sleep(3)

        ############### VENDOR EDUCATIONAL MATERIAL ###################
        param = {"identity":"veneduvid@zs.com","password":"password","name":"Vendor 02 EducationalMaterial","drug":"Propellium","services":"EducationalMaterial","activitylogs":[],"paramslogs":[]}
        signup_res = apiresponce("POST", "/api/participant/Vendor", json.dumps(param))
        if(signup_res['result']=='success'):
            print("Educational material vendor added successfully.")
            time.sleep(3)

        ############### VENDOR MEDICAL KIT ###################
        param = {"identity":"venmedkit@zs.com","password":"password","name":"Vendor 01 MedicalKit","drug":"Propellium","services":"MedicalKit","activitylogs":[],"paramslogs":[]}
        signup_res = apiresponce("POST", "/api/participant/Vendor", json.dumps(param))
        if(signup_res['result']=='success'):
            print("Medical kit vendor added successfully.")
            time.sleep(3)


        #first_name, last_name, gender, street, city, state, postcode, email, username, dob, phone, cell, someid = fakeaddr()
        
        #print(first_name, last_name, gender, street, city, state, postcode, email, username, dob, phone, cell, someid)
        #time.sleep(5)

        ############### PATIENT SATOSHI NAKAMOTO ###################
        param = {"identity":"satoshi.nakamoto","password":"password","name":"Satoshi Nakamoto","gender":"m","dob":"2009-02-01","phone":"1231231231","email":"satoshi.nakamoto@zs.com","address1":"2700  Ralph Drive","address2":"","city":"Lorain","state":"Ohio","zip":"44052","country":"Japan","idisable":True,"pname":"Hans P Jones","pspecialized":"Physician","pfacilityno":"12/56/6776","pphone":"7203120161","pnpino":"INO127845","pdeano":"DE749630","paddress1":"4222  Davis Lane","paddress2":"","pcity":"Centennial","pstate":"Colorado","pzip":"123003","pcountry":"USA","prescriptions":[],"services":["MedicalKit","CoPay"],"iagree":True,"hasCoPayCard":False,"UIDCoPayCard":"","hasMedicalKitOptIn":True,"hasEducationalMaterialOptIn":True,"hasCoPayOptIn":True,"activitylogs":[],"paramslogs":[]}
        signup_res = apiresponce("POST", "/api/participant/Patient", json.dumps(param))
        if(signup_res['result']=='success'):
            print("Satoshi Nakamoto added successfully.")
            time.sleep(3)

        ############### PATIENT EVA.G ###################
        param = {"identity":"eva.g","password":"password","name":"Eva Green","gender":"m","dob":"1990-01-15","phone":"5159256574","email":"eva.g@zs.com","address1":"4007  Heavner Court","address2":"","city":"Lone Rock","state":"Iowa","zip":"50559","country":"USA","idisable":True,"pname":"Tonya C Beebe","pspecialized":"Physician","pfacilityno":"259-17-6366","pphone":"4042215582","pnpino":"558538970","pdeano":"DE749630","paddress1":"2125  Post Farm Road","paddress2":"","pcity":"Atlanta","pstate":"Georgia","pzip":"30303","pcountry":"USA","prescriptions":[],"services":[],"iagree":True,"hasCoPayCard":False,"UIDCoPayCard":"","hasMedicalKitOptIn":False,"hasEducationalMaterialOptIn":False,"hasCoPayOptIn":False,"activitylogs":[],"paramslogs":[]}
        signup_res = apiresponce("POST", "/api/participant/Patient", json.dumps(param))
        if(signup_res['result']=='success'):
            print("Eva Green added successfully.")
            time.sleep(3)

        ############### PATIENT TONYA.C ###################
        param = {"identity":"tonya.c","password":"password","name":"Tonya C Beebe","gender":"f","dob":"1978-02-20","phone":"4042215582","email":"tonya.c@zs.com","address1":"2125  Post Farm Road","address2":"","city":"Atlanta","state":"Georgia","zip":"30303","country":"USA","idisable":False,"pname":"Deborah B Perez","pspecialized":"Physician","pfacilityno":"259-17-6366","pphone":"4042215582","pnpino":"558538970","pdeano":"DE749630","paddress1":"2125  Post Farm Road","paddress2":"","pcity":"Atlanta","pstate":"Georgia","pzip":"30303","pcountry":"USA","prescriptions":[],"services":[],"iagree":True,"hasCoPayCard":False,"UIDCoPayCard":"","hasMedicalKitOptIn":False,"hasEducationalMaterialOptIn":False,"hasCoPayOptIn":False,"activitylogs":[],"paramslogs":[]}
        signup_res = apiresponce("POST", "/api/participant/Patient", json.dumps(param))
        if(signup_res['result']=='success'):
            print("Tonya C added successfully.")
            time.sleep(3)

        #################### RAISE SERVICE REQUEST #####################
        #################### PATIENT LOGIN #####################
        #################### LOGIN satoshi.nakamoto #####################

        patients = ["satoshi.nakamoto", "eva.g", "tonya.c"] 
        ### patients = []
        for i in patients: 
            print(i + " singin initialized.") 
            patient = i
            param = {"identity":patient,"password":"password","registry":"Patient"}
            print(param)
            signup_res = apiresponce("POST", "/api/participant/signin", json.dumps(param))
            if(signup_res['result']=='success'):
                print(signup_res['actor']['actorId'] + " Login successfully.")
                TOKEN=signup_res['actor']['token']
                print(TOKEN)
                time.sleep(3)

                param = {"patient":patient,"travelLocation":"3675  Ocala Street","city":"Orlando","state":"Florida","zip":"32808","country":"USA","travelDate":"3/14/2019","unit":10,"services":"MedicalKit"}
                apiresponce("POST", "/api/asset/services", json.dumps(param), TOKEN)
                print("SR added.")
                time.sleep(3)

                param = {"patient":patient,"travelLocation":"1342  Stroop Hill Road","city":"Atlanta","state":"Georgia","zip":"30303","country":"USA","travelDate":"11/25/2019","unit":10,"services":"EducationalMaterial"} 
                apiresponce("POST", "/api/asset/services", json.dumps(param), TOKEN)
                print("SR added..")
                time.sleep(3)

                param = {"patient":patient,"travelLocation":"846  Lincoln Drive","city":"Harrisburg","state":"Pennsylvania","zip":"17109","country":"USA","travelDate":"12/28/2019","unit":10,"services":"EducationalMaterial"}
                apiresponce("POST", "/api/asset/services", json.dumps(param), TOKEN)
                print("SR added...")
                time.sleep(3)

                
                param = {"patient":patient,"travelLocation":"2492  Snowbird Lane","city":"Omaha","state":"Nebraska","zip":"68104","country":"USA","travelDate":"7/16/2019","unit":10,"services":"EducationalMaterial"}
                apiresponce("POST", "/api/asset/services", json.dumps(param), TOKEN)
                print("SR added.")
                time.sleep(3)


                param = {"patient":patient,"travelLocation":"376  Kinney Street","city":"IBERIA","state":"Missouri","zip":"65486","country":"USA","travelDate":"10/6/2019","unit":10,"services":"CoPay"}
                apiresponce("POST", "/api/asset/services", json.dumps(param), TOKEN)
                print("SR added..")
                time.sleep(3)


                param = {"patient":patient,"travelLocation":"1928  Corpening Drive","city":"Troy","state":"Michigan","zip":"48083","country":"USA","travelDate":"10/15/2019","unit":10,"services":"CoPay"}
                apiresponce("POST", "/api/asset/services", json.dumps(param), TOKEN)
                print("SR added...")
                time.sleep(3)




        '''#################### MANUFACTURERE LOGIN #####################
        param = {"identity":"manufacturer@zs.com","password":"password","registry":"Manufacturer"}
        signup_res = apiresponce("POST", "/api/participant/signin", json.dumps(param))
        if(signup_res['result']=='success'):
            print(signup_res['actor']['actorId'])
            TOKEN=signup_res['actor']['token']
            print(TOKEN)'''

        #fakeuser2 = fakeaddr()
        #print(fakeuser2.first_name, fakeuser2.last_name)
    else:
        print("Hyperledger is not running.")
else:
    print("Express server is not running. Kindly check out.")

EOF

python inject.py

rm inject.py



pwd
> mefn.sh && cat mefn.sh
