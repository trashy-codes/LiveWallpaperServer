package main

import (
	"LiveWallpaperServer/upupoo"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration

func getTags(writer http.ResponseWriter, request *http.Request) {
	tags, err := upupoo.GetTags()
	json, err := json.Marshal(tags)
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
}

func getWallpapers(writer http.ResponseWriter, request *http.Request) {
	upupoo.GetWallpapers()
	writer.Header().Set("Content-Type", "application/json")
}

func getSorts(writer http.ResponseWriter, request *http.Request) {
	sorts, err := upupoo.GetSorts()
	json, err := json.Marshal(sorts)
	if err != nil {
		log.Fatal(err)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

func main() {
	loadConfig()
	mux := http.NewServeMux()
	mux.HandleFunc("/tags", getTags)
	mux.HandleFunc("/sorts", getSorts)
	mux.HandleFunc("/wallpapers", getWallpapers)
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
