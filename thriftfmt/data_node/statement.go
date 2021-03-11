package data_node

import "bytes"

type StandardStatement struct {
	Statement string // 归一化后不包含注释的语句
	Comment   string // 注释语句
	Level     int    // 层级
}

func (s *StandardStatement) Format(buff *bytes.Buffer, maxLen int) {
	for i := 0; i < s.Level; i++ {
		buff.WriteString("    ")
	}
	buff.WriteString(s.Statement)
	if len(s.Comment) != 0 && len(s.Statement) != 0 {
		for i := 0; i < maxLen-len(s.Statement); i++ {
			buff.WriteByte(' ')
		}
		buff.WriteByte(' ')
	}

	buff.WriteString(s.Comment)
}

func defaultTail(level int) *StandardStatement {
	return &StandardStatement{
		Statement: "};",
		Comment:   "",
		Level:     level,
	}
}
