package crud

import (
	"fmt"
	"github.com/mergermarket/raymond"
	"io/ioutil"
	"os"
	"path"
	"runtime"
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
		contents, err := t.templateContents(name)
		if err != nil {
			panic(err)
		}

		tpl, err := raymond.Parse(contents)
		if err != nil {
			panic(err)
		}
		t.list[name] = tpl
	}
}

func (t *templateService) templateContents(name string) (string, error) {
	contents, err := readFileFromPathDefinedFromRuntimeCaller(name)
	if err != nil {
		contents, err = readFileFromPathDefinedFromWorkingDirectory(name)
		if err != nil {
			return "", fmt.Errorf("template \"%s\" not found in filesystem: %s", name, err)
		}
	}
	return string(contents), nil
}

func readFileFromPathDefinedFromRuntimeCaller(name string) ([]byte, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("failed to identify runtime caller to parse templates")
	}
	fp := path.Join(path.Dir(filename), fmt.Sprintf(TEMPLATE_PATH, name))
	contents, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to load from runtime caller - fp: \"%s\", err: \"%s\"\n", fp, err)
	}
	return contents, err
}

func readFileFromPathDefinedFromWorkingDirectory(name string) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fp := path.Join(dir, fmt.Sprintf(TEMPLATE_PATH, name))
	contents, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("failed to load from working directory - fp: \"%s\", err: \"%s\"\n", fp, err)
	}
	return contents, err
}

func registerTemplateHelpers() {
	_ = raymond.RegisterHelper("listColumnHeadings", ListColumnHeadings)
	_ = raymond.RegisterHelper("listRows", ListRows)
	_ = raymond.RegisterHelper("listCells", ListCells)
	_ = raymond.RegisterHelper("formElements", FormElements)
}
