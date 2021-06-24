package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Person struct {
	ID     int    `json:"id"`
	Nombre string `json:"name"`
	DNI    int    `json:"dni"`
}

type allPersons []Person

var Persons = allPersons{
	{
		ID:     1,
		Nombre: "Nex xd",
		DNI:    1,
	},
}

func getPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Persons)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var newPerson Person
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Inserte datos validos")
	}
	json.Unmarshal(reqBody, &newPerson)

	newPerson.ID = len(Persons) + 1
	Persons = append(Persons, newPerson)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPerson)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	personid, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Id invalido :(")
		return
	}

	for _, person := range Persons {
		if person.ID == personid {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(person)

		}
	}
}
func deletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid id")
		return
	}
	for i, p := range Persons {
		if p.ID == personID {
			Persons = append(Persons[:i], Persons[i+1:]...)
			fmt.Fprintf(w, "The persons with ID %v has been deleted successfully", personID)
		}
	}
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["id"])
	var updatedPerson Person
	if err != nil {
		fmt.Fprintf(w, " Invalid ID")
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter valid data")
	}

	json.Unmarshal(reqBody, &updatedPerson)

	for i, P := range Persons {
		if P.ID == personID {
			Persons = append(Persons[:i], Persons[i+1:]...)
			updatedPerson.ID = personID
			Persons = append(Persons, updatedPerson)
			fmt.Fprintf(w, "The person with ID %v has been updated successfuly", personID)
		}
	}

}
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi API")
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexRoute)
	r.HandleFunc("/persons", getPersons).Methods("GET")
	r.HandleFunc("/persons", createPerson).Methods("POST")
	r.HandleFunc("/persons/{id}", getPerson).Methods("GET")
	r.HandleFunc("/persons/{id}", deletePerson).Methods("DELETE")
	r.HandleFunc("/persons/{id}", updatePerson).Methods("PUT")
	// r.HandleFunc("/api/persons/{id}", DeletePersonsHandler).Methods(DELETE)

	server := &http.Server{
		Addr:           ":8000",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening...")
	server.ListenAndServe()

}
