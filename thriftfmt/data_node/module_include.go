package data_node

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/trigger3/toy/thriftfmt/util"
)

type ModuleInclude struct {
	nodes      *ExistNodes
	moduleName string
	comment    string

	moduleFieldLen int // 7
}

func NewModuleInclude(nodes *ExistNodes) Node {
	return &ModuleInclude{
		nodes:          nodes,
		moduleFieldLen: 7,
	}
}

func (m *ModuleInclude) Reset() {
	m.moduleName = ""
	m.comment = ""
}

func (m *ModuleInclude) Parse(terms []string, isEnd bool) error {
	// eg #include "Common.tars"
	moduleTermsLen := len(terms)
	if moduleTermsLen < m.moduleFieldLen {
		return errors.New("module including format invalid")
	}

	if terms[1] != "include" || terms[2] != "\"" || terms[4] != "." || terms[5] != "tars" || terms[6] != "\"" {
		return errors.New("module including format invalid")
	}

	moduleName := terms[3]
	if util.IsNumeric(moduleName) {
		return ErrModuleNameCannotBeNumeric
	}
	m.moduleName = moduleName
	m.comment = parseComment(terms[m.moduleFieldLen-1:])

	m.nodes.modules.Add(moduleName)

	return nil
}

func (m *ModuleInclude) Print(buff *bytes.Buffer) {
	moduleInclude := fmt.Sprintf("#include \"%v.tars\" %v", m.moduleName, m.comment)
	buff.WriteString(moduleInclude)
	buff.WriteString("\n")

}
