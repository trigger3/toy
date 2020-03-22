package data_node

import (
	"bytes"
	"toy/mytype"
)

const (
	JOIN_CHAR = " "
)

type Node interface {
	Print(buff *bytes.Buffer)
	Reset()
	Parse(terms []string, isEnd bool) error
}

type ExistNodes struct {
	modules *mytype.Set
	structs *mytype.Set
	enums   *mytype.Set
}

func NewExistNodes() *ExistNodes {
	return &ExistNodes{
		modules: mytype.NewSet(),
		structs: mytype.NewSet(),
		enums:   mytype.NewSet(),
	}
}
