#!/bin/sh

set -x

# Chaincode Install
peer chaincode install -p github.com/zs-papandas/serialization -n mycc -v 0

# Chaincode Instatiate
peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc

sleep 2
# =============================================================================
# Create User Account
# =============================================================================

# Adding Manufacturer
peer chaincode invoke -n mycc -c '{"Args":["createAccount", "a", "Papan","Das","1987-11-14","papan.das@zs.com","9641443962","ZS Associates India Pvt Ltd","manufacturer"]}' -C myc

sleep 2
# Adding Wholesaler
peer chaincode invoke -n mycc -c '{"Args":["createAccount", "b", "Wholesaler","Wholesaler","1977-10-13","rakesh.kumar@zs.com","","ZS Associates India Pvt Ltd","wholesaler"]}' -C myc

sleep 2
# Adding Retailer
peer chaincode invoke -n mycc -c '{"Args":["createAccount", "c", "Retailer","Retailer","1977-10-13","rakesh.kumar@zs.com","","ZS Associates India Pvt Ltd","retailer"]}' -C myc

sleep 2
# Adding Patient
peer chaincode invoke -n mycc -c '{"Args":["createAccount", "d", "Rakesh","Kumar","1977-10-13","rakesh.kumar@zs.com","","ZS Associates India Pvt Ltd","patient"]}' -C myc

sleep 2
# =============================================================================
# Create Product
# =============================================================================

# Create pallet
peer chaincode invoke -n mycc -c '{"Args":["createProduct", "a", "CROSINpallet","expire","gtin-102030","lotnum/23/45/as","1000","2","pallet",""]}' -C myc

sleep 2

# Create Box
peer chaincode invoke -n mycc -c '{"Args":["createProduct", "a", "CROSINbox","expire","gtin-102030","lotnum/23/45/as","1000","3","box","T7txWduqtCDMLxiD"]}' -C myc

sleep 2
# Create Packet
peer chaincode invoke -n mycc -c '{"Args":["createProduct", "a", "CROSINpacket","expire","gtin-102030","lotnum/23/45/as","1000","4","packet","ynhwJyzAHyfjXUlr"]}' -C myc

sleep 2

# Create Item
peer chaincode invoke -n mycc -c '{"Args":["createProduct", "a", "CROSINitem","expire","gtin-102030","lotnum/23/45/as","1000","0","item","WIXnkuHMYZL5fGaE"]}' -C myc

sleep 2
