package controllers

import (
	"fmt"
	"net/http"
	"text/template"
)

var getuploadForm = func(w http.ResponseWriter, r *http.Request) {

	title := "Uploading file"
	t, _ := template.ParseFiles("./templates/upload.html")

	t.Execute(w, title)

}

var getStatistic = func(w http.ResponseWriter, r *http.Request) {
	fmt.Print("second")
}
