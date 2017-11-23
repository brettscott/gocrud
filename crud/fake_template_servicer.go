package crud

func NewFakeTemplateServicer() *fakeTemplateServicer {
	return &fakeTemplateServicer{}
}

type fakeTemplateServicer struct {
	execHtml  string
	execError error
	execTmplName string
	execContext map[string]interface{}
}

func (f *fakeTemplateServicer) exec(tmplName string, ctx map[string]interface{}) (html string, err error) {
	f.execTmplName = tmplName
	f.execContext = ctx
	return f.execHtml, f.execError
}
