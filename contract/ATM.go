package main

import (
	"encoding/json"
	"fmt"
	"time"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SimpleChaincode struct {
	contractapi.Contract
}

// WS ATM struct
type ATM struct {
	ObjectType string `json:"docType"`
	ATMID string `json:"atmid"`
	CID string `json:"cid"`
	TBal int `json:"total"`
}

// type Customer struct {
// 	ObjectType string `json:"docType"`
// 	Customer_ID string `json:"customer_id"`
// 	Total_balance int `json:"total_balance"`
// 	Value int `json:"value"`
// }

type HistoryQueryResult struct {
	Record    *ATM    `json:"record"`
	TxId     string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}



// ReadATM 함수
func (t *SimpleChaincode) READ(ctx contractapi.TransactionContextInterface, atmid string) (*ATM, error) {
	ATMJSONasBytes, err := ctx.GetStub().GetState(atmid)
	if err != nil {
		return nil, fmt.Errorf("failed to get atm %s: %v", atmid, err)
	}
	if ATMJSONasBytes == nil {
		return nil, fmt.Errorf("atm %s does not exist", atmid)
	}

	var atm ATM
	err = json.Unmarshal(ATMJSONasBytes, &atm)
	if err != nil {
		return nil, err
	}

	return &atm, nil
}
// Input_ATM 함수
func (t *SimpleChaincode) INPUT(ctx contractapi.TransactionContextInterface, atmid string, cid string,  value int) error {
	fmt.Println("- start input money")
	// ATM 입금
	ATMAsBytes, err := ctx.GetStub().GetState(atmid)
	if err != nil {
		return fmt.Errorf("Failed to get money: " +err.Error())
	}
	atm := ATM{ObjectType:"ATM"}
	if ATMAsBytes != nil {
		//return fmt.Errorf("This money does not exists: " + atmid)
		_ = json.Unmarshal(ATMAsBytes, &atm)
	} else {
		atm.ATMID=atmid
	}

	// 검증 BAL + VALUE > MAX 오류

	// atm 증액
	atm.CID = cid
	atm.TBal += value

	ATMJSONasBytes, err := json.Marshal(atm)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(atmid, ATMJSONasBytes)
	if err != nil {
		return err
	}
	return nil
}

// Output_ATM 함수
func (t *SimpleChaincode) OUTPUT(ctx contractapi.TransactionContextInterface, atmid string, cid string,  value int) error {
	fmt.Println("- start output money")
	// ATM 출금
	ATMAsBytes, err := ctx.GetStub().GetState(atmid)
	if err != nil {
		return fmt.Errorf("Failed to get money: " +err.Error())
	}
	atm := ATM{ObjectType:"ATM"}
	if ATMAsBytes != nil {
		//return fmt.Errorf("This money does not exists: " + atmid)
		_ = json.Unmarshal(ATMAsBytes, &atm)
	} else {
		atm.ATMID=atmid
	}

	// 검증 BAL + VALUE < MIN 오류

	// atm 감액
	atm.CID = cid
	atm.TBal -= value

	ATMJSONasBytes, err := json.Marshal(atm)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(atmid, ATMJSONasBytes)
	if err != nil {
		return err
	}
	return nil
}

// ATMHistory 함수
func (t *SimpleChaincode) ATMHistory(ctx contractapi.TransactionContextInterface, atmid string) ([]HistoryQueryResult, error) {
	log.Printf("ATMHistory: ID %v", atmid)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(atmid)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var atm ATM
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &atm)
			if err != nil {
				return nil, err
			}
		} else {
			atm = ATM{
				ATMID: atmid,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &atm,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating atm chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting atm chaincode: %v", err)
	}
}