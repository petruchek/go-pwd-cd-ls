package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getParam(r *http.Request, param string) (string, error) {
	values, ok := r.URL.Query()[param]

	if !ok || len(values[0]) < 1 {
		return "", errors.New(fmt.Sprintf("Url Param '%s' is missing", param))
	}

	return values[0], nil
}

func pwd(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(dir)
}

func cd(w http.ResponseWriter, r *http.Request) {
	dir, err := getParam(r, "dir")
	if err != nil {
		log.Fatal(err)
		return
	}

	os.Chdir(dir)
	newDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(newDir)
}

func ls(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
		return
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
		fmt.Println(file.Name())
	}
	json.NewEncoder(w).Encode(filenames)
}

func handleRequests() {
	http.HandleFunc("/ls", ls)
	http.HandleFunc("/pwd", pwd)
	http.HandleFunc("/cd", cd)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
