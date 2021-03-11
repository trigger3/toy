package line_parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/trigger3/toy/thriftfmt/key_words"
)

func Test_parserImp_getComment(t *testing.T) {
	keyWordsMgr := key_words.NewKeyWordsMgr()
	type fields struct {
		keyWordMgr *key_words.KeyWordsMgr
	}
	type args struct {
		term []rune
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantTerms []string
	}{
		// TODO: Add test cases.
		{
			name:      `test "//"`,
			fields:    fields{keyWordMgr: keyWordsMgr},
			args:      args{[]rune("// test")},
			wantTerms: []string{"//", "test"},
		},
		{
			name:      `test "/*"`,
			fields:    fields{keyWordMgr: keyWordsMgr},
			args:      args{[]rune("/* test */")},
			wantTerms: []string{"/*", "test", "*/"},
		},
		{
			name:      `test "/*"`,
			fields:    fields{keyWordMgr: keyWordsMgr},
			args:      args{[]rune("/* test */")},
			wantTerms: []string{"/*", "test", "*/"},
		},
		{
			name:      `test "/"`,
			fields:    fields{keyWordMgr: keyWordsMgr},
			args:      args{[]rune("/ test")},
			wantTerms: []string{"//", "test"},
		},
		{
			name:      `test "/****/"`,
			fields:    fields{keyWordMgr: keyWordsMgr},
			args:      args{[]rune("/***test***/")},
			wantTerms: []string{"/*", "**test**", "*/"},
		},
		{
			name:      `test "/****test/"`,
			fields:    fields{keyWordMgr: keyWordsMgr},
			args:      args{[]rune("/***test/")},
			wantTerms: []string{"/*", "**test", "*/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parserImp{
				keyWordMgr: tt.fields.keyWordMgr,
			}
			if gotTerms := p.getComment(tt.args.term); !reflect.DeepEqual(gotTerms, tt.wantTerms) {
				t.Errorf("getComment() = %v, want %v", gotTerms, tt.wantTerms)
			}
		})
	}
}

func Test_parserImp_parseOneWord(t *testing.T) {
	keyWordsMgr := key_words.NewKeyWordsMgr()

	type fields struct {
		keyWordMgr *key_words.KeyWordsMgr
	}
	type args struct {
		term []rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
		{
			name:   `test include`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune(`"Common.tars"`)},
			want:   []string{`"`, "Common", ".", "tars", `"`},
		},
		{
			name:   `test standard module`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("common::head")},
			want:   []string{"common", "::", "head"},
		},
		{
			name:   `test error module1`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("common:head")},
			want:   []string{"common", "::", "head"},
		},
		{
			name:   `test error module2`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("common::")},
			want:   []string{"common", "::"},
		},
		{
			name:   `test error module3`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("common:")},
			want:   []string{"common", "::"},
		},
		{
			name:   `test comment"`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("/* test */")},
			want:   []string{"/*", "test", "*/"},
		},
		{
			name:   `test vector`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("vector<int>")},
			want:   []string{"vector", "<", "int", ">"},
		},
		{
			name:   `test block end"`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("};")},
			want:   []string{"}", ";"},
		},
		{
			name:   `test enum end2`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("3,")},
			want:   []string{"3", ","},
		},
		{
			name:   `test element defination end`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("elememt;")},
			want:   []string{"elememt", ";"},
		},
		{
			name:   `test include `,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{[]rune("#include")},
			want:   []string{"#", "include"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parserImp{
				keyWordMgr: tt.fields.keyWordMgr,
			}
			if got := p.parseOneWord(tt.args.term); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOneWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserImp_parseBySpace(t *testing.T) {
	keyWordsMgr := key_words.NewKeyWordsMgr()

	type fields struct {
		keyWordMgr *key_words.KeyWordsMgr
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
		{
			name:   `test include `,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`#include "test.tars"`},
			want:   []string{"#include", `"test.tars"`},
		},
		{
			name:   `test struct`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`struct test{`},
			want:   []string{"struct", "test{"},
		},
		{
			name:   `test comment1`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{"// test"},
			want:   []string{"// test"},
		},
		{
			name:   `test comment2`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`0 require long id;              // id 更新必填,创建不需要`},
			want:   []string{"0", "require", "long", "id;", "// id 更新必填,创建不需要"},
		},
		{
			name:   `test comment3`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`0 require long id;              /* id 更新必填,创建不需要*/`},
			want:   []string{"0", "require", "long", "id;", "/* id 更新必填,创建不需要*/"},
		},
		{
			name:   `test comment4`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`/* test */`},
			want:   []string{"/* test */"},
		},
		{
			name:   `test end`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`};`},
			want:   []string{"};"},
		},
		{
			name:   `test enum`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`ADMIN = 1, //管理员`},
			want:   []string{"ADMIN", "=", "1,", "//管理员"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parserImp{
				keyWordMgr: tt.fields.keyWordMgr,
			}
			if got := p.parseBySpace(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				for idx, str := range got {
					fmt.Printf("idx:%v, str:%v\n", idx, str)
				}
				t.Errorf("parseBySpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserImp_ParseOneLine(t *testing.T) {
	keyWordsMgr := key_words.NewKeyWordsMgr()

	type fields struct {
		keyWordMgr *key_words.KeyWordsMgr
	}
	type args struct {
		line string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
		{
			name:   `test include`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`#include "Common.tars"`},
			want:   []string{"#", "include", `"`, "Common", ".", "tars", `"`},
		},
		{
			name:   `test struct`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`struct AlbumInfo{`},
			want:   []string{"struct", "AlbumInfo", "{"},
		},
		{
			name:   `test comment1`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{"// test"},
			want:   []string{"//", "test"},
		},
		{
			name:   `test comment2`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`0 require long id;              // id 更新必填,创建不需要`},
			want:   []string{"0", "require", "long", "id", ";", "//", "id 更新必填,创建不需要"},
		},
		{
			name:   `test comment3`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`0 require long id;              /* id 更新必填,创建不需要*/`},
			want:   []string{"0", "require", "long", "id", ";", "/*", "id 更新必填,创建不需要", "*/"},
		},
		{
			name:   `test comment4`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`/* test */`},
			want:   []string{"/*", "test", "*/"},
		},
		{
			name:   `test end`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`};`},
			want:   []string{"}", ";"},
		},
		{
			name:   `test enum`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`test = 1, //管理员`},
			want:   []string{"test", "=", "1", ",", "//", "管理员"},
		},
		{
			name:   `test module`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`0 require Common::RequestHead reqHead;`},
			want:   []string{"0", "require", "Common", "::", "RequestHead", "reqHead", ";"},
		},
		{
			name:   `test container`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`0 require vector<string> uids;`},
			want:   []string{"0", "require", "vector", "<", "string", ">", "uids", ";"},
		},
		{
			name:   `test comment`,
			fields: fields{keyWordMgr: keyWordsMgr},
			args:   args{`int Value (Key k);//test`},
			want:   []string{"int", "Value", "(", "Key", "k", ")", ";", "//", "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parserImp{
				keyWordMgr: tt.fields.keyWordMgr,
			}
			if got := p.ParseOneLine(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseOneLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
