package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Employee struct {
	ID           int    `json:"id"`
	Lastname     string `json:"lastname"`
	Firstname    string `json:"firstname"`
	Monday       string `json:"monday"`
	Tuesday      string `json:"tuesday"`
	Wednesday    string `json:"wednesday"`
	Thursday     string `json:"thursday"`
	Friday       string `json:"friday"`
	Saturday     string `json:"saturday"`
	OutOfOffice  bool   `json:"out_of_office"`
	Sick		 bool	`json:"sick"`
}

func initDB() {
	var err error
	// Replace with your MySQL credentials
	dsn := "root:42143280MM@tcp(127.0.0.1:3306)/mySchedule"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
}

func main() {
	// Initialize the database connection
	initDB()

	// Create a new router
	router := mux.NewRouter()

	// Register handlers
	router.HandleFunc("/employees", getEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", getEmployee).Methods("GET")
	router.HandleFunc("/employees", createEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	// Serve static files from the 'static' directory
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Start the server
	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// Get all employees
func getEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM employee")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.Lastname, &emp.Firstname, &emp.Monday, &emp.Tuesday, &emp.Wednesday, &emp.Thursday, &emp.Friday, &emp.Saturday, &emp.OutOfOffice, &emp.Sick)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

// Get an employee by ID
func getEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var emp Employee
	err := db.QueryRow("SELECT * FROM employee WHERE id = ?", id).Scan(&emp.ID, &emp.Lastname, &emp.Firstname, &emp.Monday, &emp.Tuesday, &emp.Wednesday, &emp.Thursday, &emp.Friday, &emp.Saturday, &emp.OutOfOffice, &emp.Sick)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

// Create a new employee
func createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO employee (lastname, firstname, monday, tuesday, wednesday, thursday, friday, saturday, out_of_office, sick) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		emp.Lastname, emp.Firstname, emp.Monday, emp.Tuesday, emp.Wednesday, emp.Thursday, emp.Friday, emp.Saturday, emp.OutOfOffice, emp.Sick)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update an employee
func updateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE employee SET lastname = ?, firstname = ?, monday = ?, tuesday = ?, wednesday = ?, thursday = ?, friday = ?, saturday = ?, out_of_office = ?, sick = ? WHERE id = ?",
		emp.Lastname, emp.Firstname, emp.Monday, emp.Tuesday, emp.Wednesday, emp.Thursday, emp.Friday, emp.Saturday, emp.OutOfOffice, emp.Sick, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete an employee
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM employee WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
