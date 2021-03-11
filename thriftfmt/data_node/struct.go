package data_node

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/trigger3/toy/thriftfmt/key_words"
	"github.com/trigger3/toy/thriftfmt/util"
)

type Struct struct {
	nodes      *ExistNodes
	header     *StandardStatement
	statements []StandardStatement
	tail       *StandardStatement

	keyWordsMgr *key_words.KeyWordsMgr
	isBegin     bool

	headerMinLen int
	stateMinLen  int
	tailMinLen   int // 1

	elementFlag int64 // 初始化-1
}

func NewStructNode(nodes *ExistNodes, keyWordsMgr *key_words.KeyWordsMgr) Node {
	return &Struct{
		nodes:        nodes,
		keyWordsMgr:  keyWordsMgr,
		isBegin:      true,
		headerMinLen: 3,
		stateMinLen:  4,
		tailMinLen:   1,
		elementFlag:  -1,
	}
}
func (s *Struct) Reset() {
	s.header = nil
	s.isBegin = true
	s.statements = s.statements[:0]
}

func (s *Struct) Parse(terms []string, isEnd bool) error {
	if isEnd {
		return s.parseTail(terms)
	}
	if s.isBegin {
		return s.parseHeader(terms)
	}
	return s.parseBody(terms)

}

func (s *Struct) parseHeader(terms []string) error {
	if len(terms) < s.headerMinLen {
		return fmt.Errorf("%v fromat invalid", terms[0])
	}

	structName := terms[1]
	if util.IsNumeric(structName) {
		return ErrStructNameCannotBeNumeric
	}
	s.nodes.structs.Add(structName)

	state := &StandardStatement{
		Level:     1,
		Statement: fmt.Sprintf("struct %v {", structName),
		Comment:   parseComment(terms[s.headerMinLen-1:]),
	}

	s.header = state
	s.isBegin = false

	return nil
}

func (s *Struct) parseTail(terms []string) error {
	state := &StandardStatement{
		Statement: "};",
		Level:     1,
		Comment:   parseComment(terms[s.tailMinLen-1:]),
	}

	s.tail = state
	return nil
}

func (s *Struct) parseBody(terms []string) error {
	s.isBegin = false

	if s.keyWordsMgr.IsCommentWord(terms[0]) {
		return s.parseComment(terms)
	}

	termsLen := len(terms)
	if termsLen < s.stateMinLen {
		return ErrSyntexInvaild
	}

	// 将类型位归一化，将得到类型只占一个term的数组
	newTerms, err := s.formatElementType(terms)
	if err != nil {
		return err
	}

	return s.formatNewTerm(newTerms)
}

func (s *Struct) parseComment(terms []string) error {
	state := StandardStatement{
		Statement: "",
		Comment:   parseComment(terms),
		Level:     2,
	}
	s.statements = append(s.statements, state)

	return nil
}

func (s *Struct) formatElementType(terms []string) ([]string, error) {
	dataType := terms[2]
	if s.keyWordsMgr.IsBasicType(dataType) {
		return terms, nil
	} else if dataType == "unsigned" {
		return s.formatUnsigned(terms)
	} else if s.nodes.structs.Contains(dataType) || s.nodes.enums.Contains(dataType) {
		return terms, nil
	} else if s.keyWordsMgr.IsContainerType(terms[2]) {
		// 15 require vector<int> head;  // 更新时间
		return s.formatContainer(terms)
	} else if s.nodes.modules.Contains(terms[2]) {
		return s.formatModule(terms)
	} else {
		return nil, ErrSyntexInvaild
	}
}

func (s *Struct) formatNewTerm(terms []string) (err error) {
	if len(terms) < s.stateMinLen {
		return ErrSyntexInvaild
	}
	// 15 require common::head head;  // 更新时间
	// 调整元素代码
	elementFlag, err := strconv.ParseInt(terms[0], 10, 64)
	if err != nil {
		return ErrSyntexInvaild
	}
	if elementFlag <= s.elementFlag {
		elementFlag = s.elementFlag + 1
	}
	// 元素标识
	terms[0] = strconv.FormatInt(elementFlag, 10)
	s.elementFlag = elementFlag
	// 规则
	if !s.keyWordsMgr.IsElementRule(terms[1]) {
		return ErrElemetRuleInvalid
	}
	// 元素名称
	if util.IsNumeric(terms[3]) {
		return ErrVariableNameCannotBeNumeric
	}

	state := StandardStatement{
		Statement: strings.Join(terms[:s.stateMinLen], JOIN_CHAR) + ";",
		Comment:   parseComment(terms[s.stateMinLen-1:]),
		Level:     2,
	}

	s.statements = append(s.statements, state)

	return nil

}

// 0 require vector<int> id; // test
func (s *Struct) formatContainer(terms []string) ([]string, error) {
	if terms[3] != "<" {
		return nil, ErrSyntexInvaild
	}
	var (
		hasEnd bool
		endIdx int
	)
	for idx, t := range terms[3:] {
		if t == ">" {
			hasEnd = true
			endIdx = idx + 3
			break
		}
	}
	if !hasEnd || len(terms) == endIdx+1 {
		return nil, ErrSyntexInvaild
	}
	newType := strings.Join(terms[2:endIdx+1], "")

	return append(terms[:2], append([]string{newType}, terms[endIdx+1:]...)...), nil

}

func (s *Struct) formatModule(terms []string) ([]string, error) {
	// 15 require common::head head;  // 更新时间
	if terms[3] != "::" || len(terms) < 7 {
		return nil, ErrSyntexInvaild
	}

	formatType := terms[2] + "::" + terms[4]
	return append(terms[:2], append([]string{formatType}, terms[5:]...)...), nil
}

func (s *Struct) formatUnsigned(terms []string) ([]string, error) {
	//  0 require unsigned  int a;
	basicType := "unsigned " + terms[3]
	if !s.keyWordsMgr.IsBasicType(basicType) {
		return nil, ErrSyntexInvaild
	}
	return append(terms[:2], append([]string{basicType}, terms[4:]...)...), nil
}

func (s *Struct) Print(buff *bytes.Buffer) {
	s.header.Format(buff, len(s.header.Statement))
	buff.WriteByte('\n')

	var maxStateLen int
	for _, state := range s.statements {
		if len(state.Comment) > 0 && len(state.Statement) > maxStateLen {
			maxStateLen = len(state.Statement)
		}
	}
	for _, state := range s.statements {
		state.Format(buff, maxStateLen)
		buff.WriteByte('\n')
	}
	if s.tail == nil {
		s.tail = defaultTail(1)
	}
	s.tail.Format(buff, len(s.tail.Statement))
	buff.WriteByte('\n')
}
