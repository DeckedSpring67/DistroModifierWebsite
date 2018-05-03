package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/cors"
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

	fmt.Fprintf(w, "<table>%s</table>", response)
}

func createDistro(w http.ResponseWriter, r *http.Request) {
	//SETS THE HEADER
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	r.ParseForm()       //Parse url parameters passed, then parse the response packet for the POST body (request body)
	fmt.Println(r.Form) // print information on server side.
	//Initial name contains distro and date
	programs := "xorg-server xterm xf86-video-intel xf86-video-nouveau xf86-video-amdgpu xf86-video-ati xf86-video-fbdev xf86-input-libinput "
	name := fmt.Sprintf("arch-%s.", time.Now().Format("2006-01-02"))
	percentage = "0%"

	// For every value in form we get the data and we put it into a struct, this will
	// create an univoque name for the ISO
	count := 0
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			count++
			value := strings.Join(v, "")
			fmt.Println("key:", k)
			fmt.Println("val:", value)
			name += fmt.Sprintf("%s%s", k, value)
			fmt.Println(name)
			fmt.Fprintf(w, "%s:%s <br>", k, value)
			programs += value + " "
			percentage = fmt.Sprintf("%d", count) + "%"
		}
		distros = append(distros, distro{name: name, created: false})
		found := false
		percentage = "15%"
		found = checkDistro(name)
		if !found {
			percentage = "20% (This will be long)"
			fmt.Fprintf(w, "Creating Distro: %s<br>", name)
			cmd := exec.Command("sh", "../createroot.sh", programs)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out
			err := cmd.Run()
			fmt.Printf("Log Script: %s\n", out.String())
			percentage = "80% (This might take a bit)"
			cmd = exec.Command("sh", "../createsquashfs.sh", name)
			cmd.Stdout = &out
			cmd.Stderr = &out
			cmd.Run()
			fmt.Printf("Log Script: %s\n", out.String())
			percentage = "90%"
			if err != nil {
				fmt.Fprintf(w, "Distro %s created at <a href=deckedhost.ns0.it/%s.iso></a>", name, name)
			} else {
				fmt.Fprintf(w, "Something went wrong")
			}
		} else {
			fmt.Fprintf(w, "Distro already existing at <a href=deckedhost.ns0.it/%s.iso></a>, if there are some problems contact me at deckedspring@gmail.com", name)
		}
	} else {
		fmt.Println("Error parsing form data")
		fmt.Fprintln(w, "Error parsing form data")
	}
	fmt.Printf("End")
	percentage = "100%"
}

func checkDistro(name string) bool {
	found := false
	for i, v := range distros {
		if i != len(distros)-1 {
			vSplitted := strings.Split(v.name, ".")
			currentSplitted := strings.Split(name, ".")
			if vSplitted[1] == currentSplitted[1] {
				found = true
				break
			}
		}
	}
	return found
}

func refreshResponse() {
	for true {
		for _, x := range distros {
			created := "false"
			if x.created {
				created = "true"
			}
			response = fmt.Sprintf("%s<tr><td>%s</td><td>created:%s</td></tr>\n", response, x.name, created)
		}
	}
}

func getPercentage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "%s", percentage)
}

//Global array of distros
var distros []distro
var response string
var percentage string

func main() {
	response = ""

	r := http.NewServeMux()
	r.HandleFunc("/get_schedule", getSchedule)
	r.HandleFunc("/distro_creator", createDistro)
	r.HandleFunc("/current_percentage", getPercentage)
	go refreshResponse()
	log.Println("Server Started at port :8080")
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
