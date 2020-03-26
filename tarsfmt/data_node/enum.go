package data_node

import (
	"bytes"
	"fmt"

	"github.com/trigger3/toy/tarsfmt/key_words"
	"github.com/trigger3/toy/tarsfmt/util"
)

type Enum struct {
	nodes      *ExistNodes
	header     *StandardStatement
	tail       *StandardStatement
	statements []StandardStatement

	keyWordsMgr *key_words.KeyWordsMgr
	isBegin     bool

	headerMinLen int
	stateMinLen  int
	tailMinLen   int // 1
}

func NewEnumNode(nodes *ExistNodes, keyWordsMgr *key_words.KeyWordsMgr) Node {
	return &Enum{
		nodes:        nodes,
		keyWordsMgr:  keyWordsMgr,
		isBegin:      true,
		headerMinLen: 3,
		stateMinLen:  3,
		tailMinLen:   1,
	}
}

func (e *Enum) Reset() {
	panic("implement me")
}

func (e *Enum) Parse(terms []string, isEnd bool) error {
	if isEnd {
		return e.parseTail(terms)
	}
	if e.isBegin {
		return e.parseHeader(terms)
	}
	return e.parseBody(terms)
}

func (e *Enum) parseHeader(terms []string) error {
	if len(terms) < e.headerMinLen {
		return fmt.Errorf("%v fromat invalid", terms[0])
	}

	enumName := terms[1]
	if util.IsNumeric(enumName) {
		return ErrStructNameCannotBeNumeric
	}
	e.nodes.enums.Add(enumName)

	state := &StandardStatement{
		Level:     1,
		Statement: fmt.Sprintf("enum %v {", enumName),
		Comment:   parseComment(terms[e.headerMinLen-1:]),
	}

	e.header = state
	e.isBegin = false

	return nil
}

func (e *Enum) parseTail(terms []string) error {
	state := &StandardStatement{
		Statement: "};",
		Level:     1,
		Comment:   parseComment(terms[e.tailMinLen-1:]),
	}

	e.tail = state
	return nil
}

// WEIBO = 3,
func (e *Enum) parseBody(terms []string) error {
	e.isBegin = false
	if e.keyWordsMgr.IsCommentWord(terms[0]) {
		return e.parseComment(terms)
	}

	if len(terms) < e.stateMinLen {
		return ErrSyntexInvaild
	}
	enumKey, enumValue := terms[0], terms[2]
	if util.IsNumeric(enumKey) {
		return ErrEnumNameCannotBeNumeric
	}
	if !util.IsNumeric(enumValue) {
		return ErrEnumValueMustBeNumeric
	}
	state := StandardStatement{
		Statement: fmt.Sprintf("%v = %v,", enumKey, enumValue),
		Comment:   parseComment(terms[e.stateMinLen-1:]),
		Level:     2,
	}

	e.statements = append(e.statements, state)

	return nil
}

func (e *Enum) parseComment(terms []string) error {
	state := StandardStatement{
		Statement: "",
		Comment:   parseComment(terms),
		Level:     2,
	}
	e.statements = append(e.statements, state)

	return nil
}

func (e *Enum) Print(buff *bytes.Buffer) {
	e.header.Format(buff, len(e.header.Statement))
	buff.WriteByte('\n')

	var maxStateLen int
	for _, state := range e.statements {
		if len(state.Comment) > 0 && len(state.Statement) > maxStateLen {
			maxStateLen = len(state.Statement)
		}
	}
	for _, state := range e.statements {
		state.Format(buff, maxStateLen)
		buff.WriteByte('\n')
	}

	e.tail.Format(buff, len(e.tail.Statement))
	buff.WriteByte('\n')
}
