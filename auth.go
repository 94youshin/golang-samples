package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	auth := "http://10.124.210.160:58000/iam/token"
	user := "container_devops"
	pwd := "g0904gfldan32kt5409sf"

	userMsg := user + ":" + pwd
	base64Msg := base64.StdEncoding.EncodeToString([]byte(userMsg))
	request, err := http.NewRequest("POST", auth + "?grant_type=client_credentials", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	request.Header.Add("Authorization", "Basic " + base64Msg)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	token, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(token))

}
