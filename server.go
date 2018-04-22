package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type distro struct {
	name    string
	created bool
	date    string
}

func startBuild(distros []distro) {
	for _, d := range distros {
		if d.created == false {

		}
	}
}

func createDistro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //Parse url parameters passed, then parse the response packet for the POST body (request body)
	fmt.Println(r.Form) // print information on server side.
	name := fmt.Sprintf("arch-%s", time.Date())
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Fprintf(w, "Hai selezionato:\n")
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		fmt.Fprintf(w, "%s:%s \n", k, strings.Join(v, ""))
	}
}

//Global array of distros
var distros []distro

func main() {

	http.HandleFunc("/distro_creator", createDistro)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
