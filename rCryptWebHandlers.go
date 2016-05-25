package main

import (
	"github.com/boltdb/bolt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func rCryptWebIndex(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("index").ParseFiles("assets/templates/index.html")
	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}

}

func rCryptWebAbout(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("about").ParseFiles("assets/templates/about.html")

	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}

}

func rCryptWebFAQ(w http.ResponseWriter, r *http.Request, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("faq").ParseFiles("assets/templates/faq.html")

	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}

}

func rCryptWebView(w http.ResponseWriter, r *http.Request, db *bolt.DB, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	tmpl, err := template.New("view").ParseFiles("assets/templates/view.html")

	var crypt []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("crypts"))
		crypt = b.Get([]byte("test"))

		//fmt.Printf("The answer is: %s\n", v)
		return nil
	})

	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}

}

func rCryptWebCreate(w http.ResponseWriter, r *http.Request, db *bolt.DB, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}
	tmpl, err := template.New("create").ParseFiles("assets/templates/create.html")

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("crypts"))
		err := b.Put([]byte("hash"), []byte("test"))
		return err
	})

	err = tmpl.Execute(w, templateData)

	if err != nil {
		log.Fatal(err)
	}

}
