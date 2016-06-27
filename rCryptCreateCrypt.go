package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Takes a HTTP POST payload and creates a locally stored crypt.
// This function should set sensible defaults if their values are not present
func rCryptCreateCrypt(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	w.Header().Set("Content-Type", "application/json")
	if HSTSENABLED == true {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.FormatInt(HSTSEXPIRY, 10)+"; includeSubdomains")
	}

	response := APICryptResponse{StatusCode: 200, Version: APIVERSION, Success: true}

	if r.Method == "POST" {
		//The created time is stored
		createdTime := time.Now().Unix()
		md5Hash := GetMD5Hash(strconv.FormatInt(createdTime, 10))

		//r.ParseForm()
		decoder := json.NewDecoder(r.Body)
		var clientRequest ClientRequest
		var checkInDuration, missCount int64
		var byteSize int

		err := decoder.Decode(&clientRequest)

		//If we successfully decoded the JSON then we can process it
		if err == nil {

			_, err = CheckChallenge(clientRequest.Challenge, clientRequest.UserID, clientRequest.ChallengeID, db)

			//If authentication was successful lets continue
			if err == nil {

				byteSize, err = CheckByteBalance(clientRequest.CryptContent, clientRequest.UserID, db)

				//If they've got enough in their account we can store the crypt
				if err == nil {

					_, err = IncreaseUserByteUsage(byteSize, clientRequest.UserID, db)

					if err != nil {
						log.Printf("User %d wasn't charged for %d bytes of usage!", clientRequest.UserID, byteSize)
					}

					//Create a new crypt
					if clientRequest.CheckInDuration == 0 {
						checkInDuration = DEFAULT_CHECKINDURATION
					} else {
						checkInDuration = clientRequest.CheckInDuration
					}

					if clientRequest.MissCount == 0 {
						missCount = DEFAULT_MISSCOUNT
					} else {
						missCount = clientRequest.MissCount
					}

					if checkInDuration >= 3600 {

						crypt := Crypt{UserID: clientRequest.UserID,
							CryptID:         md5Hash,
							CreateTimeStamp: createdTime,
							IsDestroyed:     false,
							CheckInDuration: checkInDuration,
							MissCount:       missCount,
							Description:     clientRequest.Description,
							LastCheckIn:     createdTime,
							CipherText:      clientRequest.CryptContent}

						buf, err := json.Marshal(crypt)

						if err == nil {
							err = db.Update(func(tx *bolt.Tx) error {
								b := tx.Bucket([]byte("crypts"))
								return b.Put([]byte(md5Hash), buf)
							})

							//Store a reference to this crypt against the user
							//blah

							if err != nil {
								log.Fatal(err)
								response.StatusCode = http.StatusInternalServerError
								response.Success = false
								response.Message = err.Error()
							} else {
								response.StatusCode = http.StatusCreated
								response.Message = "Crypt successfully created!"
								response.CryptPayload = crypt
							}
						} else {
							response.StatusCode = http.StatusInternalServerError
							response.Success = false
							response.Message = "There was an issue storing your crypt"
						}

					} else {
						response.StatusCode = http.StatusConflict
						response.Success = false
						response.Message = "Check in durations must be at least one hour (60 minutes, 3600 seconds)"
					}
				} else {
					response.StatusCode = http.StatusRequestEntityTooLarge
					response.Success = false
					response.Message = "Insuffcient byte balance - try topping up or reducing the byte size of your crypt"
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

		/*} else if r.Method == "GET" {
			//Show them the GET stuff
			decoder := json.NewDecoder(r.Body)
			var clientRequest ClientRequest
			err := decoder.Decode(&clientRequest)

			//If we successfully decoded the JSON then we can process it
			if err == nil {

				_, err = CheckChallenge(clientRequest.Challenge, clientRequest.UserID)

				//If authentication was successful lets continue
				if err == nil {

					//EXPORT THE CRYPT
					//EXPORT THE CRYPT
					//EXPORT THE CRYPT
					//EXPORT THE CRYPT
					//EXPORT THE CRYPT
					//EXPORT THE CRYPT
					//EXPORT THE CRYPT

				} else {
					response.StatusCode = 403
					response.Success = false
					response.Message = "Authentication Failed"
				}

			} else {
				response.StatusCode = 500
				response.Success = false
				response.Message = "There was an issue decoding the JSON"
			}

		} else if r.Method == "DELETE" {
			//Delete the crypt
			decoder := json.NewDecoder(r.Body)
			var clientRequest ClientRequest
			err := decoder.Decode(&clientRequest)

			//If we successfully decoded the JSON then we can process it
			if err == nil {

				_, err = CheckChallenge(clientRequest.Challenge, clientRequest.UserID)

				//If authentication was successful lets continue
				if err == nil {

					//DELTE THE CRYPT
					//DELTE THE CRYPT
					//DELTE THE CRYPT
					//DELTE THE CRYPT
					//DELETE THE CRYPT
					//DELETE THE CRYPT

				} else {
					response.StatusCode = 403
					response.Success = false
					response.Message = "Authentication Failed"
				}

			} else {
				response.StatusCode = 500
				response.Success = false
				response.Message = "There was an issue decoding the JSON"
			}*/

	} else {
		response.StatusCode = http.StatusBadRequest
		response.Success = false
		response.Message = "Invalid method (/1/create only accepts POST)"
	}

	js, err := json.Marshal(response)
	if err == nil {
		w.WriteHeader(response.StatusCode)
		w.Write(js)
	} else {
		http.Error(w, err.Error(), 500)
	}

}
