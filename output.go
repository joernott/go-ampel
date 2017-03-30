// Output HTML or json
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"text/template"
)

// Renders first the body and and passes it to outputPage
func renderPage(w http.ResponseWriter, serverName string, svc *Service, tpl string, title string) {
	var body bytes.Buffer
	template, err := template.New("svc").Parse(tpl)
	if err != nil {
		panic(err)
	}
	svc.FUntil = svc.Until.Format("2006-01-02 15:04:05")
	err = template.Execute(&body, svc)
	if err != nil {
		panic(err)
	}
	outputPage(w, "Service "+serverName+title, body.String())
}

func renderList(w http.ResponseWriter, title string, templateName string, templateContent string, data interface{}) {
	var body bytes.Buffer

	tmpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(&body, data)
	if err != nil {
		panic(err)
	}
	outputPage(w, title, body.String())
}

//Writes a HTML page to the provided ResponseWriter
func outputPage(w http.ResponseWriter, title string, body string) {
	page_tmpl, err := template.New("html_page").Parse(html_page)
	if err != nil {
		panic(err)
	}
	err = page_tmpl.Execute(w, Page{Title: title, Body: body})
	if err != nil {
		panic(err)
	}
}

func outputJSON(w http.ResponseWriter, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
