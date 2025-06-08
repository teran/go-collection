package helm

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type Resource map[string]any

func (r Resource) get(key string) gjson.Result {
	v, err := json.Marshal(r)
	if err != nil {
		panic("internal error: incompatible structure received")
	}

	return gjson.GetBytes(v, key)
}

func (r Resource) GetString(key string) string {
	return r.get(key).String()
}

func (r Resource) GetNumber(key string) float64 {
	return r.get(key).Float()
}

func (r Resource) GetBoolean(key string) bool {
	return r.get(key).Bool()
}

func (r Resource) IsExists(key string) bool {
	return r.get(key).Exists()
}

func (r Resource) GetStruct(key string, in any) error {
	return errors.Wrap(
		json.NewDecoder(strings.NewReader(r.get(key).Raw)).Decode(in),
		"error unmarshaling the raw data",
	)
}
