package parser

import (
	"encoding/json"

	"github.com/smomara/gossamer/router"
)

func ParseJSONBody(r *router.Request, v interface{}) error {
	return json.Unmarshal(r.Body, v)
}
