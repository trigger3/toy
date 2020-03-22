package data_node

import "errors"

var (
	ErrModuleNameCannotBeNumeric        = errors.New("module name is numeric")
	ErrStructNameCannotBeNumeric        = errors.New("struct name is numeric")
	ErrSyntexInvaild                    = errors.New("syntex invaild")
	ErrElemetRuleInvalid                = errors.New("element rule invalid, [require|optional]")
	ErrInterfaceRespNotDefine           = errors.New("interface response struct not defined")
	ErrInterfaceReqNotDefine            = errors.New("interface request struct not defined")
	ErrInterfaceReqNameCannotBeNumeric  = errors.New("interface request name is numeric")
	ErrInterfaceFuncNameCannotBeNumeric = errors.New("interface function name is numeric")
	ErrCodeBlockNotEnd                  = errors.New("last code block is not end")
	ErrCodeBlockNotBegin                = errors.New("last code block is not begin")
	ErrVariableNameCannotBeNumeric      = errors.New("variable name is numeric")
	ErrEnumNameCannotBeNumeric          = errors.New("enum name is numeric")
	ErrEnumValueMustBeNumeric           = errors.New("enum value is not numeric")
)
