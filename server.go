package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/cors"
)

type distro struct {
	name     string
	created  bool
	programs string
}

func startBuild() {
	for true {
		for _, d := range distros {
			if d.created == false {
				currentname = d.name
				percentage = "20% (This will be long)"
				log.Printf("Creating Distro: %s", d.name)
				cmd := exec.Command("sh", "createroot.sh", d.programs)
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &out
				cmd.Run()
				fmt.Printf("Log Script: %s\n", out.String())
				percentage = "80% (This might take a bit)"
				log.Printf("Creating squashfs and iso")
				cmd = exec.Command("sh", "createsquashfs.sh", d.name)
				cmd.Stdout = &out
				cmd.Stderr = &out
				cmd.Run()
				fmt.Printf("Log Script: %s\n", out.String())
				percentage = "90%"
				log.Printf("Distro created")
				fmt.Printf("End")
				percentage = "100%, ready to take another job"
			}
		}
	}
}

func createDistro(w http.ResponseWriter, r *http.Request) {
	//SETS THE HEADER
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()                         //Parse url parameters passed, then parse the response packet for the POST body (request body)
	fmt.Printf("Appending: %s\n", r.Form) // print information on server side.
	//Initial name contains distro and date
	programmi := "xorg-server xterm xf86-video-intel xf86-video-nouveau xf86-video-amdgpu xf86-video-ati xf86-video-fbdev xf86-input-libinput xorg-xinit "
	name := fmt.Sprintf("arch-%s.", time.Now().Format("2006-01-02"))
	currentname = name
	// For every value in form we get the data and we put it into a struct, this will
	// create an univoque name for the ISO
	count := 0

	if len(r.Form) > 0 {
		for k, v := range r.Form {
			count++
			value := strings.Join(v, "")
			fmt.Println("key:", k)
			fmt.Println("val:", value)
			name += fmt.Sprintf("_%s", value)
			fmt.Fprintf(w, "Programs selected:<br>%s:%s <br>", k, value)
			programmi += value + " "
		}
		fmt.Println(name)
		found := checkDistro(name)
		if !found {
			distros = append(distros, distro{name: name, created: false, programs: programmi})
			fmt.Fprintf(w, "Appending: %s to the schedule", name)
		} else {
			fmt.Fprintf(w, "Distro already in creation or existing, if there are some problems contact me at deckedspring@gmail.com")
		}
	} else {
		fmt.Println("Error parsing form data")
		fmt.Fprintln(w, "Error parsing form data")
	}
}

func checkDistro(name string) bool {
	found := false
	for _, v := range distros {
		vSplitted := strings.Split(v.name, ".")
		currentSplitted := strings.Split(name, ".")
		if vSplitted[1] == currentSplitted[1] {
			found = true
			break
		}
	}
	return found
}

func refreshResponse() {
	response = "<br>Now Creating: <br>" + currentname + "<br>Status: " + percentage + "<br>"
	for _, x := range distros {
		created := "<span style=\"color:red\">false</span>"
		if x.created {
			created = "<span style=\"color:green\">true</span>"
		}
		response = fmt.Sprintf("%s<tr><td>%s<br>Created:&nbsp; %s</td></tr>\n", response, x.name, created)
	}
}

func getSchedule(w http.ResponseWriter, r *http.Request) {
	refreshResponse()
	fmt.Fprintf(w, "<table border sytle=\"color:gray;text-align:center\">%s</table>", response)
}

func getCurrentDistro(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating distro: %s", currentname)
}

func getDistros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	files, err := ioutil.ReadDir("/var/www/localhost/htdocs/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		filenamearr := strings.Split(f.Name(), ".")
		if len(filenamearr) > 2 {
			if filenamearr[2] == "iso" {
				fmt.Fprintf(w, "<a href=%s>%s</a><br>", f.Name(), f.Name())
			}
		}
	}
}

func appendExistingDistros() {
	files, err := ioutil.ReadDir("/var/www/localhost/htdocs/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		filenamearr := strings.Split(f.Name(), ".")
		if len(filenamearr) > 2 {
			if filenamearr[2] == "iso" {
				distros = append(distros, distro{name: f.Name(), created: true, programs: ""})
			}
		}
	}
}

//Global array of distros
var distros []distro
var response string
var percentage string
var currentname string

func main() {
	appendExistingDistros()
	response = ""
	percentage = "No distros to create yet"
	go startBuild()
	r := http.NewServeMux()
	r.HandleFunc("/get_schedule", getSchedule)
	r.HandleFunc("/distro_creator", createDistro)
	r.HandleFunc("/get_current_distro", getCurrentDistro)
	r.HandleFunc("/getDistros", getDistros)
	log.Println("Server Started at port :8080")
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
