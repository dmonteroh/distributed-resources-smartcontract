package internal

import "encoding/json"

// INVENTORY ASSET
type Asset struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Owner      string            `json:"owner"`
	Type       int               `json:"type"`       //[0: Server, 1: Sensor, 2: Robot]
	State      int               `json:"state"`      //[0: Disabled, 1: Enabled]
	Properties map[string]string `json:"properties"` //{GPU: TRUE ...}
}

func (d Asset) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func JsonToAsset(v string) (asset Asset) {
	json.Unmarshal([]byte(v), &asset)
	return asset
}
