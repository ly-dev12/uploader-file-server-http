package main

import (
	"io/ioutil"

	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	file, handle, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	defer file.Close()
	mimetype := handle.Header.Get("Content-type")
	switch mimetype {
	case "image/jpeg":
		saveFile(w, file, handle)
	case "image/png":
		saveFile(w, file, handle)
	default:
		jsonResponse(w, http.StatusCreated, "The format file is not valid")
	}

}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = ioutil.WriteFile("./"+handle.Filename, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	jsonResponse(w, http.StatusCreated, "File uploaded successfully")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, message)
}

func main() {
	http.HandleFunc("/upload", uploadFile)
	log.Println("Running")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
