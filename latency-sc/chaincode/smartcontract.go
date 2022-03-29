package chaincode

import (
	"fmt"
	"time"

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

	// RUN VALIDATIONS
	if len(asset.Results) == 0 {
		return fmt.Errorf("no latency results were posted, ignored")
	}
	if asset.ID == "" {
		return fmt.Errorf("latency results was posted without ID, ignored")
	}

	exists, err := s.AssetExists(ctx, asset.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Asset for %s already exists", asset.ID)
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

	// RUN VALIDATIONS
	if len(asset.Results) == 0 {
		return fmt.Errorf("no latency results were posted, ignored")
	}
	if asset.ID == "" {
		return fmt.Errorf("latency results was posted without ID, ignored")
	}

	exists, err := s.AssetExists(ctx, asset.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the Asset for %s already exists", asset.ID)
	}

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

func stringQuery(ctx contractapi.TransactionContextInterface, queryString string) ([]internal.LatencyAsset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return iteratorSlicer(resultsIterator)
}

func (s *SmartContract) GetAssetListTimeSource(ctx contractapi.TransactionContextInterface, source string, minutes int) ([]internal.LatencyAsset, error) {
	timeStart := time.Now()
	timeEnd := timeStart.Add(time.Duration(-time.Duration(minutes) * time.Minute))
	assetQuery := fmt.Sprintf(`{"selector":{"source":"%s","$timestamp.timeSeconds":{"$lt": %d,"$gte": %d}}}`, source, timeStart.Unix(), timeEnd.Unix())
	return stringQuery(ctx, assetQuery)
}

func (s *SmartContract) GetAssetListTimeTarget(ctx contractapi.TransactionContextInterface, target string, minutes int) ([]internal.LatencyAsset, error) {
	timeStart := time.Now()
	timeEnd := timeStart.Add(time.Duration(-time.Duration(minutes) * time.Minute))
	assetQuery := fmt.Sprintf(`{"selector":{"results":{"$elemMatch": {"hostname": "%s"}},"$timestamp.timeSeconds":{"$lt": %d,"$gte": %d}}}`, target, timeStart.Unix(), timeEnd.Unix())
	return stringQuery(ctx, assetQuery)
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

func (s *SmartContract) GetRobotAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	params := []string{"GetRobotAssets"}
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

func (s *SmartContract) GetRobotAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	params := []string{"GetRobotAssetsExceptId", excludeId}
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

func (s *SmartContract) GetSensorAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	params := []string{"GetSensorAssets"}
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

func (s *SmartContract) GetSensorAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	params := []string{"GetSensorAssetsExceptId", excludeId}
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

func (s *SmartContract) GetSensorAndRobotAssets(ctx contractapi.TransactionContextInterface) ([]internal.Asset, error) {
	params := []string{"GetSensorAndRobotAssets"}
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

func (s *SmartContract) GetSensorAndRobotAssetsExceptId(ctx contractapi.TransactionContextInterface, excludeId string) ([]internal.Asset, error) {
	params := []string{"GetSensorAndRobotAssetsExceptId", excludeId}
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
