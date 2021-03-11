package formator

import (
	"bytes"
	"fmt"
	"io"

	"github.com/trigger3/toy/thriftfmt/data_node"
	"github.com/trigger3/toy/thriftfmt/key_words"
	"github.com/trigger3/toy/thriftfmt/util"
)

type Formator interface {
	Format(line []string) error
	Print(buff *bytes.Buffer)
}

func NewFormator() Formator {
	return &formatorImp{
		writer:      nil,
		nodes:       []data_node.Node{},
		keyWordsMgr: key_words.NewKeyWordsMgr(),
		existNodes:  data_node.NewExistNodes(),
		level:       0,
	}
}

type formatorImp struct {
	writer io.Writer
	nodes  []data_node.Node

	lastNode    data_node.Node
	keyWordsMgr *key_words.KeyWordsMgr
	existNodes  *data_node.ExistNodes

	level int
}

func (a *formatorImp) Format(terms []string) (err error) {
	util.PrintArray(terms)
	dataType := a.getCodeType(terms)
	if err := a.adjustLevel(dataType); err != nil {
		return err
	}

	err = a.format(terms, dataType)
	if err != nil {
		//panic(err.Error())
		return err

	}
	return err
}

func (a *formatorImp) Print(buff *bytes.Buffer) {
	for _, node := range a.nodes {
		node.Print(buff)
		buff.WriteByte('\n')
	}
}

func (a *formatorImp) format(terms []string, dataType int8) error {
	// level不为零情况情况情况下，类型一定还为上面一个类型
	if a.level != 0 {
		if a.lastNode == nil {
			return data_node.ErrSyntexInvaild
		}
		return a.lastNode.Parse(terms, false)
	}
	// level为0时，有两种情况
	// 1. lastNode不为空时，表示为lastNode结束
	// 2. lastNode为空时，表示为新的一个Node，如 commnt, include
	if a.lastNode != nil {
		if err := a.lastNode.Parse(terms, true); err != nil {
			return err
		}
		a.nodes = append(a.nodes, a.lastNode)
		a.lastNode = nil
		return nil
	}

	if dataType == key_words.TYPE_NIL {
		return data_node.ErrSyntexInvaild
	}
	node, err := a.getNodeParser(dataType)
	if err != nil {
		return err
	}
	if err := node.Parse(terms, false); err != nil {
		return err
	}

	// 新创建一个node时，有两种情况，
	// 1. dataType为module，此时将module记为 lastNode，因为其要持续解析
	// 2. dataType为include或者comment，此时此node立即结束，热庵后将其保定到node列表里
	if dataType == key_words.TYPE_MODULE {
		a.level++
		a.lastNode = node
	} else {
		a.nodes = append(a.nodes, node)
	}

	return nil
}

func (a *formatorImp) getCodeType(terms []string) int8 {
	return a.keyWordsMgr.StatementType(terms[0])

}

func (a *formatorImp) adjustLevel(dataType int8) error {
	if a.lastNode != nil {
		if a.keyWordsMgr.IsTarsType(dataType) {
			a.level++
		} else if dataType == key_words.TYPE_END {
			a.level--
		}
	}
	if a.level > 2 {
		return data_node.ErrSyntexInvaild
	}

	return nil
}

func (a *formatorImp) getNodeParser(dataType int8) (data_node.Node, error) {
	switch dataType {
	case key_words.TYPE_MODULE_INCLUDE:
		return data_node.NewModuleInclude(a.existNodes), nil
	case key_words.TYPE_MODULE:
		return data_node.NewModule(a.existNodes, a.keyWordsMgr), nil
	case key_words.TYPE_COMMENT:
		return data_node.NewComment(0), nil
	default:
		// TODO
		return nil, fmt.Errorf("get node parse fail, err:%w, dataType:%v", data_node.ErrSyntexInvaild, dataType)
	}
}
