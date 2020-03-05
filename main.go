// Copyright 2020 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vinymeuh/gnucash-graphql/gnucash"
	"github.com/vinymeuh/gnucash-graphql/graphql"
)

var (
	buildVersion string
	buildDate    string
	db           *gnucash.Database
	appStatus    appStatusType
)

type logWriter struct {
	io.Writer
}

func (w logWriter) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format("2006-01-02T15:04:05-07:00 ")), b...))
}

type appStatusType struct {
	Version string
	Build   string
	GnuCash gnucash.DBStatistics
}

func rootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/" {
			resp, _ := json.Marshal(appStatus)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", resp)
			return
		}

		log.Printf("%s %s 404 Not Found", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
	})
}

func main() {
	var err error

	log.SetFlags(log.Lshortfile)
	log.SetOutput(logWriter{os.Stdout})

	// load configuration file
	configFile := os.Getenv("GNCGQL_CONF")
	if configFile == "" {
		configFile = "config.yml"
	}
	log.Printf("Configuration file is '%s'", configFile)

	config := NewConfig()
	err = config.LoadFromFile(configFile)
	if err != nil {
		log.Fatalf("error loading configuration file - %s", err)
	}

	// load GnuCash database
	log.Printf("GnuCash data loaded from file '%s'", config.Gnucash.File)
	db, err = gnucash.LoadFile(config.Gnucash.File)
	if err != nil {
		log.Fatalf("error loading gnucash file - %s", err)
	}

	// initialize appStatus
	appStatus.Version = buildVersion
	appStatus.Build = buildDate
	appStatus.GnuCash = db.Statistics

	// start HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", rootHandler())
	mux.Handle("/graphql", graphql.NewHandler(db))

	addr := fmt.Sprintf("%s:%s", config.Listen.Host, config.Listen.Port)
	log.Printf("Starting HTTP server on http://%s", addr)
	http.ListenAndServe(addr, mux)
}
