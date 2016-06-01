package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//HEAD    - GET Meta data status
//GET     - get the crypt contents (requires decryption)
//POST    - Checkin
//DELETE  - Kills the crypt
func rCryptManage(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	w.Header().Set("Content-Type", "application/json")
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	response := APICryptResponse{StatusCode: 200, Version: APIVERSION, Success: true}

	var clientRequest ClientRequest
	var crypt Crypt
	var CryptID string
	timeNow := time.Now().Unix()

	url := strings.Split(r.URL.String(), "/")
	if len(url) >= 3 {
		CryptID = url[3]
	} else {
		CryptID = ""
	}

	if CryptID == "" {
		response.StatusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "A CryptID must be passed as part of the URL e.g. /1/crypt/XXXXXXXXXXXXXXXXXXXX"
		js, err := json.Marshal(response)
		if err == nil {
			w.WriteHeader(response.StatusCode)
			w.Write(js)
		} else {
			http.Error(w, err.Error(), 500)
		}

		return
	}

	//GET does not require authentication (the whole point is that people can see that the secret has expired)
	if r.Method == "GET" {
		//Show them the GET stuff
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("crypts"))
			cryptJSON := b.Get([]byte(CryptID))

			err := json.Unmarshal(cryptJSON, &crypt)

			//Check if the JSON decode worked OK
			if err != nil {
				return err
			} else {
				return nil
			}

		})

		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Success = false
			response.Message = err.Error()
		} else {
			if crypt.IsDestroyed {
				response.StatusCode = http.StatusGone
				response.Message = "Crypt metadata retrived but crypt has been destroyed"
			} else {
				response.StatusCode = http.StatusOK
				response.Message = "Crypt Retrieval Successful!"
			}
			response.CryptPayload = crypt
		}
	} else if r.Method == "HEAD" {
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("crypts"))
			cryptJSON := b.Get([]byte(CryptID))

			err := json.Unmarshal(cryptJSON, &crypt)

			//Check if the JSON decode worked OK
			if err != nil {
				return err
			} else {
				return nil
			}

		})

		if err != nil {
			response.StatusCode = http.StatusInternalServerError
		} else {
			if crypt.IsDestroyed {
				response.StatusCode = http.StatusGone
			} else {
				response.StatusCode = http.StatusOK
			}
		}

		w.WriteHeader(response.StatusCode)
		return

	} else {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&clientRequest)

		//If we successfully decoded the JSON then we can process it
		if err == nil {

			_, err = CheckChallenge(clientRequest.Challenge, clientRequest.UserID, clientRequest.ChallengeID, db)

			//If authentication was successful lets continue
			if err == nil {
				if r.Method == "POST" {
					err = db.Update(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte("crypts"))
						//Get the crypt
						cryptJSON := b.Get([]byte(CryptID))

						err = json.Unmarshal(cryptJSON, &crypt)

						//Check if the JSON decode worked OK
						if err != nil {
							log.Print("Error during crypt ID query: " + err.Error())
							return err
						}

						if crypt.IsDestroyed == true {
							return errors.New("Crypt retrieval was succesful but the crypt has been destroyed (either checkin timeout or by request) and is subsequently locked from further updates")
						} else {
							crypt.LastCheckIn = timeNow
							buf, err := json.Marshal(crypt)
							if err == nil {
								return b.Put([]byte(CryptID), buf)
							} else {
								log.Print("Error writing the checkin for the crypt" + err.Error())
								return err
							}
						}
					})

					if err != nil {
						if crypt.IsDestroyed == true {
							response.StatusCode = http.StatusGone
							response.Success = false
							response.Message = err.Error()
							response.CryptPayload = crypt
						} else {
							response.StatusCode = http.StatusInternalServerError
							response.Success = false
							response.Message = err.Error()
						}
					} else {
						response.Message = "Crypt Checkin Successful!"
						response.CryptPayload = crypt
					}
				} else if r.Method == "DELETE" {
					//Delete the crypt
					//DELTE THE CRYPT
					//DELTE THE CRYPT
					//DELTE THE CRYPT
					//DELTE THE CRYPT
					//DELETE THE CRYPT
					//DELETE THE CRYPT

				} else {
					response.StatusCode = http.StatusBadRequest
					response.Success = false
					response.Message = "Invalid method (/1/create only accepts POST/GET/DELETE/HEAD)"
				}

			} else {
				response.StatusCode = http.StatusUnauthorized
				response.Success = false
				response.Message = "Authentication Failed"
				log.Print(err.Error())
			}

		} else {
			response.StatusCode = http.StatusInternalServerError
			response.Success = false
			response.Message = "There was an issue decoding the JSON" + err.Error()
		}
	}

	js, err := json.Marshal(response)
	if err == nil {
		w.WriteHeader(response.StatusCode)
		w.Write(js)
	} else {
		http.Error(w, err.Error(), 500)
	}

}
