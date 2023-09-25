package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type person struct {
	DNI  string `json:"DNI"`
    NOMBRES string `json:"NOMBRES"`
    APELLIDOS  string `json:"APELLIDOS"`
	FECHANAC  string `json:"FECHANAC"`
	EDAD  string `json:"EDAD"`
	CIUDAD string `json:"CIUDAD"`
}

func coneccionMysql()(conexion *sql.DB) {
	Driver:="mysql"
	Usuario:="root"
	Contraseña:="Renatho2023."
	Nombre:="COMPARTAMOS"

	conexion, err:= sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)

	if err!=nil {
		panic(err.Error())
	}

	return conexion
}

func getPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var personas []person
	
	conexion:= coneccionMysql()
	listaPersonas, err:= conexion.Query("select * from PERSONA")
	if err!=nil {
		panic(err.Error())
	}

	for listaPersonas.Next(){
		var persona person
		  err = listaPersonas.Scan(&persona.DNI, &persona.NOMBRES, &persona.APELLIDOS, &persona.FECHANAC, &persona.EDAD, &persona.CIUDAD);
		  personas = append(personas, persona)
	}
	json.NewEncoder(w).Encode(personas)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var persona person
	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["dni"]) 

	if err!= nil {
		fmt.Fprintf(w, "No existe DNI")
		return
	}
	conexion:= coneccionMysql()
	err = conexion.QueryRow("SELECT * FROM PERSONA WHERE DNI = ?", personID).Scan(&persona.DNI, &persona.NOMBRES, &persona.APELLIDOS, &persona.FECHANAC, &persona.EDAD, &persona.CIUDAD)

	json.NewEncoder(w).Encode(persona)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	personID, err := strconv.Atoi(vars["dni"]) 

	if err!= nil {
		fmt.Fprintf(w, "No existe DNI")
		return
	}
	conexion:= coneccionMysql()
	_, err = conexion.Exec("DELETE FROM PERSONA WHERE DNI = ?", personID)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
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

	decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&updatePersona); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	conexion:= coneccionMysql()
	_, err = conexion.Exec("UPDATE PERSONA SET NOMBRES = ?, APELLIDOS = ?, FECHANAC = ?, EDAD = ?, CIUDAD = ? WHERE DNI = ?", 
	updatePersona.NOMBRES, updatePersona.APELLIDOS, updatePersona.FECHANAC, updatePersona.EDAD, updatePersona.CIUDAD, personID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

}

func createPerson(w http.ResponseWriter, r *http.Request) {  
	var newPerson person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newPerson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conexion:= coneccionMysql()
	insertarPersonas, err:= conexion.Prepare("INSERT INTO PERSONA (DNI, NOMBRES, APELLIDOS, FECHANAC, EDAD, CIUDAD) VALUES (?,?,?,?,?,?)")

	if err!=nil {
		panic(err.Error())
	}
	_, err = insertarPersonas.Exec(newPerson.DNI, newPerson.NOMBRES, newPerson.APELLIDOS, newPerson.FECHANAC, newPerson.EDAD, newPerson.CIUDAD)

	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	fmt.Fprintln(w, "Persona insertada exitosamente.")
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

