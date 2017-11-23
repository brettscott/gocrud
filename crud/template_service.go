package crud

import (
	"fmt"
	"github.com/aymerick/raymond"
	"io/ioutil"
	"path/filepath"
)

const TEMPLATE_PATH string = "./templates/%s.hbs"

// TODO read files from filesystem instead of having to register them
var templateNames []string = []string{
	"root",
	"list",
	"form",
}

func newTemplateService() *templateService {
	tmpls := &templateService{}
	tmpls.parseTemplates()
	return tmpls
}

type templateService struct {
	list map[string]*raymond.Template
}

func (t *templateService) exec(tplName string, ctx map[string]interface{}) (html string, err error) {
	if tpl, ok := t.list[tplName]; ok {
		html, err = tpl.Exec(ctx)
		return
	}
	return "", fmt.Errorf("template \"%s\" not registered", tplName)
}

func (t *templateService) parseTemplates() {
	t.list = map[string]*raymond.Template{}

	for _, name := range templateNames {
		filename, _ := filepath.Abs(fmt.Sprintf(TEMPLATE_PATH, name))
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(fmt.Sprintf("template \"%s\" not found in filesystem: %s", name, err))
		}
		tpl, err := raymond.Parse(string(contents))
		if err != nil {
			panic(err)
		}
		t.list[name] = tpl
	}
}

func registerTemplateHelpers() {
	raymond.RegisterHelper("listColumnHeadings", ListColumnHeadings)
	raymond.RegisterHelper("listRows", ListRows)
	raymond.RegisterHelper("listCells", ListCells)
}
