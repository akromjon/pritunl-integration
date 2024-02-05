package main

import (
	"fmt"
	"net/http"
	"net/url"
	"errors"
	"os"
	"strings"
)

func main() {
	
	action, err:=getArguments()

	if(err!=nil){
		fmt.Println("Error getting action:", err)
		return
	}

	MakeRequest(action)
}

func getArguments() (map[string]string, error) {

	if len(os.Args) < 6 {
		
		fmt.Println("Usage: go run main.go {action} {PritunlID} {ClientUUID}")
		
		return nil, errors.New("Invalid number of arguments")
	}

	state := os.Args[2]

	if state != "connected" && state != "disconnected" {
		fmt.Println("Invalid State")
		return nil, errors.New("Invalid State")
	}

	result := map[string]string{
		"url":        os.Args[1],
		"state":     state,
		"token": 	os.Args[3],
		"pritunl_id":  os.Args[4],
		"client_uuid": os.Args[5],
	}

	return result, nil
}

func MakeRequest(action map[string]string) {
	
	params := url.Values{}
	
	params.Add("state", action["state"])
	
	params.Add("token", action["token"])

	params.Add("pritunl_id", action["pritunl_id"])

	params.Add("client_uuid", action["client_uuid"])
	
	query := params.Encode()

	req, err := http.NewRequest("POST", action["url"], strings.NewReader(query))

	if err != nil {

		fmt.Println("Error creating request:", err)

		return
	}

	req.Header.Set("Key", action["token"])

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	response, err := client.Do(req)

	if err != nil {

		fmt.Println("Error making request:", err)

		return
	}

	defer response.Body.Close()	

	if err != nil {

		fmt.Println("Error reading response body:", err)

		return
	}

	if response.StatusCode!=200 {
		
		fmt.Println("Error:", response.Status)
		
		return
	}

	
}
