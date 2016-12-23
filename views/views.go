package views

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var LayoutDir string = "views/templates/"

type View struct {
	Template *template.Template
	Layout   string
}

type SessionData struct {
	Data interface{}
}

func layoutfiles() []string {
	files, err := filepath.Glob(LayoutDir + "/*.html")
	if err != nil {
		panic(err)
	}
	return files
}
func NewView(layout string, files ...string) *View {
	files = append(files, layoutfiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	packageData := &SessionData{
		Data: data,
	}
	fmt.Println(packageData)
	return v.Template.ExecuteTemplate(w, v.Layout, &packageData)
}
