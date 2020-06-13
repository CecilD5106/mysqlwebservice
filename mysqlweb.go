package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	ID    string
	FName string
	LName string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "cgdavis"
	dbPass := "DzftXvz$eR7VpY^h"
	dbServer := "tcp(172.18.105.227:3306)"
	dbName := "people"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbServer+"/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM person")
	if err != nil {
		panic(err.Error())
	}

	indiv := Person{}
	res := []Person{}
	for selDB.Next() {
		var id, fname, lname string
		err = selDB.Scan(&id, &fname, &lname)
		if err != nil {
			panic(err.Error())
		}
		indiv.ID = id
		indiv.FName = fname
		indiv.LName = lname
		res = append(res, indiv)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM person WHERE person_id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	indiv := Person{}
	for selDB.Next() {
		var id, fname, lname string
		err = selDB.Scan(&id, &fname, &lname)
		if err != nil {
			panic(err.Error())
		}
		indiv.ID = id
		indiv.FName = fname
		indiv.LName = lname
	}
	tmpl.ExecuteTemplate(w, "Show", indiv)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM person WHERE person_id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	indiv := Person{}
	for selDB.Next() {
		var id, fname, lname string
		err = selDB.Scan(&id, &fname, &lname)
		if err != nil {
			panic(err.Error())
		}
		indiv.ID = id
		indiv.FName = fname
		indiv.LName = lname
	}
	tmpl.ExecuteTemplate(w, "Edit", indiv)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		insForm, err := db.Prepare("INSERT INTO person(first_name, last_name) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fname, lname)
		log.Println("INSERT: First Name: " + fname + " | Last Name: " + lname)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		fname := r.FormValue("fname")
		lname := r.FormValue("lname")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE person SET first_name=?, last_name=? WHERE person_id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fname, lname, id)
		log.Println("UPDATE: First Name: " + fname + " | Last Name: " + lname)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	indiv := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM person WHERE person_id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(indiv)
	log.Println("DELETE")
	defer db.Close()
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
