package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	var cmd string

	for {
		fmt.Scanln(&cmd)
		resp, err := http.Get("http://localhost:8080/" + cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(string(body))
	}
}
