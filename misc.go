package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DoPost (url string, v interface{})([]byte, error){



	jsn, err := json.Marshal(v)

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsn))

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	token := "Bearer " + access_token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	client := &http.Client{}

	//req.Header.Add("Authorization", string("Bearer " + c.token))

	resp, err := client.Do(req)

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	return body,err

}
