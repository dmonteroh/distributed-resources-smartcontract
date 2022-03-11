package internal

import (
	"encoding/json"

	"github.com/wI2L/jettison"
)

// INVENTORY ASSET
type Asset struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Owner      string            `json:"owner"`
	Type       int               `json:"type"`       //[0: Server, 1: Robot, 2: Sensor]
	State      int               `json:"state"`      //[0: Disabled, 1: Enabled]
	Properties map[string]string `json:"properties"` //{GPU: TRUE ...}
}

func (d Asset) String() string {
	s, _ := jettison.MarshalOpts(d, jettison.NilMapEmpty(), jettison.NilSliceEmpty())
	return string(s)
}

func JsonToAsset(v string) (asset Asset, err error) {
	err = json.Unmarshal([]byte(v), &asset)
	return asset, err
}

func JsonToAssetArray(v string) (assets []Asset, err error) {
	err = json.Unmarshal([]byte(v), &assets)
	return assets, err
}
