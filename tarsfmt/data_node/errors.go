package data_node

import "errors"

var (
	ErrModuleNameCannotBeNumeric       = errors.New("module name is numeric")
	ErrStructNameCannotBeNumeric       = errors.New("struct name is numeric")
	ErrSyntexInvaild                   = errors.New("syntex invaild")
	ErrElemetRuleInvalid               = errors.New("element rule invalid, [require|optional]")
	ErrInterfaceRespNotDefine          = errors.New("interface response struct not defined")
	ErrInterfaceReqNotDefine           = errors.New("interface request struct not defined")
	ErrInterfaceReqNameCannotBeNumeric = errors.New("interface request name is numeric")

	ErrInterfaceFuncNameCannotBeNumeric = errors.New("interface function name is numeric")
)

var ErrCodeBlockNotBegin error = errors.New("last code block is not begin")
var ErrCodeBlockNotEnd error = errors.New("last code block is not end")
var ErrVariableNameCannotBeNumeric = errors.New("variable name is numeric")
var ErrEnumNameCannotBeNumeric = errors.New("enum name is numeric")
var ErrEnumValueMustBeNumeric = errors.New("enum value is not numeric")
