package wildlink

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestClients(t *testing.T) {

	key, ok := os.LookupEnv("APPKEY")
	if !ok {
		t.Skip()
		return
	}

	appID, ok := os.LookupEnv("APPID")
	if !ok {
		t.Skip()
		return
	}

	id, err := strconv.ParseUint(appID, 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dKey, ok := os.LookupEnv("DEVICEKEY")
	if !ok {
		t.Skip()
		return
	}

	// New Client
	client := NewClient(&http.Client{}).SetAppID(id).SetAppKey(key)

	// If you would like ot keep the same device to track stats and sales
	// You will need to save your device's Key from a previous run and use it before `Connect()`.
	// Please Note the Device Key should be kept Secret, on the server side, or in a clients secure storage
	// This is a testing key and has no transactions or value.
	client.SetDevice(&Device{Key: dKey})
	// If you do not set up your device `Connect()` will make a new one.
	err = client.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//device := client.Device()
	//fmt.Println(device)

	iter, err := client.ConceptService.List(&ConceptListParams{Limit: 5})
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

	arIter, err := client.NLPService.ListAnalysis(&AnalyzeParams{Content: "marriott in san diego"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer arIter.Close()
	contentWithLinks := make([]string, 0)
	for arIter.Next() {
		result, err := arIter.Scan()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		contentWithLinks = append(contentWithLinks, result.ContentPart...)
		if result.URL != "" {
			contentWithLinks = append(contentWithLinks, fmt.Sprintf("(%v)", result.URL))
		}
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(strings.Join(contentWithLinks, " "))
}
