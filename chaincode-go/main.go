/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/thanhdaonguyen/TEDU-CertiBlock/certi-abac/chaincode-go/chaincode"
)

func main() {
	certSmartContract, err := contractapi.NewChaincode(&chaincode.TEduCertContract{})
	if err != nil {
		log.Panicf("Error creating certificate chaincode: %v", err)
	}

	if err := certSmartContract.Start(); err != nil {
		log.Panicf("Error starting certificate chaincode: %v", err)
	}
}
