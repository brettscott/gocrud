package crud

func NewFakeTemplateServicer() *fakeTemplateServicer {
	return &fakeTemplateServicer{}
}

type fakeTemplateServicer struct {
	execHtml  string
	execError error
}

func (f *fakeTemplateServicer) exec(tmplName string, ctx map[string]interface{}) (html string, err error) {
	return f.execHtml, f.execError
}
