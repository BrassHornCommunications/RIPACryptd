package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handles requests to /
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

// Handles requests to /about/
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

// Handles requests to /faq/
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

// This function will extract a CryptID from the URL path (/view/CRYPTID/) and
// query the DB for the crypt.
// Once found it will display various information in a nicely formatted way.
func rCryptWebView(w http.ResponseWriter, r *http.Request, db *bolt.DB, templateData TemplateConf) {
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	var CryptID string
	var thisCrypt Crypt
	type ViewTemplateConf struct {
		FQDN        string
		ListenPort  int64
		Crypt       Crypt
		LastCheckIn string
	}

	thisTemplateConf := ViewTemplateConf{FQDN: templateData.FQDN, ListenPort: templateData.ListenPort}

	url := strings.Split(r.URL.String(), "/")
	if len(url) >= 2 {
		CryptID = url[2]
	} else {
		CryptID = ""
	}

	if CryptID == "" {
		http.Error(w, "A CryptID must be passed as part of the URL e.g. /view/XXXXXXXXXXXXXXXXXXXX", 400)

		return
	} else {
		log.Println("Found CryptID: " + CryptID)
	}

	tmpl, err := template.New("view").ParseFiles("assets/templates/view.html")

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("crypts"))
		cryptJSON := b.Get([]byte(CryptID))

		err := json.Unmarshal(cryptJSON, &thisCrypt)

		log.Println(string(cryptJSON))
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		thisTemplateConf.Crypt = thisCrypt

		tm := time.Unix(thisCrypt.LastCheckIn, 0)
		thisTemplateConf.LastCheckIn = tm.Format(time.RFC822)

		err = tmpl.Execute(w, thisTemplateConf)

		log.Println("CryptID:")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// This will probably never be exposed properly
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
		log.Println(err)
	}

}
