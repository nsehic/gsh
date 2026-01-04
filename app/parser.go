package main

import "strings"

type Parser struct {
	result          []string
	buffer          strings.Builder
	input           string
	singleQuoteMode bool
	doubleQuoteMode bool
	escapeMode      bool
	concatMode      bool
}

func (p *Parser) Parse(input string) (string, []string) {
	p.input = input
	defer p.Reset()

	for pos, char := range input {
		switch char {
		case '\'':
			if p.escapeMode {
				p.buffer.WriteRune(char)
				p.escapeMode = false
				continue
			}
			if p.doubleQuoteMode {
				p.buffer.WriteRune(char)
				continue
			}
			if p.singleQuoteMode {
				if p.concatMode {
					p.concatMode = false
				} else if p.getNextChar(pos) != ' ' {
					p.concatMode = true
				} else {
					p.singleQuoteMode = false
					p.flushBuffer()
				}
			} else {
				p.singleQuoteMode = true
			}
		case '"':
			if p.escapeMode {
				p.buffer.WriteRune(char)
				p.escapeMode = false
				continue
			}
			if p.singleQuoteMode {
				p.buffer.WriteRune(char)
				continue
			}
			if p.doubleQuoteMode {
				if p.concatMode {
					p.concatMode = false
				} else if p.getNextChar(pos) != ' ' {
					p.concatMode = true
				} else {
					p.doubleQuoteMode = false
					p.flushBuffer()
				}
			} else {
				p.doubleQuoteMode = true
			}
		case ' ':
			if p.escapeMode {
				p.buffer.WriteRune(char)
				p.escapeMode = false
				continue
			}
			if p.singleQuoteMode || p.doubleQuoteMode {
				p.buffer.WriteRune(char)
			} else {
				p.flushBuffer()
			}
		case '\\':
			nextChar := p.getNextChar(pos)
			if p.escapeMode {
				p.buffer.WriteRune(char)
				p.escapeMode = false
			} else if p.doubleQuoteMode && (nextChar == '\\' || nextChar == '"') {
				p.escapeMode = true
			} else {
				p.buffer.WriteRune(char)
			}
		default:
			if p.escapeMode {
				p.escapeMode = false
			}
			p.buffer.WriteRune(char)
		}
	}
	p.flushBuffer()
	if len(p.result) > 0 {
		return p.result[0], p.result[1:]
	}
	return "", []string{}
}

func (p *Parser) getNextChar(pos int) rune {
	r := []rune(p.input)
	if pos >= len(r)-1 {
		return 0
	}
	return r[pos+1]
}

func (p *Parser) Reset() {
	p.input = ""
	p.result = []string{}
	p.buffer.Reset()
	p.singleQuoteMode = false
	p.doubleQuoteMode = false
	p.concatMode = false
	p.escapeMode = false
}

func (p *Parser) flushBuffer() {
	if p.buffer.Len() > 0 {
		p.result = append(p.result, p.buffer.String())
		p.buffer.Reset()
	}
}
