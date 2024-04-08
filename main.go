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

var db *sql.DB

type ResponseMessage struct {
	Message string `json:"message"`
}

func initDB() {
	var err error
	db, err = sql.Open("mysql", "jasper:admin123@tcp(db:3306)/test")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func getUsers() []*User {
	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err) // Sử dụng log.Fatal trong ví dụ này, nhưng bạn nên xử lý lỗi một cách tử tế hơn trong sản phẩm thực
	}
	defer results.Close()

	var users []*User
	for results.Next() {
		var u User

		err = results.Scan(&u.ID, &u.Name)
		if err != nil {
			panic(err.Error()) // proper error handling
		}

		users = append(users, &u)
	}

	return users
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)

	stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(newUser.Name)
	if err != nil {
		log.Fatal(err)
	}

	response := ResponseMessage{Message: "Added successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var updateUser User
	json.NewDecoder(r.Body).Decode(&updateUser)

	stmt, err := db.Prepare("UPDATE users SET name = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(updateUser.Name, id)
	if err != nil {
		log.Fatal(err)
	}

	response := ResponseMessage{Message: "Updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	response := ResponseMessage{Message: "Deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Fprintf(w, "User with ID = %s was deleted", params["id"])
}

func userPage(w http.ResponseWriter, r *http.Request) {
	users := getUsers()

	fmt.Println("Endpoint Hit: usersPage")
	json.NewEncoder(w).Encode(users)
}

func main() {

	initDB()

	router := mux.NewRouter()

	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/users", userPage).Methods("GET")
	router.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/", homePage).Methods("GET")

	log.Fatal(http.ListenAndServe(":8083", router))
}
