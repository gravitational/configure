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
package configure

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/hzakher/configure/cstrings"
	//"gopkg.in/yaml.v2"
	"github.com/hzakher/yaml/v2"
	// "github.com/narisongdev/yaml"
)

// ParseYAML parses yaml-encoded byte string into the struct
// passed to the function.
// EnableTemplating() argument allows to treat configuration file as a template
// for example, it will support {{env "VAR"}} - that will substitute
// environment variable "VAR" and pass it to YAML file parser
func ParseYAML(data []byte, cfg interface{}, funcArgs ...ParseOption) error {
	var opts parseOptions
	// yaml.Tag = "config"
	for _, fn := range funcArgs {
		fn(&opts)
	}
	var err error
	if opts.templating {
		if data, err = renderTemplate(data); err != nil {
			return err
		}
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}

type parseOptions struct {
	// templating turns on templating mode when
	// parsing yaml file
	templating bool
}

// ParseOption is a functional argument type
type ParseOption func(p *parseOptions)

// EnableTemplating allows to treat configuration file as a template
// for example, it will support {{env "VAR"}} - that will substitute
// environment variable "VAR" and pass it to YAML file parser
func EnableTemplating() ParseOption {
	return func(p *parseOptions) {
		p.templating = true
	}
}

func renderTemplate(data []byte) ([]byte, error) {
	t := template.New("tpl")
	c, err := newCtx()
	if err != nil {
		return nil, err
	}
	t.Funcs(map[string]interface{}{
		"env":  c.Env,
		"file": c.File,
	})
	t, err = t.Parse(string(data))
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func newCtx() (*ctx, error) {
	values := os.Environ()
	c := &ctx{
		env: make(map[string]string, len(values)),
	}
	for _, v := range values {
		vals := strings.SplitN(v, "=", 2)
		if len(vals) != 2 {
			return nil, fmt.Errorf("failed to parse variable: '%v'", v)
		}
		c.env[vals[0]] = vals[1]
	}
	return c, nil
}

type ctx struct {
	env map[string]string
}

func (c *ctx) File(path string) (string, error) {
	o, err := ioutil.ReadFile(path)
	if err != nil {
		return "", cstrings.Wrap(err, fmt.Sprintf("reading file: %v", path))
	}
	return string(o), nil
}

func (c *ctx) Env(key string) (string, error) {
	v, ok := c.env[key]
	if !ok {
		return "", fmt.Errorf("environment variable '%v' is not set", key)
	}
	values := cstrings.SplitComma(v)
	out := make([]string, len(values))
	for i, p := range values {
		out[i] = quoteYAML(p)
	}
	return strings.Join(out, ","), nil
}

func quoteYAML(val string) string {
	if len(val) == 0 {
		return val
	}
	if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
		return val
	}
	if strings.ContainsAny(val, ":") {
		return "'" + val + "'"
	}
	return val
}
