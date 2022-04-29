package db

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/archi-dex/ingester/pkg/util"
	trace "github.com/hans-m-song/go-stacktrace"
)

var (
	ErrorInvalidEntity   = trace.New("ERROR_INVALID_ENTITY")
	ErrorMarshalEntity   = trace.New("ERROR_MARSHAL_ENTITY")
	ErrorUnmarshalEntity = trace.New("ERROR_UNMARSHAL_ENTITY")
)

type EntityMeta struct {
	Indicies []string
}

type EntityAttributes struct {
	Path       string                 `bson:"path"       json:"path"`
	Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
}

func EntityAttributesFromBytes(input []byte) (*EntityAttributes, error) {
	var parsed EntityAttributes
	if err := json.Unmarshal(input, &parsed); err != nil {
		return nil, ErrorUnmarshalEntity.Trace(err).Add("input", input)
	}

	return &parsed, nil
}

type Entity struct {
	ID         string                 `bson:"_id"        json:"_id,omitempty"`
	Dir        string                 `bson:"dir"        json:"dir"`
	Base       string                 `bson:"base"       json:"base"`
	Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updated_at"`
}

func NewEntity(path string, attributes map[string]interface{}) *Entity {
	now := time.Now().UTC()
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	return &Entity{
		ID:         util.NewHash(base),
		Dir:        dir,
		Base:       base,
		Attributes: attributes,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
