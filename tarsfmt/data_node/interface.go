package data_node

import (
	"bytes"
	"fmt"

	"github.com/trigger3/toy/tarsfmt/key_words"
	"github.com/trigger3/toy/tarsfmt/util"
)

type Interface struct {
	nodes      *ExistNodes
	header     *StandardStatement
	statements []StandardStatement
	tail       *StandardStatement

	keyWordsMgr *key_words.KeyWordsMgr
	isBegin     bool

	headerMinLen int
	stateMinLen  int
	tailMinLen   int // 1
}

func NewInterface(nodes *ExistNodes, keyWordsMgr *key_words.KeyWordsMgr) Node {
	return &Interface{
		nodes:        nodes,
		keyWordsMgr:  keyWordsMgr,
		isBegin:      true,
		headerMinLen: 3,
		stateMinLen:  6,
		tailMinLen:   1,
	}
}

func (i *Interface) Parse(terms []string, isEnd bool) error {
	if isEnd {
		return i.parseTail(terms)
	}
	if i.isBegin {
		return i.parseHeader(terms)
	}
	return i.parseBody(terms)

}

func (i *Interface) parseHeader(terms []string) error {
	if len(terms) < i.headerMinLen {
		return fmt.Errorf("%v fromat invalid", terms[0])
	}

	enumName := terms[1]
	if util.IsNumeric(enumName) {
		return ErrStructNameCannotBeNumeric
	}
	i.nodes.enums.Add(enumName)

	state := &StandardStatement{
		Level:     1,
		Statement: fmt.Sprintf("interface %v {", enumName),
		Comment:   parseComment(terms[i.headerMinLen-1:]),
	}

	i.header = state
	i.isBegin = false

	return nil
}

func (i *Interface) parseTail(terms []string) error {
	state := &StandardStatement{
		Statement: "};",
		Level:     1,
		Comment:   parseComment(terms[i.tailMinLen-1:]),
	}

	i.tail = state
	return nil
}

// Resp DoReq(Req req); //xxx
func (i *Interface) parseBody(terms []string) error {
	if i.keyWordsMgr.IsCommentWord(terms[0]) {
		return i.parseComment(terms)
	}
	if len(terms) < i.stateMinLen {
		return ErrSyntexInvaild
	}
	if terms[2] != "(" || terms[5] != ")" {
		return ErrSyntexInvaild
	}
	intfResp, intfFunc, intfReq, intfReqName := terms[0], terms[1], terms[3], terms[4]
	if !i.nodes.structs.Contains(intfResp) && !i.keyWordsMgr.IsBasicType(intfResp) {
		return ErrInterfaceRespNotDefine
	}
	if util.IsNumeric(intfFunc) {
		return ErrInterfaceFuncNameCannotBeNumeric
	}
	if !i.nodes.structs.Contains(intfReq) && !i.keyWordsMgr.IsBasicType(intfReq) {
		return ErrInterfaceReqNotDefine
	}
	if util.IsNumeric(intfReqName) {
		return ErrInterfaceReqNameCannotBeNumeric
	}
	state := StandardStatement{
		Statement: fmt.Sprintf("%v %v(%v %v);", intfResp, intfFunc, intfReq, intfReqName),
		Comment:   parseComment(terms[i.stateMinLen-1:]),
		Level:     2,
	}

	i.statements = append(i.statements, state)

	return nil
}

func (i *Interface) parseComment(terms []string) error {
	state := StandardStatement{
		Statement: "",
		Comment:   parseComment(terms),
		Level:     2,
	}
	i.statements = append(i.statements, state)

	return nil
}

//func (i *Interface) formatNewTerm(terms []string) []string {
//	// Resp DoReq(Req req); //xxx
//	basicType := terms[3]
//	if basicType == "unsigned" {
//		if i.keyWordsMgr.IsBasicType(terms[4]) {
//			basicType = basicType + " " + terms[4]
//			return append(terms[:2], append([]string{basicType}, terms[5:]...)...)
//
//		}
//	}
//}

func (i *Interface) Print(buff *bytes.Buffer) {
	i.header.Format(buff, len(i.header.Statement))
	buff.WriteByte('\n')

	var maxStateLen int
	for _, state := range i.statements {
		if len(state.Comment) > 0 && len(state.Statement) > maxStateLen {
			maxStateLen = len(state.Statement)
		}
	}
	for _, state := range i.statements {
		state.Format(buff, maxStateLen)
		buff.WriteByte('\n')
	}

	i.tail.Format(buff, len(i.tail.Statement))
	buff.WriteByte('\n')
}

func (i *Interface) Reset() {
	panic("implement me")
}
