package handler

type Node struct {
	ID    int      `json:"id"` // Altere para um inteiro incremental
	Label string   `json:"label"`
	Links []string `json:"-"`
}

type Edge struct {
	From int `json:"from"` // ID do nó de origem
	To   int `json:"to"`   // ID do nó de destino
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
