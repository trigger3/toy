package line_parser

import (
	"strings"

	"github.com/trigger3/toy/tarsfmt/key_words"
)

type Parser interface {
	ParseOneLine(string) []string
}

func NewLineParse(keyWordsMgr *key_words.KeyWordsMgr) Parser {
	return &parserImp{keyWordsMgr}
}

type parserImp struct {
	keyWordMgr *key_words.KeyWordsMgr
}

func (p *parserImp) ParseOneLine(line string) []string {
	if line == "//" {
		return nil
	}
	terms := p.parseBySpace(line)
	//util.PrintArray(terms)
	return p.parseBySep(terms)
}

func (p *parserImp) parseBySpace(line string) []string {
	terms := make([]string, 0)
	runeSlice := []rune(line)

	var (
		begin                   int
		lastIsBlock, curIsBlock bool
	)

	for idx, word := range runeSlice {
		if p.keyWordMgr.IsCommentWord(word) {
			// 可能出现test;//场景，这个时候 test;需要单独分
			//if begin < idx && !lastIsBlock {
			if !lastIsBlock {
				terms = append(terms, string(runeSlice[begin:]))
			} else {
				terms = append(terms, string(runeSlice[idx:]))
			}
			//terms = append(terms, string(runeSlice[idx:]))
			return terms
		}

		curIsBlock = p.isBlock(word)
		// 当上个字符和当前字符都为空格 或者不为空格时，continue
		if (curIsBlock && lastIsBlock) || (!curIsBlock && !lastIsBlock) {
			continue
		}

		// 当上个字符是空格，当前不为空格时，设置begin为当前idx
		// 当所有的comment，包括标识符，作为一个单词
		if lastIsBlock {
			begin = idx
			lastIsBlock = false
			continue
		}

		// 当上个字符不为空格，当前为空格时，获取一个单词，档次为runeSlice[begin : idx]
		if curIsBlock {
			term := runeSlice[begin:idx]
			terms = append(terms, string(term))
			lastIsBlock = true
		}
	}

	if !lastIsBlock {
		term := runeSlice[begin:]
		terms = append(terms, string(term))
	}

	return terms
}

func (p *parserImp) parseBySep(words []string) []string {
	terms := make([]string, 0)
	for _, word := range words {
		if len(word) == 1 {
			terms = append(terms, string(word))
			continue
		}

		ts := p.parseOneWord([]rune(word))
		terms = append(terms, ts...)
	}

	return terms
}

// 需要注意的是 注释"// xxx"、"//xxx" 及 "/* */"、"/**/" 都被当成了一个word，但是不可能是"xxx;//"
func (p *parserImp) parseOneWord(term []rune) []string {
	if term[0] == '/' {
		return p.getComment(term)
	}

	terms := make([]string, 0)
	var (
		begin         int
		lastIsSepWord = true
		curIsSepWord  bool
	)
	for idx, word := range term {
		curIsSepWord = p.keyWordMgr.IsSepWord(word)
		// 上一个是，当前不是，代表一个变量正常单词开始了,
		// 上一个和当前都不为sep，直接continu
		if !curIsSepWord {
			if lastIsSepWord { // 当前不是sep，上一个是sep；新单词开始
				begin = idx
			}
			// 当期那不是sep，上一个不是sep；代表一个未完整的单词
			lastIsSepWord = false
			continue
		}

		// 先算前部分
		if !lastIsSepWord { // 当前是sep，上一个不是sep；代表一个单词结束
			// 上一个不是，当前是，代表一个变量正常单词结束了
			lastIsSepWord = true
			t := term[begin:idx]
			terms = append(terms, string(t))
		}

		// 再算后一部分
		// 判断是否为特殊的双分割符。如"::"、"//"、"/*"
		if p.keyWordMgr.IsDoubleSepWord(word) {
			newTerms := p.getDoubleSepResult(term[idx:])
			terms = append(terms, newTerms...)
			return terms
		}

		terms = append(terms, string(word))
		continue
	}

	if !lastIsSepWord {
		t := term[begin:]
		terms = append(terms, string(t))
	}

	return terms
}

func (p *parserImp) getDoubleSepResult(terms []rune) []string {
	if terms[0] == '/' {
		return p.getComment(terms)
	}
	if terms[0] == ':' {
		return p.getModule(terms)
	}

	return []string{}
}

// 对comment做了修正
func (p *parserImp) getComment(term []rune) (terms []string) {
	termLen := len(term)

	if termLen <= 2 {
		return
	}
	if term[1] == '/' {
		return []string{"//", p.TrimSpace(string(term[2:]))}
	}
	if term[1] == '*' {
		if termLen <= 4 {
			return
		}
		if term[termLen-1] == '/' {
			if term[termLen-2] == '*' {
				return []string{"/*", p.TrimSpace(string(term[2 : termLen-2])), "*/"}
			} else {
				return []string{"/*", p.TrimSpace(string(term[2 : termLen-1])), "*/"}
			}
		} else {
			if term[termLen-1] == '*' {
				return []string{"/*", p.TrimSpace(string(term[2 : termLen-1])), "*/"}
			} else {
				return []string{"/*", p.TrimSpace(string(term[2:termLen])), "*/"}
			}
		}
	} else { // 不为* 也不为/
		return []string{"//", p.TrimSpace(string(term[1:]))}
	}
}

// 修正common:head情况
func (p *parserImp) getModule(term []rune) (terms []string) {
	switch len(term) {
	case 0:
		return
	case 1:
		return []string{"::"}
	default:
		if term[1] != ':' {
			return []string{"::", string(term[1:])}
		}
		if len(term) > 2 {
			return []string{"::", string(term[2:])}
		} else {
			return []string{"::"}
		}
	}
}

func (p *parserImp) isTwo(word string) bool {
	return false
}

func (p *parserImp) isBlock(b int32) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func (p *parserImp) TrimSpace(word string) string {
	return strings.Trim(word, " \t\r\n")
}
