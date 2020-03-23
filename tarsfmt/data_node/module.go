package data_node

import (
	"bytes"
	"fmt"

	"github.com/trigger3/toy/tarsfmt/key_words"
	"github.com/trigger3/toy/tarsfmt/util"
)

type Module struct {
	header *StandardStatement
	tail   *StandardStatement
	nodes  []Node

	nodesName *ExistNodes

	keyWordsMgr *key_words.KeyWordsMgr
	isBegin     bool

	headerMinLen int
	stateMinLen  int
	tailMinLen   int // 1

	lastNode Node

	curDataType int8 // 代码数据类型
}

func NewModule(nodes *ExistNodes, mgr *key_words.KeyWordsMgr) Node {
	return &Module{
		keyWordsMgr: mgr,
		nodesName:   nodes,
		isBegin:     true,

		headerMinLen: 3,
		tailMinLen:   1,
	}
}

func (m *Module) Reset() {
	panic("implement me")
}

func (m *Module) Parse(terms []string, isEnd bool) error {
	if isEnd {
		return m.parseTail(terms)
	}
	if m.isBegin {
		return m.parseHeader(terms)
	}

	return m.parseBody(terms)

}

func (m *Module) parseHeader(terms []string) error {
	if len(terms) < m.headerMinLen {
		return fmt.Errorf("%v fromat invalid", terms[0])
	}

	moduleName := terms[1]
	if util.IsNumeric(moduleName) {
		return ErrStructNameCannotBeNumeric
	}

	state := &StandardStatement{
		Level:     0,
		Statement: fmt.Sprintf("module %v {", moduleName),
		Comment:   parseComment(terms[m.headerMinLen-1:]),
	}

	m.header = state
	m.isBegin = false

	return nil
}

func (m *Module) parseTail(terms []string) error {
	state := &StandardStatement{
		Statement: "};",
		Level:     0,
		Comment:   parseComment(terms[m.tailMinLen-1:]),
	}

	m.tail = state
	return nil
}

func (m *Module) Print(buff *bytes.Buffer) {
	stateLen := len(m.header.Statement)
	m.header.Format(buff, stateLen)
	buff.WriteByte('\n')

	for i, node := range m.nodes {
		node.Print(buff)
		if i != len(m.nodes)-1 {
			buff.WriteByte('\n')
		}
	}

	m.tail.Format(buff, len(m.tail.Statement))
	buff.WriteByte('\n')
}

func (m *Module) parseBody(terms []string) error {
	dataType := m.getCodeType(terms)
	// dataType 必然不为0
	if m.lastNode == nil {
		//m.lastNode = m.getNodeParser(dataType)
		node := m.getNodeParser(dataType)
		if err := node.Parse(terms, false); err != nil {
			return err
		}
		if dataType == key_words.TYPE_COMMENT {
			m.nodes = append(m.nodes, node)
		} else {
			m.lastNode = node
		}
		return nil
	}

	if dataType != key_words.TYPE_END {
		return m.lastNode.Parse(terms, false)
	}

	if err := m.lastNode.Parse(terms, true); err != nil {
		return err
	}
	m.nodes = append(m.nodes, m.lastNode)
	m.lastNode = nil
	return nil
}

func (m *Module) getCodeType(terms []string) int8 {
	codeType := m.keyWordsMgr.StatementType(terms[0])
	if codeType == key_words.TYPE_NIL {
		return m.curDataType
	}

	if codeType == key_words.TYPE_END {
		m.curDataType = key_words.TYPE_NIL
	} else {
		m.curDataType = codeType
	}

	return codeType
}

func (m *Module) getNodeParser(dataType int8) Node {
	switch dataType {
	case key_words.TYPE_STRUCT:
		return NewStructNode(m.nodesName, m.keyWordsMgr)
	case key_words.TYPE_ENUM:
		return NewEnumNode(m.nodesName, m.keyWordsMgr)
	case key_words.TYPE_INTERFACE:
		return NewInterface(m.nodesName, m.keyWordsMgr)
	case key_words.TYPE_COMMENT:
		return NewComment(1)
	default:
		return nil
	}
}
