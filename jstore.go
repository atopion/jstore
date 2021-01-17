package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"

	uuid "github.com/nu7hatch/gouuid"
)

// Flags
var url string
var port string
var folder string

func init() {
	flag.StringVar(&url, "u", "http://localhost:8080", "The url of the current machine. To be displayed in the returned identifier.")
	flag.StringVar(&port, "p", "8080", "The port this application should run on.")
	flag.StringVar(&folder, "f", "/store", "Folder to store the files in. (absolut)")
	flag.Parse()
}

/*Document  :
A type that represents a stored Document based on title and body (content)
*/
type Document struct {
	Title string
	Body  []byte
}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9\\-]+)(.json)?$")

func (d *Document) save() error {
	filename := d.Title + ".json"
	return ioutil.WriteFile(filepath.Join(folder, filename), d.Body, 0600)
}

func loadDocument(title string) (*Document, error) {
	filename := title + ".json"
	body, err := ioutil.ReadFile(filepath.Join(folder, filename))
	if err != nil {
		return nil, err
	}
	return &Document{Title: title, Body: body}, nil
}

func getTitleFromPath(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Document Identifier")
	}
	return m[1], nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitleFromPath(w, r)
	if err != nil {
		return
	}
	d, err := loadDocument(title)
	if err != nil {
		log.Printf("Cannot find file: %s\n", err.Error())
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(d.Body)
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read request's body: %s\n", err.Error())
		http.Error(w, "Cannot read the request", http.StatusBadRequest)
		return
	}
	u4, err := uuid.NewV4()
	if err != nil {
		log.Printf("Problems while generating a new UUID: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	title := u4.String()
	d := &Document{Title: title, Body: []byte(body)}
	err = d.save()
	if err != nil {
		log.Printf("Cannot save the file: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("location", path.Join(url, title+".json"))
	w.Write([]byte(path.Join(url, title+".json")))
}

func modifyHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitleFromPath(w, r)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read request's body: %s\n", err.Error())
		http.Error(w, "Cannot read the request", http.StatusBadRequest)
		return
	}
	d, err := loadDocument(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	d.Body = body
	err = d.save()
	if err != nil {
		log.Printf("Cannot save the file: %s\n", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(path.Join(url, title+".json")))
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func corsHeaderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
}

func handlersSwitch(w http.ResponseWriter, r *http.Request) {
	corsHeaderHandler(w, r)

	switch r.Method {
	case http.MethodGet:
		viewHandler(w, r)
	case http.MethodPost:
		storeHandler(w, r)
	case http.MethodPut:
		modifyHandler(w, r)
	case http.MethodOptions:
		optionsHandler(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusNotImplemented)
	}
}

func main() {
	log.Printf("Start server on port %s ...\n", port)
	http.HandleFunc("/", handlersSwitch)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
