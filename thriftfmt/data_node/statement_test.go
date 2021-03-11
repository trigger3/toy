package data_node

import (
	"bytes"
	"testing"
)

func TestStandardStatement_format(t *testing.T) {
	type fields struct {
		Statement string
		Comment   string
		Level     int
	}
	type args struct {
		buff            *bytes.Buffer
		maxStatementLen int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "test level",
			fields: fields{
				Statement: "2 require map<int,int> id;",
				Comment:   "test",
				Level:     1,
			},
			args: args{
				buff:            bytes.NewBuffer(nil),
				maxStatementLen: 40,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//ss := &StandardStatement{
			//	Statement: tt.fields.Statement,
			//	Comment:   tt.fields.Comment,
			//	Level:     tt.fields.Level,
			//}
		})
	}
}
