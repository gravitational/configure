package cxml

import (
	"encoding/xml"
	"io"
	"strings"
)

// TransofrmFunc is a function that can transform incoming token
// into a series of outgoing tokens when traversing XML tree
type TransformFunc func(in xml.Token) []xml.Token

// TransformXML parses the XML tree, traverses it and calls TransformFunc
// on each XML token, writing the output to the writer, resulting in a
// transformed XML tree
func TransformXML(r io.Reader, w io.Writer, fn TransformFunc, indent bool) error {
	parser := xml.NewDecoder(r)
	encoder := xml.NewEncoder(w)
	if indent {
		encoder.Indent("  ", "    ")
	}
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		for _, t := range fn(token) {
			if err := encoder.EncodeToken(t); err != nil {
				return err
			}
		}
	}
	encoder.Flush()
	return nil
}

// SetAttribute is XML helper that allows to set attribute on a node
func SetAttribute(e xml.StartElement, name, value string) xml.StartElement {
	if len(e.Attr) != 0 {
		for i := range e.Attr {
			if e.Attr[i].Name.Local == name {
				e.Attr[i].Value = value
				return e
			}
		}
	} else {
		e.Attr = append(e.Attr, xml.Attr{Name: xml.Name{Local: name}, Value: value})
	}
	return e
}

// ReplaceCDATAIf replaces CDATA value of the matched node
// if the parent node name matches the name
func ReplaceCDATA(el xml.Name, value string) TransformFunc {
	var prev xml.Token
	return func(in xml.Token) []xml.Token {
		switch t := in.(type) {
		case xml.StartElement:
			prev = t
		case xml.CharData:
			oldPrev := prev
			prev = nil
			if oldPrev != nil {
				if prevStart, ok := oldPrev.(xml.StartElement); ok {
					if prevStart.Name.Local == el.Local && prevStart.Name.Space == el.Space {
						return []xml.Token{xml.CharData(value)}
					}
				}
			}
		default:
			prev = nil
		}
		return []xml.Token{in}
	}
}

// InjectNodes injects nodes at the end of the tag that matches name
func InjectNodes(name xml.Name, nodes []xml.Token) TransformFunc {
	return func(in xml.Token) []xml.Token {
		switch t := in.(type) {
		case xml.EndElement:
			e := xml.EndElement(t)
			if e.Name.Local == name.Local && e.Name.Space == name.Space {
				return append(nodes, in)
			}
		}
		return []xml.Token{in}
	}
}

// ReplaceAttribute replaces the attribute of the first node that matches the name
func ReplaceAttributeIf(attrName, attrValue string, match func(e xml.Token) bool) TransformFunc {
	return func(in xml.Token) []xml.Token {
		switch t := in.(type) {
		case xml.StartElement:
			if match(in) {
				e := xml.StartElement(t)
				return []xml.Token{SetAttribute(e, attrName, attrValue)}
			}
		}
		return []xml.Token{in}
	}
}

// TrimSpace is a transformer function that replaces CDATA with blank
// characters with empty strings
func TrimSpace(in xml.Token) []xml.Token {
	switch t := in.(type) {
	case xml.CharData:
		if strings.TrimSpace(string(t)) == "" {
			return []xml.Token{xml.CharData("")}
		}
	}
	return []xml.Token{in}
}

// Combine takes a list of TransformFuncs and converts them
// into a single transform function applying all functions one by one
func Combine(funcs ...TransformFunc) TransformFunc {
	return func(in xml.Token) []xml.Token {
		out := []xml.Token{in}
		for _, f := range funcs {
			new := []xml.Token{}
			for _, t := range out {
				new = append(new, f(t)...)
			}
			out = new
		}
		return out
	}
}
