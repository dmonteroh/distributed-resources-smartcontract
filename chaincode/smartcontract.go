package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/dmonteroh/fabric-distributed-resources/internal"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	landings := []internal.StatsObject{}

	for _, landing := range landings {
		landingJSON, err := json.Marshal(landing)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(landing.Dronename, landingJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateLanding(ctx contractapi.TransactionContextInterface, dronename string, landingxposition float64, landingyposition float64) error {
	exists, err := s.LandingExists(ctx, dronename)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the landing %s already exists", dronename)
	}

	landing := internal.StatsObject{}
	landingJSON, err := json.Marshal(landing)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dronename, landingJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadLanding(ctx contractapi.TransactionContextInterface, dronename string) (*internal.StatsObject, error) {
	landingJSON, err := ctx.GetStub().GetState(dronename)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if landingJSON == nil {
		return nil, fmt.Errorf("the landing %s does not exist", dronename)
	}

	var landing internal.StatsObject
	err = json.Unmarshal(landingJSON, &landing)
	if err != nil {
		return nil, err
	}

	return &landing, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateLanding(ctx contractapi.TransactionContextInterface, dronename string, landingxposition float64, landingyposition float64) error {
	exists, err := s.LandingExists(ctx, dronename)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the landing %s does not exist", dronename)
	}

	// overwriting original asset with new asset
	landing := internal.StatsObject{}
	landingJSON, err := json.Marshal(landing)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dronename, landingJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteLanding(ctx contractapi.TransactionContextInterface, dronename string) error {
	exists, err := s.LandingExists(ctx, dronename)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dronename)
	}

	return ctx.GetStub().DelState(dronename)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) LandingExists(ctx contractapi.TransactionContextInterface, dronename string) (bool, error) {
	landingJSON, err := ctx.GetStub().GetState(dronename)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return landingJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferLanding(ctx contractapi.TransactionContextInterface, dronename string, AnotherDrone string) error {
	landing, err := s.ReadLanding(ctx, dronename)
	if err != nil {
		return err
	}

	landing.Dronename = AnotherDrone
	landingJSON, err := json.Marshal(landing)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dronename, landingJSON)
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

	var landings []*internal.StatsObject
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var landing internal.StatsObject
		err = json.Unmarshal(queryResponse.Value, &landing)
		if err != nil {
			return nil, err
		}
		landings = append(landings, &landing)
	}

	return landings, nil
}
