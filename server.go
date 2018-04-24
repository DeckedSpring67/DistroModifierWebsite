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
}

type options struct {
	Browser string
	Email   string
	Media   string
	Office  string
}

func startBuild(distros []distro) {
	for _, d := range distros {
		if d.created == false {

		}
	}
}

func getSchedule(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s</table>", response)
}

func createDistro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //Parse url parameters passed, then parse the response packet for the POST body (request body)
	fmt.Println(r.Form) // print information on server side.
	//Initial name contains distro and date
	name := fmt.Sprintf("arch-%s", time.Now().Format("2006-01-02"))

	// For every value in form we get the data and we put it into a struct, this will
	// create an univoque name for the ISO
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			value := strings.Join(v, "")
			fmt.Println("key:", k)
			fmt.Println("val:", value)
			name = fmt.Sprintf("%s-%s", string(k[0]), string(value[0]))
			fmt.Fprintf(w, "%s:%s <br>", k, strings.Join(v, ""))
		}
		distros = append(distros, distro{name: name, created: false})
	} else {
		fmt.Println("Error parsing form data")
	}
	found := checkDistro(name)
	if !found {
		fmt.Fprintf(w, "Creating Distro: %s<br>", name)
	} else {
		fmt.Fprint(w, "Distro already existing, if there are some problems contact me at deckedspring@gmail.com")
	}
}

func checkDistro(name string) bool {
	found := false
	for _, v := range distros {
		vSplitted := strings.Split(v.name, "-")
		currentSplitted := strings.Split(name, "-")
		if vSplitted[2] == currentSplitted[2] {
			found = true
			break
		}
	}
	return found
}

//Global array of distros
var distros []distro
var response string

func main() {
	response = "<table>"
	http.HandleFunc("/get_schedule", getSchedule)
	http.HandleFunc("/distro_creator", createDistro)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
