package main

import (
	"html/template"
	"net/http"
	"fmt"
	"bytes"
)

// layoutFuncs will be the function to execute the html defined content to layout.html
// This will be the default template.FuncMap if there's no passed template.FuncMap instance
// to layout.html
var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called inappropriately")
	},
}

// layout is the main html layout where other templates will be executed
var layout = template.Must(
	template.New("layout.html").
		Funcs(layoutFuncs).ParseFiles("templates/layout.html"))

// templates is the main instance of all templates in gophr
var templates = template.Must(template.New("t").ParseGlob("templates/**/*.html"))

// errorTemplate is an html that will display to the client browser when something error
// occur.
var errorTemplate = `
<html>
	<body>
		<h1>Error rendering template %s\n</h1>
		<p>%s</p>
	</body>
</html>
`

// RenderTemplate is an helper function for displaying the html template.
// It will use a template.FuncMap to load the content of the defined HTML template.
func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["CurrentUser"] = RequestUser(r)
	data["Flash"] = r.URL.Query().Get("flash")
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			templates.ExecuteTemplate(buf, name, data)
			return template.HTML(buf.String()), nil
		},
	}

	cloneLayout, _ := layout.Clone()
	cloneLayout.Funcs(funcs)

	err := cloneLayout.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf(errorTemplate, name, err), http.StatusInternalServerError)
	}
}
