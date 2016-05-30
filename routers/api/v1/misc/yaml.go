// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package misc

import (
	api "github.com/gogits/go-gogs-client"

	"github.com/gogits/gogs/modules/context"
	"github.com/gogits/gogs/modules/yaml"
)

// https://github.com/gogits/go-gogs-client/wiki/Miscellaneous#render-an-arbitrary-markdown-document
func Yaml(ctx *context.APIContext, form api.YamlOption) {
	if ctx.HasApiError() {
		ctx.Error(422, "", ctx.GetErrMsg())
		return
	}

	if len(form.Text) == 0 {
		ctx.Write([]byte(""))
		return
	}

	var yamlHtml, yamlBody []byte
	yamlHtml = yaml.RenderYamlHtmlTable([]byte(form.Text))
	switch form.Mode {
	case "gfm":
		yamlBody = yaml.Render(yaml.StripYamlFromText([]byte(form.Text)), form.Context, nil)
	default:
		yamlBody = yaml.RenderRaw([]byte(form.Text), "")
	}
	ctx.Write(append(yamlHtml, yamlBody...))
}

// https://github.com/gogits/go-gogs-client/wiki/Miscellaneous#render-a-markdown-document-in-raw-mode
func YamlRaw(ctx *context.APIContext) {
	body, err := ctx.Req.Body().Bytes()
	if err != nil {
		ctx.Error(422, "", err)
		return
	}
	ctx.Write(yaml.RenderRaw(body, ""))
}
