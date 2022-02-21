package chaincode

import (
	"fmt"

	"github.com/dmonteroh/distributed-resources-smartcontract/latency-sc/internal"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []internal.LatencyAsset{}

	for _, asset := range assets {
		err := ctx.GetStub().PutState(asset.ID, []byte(asset.String()))
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetKey string) (internal.LatencyAsset, error) {
	assetJson, err := ctx.GetStub().GetState(assetKey)
	if err != nil {
		return internal.LatencyAsset{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJson == nil {
		return internal.LatencyAsset{}, fmt.Errorf("the Asset with key: %s does not exist", assetKey)
	}

	asset, err := internal.LatencyAssetJsonToStruct(string(assetJson))
	if err != nil {
		return internal.LatencyAsset{}, err
	}

	return asset, nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetJson string) error {
	asset, err := internal.LatencyAssetJsonToStruct(assetJson)
	if err != nil {
		return err
	}
	exists, err := s.AssetExists(ctx, asset.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Asset for %s already exists", asset.ID)
	}

	// RUN VALIDATIONS
	if len(asset.Results) == 0 {
		return fmt.Errorf("no latency results were posted, ignored")
	}
	validJson := []byte(asset.String())

	return ctx.GetStub().PutState(asset.ID, validJson)
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetJson string) error {
	asset, err := internal.LatencyAssetJsonToStruct(assetJson)
	if err != nil {
		return err
	}
	exists, err := s.AssetExists(ctx, asset.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Asset for %s already exists", asset.ID)
	}

	// RUN VALIDATION
	validJson := []byte(asset.String())

	return ctx.GetStub().PutState(asset.ID, validJson)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetKey string) error {
	exists, err := s.AssetExists(ctx, assetKey)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Stats for %s do not exist", assetKey)
	}

	return ctx.GetStub().DelState(assetKey)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, assetKey string) (bool, error) {
	statJSON, err := ctx.GetStub().GetState(assetKey)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return statJSON != nil, nil
}

// CURRENTLY IN TO-DO (NO OWNERSHIP REQUIRED ATM)
// TransferAsset updates the owner field of asset with given id in world state.
// func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, statIP string, newStatIP string) (string, error) {
// 	statObject, err := s.ReadAsset(ctx, statIP)
// 	if err != nil {
// 		return "", err
// 	}

// 	statObject.ID = newStatIP

// 	return statObject.String(), ctx.GetStub().PutState(newStatIP, []byte(statObject.String()))
// }

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]internal.LatencyAsset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return iteratorSlicer(resultsIterator)
}

func iteratorSlicer(resultsIterator shim.StateQueryIteratorInterface) ([]internal.LatencyAsset, error) {
	var assets []internal.LatencyAsset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		asset, err := internal.LatencyAssetJsonToStruct(string(queryResponse.Value))
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

// INVETORY SMART CONTRACT INVOKATION
func (s *SmartContract) GetServerAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	params := []string{"GetServerAssets"}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("inventory-sc", queryArgs, "mychannel")
	if response.Status != shim.OK {
		return nil, fmt.Errorf("failed to query chaincode. Error %s", response.Payload)
	}

	assetArray, err := internal.JsonToAssetArray(string(response.GetPayload()))
	if err != nil {
		return nil, fmt.Errorf("failed to query chaincode. Error %s", err)
	}
	return assetArray, nil
}

func (s *SmartContract) GetServerAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	params := []string{"GetServerAssetsExceptId", excludeId}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("inventory-sc", queryArgs, "mychannel")
	if response.Status != shim.OK {
		return nil, fmt.Errorf("failed to query chaincode. Error %s", response.Payload)
	}

	assetArray, err := internal.JsonToAssetArray(string(response.GetPayload()))
	if err != nil {
		return nil, fmt.Errorf("failed to query chaincode. Error %s", err)
	}
	return assetArray, nil
}

// func iteratorSlicerAsset(resultsIterator shim.StateQueryIteratorInterface) ([]internal.Asset, error) {
// 	var assets []internal.Asset
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}
// 		asset, err := internal.JsonToAsset(string(queryResponse.Value))
// 		if err != nil {
// 			return nil, err
// 		}
// 		assets = append(assets, asset)
// 	}

// 	return assets, nil
// }
