package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"log"
)

const HSTSEXPIRY = 94670856
const APIVERSION = 1
const DEFAULT_CHECKINDURATION = 86400
const DEFAULT_MISSCOUNT = 3
const HSTSENABLED = false
const BRASSHORNCOMMSPUBLICKEY = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
=4UbB
-----END PGP PUBLIC KEY BLOCK-----`

func GetMD5Hash(text string) string {
	//log.Printf("GetMD5Hash received %s", text)
	hash := md5.Sum([]byte(text))
	//log.Printf("Returned %s", hex.EncodeToString(hash[:]))
	return hex.EncodeToString(hash[:])
}

// Takes the challenge (which would have to have been decrypted), the user ID and the challenge ID
// Queries bolt for the challenge string using the challenge ID
// Compares the challenge strings and the user IDs
func CheckChallenge(challengeStr string, userID uint64, challengeID uint64, db *bolt.DB) (success bool, err error) {
	if challengeStr == "" {
		success = false
		err = errors.New("Empty challenge provided")
	} else {
		var challenge Challenge
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("challenges"))
			challengeJSON := b.Get(itob(challengeID))

			json.Unmarshal(challengeJSON, &challenge)
			return nil
		})

		if err == nil {
			if challenge.Challenge == challengeStr {
				if challenge.UserID == userID {
					//All OK
					success = true
					err = nil

					//Delete the challenge from the DB
					err = db.Update(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte("challenges"))

						return b.Delete(itob(challengeID))

					})

					if err != nil {
						log.Print(err.Error())

					}
				} else {
					success = false
					err = errors.New("This is not a challenge for this user")
				}
			} else {
				success = false
				err = errors.New("Challenge String Failed")
			}
		} else {
			success = false
			err = errors.New("Failed to successfully query the DB for your challenge")
		}

	}

	return success, err
}

func CheckChallengeRequest(clientRequest ClientRequest, db *bolt.DB) (success bool, err error) {
	if clientRequest.UserID != 0 {
		return true, nil
	} else {
		return false, errors.New("Invalid request")
	}
}

//TODO Do we want to charge per byte/crypt etc?
func CheckByteBalance(CryptContent string, UserID uint64, db *bolt.DB) (int, error) {
	var account Account
	byteCount := bytes.NewBufferString(CryptContent)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		accountJSON := b.Get(itob(UserID))

		json.Unmarshal(accountJSON, &account)
		return nil
	})

	if err == nil {
		if int64(byteCount.Len()) < account.ByteBudget {
			log.Printf("Permitting a crypt that's %d bytes", byteCount.Len())

			return byteCount.Len(), nil
		} else {
			log.Printf("Denying a crypt that's %d bytes on an account (%d) whose remaining budget is %d bytes (total stored: %d bytes)", byteCount.Len(), UserID, account.ByteBudget, account.ByteSize)
			return byteCount.Len(), errors.New("Crypt size exceeds your budget")
		}
	} else {
		return 0, err
	}

	/*if byteCount.Len() > 2048 {
		return 0, errors.New("Crypt size exceeded 2Kb maximum (" + strconv.Itoa(byteCount.Len()) + "bytes)")
	} else {
		log.Printf("Storing a crypt that's %d bytes", byteCount.Len())

		return 0, nil
	}*/

}

//Adjusts a users Byte balance
func IncreaseUserByteUsage(byteCount int, userID uint64, db *bolt.DB) (int64, error) {

	//Populate the users account
	var account Account
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		accountJSON := b.Get(itob(userID))

		json.Unmarshal(accountJSON, &account)
		return nil
	})

	if err == nil {
		account.ByteBudget = (account.ByteBudget - int64(byteCount))
		account.ByteSize = (account.ByteSize + int64(byteCount))

		//Update the users account
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("users"))

			buf, err := json.Marshal(account)
			if err != nil {
				return err
			}

			//Store the user account
			return b.Put(itob(account.UserID), buf)
		})

		return account.ByteSize, nil
	} else {
		return 0, err
	}
}

//Used when a user pays for more storage
func AddUserByteBalance(byteCount int64, userID uint64, db *bolt.DB) (byteBudget int64, err error) {
	//Populate the users account
	var account Account
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		accountJSON := b.Get(itob(userID))

		json.Unmarshal(accountJSON, &account)
		return nil
	})

	if err == nil {
		account.ByteBudget = (account.ByteBudget + byteCount)

		//Update the users account
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("users"))

			buf, err := json.Marshal(account)
			if err != nil {
				return err
			}

			//Store the user account
			return b.Put(itob(account.UserID), buf)
		})

		return account.ByteBudget, nil
	} else {
		return 0, err
	}
}

func EncryptForUser(secret string, userID uint64, db *bolt.DB) (cipherText string, err error) {

	var account Account
	var entityList openpgp.EntityList
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		accountJSON := b.Get(itob(userID))

		json.Unmarshal(accountJSON, &account)
		return nil
	})

	if err == nil {
		keyBuffer := bytes.NewBufferString(account.PublicKey)
		entityList, err = openpgp.ReadArmoredKeyRing(keyBuffer)
		if err == nil {
			buf := new(bytes.Buffer)
			w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
			if err != nil {
				return "", err
			}
			_, err = w.Write([]byte(secret))
			if err != nil {
				return "", err
			}
			err = w.Close()
			if err != nil {
				return "", err
			}

			// Encode to base64
			bytes, err := ioutil.ReadAll(buf)
			if err != nil {
				return "", err
			}
			encStr := base64.StdEncoding.EncodeToString(bytes)

			// Output encrypted/encoded string
			log.Println("Encrypted Secret:", encStr)

			return encStr, nil
		} else {
			return "", err
		}

	} else {
		return "", errors.New("Unable to find the users GPG key in the keystore")
	}
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func VerifyGPGPublicKey(PublicKey string) (string, error) {
	keyBuffer := bytes.NewBufferString(PublicKey)
	entityList, err := openpgp.ReadArmoredKeyRing(keyBuffer)

	if err != nil {
		return "", err
	} else {
		return entityList[0].PrimaryKey.KeyIdString(), nil
	}
}
