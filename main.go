package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("v3ry-t0p_s3cr3t"))

func getWorkingDirectory(session *sessions.Session) (string, error) {
	untyped, ok := session.Values["directory"]
	if !ok {
		default_directory, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return default_directory, nil;
	} else {
		directory, ok := untyped.(string)
		if !ok {
			return "", errors.New("Directory is not a string?")
		}
		return directory, nil;
	}
}

func handlePwdRequest(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	directory, err := getWorkingDirectory(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(directory))
}

func handleCdRequest(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	dir := r.URL.Query().Get("dir")
	if dir == "" {
        http.Error(w, "Missing <dir> param", 400)
        return
	}

	directory, err := getWorkingDirectory(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.Chdir(directory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	err = os.Chdir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	current_directory, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	session.Values["directory"] = current_directory
	session.Save(r, w)
	w.Write([]byte(current_directory))
}

func handleLsRequest(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	directory, err := getWorkingDirectory(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	json.NewEncoder(w).Encode(filenames)
}

func handleRequests() {
	http.HandleFunc("/ls", handleLsRequest)
	http.HandleFunc("/pwd", handlePwdRequest)
	http.HandleFunc("/cd", handleCdRequest)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}
