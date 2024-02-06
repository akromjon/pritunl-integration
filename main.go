package main

import (
	"fmt"
	"net/http"
	"net/url"
	"errors"
	"os"
	"strings"
	"time"
)

var ErrorFile string = "error.log"

func main() {	
	
	action, err:=getArguments()

	if(err != nil){
		
		errMessage := fmt.Sprintf("Error: %s", err)
		
		writeErrorToFile(ErrorFile, errMessage)	

		os.Exit(1)
	}	

	MakeRequest(action)
}

func getArguments() (map[string]string, error) {

	if len(os.Args) < 6 {		

		errMessage := fmt.Sprintf("Usage: go run main.go {url} {connected/disconnected} {PritunlID} {ClientUUID}")
		
		fmt.Println(errMessage)
		
		writeErrorToFile(ErrorFile, errMessage)	
		
		return nil, errors.New(errMessage)
	}

	state := os.Args[3]

	if state != "connected" && state != "disconnected" {

		errMessage := fmt.Sprintf("Invalid State")
		
		fmt.Println(errMessage)

		writeErrorToFile(ErrorFile, errMessage)
		
		return nil, errors.New("Invalid State")
	}

	result := map[string]string{

		"url":        os.Args[1],

		"token": 	os.Args[2],

		"state":     state,

		"pritunl_user_id":  os.Args[4],
		
		"client_uuid": os.Args[5],
	}

	return result, nil
}

func MakeRequest(action map[string]string) {
	
	params := url.Values{}
	
	params.Add("state", action["state"])	

	params.Add("pritunl_user_id", action["pritunl_user_id"])

	params.Add("client_uuid", action["client_uuid"])
	
	query := params.Encode()

	req, err := http.NewRequest("POST", action["url"], strings.NewReader(query))

	if err != nil {

		errMessage := fmt.Sprintf("Error: %s", err)

		fmt.Println(errMessage)

		writeErrorToFile(ErrorFile, errMessage)

		return
	}

	req.Header.Set("Token", action["token"])

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	response, err := client.Do(req)

	if err != nil {

		errMessage := fmt.Sprintf("Error: %s", err)

		writeErrorToFile(ErrorFile, errMessage)

		fmt.Println(errMessage)

		return
	}

	defer response.Body.Close()	

	if err != nil {

		errMessage := fmt.Sprintf("Error: %s", err)
		
		writeErrorToFile(ErrorFile, errMessage)	

		return
	}

	if response.StatusCode!=200 {

		errMessage := fmt.Sprintf("Error: %s", response.Status)
		
		writeErrorToFile(ErrorFile, errMessage)	
		
		return
	}

	
}


func writeErrorToFile(filename, message string) error {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	
	if err != nil {
	
		return err
	
	}
	
	defer file.Close()

	currentTime := time.Now()

	formattedMessage := fmt.Sprintf("[%s] %s\n", currentTime.Format("02-01-2006 15:04:05"), message)

	_, err = file.WriteString(formattedMessage)
	
	if err != nil {
	
		return err
	
	}

	return nil
}
