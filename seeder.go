package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
	"github.com/revrost/fake"
)

func helpString() string {
	return "Usage: ./seeder <username> <password> <db> <definition.json>"
}

var username string
var password string
var dbName string
var definitionFile string
var db *dbr.Connection
var seedDefinition map[string]map[string]string

func isInitialiseOk() bool {
	if len(os.Args) < 4 {
		fmt.Println(helpString())
		return false
	}
	username = os.Args[1]
	password = os.Args[2]
	dbName = os.Args[3]
	definitionFile = os.Args[4]

	// Check definition file is ok!
	if _, err := os.Stat(definitionFile); os.IsNotExist(err) {
		fmt.Println("Unable to find definition file. ", err)
		return false
	}
	definition, err := os.Open(definitionFile)
	if err != nil {
		fmt.Println("Failed to open definition file.", err)
		return false
	}
	jsonReader := json.NewDecoder(definition)
	seedDefinition = make(map[string]map[string]string)
	if err := jsonReader.Decode(&seedDefinition); err != nil {
		fmt.Println("Unable to parse config file %#v\n", err)
		return false
	}

	// Check db connection is ok
	db, err = dbr.Open("postgres", "user="+username+" password="+password+" dbname="+dbName+" sslmode=disable", nil)
	if err != nil {
		fmt.Println("Unable to establish connection to the database. goodbye, Error:", err)
		return false
	}
	return true
}

func giveRandomValue(valueType string) interface{} {
	var val interface{}
	switch valueType {
	case "Number.Tenth":
		val = fake.Number.Number(9)
	case "Number.Hundredth":
		val = fake.Number.Number(99)
	case "String":
		val = "F"
	case "Pharmacy.FullDrug":
		val = fake.Pharmacy.FullDrug()
	}
	return val
}

func generateRandomValue(valueTypes []string) []interface{} {
	values := make([]interface{}, len(valueTypes))
	for i := 0; i < len(valueTypes); i++ {
		values[i] = giveRandomValue(valueTypes[i])
	}
	return values
}

func seedTable(table string, count int) error {
	fmt.Printf("You've chosen %s, with %d entries\n", table, count)
	sess := db.NewSession(nil)
	schema := seedDefinition[table]
	var columns []string
	var valueTypes []string
	for k, v := range schema {
		columns = append(columns, k)
		valueTypes = append(valueTypes, v)
	}

	fmt.Println(columns)
	for i := 0; i < count; i++ {
		var values []interface{}
		// Need to generate the values
		values = generateRandomValue(valueTypes)
		fmt.Printf("%v\n", values)
		query := sess.InsertInto(table).Columns(columns...).Values(values...)
		fmt.Println(query.ToSql())
		test, err := query.Exec()

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(test)
	}
	return nil
}

func main() {
	if !isInitialiseOk() {
		return
	}

	var tables []string
	for k, _ := range seedDefinition {
		fmt.Printf("Definition for table [%s] found!\n", k)
		tables = append(tables, k)
	}

	fmt.Println("------- Seeder -------")
	fmt.Printf("Successfully logged in as %s to %s\n", username, dbName)
	for i, k := range tables {
		fmt.Printf("[%d] %s\n", i, k)
	}
	var index int
	fmt.Println("Which table would you like to seed?")
	_, err := fmt.Scanf("%d", &index)
	if err != nil {
		fmt.Println("Invalid integer input ", err)
		return
	}
	chosenTable := tables[index]
	fmt.Println("How many entries?")
	_, err = fmt.Scanf("%d", &index)
	if err != nil {
		fmt.Println("Invalid integer input ", err)
		return
	}
	entriesNo := index

	seedTable(chosenTable, entriesNo)
	return
}
