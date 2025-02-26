package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {

	// client := init_aws_client()
	// bucketName := os.Getenv("BUCKET_NAME")
	// folderPath := os.Getenv("FOLDER_PATH")

	// Calculate folder size
	// totalSize, err := getFolderSize(client, bucketName, folderPath)
	// if err != nil {
	// 	log.Fatalf("Error calculating folder size: %v", err)
	// }

	// fmt.Printf("Total size of folder '%s' in bucket '%s': %d bytes\n", folderPath, bucketName, totalSize)

	// only_2024, err := getListClientsFiles(client, bucketName, folderPath)
	// if err != nil {
	// 	log.Fatalf("Error calculating folder size: %v", err)
	// }

	// fmt.Printf("ONLY 2024 LIST CLIENTS ENDPOINT :) '%s' in bucket '%s': %d bytes\n", folderPath, bucketName, only_2024)
	read_json()
}

func read_json() {
	// Open our jsonFile
	jsonFile, err := os.Open("list_clients--2025-02-07--05-51-42.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened json")
	// // defer the closing of our jsonFile so that we can parse it later on
	// defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// fmt.Println(byteValue)

	// entries := Entries{}
	var entries []Entry

	json.Unmarshal(byteValue, &entries)

	fmt.Println(len(entries))

	for i := 0; i < len(entries); i++ {
		fmt.Println("Entry: " + entries[i].Mac)
		// fmt.Println("User Type: " + users.Users[i].Type)
		// fmt.Println("User Age: " + strconv.Itoa(users.Users[i].Age))
		// fmt.Println("User Name: " + users.Users[i].Name)
		// fmt.Println("Facebook Url: " + users.Users[i].Social.Facebook)
	}

	fmt.Println("all done")

}

// type Entries struct {
// 	Entries []Entry `json:"Entries"`
// }

type Entry struct {
	Site_id string `json:"Site_id"`
	Mac     string `json:"mac"`
}
