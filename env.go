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
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/hzakher/configure/cstrings"
)

// ParseEnv takes a pointer to a struct and attempts
// to initialize it from environment variables.
func ParseEnv(v interface{}) error {
	env, err := parseEnvironment()
	if err != nil {
		return err
	}
	s := reflect.ValueOf(v).Elem()
	return setEnv(s, env, "")
}

// Setter is an interface that properties of struct can implement
// to initialize the value of any struct from string
type EnvSetter interface {
	SetEnv(string) error
}

// StringSetter is an interface that allows to set variable from any string
type StringSetter interface {
	Set(string) error
}

func setEnv(v reflect.Value, env map[string]string, prefix string) error {
	// for structs, walk every element and parse
	vType := v.Type()
	if v.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < v.NumField(); i++ {
		structField := vType.Field(i)
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		kind := field.Kind()
		// if kind == reflect.Struct {
		// 	if err := setEnv(field, env, prefix); err != nil {
		// 		return cstrings.Wrap(err,
		// 			fmt.Sprintf("failed parsing struct field %v",
		// 				structField.Name))
		// 	}
		// }
		envKey := structField.Tag.Get(Tag)
		envSkipKey := structField.Tag.Get("env")
		if envKey == "" || envSkipKey == "-" {
			continue
		}

		//if cli flag is set to - on the struct, then reset the prefix
		if envKey == "-" && kind == reflect.Struct {
			prefix = ""
		} else {
			if prefix == "" && kind == reflect.Struct {
				prefix = envKey
			} else if prefix != "" {
				envKey = strings.Join([]string{prefix, envKey}, ".")
			}
		}

		_, ok := env[envKey]
		envDefault := structField.Tag.Get("default")
		// 		sRequired := structField.Tag.Get("required")
		// 		isRequired, err := strconv.ParseBool(sRequired)
		//
		// 		if err != nil && isRequired {
		// 			if !ok && kind != reflect.Struct {
		// 				//this is required and missing, return error
		// 				err := fmt.Errorf("required field  %v is missing", envKey)
		// 				return err
		// 			}
		// 		}

		if !ok && kind != reflect.Struct {
			if envDefault != "" {
				env[envKey] = envDefault
			}
		}

		val, ok := env[envKey]
		if !ok || val == "" { // assume defaults
			continue
		}

		if field.CanAddr() {
			if s, ok := field.Addr().Interface().(EnvSetter); ok {
				if err := s.SetEnv(val); err != nil {
					return err
				}
				continue
			}
			if s, ok := field.Addr().Interface().(StringSetter); ok {
				if err := s.Set(val); err != nil {
					return err
				}
				continue
			}
		}
		switch kind {
		case reflect.Struct:
			if err := setEnv(field, env, prefix); err != nil {
				return cstrings.Wrap(err,
					fmt.Sprintf("failed parsing struct field %v",
						structField.Name))
			}
		case reflect.Slice:
			if _, ok := field.Interface().([]map[string]string); ok {
				var kv KeyValSlice
				if err := kv.SetEnv(val); err != nil {
					return cstrings.Wrap(err, "error parsing key value list")
				}
				field.Set(reflect.ValueOf([]map[string]string(kv)))
			}
		case reflect.Map:
			if _, ok := field.Interface().(map[string]string); ok {
				var kv KeyVal
				if err := kv.SetEnv(val); err != nil {
					return cstrings.Wrap(err, "error parsing key value list")
				}
				field.Set(reflect.ValueOf(map[string]string(kv)))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(val, 0, field.Type().Bits())
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		case reflect.String:
			field.SetString(val)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				return cstrings.Wrap(
					err,
					fmt.Sprintf("failed parsing struct field %v, expected bool, got '%v'",
						structField.Name, val))
			}
			field.SetBool(boolVal)
		}
	}
	return nil
}

func parseEnvironment() (map[string]string, error) {
	values := os.Environ()
	env := make(map[string]string, len(values))
	for _, v := range values {
		vals := strings.SplitN(v, "=", 2)
		if len(vals) != 2 {
			return nil, fmt.Errorf("failed to parse variable: '%v'", v)
		}
		env[vals[0]] = vals[1]
	}
	return env, nil
}
