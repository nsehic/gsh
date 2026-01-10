package main

import (
	"strings"
)

type TokenKind uint8

const (
	TokenKindLiteral TokenKind = iota
	TokenKindSymbol
)

type Token struct {
	Kind  TokenKind
	Value string
}

type Parser struct {
	tokens          []Token
	buffer          strings.Builder
	input           string
	singleQuoteMode bool
	doubleQuoteMode bool
	escapeMode      bool
	concatMode      bool
	symbolMode      bool
}

type ParseResult struct {
	Command            string
	Args               []string
	StdoutRedirectPath string
	StderrRedirectPath string
}

func (p *Parser) Parse(input string) *ParseResult {
	p.input = input
	defer p.Reset()

	for pos, char := range input {
		switch char {
		case '\'':
			if p.escapeMode || p.doubleQuoteMode {
				p.writeChar(char)
			} else if p.singleQuoteMode {
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
			if p.escapeMode || p.singleQuoteMode {
				p.writeChar(char)
			} else if p.doubleQuoteMode {
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
			if p.escapeMode || p.singleQuoteMode || p.doubleQuoteMode {
				p.writeChar(char)
			} else {
				p.flushBuffer()
			}
		case '>':
			prevChar := p.getPrevChar(pos)
			if p.escapeMode || p.singleQuoteMode || p.doubleQuoteMode {
				p.writeChar(char)
			} else if prevChar == ' ' || prevChar == '1' || prevChar == '2' {
				p.writeSymbol(char)
			} else {
				p.writeChar(char)
			}
		case '\\':
			nextChar := p.getNextChar(pos)

			if p.escapeMode || p.singleQuoteMode {
				p.writeChar(char)
			} else if !p.doubleQuoteMode || (p.doubleQuoteMode && (nextChar == '\\' || nextChar == '"')) {
				p.escapeMode = true
			} else {
				p.writeChar(char)
			}
		default:
			p.writeChar(char)
		}
	}
	p.flushBuffer()
	if len(p.tokens) > 0 {
		commandToken, tokens := p.tokens[0], p.tokens[1:]
		args := []string{}
		symbolPos := -1
		for pos, token := range tokens {
			if token.Kind == TokenKindSymbol {
				symbolPos = pos
				break
			}
			args = append(args, token.Value)
		}
		res := ParseResult{
			Command: commandToken.Value,
			Args:    args,
		}

		if symbolPos != -1 && symbolPos < len(tokens)-1 {
			token := tokens[symbolPos]
			nextToken := tokens[symbolPos+1]
			switch token.Value {
			case ">":
			case "1>":
				res.StdoutRedirectPath = nextToken.Value
			case "2>":
				res.StderrRedirectPath = nextToken.Value
			}
		}
		return &res
	}
	return nil
}

func (p *Parser) writeChar(char rune) {
	p.buffer.WriteRune(char)
	p.escapeMode = false
	p.symbolMode = false
}

func (p *Parser) writeSymbol(char rune) {
	p.buffer.WriteRune(char)
	p.escapeMode = false
	p.symbolMode = true
}

func (p *Parser) getNextChar(pos int) rune {
	r := []rune(p.input)
	if pos >= len(r)-1 {
		return 0
	}
	return r[pos+1]
}

func (p *Parser) getPrevChar(pos int) rune {
	r := []rune(p.input)
	if pos <= 0 {
		return 0
	}
	return r[pos-1]
}

func (p *Parser) Reset() {
	p.input = ""
	p.tokens = []Token{}
	p.buffer.Reset()
	p.singleQuoteMode = false
	p.doubleQuoteMode = false
	p.concatMode = false
	p.escapeMode = false
	p.symbolMode = false
}

func (p *Parser) flushBuffer() {
	if p.buffer.Len() > 0 {
		kind := TokenKindLiteral
		if p.symbolMode {
			kind = TokenKindSymbol
			p.symbolMode = false
		}
		p.tokens = append(p.tokens, Token{Kind: kind, Value: p.buffer.String()})
		p.buffer.Reset()
	}
}
