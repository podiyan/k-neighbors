package main

import (
	"net/http"
	"io"
	"fmt"
	"log"
	"os"
	"html/template"
	"github.com/ericlagergren/go-kml"
	"encoding/json"
	"io/ioutil"
	"encoding/xml"
	"persistence"
	"strconv"
)

var Log = log.New(os.Stdout, "proxi: ", log.Llongfile)

// ServiceStatus is a simple struct that stores details about the proxi service
type ServiceStatus struct {
	Status string
	Backend string
	PlacemarksProcessed int
	IndexSize int
}

var repo persistence.Repository
var proxiStatus = ServiceStatus{Status:"UnInitialized", Backend:"UnInitialized", PlacemarksProcessed:0}

func main() {

	bytes, err := ioutil.ReadFile("./config.json")

	if err != nil {
		Log.Fatalf("Unable to read config. : %#v", err)
	}

	var config ProxiConfig
	json.Unmarshal(bytes, &config)

	Log.Printf("Config : %#v", config)

	proxiStatus.Backend = config.PersistenceConfig.Type

	if proxiStatus.Backend == "memory" {
		Log.Printf("Using memory persistence")
		repo = persistence.NewMemoryStore(*config.PersistenceConfig)
	}

	startWebServer(config)
}


// For the functions to be fully testable, param extraction should happen in a central place.
// Functions will then accept a param array instead of http.Request

// health-check
func statusCheck(w http.ResponseWriter, r *http.Request) {
	Log.Printf("/status")
	proxiStatus.IndexSize = repo.GetIndexSize()
	status,_ := json.Marshal(proxiStatus)
	io.WriteString(w, string(status))
}

// Index KML file from a specific URL
func indexKML(w http.ResponseWriter, r *http.Request) {
	kmlPath := r.FormValue("path")

	if kmlPath == "" {
		Log.Printf("No path provided, returning http 400")
		reportBadRequest(w,"provide \"path\" for the KML file")
		return
	}

	// If the path is ok, return http 202 immediately.
	// Indexing will happen in a go routine below
	w.WriteHeader(http.StatusAccepted)

	go func(kmlPath string) {
		Log.Printf("Reading KML file %#v", kmlPath)

		resp, err := http.Get(kmlPath)
		if err != nil {
			// handle error
			Log.Fatalf("Fatal error : %#v", err)
		}

		proxiStatus.Status = "Indexing"
		defer resp.Body.Close()
		decoder := xml.NewDecoder(resp.Body)

		for {
			// Read tokens from the XML document in a stream.
			t, _ := decoder.Token()
			if t == nil {
				break
			}
			// Inspect the type of the token just read.
			switch se := t.(type) {
			case xml.StartElement:
				// If we just read a StartElement token
				// ...and its name is "page"
				if se.Name.Local == "Placemark" {
					var p kml.Placemark
					// decode a whole chunk of following XML into the
					// variable p which is a Page (se above)
					decoder.DecodeElement(&p, &se)
					err := repo.AddPlaceMark(p)
					if err != nil {
						Log.Printf("Repo returned error. halting indexing operation")
						proxiStatus.Status = "stopped indexing"
						return
					}
					proxiStatus.PlacemarksProcessed++
					Log.Printf("Placemark : %#v", p)
				}
			}
		}

	}(kmlPath)
}

// Find closest K matches for a given location
func findNearestK(w http.ResponseWriter, r *http.Request) {

	lat,e := strconv.ParseFloat(r.FormValue("lat"), 64)
	if e != nil {
		reportBadRequest(w,"required parameter \"lat\" invalid or not specified")
		return
	}

	lon,e := strconv.ParseFloat(r.FormValue("lon"), 64)
	if e != nil {
		reportBadRequest(w,"required parameter \"lon\" invalid or not specified")
		return
	}

	size,_ := strconv.ParseInt(r.FormValue("size"),10,32)

	format := r.FormValue("format")

	if format == "" {
		format = "json"
	}

	// Find the nearest one if size isn't specified
	if size <= 0 {
		size = 1
	}

	k := repo.FindKNearestPlaceMarks(lat, lon, int(size))

	if format == "kml" {
		xs,_ := xml.Marshal(resultsToKML(k))
		io.WriteString(w, string(xs))
		return
	}
	ks,_ := json.Marshal(k)
	io.WriteString(w, string(ks))
}

// display home page
func proxiHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("proxi.html")
	if err != nil {
		Log.Printf("Template is messed up: %v", err)
	}
	t.Execute(w, nil)
}



func startWebServer(config ProxiConfig) {

	// health-check
	http.HandleFunc("/proxi/status", statusCheck)

	// Index KML file from a specific URL
	http.HandleFunc("/proxi/indexKML", indexKML)

	// Find closest match(es) for a given location
	http.HandleFunc("/proxi/find", findNearestK)

	// root, display main message and load template
	http.HandleFunc("/", proxiHome)

	host := fmt.Sprintf(":%d", config.Port)
	Log.Printf("Web server listening at %v", host)
	Log.Printf("Web server died: %v", http.ListenAndServe(host, nil))
}



