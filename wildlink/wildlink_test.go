package wildlink

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestClients(t *testing.T) {

	val, ok := os.LookupEnv("APPKEY") 
	if !ok {
		t.Skip()
		return
	}

	// New Client
	client := NewClient(&http.Client{}).SetAppID(20).SetAppKey(val)
	
	// If you would like ot keep the same device to track stats and sales
	// You will need to save your device's Key from a previous run and use it before `Connect()`.
	// Please Note the Device Key should be kept Secret, on the server side, or in a clients secure storage
	// This is a testing key and has no transactions or value.
	client.SetDevice(&Device{Key: "66626391fadb79dbbe1036a5b9cd18a84500"})
	// If you do not set up your device `Connect()` will make a new one.
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	//device := client.Device()
	//fmt.Println(device)

	iter, err := client.ConceptService.List(&ConceptListParams{ Limit:5})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for iter.Next() {
		concept, err := iter.Scan()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(concept.URL)
		// for the test shut down early
		iter.Close()
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}
}
