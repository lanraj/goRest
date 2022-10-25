package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "strconv"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
)

// Student struct (Model)
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"Name"`
	Grade int    `json:"Grade"`
}

// MSSQL DB configuration
var db *sql.DB
var servers = "trialapis.database.windows.net"
var ports = 1433
var users = "maniraj"
var passwords = "Raj#man#7548"
var databases = "TrialDB"

func GetMySQLDB() *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		servers, users, passwords, ports, databases)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)

	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	return db
}

// Init Students var as a slice Student struct
var Students []Student

// Get all Students
func getStudents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	studentList := []Student{}
	s := Student{}
	rows, err := db.Query("select * from TestSchema.student")
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		for rows.Next() {
			rows.Scan(&s.ID, &s.Name, &s.Grade)
			studentList = append(studentList, s)
		}
		json.NewEncoder(w).Encode(studentList)

	}
}

func createStudents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	s := Student{}
	json.NewDecoder(r.Body).Decode(&s)
	tsql := fmt.Sprintf("insert into TestSchema.student (sid, name, Grade) values('%d','%s','%s')", s.ID, s.Name, s.Grade)
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.LastInsertId()
		if err != nil {
			json.NewEncoder(w).Encode(s)
		} else {
			json.NewEncoder(w).Encode(s)

		}
	}

}

func updateStudents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	s := Student{}
	json.NewDecoder(r.Body).Decode(&s)
	// vars := mux.Vars(r)
	tsql := fmt.Sprintf("update TestSchema.student set name='%s', Grade='%s' where id='%d'", s.Name, s.Grade, s.ID)
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{ error: record not updated }")
		} else {
			json.NewEncoder(w).Encode(s)

		}
	}
}

func deleteStudents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	vars := mux.Vars(r)
	tsql := fmt.Sprintf("delete from TestSchema.student where id='%d';", vars["ID"])
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Fprintf(w, ""+err.Error())
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("{ error: record not Deleted }")
		} else {
			json.NewEncoder(w).Encode("{ result: record Sucessfully Deleted }")

		}
	}

}

func main() {
	db = GetMySQLDB()
	fmt.Print("testing")
	// db.Close()
	// fmt.Print("testing")
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students", createStudents).Methods("POST")
	r.HandleFunc("/students/{id}", updateStudents).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudents).Methods("DELETE")
	fmt.Println("server started")
	http.ListenAndServe(":8000", r)
}
