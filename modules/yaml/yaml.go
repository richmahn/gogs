// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package yaml

import (
	"fmt"
	"strings"
	"reflect"

	"gopkg.in/yaml.v2"
	"github.com/microcosm-cc/bluemonday"
)

// Note: this section is for purpose of increase performance and
// reduce memory allocation at runtime since they are constant literals.
var (
	space             = " "
	spaceEncoded      = "%20"
)

const (
	DIR_HORIZONTAL    = "horizontal"
	DIR_VERTICAL      = "vertial"
)

var Sanitizer = bluemonday.UGCPolicy()

func renderHorizontalHtmlTable(m yaml.MapSlice) string {
	var thead, tbody, table string
	var mi yaml.MapItem
	for _, mi = range m {
		key := mi.Key
		value := mi.Value

		if  key != nil && reflect.TypeOf(key).String() == "yaml.MapSlice" {
			key = renderHorizontalHtmlTable(key.(yaml.MapSlice))
		}
		thead += fmt.Sprintf("<th>%v</th>", key)

		if value != nil && reflect.TypeOf(value).String() == "yaml.MapSlice" {
			value = renderHorizontalHtmlTable(value.(yaml.MapSlice))
		}
		tbody += fmt.Sprintf("<td>%v</td>", value)
	}

	table = ""
	if len(thead) > 0 {
		table = fmt.Sprintf(`<table data="yaml-metadata"><thead><tr>%s</tr></thead><tbody><tr>%s</tr></table>`, thead, tbody)
	}
	return table
}

func renderVerticalHtmlTable(m yaml.MapSlice) string {
	var mi yaml.MapItem
	var table string

	table = `<table data="yaml-metadata">`
	for _, mi = range m {
		key := mi.Key
		value := mi.Value

		table += `<tr>`
		if  key != nil && reflect.TypeOf(key).String() == "yaml.MapSlice" {
			key = renderHorizontalHtmlTable(key.(yaml.MapSlice))
		}
		table += fmt.Sprintf("<td>%v</td>", key)

		if value != nil && reflect.TypeOf(value).String() == "yaml.MapSlice" {
			value = renderVerticalHtmlTable(value.(yaml.MapSlice))
		}
		table += fmt.Sprintf("<td>%v</td>", value)

		table += `</tr>`
	}
	table += `</table>`

	return table
}

func RenderYamlHtmlTable(data []byte, dir string) []byte {
	m := yaml.MapSlice{}

	if len(data) < 1 {
		return []byte("")
	}

	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 1 {
		lines = strings.Split(string(data), "\n")
	}
	if len(lines) < 1 || lines[0] != "---" {
		return []byte("")
	}

	if err := yaml.Unmarshal(data, &m); err != nil {
		return []byte("")
	}

	if dir == DIR_HORIZONTAL {
		return []byte(renderHorizontalHtmlTable(m))
	} else if dir == DIR_VERTICAL {
		return []byte(renderVerticalHtmlTable(m))
	} else {
		return data
	}
}

func StripYamlFromText(data []byte) []byte {
	m := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(data, &m); err != nil {
		return data
	}

	lines := strings.Split(string(data), "\r\n")
	if len(lines) == 1 {
		lines = strings.Split(string(data), "\n")
	}
	if len(lines) < 1 || lines[0] != "---" {
		return data
	}
	body := ""
	atBody := false
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if line == "---" {
			atBody = true
		} else if atBody {
			body += line+"\n"
		}
	}
	return []byte(body)
}

func Render(rawBytes []byte) []byte {
	result := RenderYamlHtmlTable(rawBytes, DIR_VERTICAL)
	result = Sanitizer.SanitizeBytes(result)
	return result
}

// Renders the YAML and text as a string
func RenderString(rawBytes []byte) string {
	return string(Render(rawBytes))
}
