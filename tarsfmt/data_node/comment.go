package data_node

import (
	"bytes"
	"strings"
)

type Comment struct {
	comment string
	level   int
}

// FIXME 不能允许有 正常语句前有注释的情况，如 /*******/ module M {
func NewComment(level int) Node {
	return &Comment{level: level}
}

func (c *Comment) Parse(terms []string, isEnd bool) error {
	c.comment = parseComment(terms)
	return nil
}

func (c *Comment) Print(buff *bytes.Buffer) {
	for i := 0; i < c.level; i++ {
		buff.WriteString("    ")
	}
	buff.WriteString(c.comment)
}

func (c *Comment) Reset() {
	panic("implement me")
}

func parseComment(terms []string) string {
	for idx, term := range terms {
		// 判断是不是最后一个字符
		if idx+1 == len(terms) {
			return ""
		}

		if term == "//" {
			return strings.Join(terms[idx:], JOIN_CHAR)
		} else if term == "/*" {
			if terms[len(terms)-1] == "*/" {
				return strings.Join(terms[idx:len(terms)], JOIN_CHAR)
			} else {
				return strings.Join(terms[idx:], JOIN_CHAR)
			}
		}
	}

	return ""
}
