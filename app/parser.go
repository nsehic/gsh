package main

import "strings"

type Parser struct {
	result       []string
	singleQuote  bool
	doubleQuote  bool
	escape       bool
	concatString bool
	stringBuffer strings.Builder
	input        string
}

func (p *Parser) Parse(input string) (string, []string) {
	p.input = input
	defer p.Reset()
	for i, c := range input {
		switch string(c) {
		case "'":
			if p.escape {
				p.stringBuffer.WriteRune(c)
				p.escape = false
				continue
			}
			if p.doubleQuote {
				p.stringBuffer.WriteRune(c)
				continue
			}
			if p.singleQuote {
				if p.concatString {
					p.concatString = false
				} else if p.getNextChar(i) != " " {
					p.concatString = true
				} else {
					p.singleQuote = false
					p.flushStringBuffer()
				}
			} else {
				p.singleQuote = true
			}
		case "\"":
			if p.escape {
				p.stringBuffer.WriteRune(c)
				p.escape = false
				continue
			}
			if p.singleQuote {
				p.stringBuffer.WriteRune(c)
				continue
			}
			if p.doubleQuote {
				if p.concatString {
					p.concatString = false
				} else if p.getNextChar(i) != " " {
					p.concatString = true
				} else {
					p.doubleQuote = false
					p.flushStringBuffer()
				}
			} else {
				p.doubleQuote = true
			}
		case " ":
			if p.escape {
				p.stringBuffer.WriteRune(c)
				p.escape = false
				continue
			}
			if p.singleQuote || p.doubleQuote {
				p.stringBuffer.WriteRune(c)
			} else {
				p.flushStringBuffer()
			}
		case "\\":
			if p.escape {
				p.stringBuffer.WriteRune(c)
				p.escape = false
			} else if !p.singleQuote && !p.doubleQuote {
				p.escape = true
			}
		default:
			p.stringBuffer.WriteRune(c)
		}
	}
	p.flushStringBuffer()
	return p.result[0], p.result[1:]
}

func (p *Parser) getNextChar(pos int) string {
	r := []rune(p.input)
	if pos >= len(r)-1 {
		return ""
	}
	return string(r[pos+1])
}

func (p *Parser) Reset() {
	p.input = ""
	p.result = []string{}
	p.stringBuffer.Reset()
	p.singleQuote = false
	p.doubleQuote = false
	p.concatString = false
	p.escape = false
}

func (p *Parser) flushStringBuffer() {
	if p.stringBuffer.Len() > 0 {
		p.result = append(p.result, p.stringBuffer.String())
		p.stringBuffer.Reset()
	}
}
