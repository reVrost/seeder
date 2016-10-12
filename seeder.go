package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocraft/dbr"
	//"github.com/revrost/fake"
	_ "github.com/lib/pq"
)

func helpString() string {
	return "Usage: ./seeder <username> <password> <db> <definition.json>"
}

var username string
var password string
var dbName string
var definitionFile string
var db *dbr.Connection

func main() {

	if len(os.Args) < 4 {
		fmt.Println(helpString())
		return
	}

	username = os.Args[1]
	password = os.Args[2]
	dbName = os.Args[3]
	definitionFile = os.Args[4]

	seedDefinition := make(map[string]map[string]string)

	if _, err := os.Stat(definitionFile); os.IsNotExist(err) {
		fmt.Println("Unable to find definition file. ", err)
		return
	}

	definition, err := os.Open(definitionFile)
	if err != nil {
		fmt.Println("Failed to open definition file.", err)
		return
	}

	jsonReader := json.NewDecoder(definition)
	if err := jsonReader.Decode(&seedDefinition); err != nil {
		fmt.Errorf("Unable to parse config file %#v\n", err)
	}

	//fmt.Println(seedDefinition)
	for k, _ := range seedDefinition {
		fmt.Printf("Definition for table [%s] found!\n", k)
	}

	_, err = dbr.Open("postgres", "user="+username+" password"+password+" dbname="+dbName+" sslmode=disable", nil)
	if err != nil {
		fmt.Println("Unable to establish connection to the database. goodbye, Error:", err)
		return
	}

	fmt.Println("------- Seeder -------")
	fmt.Printf("Successfully logged in as %s to %s\n", username, dbName)
	return

}
