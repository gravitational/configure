/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cxml

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func TestSchema(t *testing.T) { TestingT(t) }

type SchemaSuite struct {
}

var _ = Suite(&SchemaSuite{})

func (s *SchemaSuite) TestCases(c *C) {
	type tc struct {
		in  string
		out string
		fn  TransformFunc
	}

	testCases := []tc{
		{
			in:  `<xml><source file="before"></source></xml>`,
			out: `<xml><source file="after"></source></xml>`,
			fn: ReplaceAttributeIf("file", "after", func(el xml.Token) bool {
				e, ok := el.(xml.StartElement)
				if !ok {
					return false
				}
				return e.Name.Local == "source"
			}),
		},
		{
			in:  `<xml><node></node></xml>`,
			out: `<xml><node></node><disk type="file"></disk></xml>`,
			fn: InjectNodes(xml.Name{Local: "xml"}, []xml.Token{
				xml.StartElement{
					Name: xml.Name{Local: "disk"},
					Attr: []xml.Attr{{
						Name:  xml.Name{Local: "type"},
						Value: "file"}}},
				xml.EndElement{Name: xml.Name{Local: "disk"}},
			}),
		},
		{
			in:  `<xml><node>prev value</node><other>other</other></xml>`,
			out: `<xml><node>new value</node><other>other</other></xml>`,
			fn:  ReplaceCDATA(xml.Name{Local: "node"}, "new value"),
		},
	}

	for i, tc := range testCases {
		comment := Commentf("test #%d", i+1)

		out := &bytes.Buffer{}
		err := TransformXML(strings.NewReader(tc.in), out, tc.fn, false)
		c.Assert(err, IsNil, comment)
		c.Assert(out.String(), Equals, tc.out)
	}
}
