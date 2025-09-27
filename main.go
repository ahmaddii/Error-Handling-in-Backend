package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Custom error type
type UserNotFound struct {
	Id string
}

func (e UserNotFound) Error() string {
	return fmt.Sprintf("User of Specific Id %s not Found", e.Id)
}

// Fake DB function
func getUserDatafromDB(id string) (string, error) {
	if id != "1" { // Only user with ID=1 exists
		return "", UserNotFound{Id: id}
	}
	return "User1 Data", nil
}

// Handler with error handling
func getUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserDatafromDB("123") // pretend this is DB call

	if err != nil {
		if _, ok := err.(UserNotFound); ok {
			http.Error(w, err.Error(), http.StatusNotFound) // 404
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user) // success
}

// Panic recovery example
func GetUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	var arr []string
	fmt.Println(arr[1]) // panic: index out of range
}

func main() {
	http.HandleFunc("/getuser", getUser)
	http.HandleFunc("/panic", GetUser)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
