package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"math"
	"time"
)

func CryptWatcher(db *bolt.DB) {
	for {
		var crypt Crypt
		checkTime := time.Now().Unix()

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("crypts"))

			b.ForEach(func(cryptID, cryptJSON []byte) error {
				err := json.Unmarshal(cryptJSON, &crypt)

				log.Printf("CryptID: %s | Crypt Last Checkin: %d", cryptID, crypt.LastCheckIn)
				if !crypt.IsDestroyed {
					duration := checkTime - crypt.LastCheckIn

					if duration > crypt.CheckInDuration {
						missCount := math.Mod(float64(duration), float64(crypt.CheckInDuration))

						if missCount > float64(crypt.MissCount) {
							log.Printf("CryptID: %s | Crypt Last Checkin: %d is over threshold (%d) - destroying", cryptID, crypt.LastCheckIn, crypt.MissCount)
							err = db.Update(func(tx *bolt.Tx) error {
								b := tx.Bucket([]byte("crypts"))

								crypt.CipherText = ""
								crypt.IsDestroyed = true

								buf, err := json.Marshal(crypt)
								if err == nil {
									return b.Put([]byte(cryptID), buf)
								} else {
									log.Print("Error writing to the record" + err.Error())
									return err
								}
							})

						} else {
							log.Printf("CryptID: %s | Crypt Last Checkin: %d is nearing threshold (%f of %d)", cryptID, crypt.LastCheckIn, missCount, crypt.MissCount)
						}
					} else {
						log.Printf("CryptID: %s is within boundaries (%d / %d)", cryptID, duration, crypt.CheckInDuration)
					}
				} else {
					log.Printf("CryptID: %s has already been destroyed", cryptID)
				}
				log.Printf("----------------------")

				return err
			})
			return nil
		})

		//And now we sleep
		time.Sleep(300000 * time.Millisecond)
	}
}

// This watcher is for keeping track of users btc addresses and ensuring that
// balance changes are reflected in their usage budget
func BTCWatcher(db *bolt.DB, conf CoreConf) {
	for {
		//Check all accounts, ensure the balance is correct
		time.Sleep(300000 * time.Millisecond)
	}
}
