package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func rCryptChallenge(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	w.Header().Set("Content-Type", "application/json")
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	response := APIChallengeResponse{StatusCode: 200, Version: APIVERSION, Success: true}

	if r.Method == "POST" {
		var challenge Challenge
		var clientRequest ClientRequest
		var challengeSecret string

		challenge.Challenge = GetMD5Hash(strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(2048)) + strconv.FormatUint(clientRequest.UserID, 10))
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&clientRequest)

		if err == nil {
			_, err = CheckChallengeRequest(clientRequest, db)

			if err == nil {
				//Generate (and encrypt)the challenge
				challengeSecret, err = EncryptForUser(challenge.Challenge, clientRequest.UserID, db)

				if err == nil {
					log.Printf("Challenge is: %s", challenge.Challenge)

					//Add a double check that this challenge is only for this user
					challenge.UserID = clientRequest.UserID

					err = db.Update(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte("challenges"))
						id, _ := b.NextSequence()
						challenge.ChallengeID = uint64(id)

						buf, err := json.Marshal(challenge)
						if err != nil {
							return err
						}

						return b.Put(itob(id), buf)

					})

					if err == nil {
						response.Message = "challenge Successfully generated!"
						response.UserID = clientRequest.UserID
						response.Challenge = challengeSecret
						response.ChallengeID = challenge.ChallengeID

					} else {
						response.StatusCode = http.StatusInternalServerError
						response.Success = false
						response.Message = "There was an issue creating your secret"
						log.Print(err.Error())
					}
				} else {
					response.StatusCode = http.StatusInternalServerError
					response.Success = false
					response.Message = "There was an issue encrypting the secret with your public key"
					log.Print(err.Error())
				}

			} else {
				response.StatusCode = http.StatusNotFound
				response.Success = false
				response.Message = "Challenge requests must pass a corresponding userid and public key fingerprint"
				log.Print("Challenge generation failed due to mismatched userid/fingerprint")
			}
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
