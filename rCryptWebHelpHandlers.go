package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func rCryptWebHelp(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("help-index").ParseFiles("assets/templates/help-index.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}

func rCryptWebHelpCreate(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("help-create").ParseFiles("assets/templates/help-create.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}

func rCryptWebHelpCheckin(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("help-checkin").ParseFiles("assets/templates/help-checkin.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}

func rCryptWebHelpView(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("help-view").ParseFiles("assets/templates/help-view.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}

func rCryptWebHelpDestroy(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("help-destroy").ParseFiles("assets/templates/help-destroy.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}

func rCryptWebHelpChallenge(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("help-challenge").ParseFiles("assets/templates/help-challenge.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}

func rCryptWebHelpBitcoin(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("help-bitcoin").ParseFiles("assets/templates/help-bitcoin.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}
}
