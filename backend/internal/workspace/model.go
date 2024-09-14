package workspace

import (
	"github.com/google/uuid"
)

type NodeType string

const (
	InputNode   NodeType = "INPUT"
	OutputNode  NodeType = "OUTPUT"
	ProcessNode NodeType = "PROCESS"
)

type Node struct {
	ID       uuid.UUID              `json:"id"`
	Type     NodeType               `json:"type"`
	Label    string                 `json:"label"`
	Data     map[string]interface{} `json:"data"`
	Position Position               `json:"position"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Workspace struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Nodes []Node    `json:"nodes"`
	Edges []Edge    `json:"edges"`
}

type Edge struct {
	ID     uuid.UUID `json:"id"`
	Source uuid.UUID `json:"source"`
	Target uuid.UUID `json:"target"`
}
