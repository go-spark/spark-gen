package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Parser struct {
	r     io.Reader
	file  string
	props []string
}

type Element struct {
	Name        string               `json:"name"`
	Content     string               `json:"content,omitempty"`
	SelfClosing bool                 `json:"selfClosing"`
	Attributes  map[string]Attribute `json:"attributes,omitempty"`
	Children    []*Element           `json:"children,omitempty"`
}

type Component struct {
	File     string     `json:"file"`
	Elements []*Element `json:"elements"`
	Props    []string   `json:"props"`
}

type Attribute struct {
	Value          string
	Go             bool
	JSExec         bool
	IsDefaultValue bool
}

func NewParser(r io.Reader, file string) *Parser {
	data, _ := io.ReadAll(r)
	s := strings.Split(string(data), "\n")

	parser := &Parser{file: file}

	if len(s) > 0 && strings.HasPrefix(s[0], "!!") {
		props := strings.TrimPrefix(s[0], "!!")
		sProps := strings.Split(props, ",")
		for _, p := range sProps {
			parser.props = append(parser.props, "\""+strings.TrimSpace(p)+"\"")
		}

		s = s[1:]
	}

	parser.r = strings.NewReader(strings.Join(s, "\n"))
	return parser
}

func (p *Parser) Parse() (*Component, error) {
	tokenizer := html.NewTokenizer(p.r)
	var roots []*Element
	var stack []*Element

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return &Component{
					File:     p.file,
					Elements: roots,
					Props:    p.props,
				}, nil
			}
			return nil, err
		case html.StartTagToken:
			token := tokenizer.Token()
			element := &Element{
				Name:        token.Data,
				SelfClosing: false,
				Attributes:  p.parseAttributes(token.Attr),
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, element)
			} else {
				roots = append(roots, element)
			}
			stack = append(stack, element)
		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			element := &Element{
				Name:        token.Data,
				SelfClosing: true,
				Attributes:  p.parseAttributes(token.Attr),
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, element)
			} else {
				roots = append(roots, element)
			}
		case html.EndTagToken:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case html.TextToken:
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				current.Content += string(tokenizer.Text())
			}
		}
	}
}

func (p *Parser) parseAttributes(attr []html.Attribute) map[string]Attribute {
	attributes := make(map[string]Attribute)

	for _, a := range attr {
		att := Attribute{Value: a.Val}
		key := a.Key

		if strings.HasPrefix(key, ":") {
			att.Go = true
			key = key[1:]
		}

		if strings.HasPrefix(key, "@") {
			att.JSExec = true
			key = key[1:]
		}

		if strings.HasPrefix(key, "!") {
			att.IsDefaultValue = true
			key = key[1:]
		}

		key = camelCase(key)

		attributes[key] = att
	}

	return attributes
}
