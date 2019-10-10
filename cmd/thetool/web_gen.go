package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const templateText = `

## {{ lang .Lang .Title }}
##### {{ if gt .Year 0 }}{{ .Year }}/{{ end }}{{ .Month }}/{{ .Day }}

{{ lang .Lang .Description }}

`

var tpl = template.Must(template.New("entry").Funcs(
	template.FuncMap{
		"lang": func(lang string, key map[string]string) string {
			return key[lang]
		},
	}).Parse(templateText))

func generateMarkdowns(fl *File, folder string) error {
	bufs := make(map[string]*bytes.Buffer)
	getbuf := func(ln string) *bytes.Buffer {
		b := bufs[ln]
		if b == nil {
			b = &bytes.Buffer{}
			bufs[ln] = b
		}

		return b
	}
	for i := range fl.Events {
		key := fl.Months.Name[fl.Events[i].Month-1]["en_US"] + "-%s.md"
		en := getbuf(fmt.Sprintf(key, "en_US"))
		fa := getbuf(fmt.Sprintf(key, "fa_IR"))

		if err := tpl.Execute(en, struct {
			Lang string
			Event
		}{
			Lang:  "en_US",
			Event: fl.Events[i],
		}); err != nil {
			return err
		}

		if err := tpl.Execute(fa, struct {
			Lang string
			Event
		}{
			Lang:  "fa_IR",
			Event: fl.Events[i],
		}); err != nil {
			return err
		}
	}

	path := filepath.Join(folder, strings.ToLower(fl.Name))
	if err := os.MkdirAll(path, 0750); err != nil {
		return err
	}

	for i := range bufs {
		file := filepath.Join(path, i)
		if err := ioutil.WriteFile(file, bufs[i].Bytes(), 0600); err != nil {
			return err
		}
	}

	return nil
}
