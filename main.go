package main

import (
	//"html/template"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

func main() {
	log.Println("---------------------------------------------")

	//Grab all our command line config
	configuration := flag.String("conf", "", "path to configuration file")
	flag.Parse()
	conf := readConfig(*configuration)

	//Populate our template conf - normally passed by value to other functions
	tmplConf := TemplateConf{
		FQDN:       conf.FQDN,
		ListenPort: conf.ListenPort,
	}

	//We have our main web thread and the watcher threads
	runtime.GOMAXPROCS(3)

	//Tell people what we're up to
	if *configuration == "" {
		log.Println("Using default configuration")
	} else {
		log.Println("Using config: " + *configuration)
	}
	log.Println("DB Path is: " + conf.DbPath)
	log.Println("FQDN is:" + conf.FQDN)
	log.Println("Listening on: (" + conf.ListenIP + " / " + conf.ListenIPv6 + " ) : " + strconv.FormatInt(conf.ListenPort, 10))

	//We need a DB for holding interactions
	db, err := bolt.Open(conf.DbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DB opened, all is OK")
	}
	defer db.Close()

	//Make sure our buckets exist
	for _, bucket := range []string{"users", "crypts", "challenges"} {
		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
	}

	//Start our watcher process
	go CryptWatcher(db)

	//Human readable pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { rCryptWebIndex(w, r, tmplConf) })
	http.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) { rCryptWebAbout(w, r, tmplConf) })
	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) { rCryptWebView(w, r, db, tmplConf) })
	http.HandleFunc("/create/", func(w http.ResponseWriter, r *http.Request) { rCryptWebCreate(w, r, db, tmplConf) })
	http.HandleFunc("/faq/", func(w http.ResponseWriter, r *http.Request) { rCryptWebFAQ(w, r, tmplConf) })

	//Help
	http.HandleFunc("/help/", func(w http.ResponseWriter, r *http.Request) { rCryptWebHelp(w, r, tmplConf) })
	http.HandleFunc("/help/create/", func(w http.ResponseWriter, r *http.Request) { rCryptWebHelpCreate(w, r, tmplConf) })
	http.HandleFunc("/help/checkin/", func(w http.ResponseWriter, r *http.Request) { rCryptWebHelpCheckin(w, r, tmplConf) })
	http.HandleFunc("/help/view/", func(w http.ResponseWriter, r *http.Request) { rCryptWebHelpView(w, r, tmplConf) })
	http.HandleFunc("/help/destroy/", func(w http.ResponseWriter, r *http.Request) { rCryptWebHelpDestroy(w, r, tmplConf) })
	http.HandleFunc("/help/challenge/", func(w http.ResponseWriter, r *http.Request) { rCryptWebHelpChallenge(w, r, tmplConf) })
	http.HandleFunc("/help/newbtc/", func(w http.ResponseWriter, r *http.Request) { rCryptWebHelpBitcoin(w, r, tmplConf) })

	//Assets (css, images etc)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("assets/css/"))))
	http.Handle("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir("assets/font/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("assets/images/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js/"))))

	//Programatic interfaces
	//API 1.0
	http.HandleFunc("/1/register/", func(w http.ResponseWriter, r *http.Request) { rCryptRegister(w, r, db, conf) })
	http.HandleFunc("/1/challenge/", func(w http.ResponseWriter, r *http.Request) { rCryptChallenge(w, r, db) })
	http.HandleFunc("/1/status/", func(w http.ResponseWriter, r *http.Request) { rCryptStatus(w, r, db) })

	//The crypt itself
	//HEAD 		- GET Meta data status
	//GET 		- get the crypt contents (requires decryption)
	//POST 		- Checkin
	//DELETE 	- Kills the crypt
	http.HandleFunc("/1/crypt/", func(w http.ResponseWriter, r *http.Request) { rCryptManage(w, r, db) })

	//Creates a new crypt
	http.HandleFunc("/1/crypt/new/", func(w http.ResponseWriter, r *http.Request) { rCryptCreateCrypt(w, r, db) })

	//In production we'll probably be using nginx or a Tor HS to proxy this
	ListenPort := strconv.FormatInt(conf.ListenPort, 10)
	if conf.TLS {
		v4err := http.ListenAndServeTLS(conf.ListenIP+":"+ListenPort, conf.TLSCert, conf.TLSKey, nil)

		if v4err != nil {
			log.Fatal(v4err)
		}

		//Seems ListenAndServe doesn't support v6 yet
		//https://golang.org/src/net/http/server.go?s=65528:65583#L2777
		/*v6err := http.ListenAndServeTLS(conf.ListenIPv6+":"+ListenPort, conf.TLSCert, conf.TLSKey, nil)

		  if v6err != nil {
		    log.Fatal(v6err)
		  }*/
	} else {
		v4err := http.ListenAndServe(conf.ListenIP+":"+ListenPort, nil)
		if v4err != nil {
			log.Fatal(v4err)
		}

		//Seems ListenAndServe doesn't support v6 yet
		//https://golang.org/src/net/http/server.go?s=65528:65583#L2094
		/*v6err := http.ListenAndServe(conf.ListenIPv6+":"+ListenPort, nil)

		if v6err != nil {
			log.Fatal(v6err)
		}*/

	}
}

// Reads our JSON formatted config file
// and returns a struct
func readConfig(filename string) CoreConf {
	var conf CoreConf

	if filename == "" {
		conf.DbPath = "./rcrypt.bolt"
		conf.ListenIP = "127.0.0.1"
		conf.ListenIPv6 = "::1"
		conf.ListenPort = 8080
		conf.FQDN = "ripacrypt.download"
	} else {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal("Cannot read configuration file ", filename)
		}
		err = json.Unmarshal(b, &conf)
		if err != nil {
			log.Fatal("Cannot parse configuration file ", filename)
		}

		if conf.DbPath == "" {
			conf.DbPath = "./rcrypt.bolt"
		}

		if conf.ListenIP == "" {
			conf.ListenIP = "127.0.0.1"
		}

		if conf.ListenIPv6 == "" {
			conf.ListenIPv6 = "::1"
		}

		if conf.ListenPort == 0 {
			conf.ListenPort = 8080
		}

		if conf.FQDN == "" {
			conf.FQDN = "ripacrypt.download"
		}
	}
	return conf
}
