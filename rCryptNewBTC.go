package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"strconv"
)

// This function takes a HTTP Post, verifies the users challenge nonce and
// then requests a new bitcoin address.
// Returns a JSON payload
func rCryptNewBTC(w http.ResponseWriter, r *http.Request, db *bolt.DB, conf CoreConf) {
	w.Header().Set("Content-Type", "application/json")
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	response := APIRegisterResponse{StatusCode: 200, Version: APIVERSION, Success: true}

	if r.Method == "POST" {
		//md5Hash := GetMD5Hash(string(time.Now().Unix()))
		var account Account
		var clientRequest ClientRequest

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&clientRequest)

		//Check the users challenge
		_, challengeErr := CheckChallenge(clientRequest.Challenge, clientRequest.UserID, clientRequest.ChallengeID, db)

		if challengeErr == nil {

			//Retrieve the users account
			accountErr := db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("users"))
				accountJSON := b.Get(itob(clientRequest.UserID))

				err = json.Unmarshal(accountJSON, &account)
				return err
			})

			if accountErr == nil {

				//Get a new address
				account.BTCAddr, err = getBTCAddr(conf.BTCAddr, conf.BTCUser, conf.BTCPass, true)

				//This is not ideal, we should probably bomb out instead
				if err != nil {
					account.BTCAddr = "Not Available"
				}

				err = db.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("users"))
					id, _ := b.NextSequence()
					account.UserID = uint64(id)

					// Marshal user data into bytes.
					buf, err := json.Marshal(account)
					if err != nil {
						return err
					}

					//Store the user account
					return b.Put(itob(account.UserID), buf)
				})

				if err == nil {
					response.Message = "A new bitcoin address has been successfully generated"
					response.UserID = account.UserID
					response.BTCAddr = account.BTCAddr

				} else {
					response.StatusCode = http.StatusInternalServerError
					response.Success = false
					response.Message = "There was an issue generating a new bitcoin address"
					log.Print(err.Error())
				}
			} else {
				response.StatusCode = http.StatusBadRequest
				response.Success = false
				response.Message = "There was an issue retrieving your account"
				log.Print(accountErr.Error())

			}
		} else {
			response.StatusCode = http.StatusUnauthorized
			response.Success = false
			response.Message = "Authentication Failed"
			log.Print(challengeErr.Error())
		}
	} else {
		response.StatusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "Invalid method (/1/register only accepts POST)"
	}

	js, err := json.Marshal(response)
	if err == nil {
		w.WriteHeader(response.StatusCode)
		w.Write(js)
	} else {
		http.Error(w, err.Error(), 500)
	}
}
