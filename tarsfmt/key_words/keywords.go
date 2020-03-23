package key_words

import (
	"github.com/trigger3/toy/mytype"
)

const (
	TYPE_NIL            int8 = iota
	TYPE_MODULE_INCLUDE      // module引用
	TYPE_MODULE              // module定义
	TYPE_STRUCT              // struct
	TYPE_ENUM                // 枚举
	TYPE_INTERFACE           // 接口
	TYPE_COMMENT             // 注释
	TYPE_END                 // 块结束
)

type KeyWordsMgr struct {
	total            *mytype.Set // 全部的keyword
	codeBlock        *mytype.Set // 包含代码块的keywords，module, struct, interface, enum
	rule             *mytype.Set // 规则可以words，require, optional
	container        *mytype.Set // 容器， // map, vector
	basicType        *mytype.Set // 字段基本类型，int, long, float, string
	sepFlag          *mytype.Set // 界符， {} / * ; <> () # : ,
	commentFlag      *mytype.Set
	doubleSepFlag    *mytype.Set
	tarsType         *mytype.Set // 包括module,结构体，interface,枚举
	moduleFlag       string      // #
	codeBlockEndFlag string      // };
	typeMap          map[string]int8
	statementTypeMap map[string]int8
}

func NewKeyWordsMgr() *KeyWordsMgr {
	sepFlag := mytype.NewSet('{', '}', '/', '*', ';', '<', '>', '(', ')', ':', ',', '.', '"', '#')
	commentFlag := mytype.NewSet('/')
	doubleSepFlag := mytype.NewSet('/', ':')
	basicType := mytype.NewSet("int", "long", "float", "string", "void", "bool", "byte", "short",
		"double", "unsigned byte", "unsigned int", "unsigned short")
	rule := mytype.NewSet("require", "optional")
	container := mytype.NewSet("map", "vector")
	tarsType := mytype.NewSet(TYPE_MODULE, TYPE_STRUCT, TYPE_ENUM, TYPE_INTERFACE)
	statementTypeMap := map[string]int8{
		"#":         TYPE_MODULE_INCLUDE,
		"enum":      TYPE_ENUM,
		"struct":    TYPE_STRUCT,
		"interface": TYPE_INTERFACE,
		"module":    TYPE_MODULE,
		"//":        TYPE_COMMENT,
		"/*":        TYPE_COMMENT,
		"}":         TYPE_END,
	}

	return &KeyWordsMgr{
		sepFlag:          sepFlag,
		commentFlag:      commentFlag,
		doubleSepFlag:    doubleSepFlag,
		basicType:        basicType,
		rule:             rule,
		container:        container,
		tarsType:         tarsType,
		statementTypeMap: statementTypeMap,
	}
}

func (k *KeyWordsMgr) IsKeyWord(firstTerm string) bool {
	return k.total.Contains(firstTerm)
}

func (k *KeyWordsMgr) IsElementRule(term string) bool {
	return k.rule.Contains(term)
}

func (k *KeyWordsMgr) IsBasicType(term string) bool {
	return k.basicType.Contains(term)
}

func (k *KeyWordsMgr) IsContainerType(term string) bool {
	return k.container.Contains(term)
}

func (k *KeyWordsMgr) IsSepWord(w int32) bool {
	return k.sepFlag.Contains(w)
}

func (k *KeyWordsMgr) IsCommentWord(w int32) bool {
	return k.commentFlag.Contains(w)
}

func (k *KeyWordsMgr) IsDoubleSepWord(w int32) bool {
	return k.doubleSepFlag.Contains(w)
}

func (k *KeyWordsMgr) IsTarsType(dataType int8) bool {
	return k.tarsType.Contains(dataType)
}

func (k *KeyWordsMgr) StatementType(dataType string) int8 {
	return k.statementTypeMap[dataType]
}
