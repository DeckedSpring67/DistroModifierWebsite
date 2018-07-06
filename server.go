package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"math"

	"github.com/rs/cors"
)

type distro struct {
	name     string
	created  bool
	programs string
}

func startBuild() {
	for true {
		time.Sleep(10000 * time.Millisecond)
		for i, d := range distros {
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
				distros[i].created = true
			}
		}
	}
}

func createDistro(w http.ResponseWriter, r *http.Request) {
	//SETS THE HEADER
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()                         //Parse url parameters passed, then parse the response packet for the POST body (request body)
	if !canCreate(){
		fmt.Fprintf(w, "Can't create more than 3 distros at once") 
		return
	}
	fmt.Printf("Appending: %s\n", r.Form) // print information on server side.
	//Sets a blank space for the program's container
	programmi := " "
	//Initial name contains distro and date
	name := fmt.Sprintf("arch-%s.", time.Now().Format("2006-01-02"))
	// For every value in form we get the data and we put it into a struct, this will
	// create an univoque name for the ISO
	if len(r.Form) > 0 {
		var mult uint64
		mult = 1
 		for k, v := range r.Form {
			value := strings.Join(v, "")
			value = checkInjection(value);
			temp, _ := strconv.Atoi(k)
			if !isPrime(temp){
				fmt.Fprintf(w, "A program's identifier is not a prime number, ABORTING") 
				return
			}
			mult = mult * uint64(temp) //Multiplying each key (which is a prime number) to verify its non-duplicity. (Fundamental theorem of arithmetics)
			fmt.Println("val", value)
			fmt.Fprintf(w, "You selected:<br>%s,%s<br>", value, k)
			programmi += value + " "
		}
		fmt.Fprintf(w, "Divide the second part of the ISO's name to check if a program is already instelled<br>")
		name = fmt.Sprintf("%s%d",name,mult)
		currentname = name
		fmt.Println(name)
		found := checkDistro(name)
		if !found {
			distros = append(distros, distro{name: name, created: false, programs: programmi})
			fmt.Fprintf(w, "Appending: %s to the schedule", name)
		} else {
			fmt.Fprintf(w, "Distro already in creation or existing and is/will be located <a href=%s.iso>here</a>, if there are some problems contact me at deckedspring@gmail.com", name)
		}
	} else {
		fmt.Println("Error parsing form data")
		fmt.Fprintln(w, "Error parsing form data")
	}
}

//Check if it's a prime number
func isPrime(value int) bool {
    for i := 2; i <= int(math.Floor(math.Sqrt(float64(value)))); i++ {
        if value%i == 0 {
            return false
        }
    }
    return value > 1
}


func canCreate() bool{
	count := 0
	for _, x := range distros {
		if !x.created{
			count++	
		}
    }
	return !(count > 3)
}

//To sanitize the input
func checkInjection(value string) string {
	//These characters can compromise the program string
	specialChars := [17]string{";","<",">","|","&","`","$","@","(",")","{","}","[","]","*","%",","}	
    for _, x := range specialChars {
		index := strings.Index(value,x)
		//If the index is found (> -1) then I delete everything that follows the character
		if index > -1 {
			partial := strings.Split(value,x) 
			value = partial[0]
		}
	}
	return value 
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
		created := "<span style=\"color:#FA8072\">false</span>"
		if x.created {
			created = "<span style=\"color:#90EE90\">true</span>"
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
	percentage = "No distros created yet"
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
