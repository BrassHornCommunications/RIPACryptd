package main

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcrpcclient"
	"log"
	"time"
)

func getBTCAddr(BTCURL, BTCUser, BTCPass string, BTCDisableTLS bool) (string, error) {
	log.Printf("User: %s, Pass: %s", BTCUser, BTCPass)

	connCfg := &btcrpcclient.ConnConfig{
		Host:         BTCURL,
		Endpoint:     "ws",
		User:         BTCUser,
		Pass:         BTCPass,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnBlockConnected: func(hash *wire.ShaHash, height int32, time time.Time) {
			log.Printf("Block connected: %v (%d) %v", hash, height, time)
		},
		OnBlockDisconnected: func(hash *wire.ShaHash, height int32, time time.Time) {
			log.Printf("Block disconnected: %v (%d) %v", hash, height, time)
		},
	}
	client, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Println("Error Creating a client")
		log.Fatal(err)
	}
	address, err := client.GetNewAddress("")
	if err != nil {
		log.Println("Error getting a new address")
		return "", err
	} else {
		//return "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy", nil
		return address.String(), err
	}
}

/*func BitoinHandler(w http.ResponseWriter, r *http.Request, conf CoreConf, db *bolt.DB) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Fprintf(w, "{\"success\":true, \"btc\":\"fhsfsfhosfisfsduf8sduf\",\"mbtc\":2344}")
	var response ResponsePayload
	response.Success = false
	response.DateTime = time.Now().Format(time.RFC3339)
	response.BTCAddr, _ = getBTCAddr(conf.BTCURL, conf.BTCUser, conf.BTCPass, conf.BTCDisableTLS)
	response.Cost = conf.HourlyBTCCost

	//Return the content to the user
	js, err := json.Marshal(response)
	if err == nil {
		w.Write(js)
	} else {
		http.Error(w, err.Error(), 500)
	}
}*/
