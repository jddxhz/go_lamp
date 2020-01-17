package myweb

import (
	"net/http"
	"text/template"
)

func MyWeb(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./view/index.html")

	data := map[string]string{
		"name":    "zeta",
		"someStr": "这是一个开始",
	}

	t.Execute(w, data)

	// fmt.Fprintf(w, "这是一个开始")
}
