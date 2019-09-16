package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func uploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	r.ParseMultipartForm(1 << 40)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	f, err := os.Create("Downloads/" + handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	io.Copy(f, file)

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
