package model

type Completion struct {
	Name    string   `json:"name" yaml:"name"`
	Doc     string   `json:"doc" yaml:"doc"`
	Symbols []Symbol `json:"symbols" yaml:"symbols"`
}

type Symbol struct {
	Name   string `json:"name" yaml:"name"`
	Kind   string `json:"kind" yaml:"kind"`
	Detail string `json:"detail" yaml:"detail"`
	Doc    string `json:"doc,omitempty" yaml:"doc,omitempty"`
}
