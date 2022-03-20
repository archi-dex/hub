package db

import (
	"encoding/json"
	"path"
	"time"

	"github.com/archi-dex/ingester/pkg/util"
	trace "github.com/hans-m-song/go-stacktrace"
)

type EntityType string

const (
	FileEntityType   EntityType = "FILE"
	FolderEntityType EntityType = "FOLDER"
)

var (
	ErrorMarshalEntity   = trace.New("ERROR_MARSHAL_ENTITY")
	ErrorUnmarshalEntity = trace.New("ERROR_UNMARSHAL_ENTITY")
)

type EntityMeta struct {
	Indicies []string
}

type Entity struct {
	ID         string            `bson:"_id"        json:"_id,omitempty"`
	Path       string            `bson:"path"       json:"path"`
	Type       EntityType        `bson:"type"       json:"type"`
	Attributes map[string]string `bson:"attributes" json:"attributes"`
	CreatedAt  time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time         `bson:"updated_at" json:"updated_at"`
}

func getEntityType(filepath string) EntityType {
	if path.Ext(filepath) == "" {
		return FolderEntityType
	}

	return FileEntityType
}

func NewEntity(path string, attributes map[string]string) *Entity {
	now := time.Now().UTC()
	return &Entity{
		ID:         util.NewHash(path),
		Path:       path,
		Type:       getEntityType(path),
		Attributes: attributes,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func EntityFromBytes(input []byte) (*Entity, error) {
	var parsed Entity
	if err := json.Unmarshal(input, &parsed); err != nil {
		return nil, ErrorUnmarshalEntity.Trace(err).Add("input", input)
	}

	return NewEntity(parsed.Path, parsed.Attributes), nil
}
