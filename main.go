/*
Creator James Marshall
Date 	2022-02-16


This article was very helpful in creating this API
https://tutorialedge.net/golang/creating-restful-api-with-golang/
*/
package main

import (
	"encoding/json" // json ser/des & tagging
	"fmt"           // console printing
	"io/ioutil"     // reading data stream
	"log"           // capture of fatal messages
	"net/http"      // http response, request, and service

	"github.com/google/uuid" // lib for easy uuid
	"github.com/gorilla/mux" // lib for easy routing
)

//foo structure
type foo struct {
	// Note: Public is PascalCase and Private is camelCase. json.Unmarshal will only ser/des public members

	Name string `json:"name"`
	Id   string `json:"id"`
}

// In memory collection of foos
var Foos []foo

// Returns the index of foo in Foos, -1 being not found
func indexOfFooInFoos(key string) int {
	var retIndex int = -1

	//Loop through collection till we find matching id
	for index, fooInCollection := range Foos {
		if fooInCollection.Id == key {
			retIndex = index
			break
		}
	}
	return retIndex
}

// Handles /foo POST by creating a new foo and adding it to the foo collection.
// If the id of the foo is provided, the id is replaced with a new uuid is generated.
func createFoo(w http.ResponseWriter, r *http.Request) {
	// Requirement: 5. It should support a POST endpoint (/foo) that accepts a json foo object and responds with a 200 response code.
	//                 The value of the id field should be added by this endpoint using a generated UUID.
	var newFoo foo

	// Deserialize foo from request
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newFoo)

	// Give the foo.Id a new UUID
	newFoo.Id = uuid.New().String()
	Foos = append(Foos, newFoo)

	// Serialize foo into the response
	json.NewEncoder(w).Encode(&newFoo)
	fmt.Println("Foo created: ", newFoo)
}

// Handles /foo/{id} GET and Returns foo (status 200) or empty (status 404)
func returnFoo(w http.ResponseWriter, r *http.Request) {
	// Requirement 6. It should support a GET endpoint (foo/{id}) that responds with a 200 response code if the record is found, or a 404 response code if not found.

	// Retrieve id from url
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Getting foo with id: ", id)

	fooIndex := indexOfFooInFoos(id)

	if fooIndex > -1 {
		// Return found foo
		json.NewEncoder(w).Encode(Foos[fooIndex])
		fmt.Println("Returning foo with id: ", id)
	} else {
		// Foo not found, return 404
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("foo not found with id: ", id)
	}
}

// Handles /foo/{id} DELETE returns status 200 if successful or 404 if not found
func deleteFoo(w http.ResponseWriter, r *http.Request) {
	// Requirement 7. It should support a DELETE endpoint (foo/{id}) that responds with a 204 response code if the record is found, or a 404 response code if not found.

	// Retrieve id from url
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println("Attempting to delete foo with id: ", id)

	fooIndex := indexOfFooInFoos(id)

	if fooIndex > -1 {
		// Remove foo by appending foos before and after
		Foos = append(Foos[:fooIndex], Foos[fooIndex+1:]...)
		fmt.Println("foo removed with id: ", id)
	} else {
		// Foo not found, return 404
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("foo not found with id: ", id)
	}
}

// Starts the server and handles requests
func handleRequests() {
	fmt.Println("Starting server on port 8080")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/foo", createFoo).Methods("POST")
	router.HandleFunc("/foo/{id}", deleteFoo).Methods("DELETE")
	router.HandleFunc("/foo/{id}", returnFoo)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("Starting Rakuten Coding Exercise")
	Foos = []foo{}
	handleRequests()
}
