package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	host = "https://api.helium.io/v1/hotspots/"
	online = "online"
)

type MyEvent struct {
	Addr string `json:"addr"`
}

type DataResponse struct {
	Data struct {
		Status struct {
			Online string `json:"online"`
		} `json:"status"`
	} `json:"data"`
}

func handler(event MyEvent) (bool, error) {
	fmt.Printf("%s\n", event)

	if (MyEvent{}) == event {
		log.Println("Empty event provided!!!")
		return false, errors.New("Empty Event Provided!!!")
	}
	url := host + event.Addr
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "NvrMineApp")

	log.Printf("Using %s\n", url)

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Oops something broke!!! \n %s", err.Error())
		return false, errors.New("OOOPS!!!")
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("The following error response code was returned: %d\n", resp.StatusCode)
		return false, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var dataResp DataResponse

	json.Unmarshal(body, &dataResp)

	log.Printf("%s\n", dataResp.Data.Status.Online)

	return strings.ToLower(string(dataResp.Data.Status.Online)) == online, nil
}

func main() {
	lambda.Start(handler)
}
