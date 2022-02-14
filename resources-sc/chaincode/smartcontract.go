package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/dmonteroh/distributed-resources-smartcontract/resources-sc/internal"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stats := []internal.StoredStat{}

	for _, stat := range stats {
		err := ctx.GetStub().PutState(stat.DrcHost.HostID, []byte(stat.String()))
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, statIP string, statJSON string) error {
	exists, err := s.AssetExists(ctx, statIP)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Stats for %s already exists", statIP)
	}

	tmpStat, err := internal.DrcJsonToStruct(statJSON)
	if err != nil {
		return err
	}
	toStore := internal.ConvertToStorage(tmpStat)
	toStore.ID = statIP
	// RUN VALIDATION

	return ctx.GetStub().PutState(statIP, []byte(toStore.String()))
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, statIP string) (*internal.StoredStat, error) {
	statJSON, err := ctx.GetStub().GetState(statIP)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if statJSON == nil {
		return nil, fmt.Errorf("the Stats for %s do not exist", statIP)
	}

	var stat internal.StoredStat
	err = json.Unmarshal(statJSON, &stat)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, statIP string, statJSON string) error {
	exists, err := s.AssetExists(ctx, statIP)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Stats for %s do not exist", statIP)
	}

	tmpStat, err := internal.DrcJsonToStruct(statJSON)
	if err != nil {
		return err
	}
	toStore := internal.ConvertToStorage(tmpStat)
	toStore.ID = statIP

	return ctx.GetStub().PutState(statIP, []byte(toStore.String()))
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, statIP string) error {
	exists, err := s.AssetExists(ctx, statIP)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Stats for %s do not exist", statIP)
	}

	return ctx.GetStub().DelState(statIP)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, statIP string) (bool, error) {
	statJSON, err := ctx.GetStub().GetState(statIP)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return statJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, statIP string, newStatIP string) (string, error) {
	statObject, err := s.ReadAsset(ctx, statIP)
	if err != nil {
		return "", err
	}

	statObject.ID = newStatIP

	return statObject.String(), ctx.GetStub().PutState(newStatIP, []byte(statObject.String()))
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*internal.StoredStat, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var statObjects []*internal.StoredStat
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var statObject internal.StoredStat
		err = json.Unmarshal(queryResponse.Value, &statObject)
		if err != nil {
			return nil, err
		}
		statObjects = append(statObjects, &statObject)
	}

	return statObjects, nil
}
