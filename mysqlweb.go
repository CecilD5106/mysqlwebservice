package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	//"github.com/gin-gonic/gin"
	// _ "github.com/go-sql-driver/mysql"
)

// Person is a model of the person table
type Person struct {
	ID    string `json:"person_id"`
	FName string `json:"first_name"`
	LName string `json:"last_name"`
}

// Response is a list of person objects
type Response struct {
	People []Person `json:"result"`
}

var tmpl = template.Must(template.ParseGlob("form/*"))

// Index is the main primary website page
func Index(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("http://localhost:8000/getpeople")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	person := Person{}
	res := []Person{}
	for i := 0; i < len(responseObject.People); i++ {
		person.ID = responseObject.People[i].ID
		person.FName = responseObject.People[i].FName
		person.LName = responseObject.People[i].LName
		res = append(res, person)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
}

// Show activates a web page that show a single record
func Show(w http.ResponseWriter, r *http.Request) {
	nID := r.URL.Query().Get("id")
	sURL := "http://localhost:8000/getperson/" + nID
	response, err := http.Get(sURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	person := Person{}
	person.ID = responseObject.People[0].ID
	person.FName = responseObject.People[0].FName
	person.LName = responseObject.People[0].LName
	tmpl.ExecuteTemplate(w, "Show", person)
}

// New opens the web page to create a new person
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Edit opens the web page to edit a person record
func Edit(w http.ResponseWriter, r *http.Request) {
	nID := r.URL.Query().Get("id")
	sURL := "http://localhost:8000/getperson/" + nID
	response, err := http.Get(sURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	person := Person{}
	person.ID = responseObject.People[0].ID
	person.FName = responseObject.People[0].FName
	person.LName = responseObject.People[0].LName

	tmpl.ExecuteTemplate(w, "Edit", person)
}

// Insert is the code to add a new person to the database
func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		jsonData := map[string]string{"first_name": fname, "last_name": lname}
		jsonValue, _ := json.Marshal(jsonData)
		res, err := http.Post("http://localhost:8000/createperson", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			panic(err.Error())
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(data))
		}
		log.Println("INSERT: First Name: " + fname + " | Last Name: " + lname)
	}
	http.Redirect(w, r, "/", 301)
}

// Update data from Edit web page
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		id := r.FormValue("uid")
		jsonData := map[string]string{"person_id": id, "first_name": fname, "last_name": lname}
		jsonValue, _ := json.Marshal(jsonData)
		res, err := http.Post("http://localhost:8000/updateperson", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			panic(err.Error())
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			fmt.Print(string(data))
		}
		log.Println("UPDATE: First Name: " + fname + " | Last Name: " + lname)
	}
	http.Redirect(w, r, "/", 301)
}

// Delete activates a Delete Person call
func Delete(w http.ResponseWriter, r *http.Request) {
	nID := r.URL.Query().Get("id")
	sURL := "http://localhost:8000/deleteperson/" + nID
	_, err := http.Get(sURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	log.Println("DELETE")
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
