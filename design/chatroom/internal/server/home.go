package server

import (
	"fmt"
	"net/http"
	"text/template"
)

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles(homeTemplate)
	if err != nil {
		fmt.Fprint(w, "template parse failed!")
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "template execute failed!")
		return
	}
}
