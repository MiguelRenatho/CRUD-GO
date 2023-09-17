package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type person struct {
	DNI  int `json:"DNI"`
    NOMBRES string `json:"Nombres"`
    APELLIDOS  string `json:"APELLIDOS"`
	FECHANAC  string `json:"FECHANAC"`
	EDAD  int `json:"EDAD"`
	CIUDAD string `json:"CIUDAD"`
}

type allPersons = []person

var persons = allPersons {
	{
		DNI: 12345678,
		NOMBRES: "Miguel", 
    	APELLIDOS: "Vegas",   
		FECHANAC: "05/04/1998",    
		EDAD: 25,     
		CIUDAD: "Lima",  
	},
}

func getPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["dni"]) 

	if err!= nil {
		fmt.Fprintf(w, "No existe DNI")
		return
	}

	for _, person := range persons {
		if person.DNI == personID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(person)
		}
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["dni"]) 

	if err!= nil {
		fmt.Fprintf(w, "No existe DNI")
		return
	}

	for i, p := range persons {
		if p.DNI == personID {
			persons = append(persons[:i], persons[i + 1:]...)
			fmt.Fprintf(w, "La persona fue eliminada")
		}
	}
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["dni"]) 
	var updatePersona person

	if err!= nil {
		fmt.Fprintf(w, "No existe DNI")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserta Persona Valida")
	}
	json.Unmarshal(reqBody, &updatePersona)

	for i, p := range persons {
		if p.DNI == personID {
			persons = append(persons[:i], persons[i + 1:]...)
			updatePersona.DNI = personID
			persons = append(persons, updatePersona)

			fmt.Fprintf(w, "Se actualizo la persona")
		}
	}

}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var newPerson person
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserta Persona Valida")
	}
	json.Unmarshal(reqBody, &newPerson)
	persons = append(persons, newPerson)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPerson)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PROYECTO COMPARTAMOS")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index)
	router.HandleFunc("/persons", getPersons).Methods("GET")
	router.HandleFunc("/persons", createPerson).Methods("POST")
	router.HandleFunc("/persons/{dni}", getPerson).Methods("GET")
	router.HandleFunc("/persons/{dni}", deletePerson).Methods("DELETE")
	router.HandleFunc("/persons/{dni}", updatePerson).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", router))
}

