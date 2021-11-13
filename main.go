package main

import (
	"fmt"
	"net/http"
	"strconv"
	"os"
	"io/ioutil"
	"strings"
)

func home(w http.ResponseWriter, req *http.Request) {
	 if len(req.URL.Path) >= 3 && req.URL.Path[0:2] == "/f" {
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
		fmt.Fprintf(w, 
`<html>
<head>
	<title>Image host</title>
	<style>
	input {
		display: block;
	}
	</style>
</head>
<body>
	<p>Upload a file:</p>
	<form action="upload" method="post" enctype="multipart/form-data">
		<input type="file" name="file" id="file">
		<input type="text" name="key" id="key" value="asd123">
		<input type="submit" value="Upload" name="submit">
	</form>
</body>
</html>`)
	}
}

// TODO: more efficient solution is to load size of dir on startup
// and count from there in memory
func count() string {
	data, _ := ioutil.ReadFile("file/counter")
	text := strings.Split(string(data), "\n")
	num, _ := strconv.Atoi(text[0])

	num++

	c, _ := os.OpenFile("file/counter", os.O_WRONLY, os.ModeAppend)
	count2, _ := c.WriteString(strconv.Itoa(num) + "\n")
	count2 = count2

	c.Close()
	return strconv.Itoa(num)
}

func upload(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(10 << 20)
	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error - no file parameter")
		return
	}

	key := req.FormValue("key")
	if key == "" {
		fmt.Fprintf(w, "Error - no key parameter")
		return
	}

	num := count()

	content := make([]byte, handler.Size)
	_, err = file.Read(content)
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}

	extensions := strings.Split(handler.Filename, ".")
	extension := "txt"
	if len(extensions) == 2 {
		extension = extensions[1]
	}

	output, err := os.OpenFile("file/" + num + "." + extension, os.O_WRONLY | os.O_CREATE, 777)

	// XOR encode
	for i := 0; i <= len(content) - 1; i++ {
		code := byte(content[i]) ^ key[i % len(key)]
		output.Write([]byte{code})
	}

	output.Close()
	file.Close()

	fmt.Fprintf(w, 
	`<p>
	` + handler.Filename + `
	has been uploaded with encryption key "` + key + `"<br>
	<a href="f/` + num + `-` + key + `.` + extension + `">f/` + num + `-` + key + `.` + extension + `</a>
	</p>`)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)
	fmt.Println("Hosting on 8090")
	http.ListenAndServe(":8090", nil)
}
