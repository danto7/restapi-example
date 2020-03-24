// (c) 2019 Christian Bargmann
//
// This project serves for teaching purposes in the CloudWP with Stefan Sarstedt
// at the University of Applied Sciences in Hamburg. The project provides a basic framework for a Restful API,
// which can be used to manage courses of study via simple web calls.
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Studiengang represents a study course at HAW Hamburg
type Studiengang struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Beschreibung string       `json:"beschreibung"`
	Kontakt      *[]Professor `json:"kontakt"`
}

// Kontakt represents one of our nice profs
type Professor struct {
	Vorname  string `json:"vorname"`
	Nachname string `json:"nachname"`
}

// Init a slice studiengaenge to store our data and mock a database
var studiengaenge []Studiengang

// getStudiengaenge returns a collection of existing study courses
func getStudiengaenge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(studiengaenge)
}

// getStudiengang returns an existing studiengang
func getStudiengang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the parameters out of the http request
	params := mux.Vars(r)

	// Loop through studiengaenge and find one with the id from the params
	// I know you can do it better, guys :-)
	for _, studiengang := range studiengaenge {
		if studiengang.ID == params["id"] {
			json.NewEncoder(w).Encode(studiengang)
			return
		}
	}
	json.NewEncoder(w).Encode(&Studiengang{})
}

// createStudiengang creates a new Studiengang
func createStudiengang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var studiengang Studiengang

	err := json.NewDecoder(r.Body).Decode(&studiengang)
	if err != nil {
		fmt.Printf("An error occured while decoding the json: json was %s \n", r.Body)
	}

	// Mock our studiengang id
	// attention: ids may not be unique using this implementation but serves well as an example
	studiengang.ID = strconv.Itoa(rand.Intn(100000))
	studiengaenge = append(studiengaenge, studiengang)
	json.NewEncoder(w).Encode(studiengang)
}

// updateStudiengang updates an existing studiengang
func updateStudiengang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// again, get our params out of the request
	params := mux.Vars(r)

	// Loop through all our study courses and update matching one if found
	// Again, does this implementation does not scale well :-)
	for index, studiengang := range studiengaenge {
		if studiengang.ID == params["id"] {

			// we will simply replace the found studiengang with our new spec
			studiengaenge = append(studiengaenge[:index], studiengaenge[index+1:]...)

			var stg Studiengang

			err := json.NewDecoder(r.Body).Decode(&stg)
			if err != nil {
				fmt.Printf("An error occured while decoding the json: json was %s \n", r.Body)
			}

			stg.ID = params["id"]
			studiengaenge = append(studiengaenge, stg)
			json.NewEncoder(w).Encode(stg)
			return
		}
	}
}

// deleteStudiengang deletes an existing studiengang
func deleteStudiengang(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, item := range studiengaenge {
		if item.ID == params["id"] {
			studiengaenge = append(studiengaenge[:index], studiengaenge[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(studiengaenge)
}

// main launches our simple studiengang restful api
// First we create some sample data. Thereafter we define our api routes and corresponding functions.
func main() {
	// Init router
	r := mux.NewRouter()

	// Hardcoded data - @todo: add database
	studiengaenge = append(studiengaenge, Studiengang{
		ID:           "1",
		Name:         "Angewandte Informatik",
		Beschreibung: "Programmieren, programmieren, programmieren...",
		Kontakt: &[]Professor{
			Professor{Vorname: "Stefan", Nachname: "Sarstedt"},
			Professor{Vorname: "Christian", Nachname: "Bargmann"},
		},
	})
	studiengaenge = append(studiengaenge, Studiengang{
		ID:           "2",
		Name:         "Wirtschaftsinformatik",
		Beschreibung: "Programmieren, und ein bisschen BWL...",
		Kontakt: &[]Professor{
			Professor{Vorname: "Ulrike", Nachname: "Steffens"},
		},
	})
	studiengaenge = append(studiengaenge, Studiengang{
		ID:           "2",
		Name:         "Technische Informatik",
		Beschreibung: "Ich mag Platinen...",
		Kontakt: &[]Professor{
			Professor{Vorname: "Thomas", Nachname: "Lehmann"},
		},
	})

	// Route handles & endpoints
	setupEndpoints(r)

	// Start server
	fmt.Printf("webserver is running at 0.0.0.0:8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func setupEndpoints(r *mux.Router) {

	r.HandleFunc("/studiengaenge", getStudiengaenge).Methods("GET")
	r.HandleFunc("/studiengaenge/{id}", getStudiengang).Methods("GET")
	r.HandleFunc("/studiengaenge", createStudiengang).Methods("POST")
	r.HandleFunc("/studiengaenge/{id}", updateStudiengang).Methods("PUT")
	r.HandleFunc("/studiengaenge/{id}", deleteStudiengang).Methods("DELETE")
}
