package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"strconv"
	"strings"
	"html/template"
)

var count = 0

var htmlMain = []byte{}
var htmlUploaded = []byte{}

func home(w http.ResponseWriter, req *http.Request) {
	 if len(req.URL.Path) >= 3 && req.URL.Path[0:2] == "/f" {
		// Files are loaded in with this format <ID>-<KEY>.<EXTENSION>
		path := req.URL.Path[3:]
		extension := strings.Split(path, ".")
		if len(extension) != 2 {
			fmt.Fprintf(w, "Parsing error - no filename")
			return
		}

		ids := strings.Split(extension[0], "-")
		if len(ids) != 2 {
			fmt.Fprintf(w, "Parsing error - no ID/KEY")
			return
		}

		file := "file/" + ids[0] + "." + extension[1]

		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(w, "IO error - Bad ID", err)
			return
		}

		for i := 0; i <= len(bytes) - 1; i++ {
			code := byte(bytes[i]) ^ ids[1][i % len(ids[1])]
			w.Write([]byte{code})
		}
	} else if req.URL.Path != "/" {
		fmt.Fprintf(w, "404")
	} else {
		fmt.Fprintf(w, string(htmlMain))
	}
}

func upload(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(10 << 20)
	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error - no file parameter")
		return
	}

	if handler.Size >= 1000000 {
		fmt.Fprintf(w, "Error - Maximum file size is 5 megabytes")
		return
	}

	key := req.FormValue("key")
	if key == "" {
		fmt.Fprintf(w, "Error - no key parameter")
		return
	}

	if len(key) < 10 {
		fmt.Fprintf(w, "Error - key must be at least 10 characters")
		return
	}

	if len(key) >= 100 {
		fmt.Fprintf(w, "Error - key must be less than 100 characters")
		return
	}

	num := strconv.Itoa(count)
	count++

	content := make([]byte, handler.Size)
	_, err = file.Read(content)
	if err != nil {
		fmt.Fprintf(w, "File read error")
		return
	}

	extensions := strings.Split(handler.Filename, ".")
	extension := "txt"
	if len(extensions) == 2 {
		extension = extensions[1]
	}

	output, err := os.OpenFile("file/" + num + "." + extension, os.O_WRONLY | os.O_CREATE, 777)

	// XOR encode array into itself
	for i := 0; i <= len(content) - 1; i++ {
		content[i] = byte(content[i]) ^ key[i % len(key)]
	}

	output.Write(content)

	output.Close()
	file.Close()

	template, err := template.New("webpage").Parse(string(htmlUploaded))

	data := struct {
		Title string
		Number string
		Key string
		Extension string
	}{
		Title: handler.Filename,
		Number: num,
		Key: key,
		Extension: extension,
	}

	template.Execute(w, data)
}

func main() {
	// Set up file counter, from files in directory
	files, err := ioutil.ReadDir("file/")
	if err != nil {
		fmt.Println("Can't read file/")
		return
	}

	count = len(files)

	fmt.Println("Starting from file", count)

	htmlMain, err = ioutil.ReadFile("main.html")
	if err != nil {
		fmt.Println("Can't read main.html")
		return
	}

	htmlUploaded, err = ioutil.ReadFile("uploaded.html")
	if err != nil {
		fmt.Println("Can't read main.html")
		return
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)
	fmt.Println("Hosting on 8090")
	http.ListenAndServe(":8090", nil)
}
