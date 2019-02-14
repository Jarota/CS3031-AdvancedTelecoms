package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	client := &http.Client{}
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Contacting", r.URL.Hostname())
		r.RequestURI = ""
		resp, err := client.Do(r)
		if err != nil {
			log.Fatal(err)
		}
		resp.Write(w)
		//w.Write([]byte("RECEIVED RESPONSE"))
	}))
	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
