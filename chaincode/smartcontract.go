package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/dmonteroh/distributed-resources-smartcontract/internal"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stats := []internal.StatsObject{}

	for _, stat := range stats {
		err := ctx.GetStub().PutState(stat.Identity.IP, []byte(stat.String()))
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateStat(ctx contractapi.TransactionContextInterface, statIP string, statJSON string) error {
	exists, err := s.StatExists(ctx, statIP)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Stats for %s already exists", statIP)
	}
	return ctx.GetStub().PutState(statIP, []byte(statJSON))
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadStat(ctx contractapi.TransactionContextInterface, statIP string) (*internal.StatsObject, error) {
	statJSON, err := ctx.GetStub().GetState(statIP)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if statJSON == nil {
		return nil, fmt.Errorf("the Stats for %s do not exist", statIP)
	}

	var stat internal.StatsObject
	err = json.Unmarshal(statJSON, &stat)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateStat(ctx contractapi.TransactionContextInterface, statIP string, statJSON string) error {
	exists, err := s.StatExists(ctx, statIP)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Stats for %s do not exist", statIP)
	}
	return ctx.GetStub().PutState(statIP, []byte(statJSON))
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteLanding(ctx contractapi.TransactionContextInterface, dronename string) error {
	exists, err := s.StatExists(ctx, dronename)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Stats for %s do not exist", dronename)
	}

	return ctx.GetStub().DelState(dronename)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) StatExists(ctx contractapi.TransactionContextInterface, statIP string) (bool, error) {
	statJSON, err := ctx.GetStub().GetState(statIP)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return statJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferStat(ctx contractapi.TransactionContextInterface, statIP string, newStatIP string) (string, error) {
	statObject, err := s.ReadStat(ctx, statIP)
	if err != nil {
		return "", err
	}

	statObject.Identity.IP = newStatIP

	return statObject.String(), ctx.GetStub().PutState(newStatIP, []byte(statObject.String()))
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*internal.StatsObject, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var statObjects []*internal.StatsObject
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var statObject internal.StatsObject
		err = json.Unmarshal(queryResponse.Value, &statObject)
		if err != nil {
			return nil, err
		}
		statObjects = append(statObjects, &statObject)
	}

	return statObjects, nil
}
