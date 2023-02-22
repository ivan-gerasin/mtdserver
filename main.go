package main

import (
	"encoding/json"
	"fmt"
	mtdCore "github.com/ivan-gerasin/mtdcore"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Item struct {
	Item string `json:"item"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Match root")
		if request.URL.Path != "/" {
			return
		}
		writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		writer.WriteHeader(200)
		writer.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Match /add")
		writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		if request.Method != "POST" {
			writer.Write([]byte(`{"status":"error", "message": "Unsupported method"}`))
			writer.WriteHeader(405)
			return
		}
		body := make([]byte, request.ContentLength)
		request.Body.Read(body)

		item := Item{}
		err := json.Unmarshal(body, &item)
		if err != nil {
			writer.Write([]byte(`{"status":"error", "message": "Internal server error"}`))
			writer.WriteHeader(500)
			return
		}
		mtdCore.AddItem(item.Item, 0)

		writer.WriteHeader(201)
	})

	mux.HandleFunc("/list", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Match /list")
		writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		if request.Method != "GET" {
			writer.WriteHeader(405)
			writer.Write([]byte(`{"status":"error", "message": "Unsupported method"}`))
			return
		}
		todoList := mtdCore.List()
		responseData, err := json.Marshal(todoList)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(`{"status":"error", "message": "Internal server error"}`))
			return
		}
		writer.WriteHeader(200)
		writer.Write(responseData)
	})

	mux.HandleFunc("/done/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Match /done")
		writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		if request.Method != "POST" {
			writer.WriteHeader(405)
			writer.Write([]byte(`{"status":"error", "message": "Unsupported method"}`))
			return
		}
		path := request.URL.Path
		idRegEx := regexp.MustCompile(`/done/(?P<id>\d+)`)
		result := idRegEx.FindStringSubmatch(path)
		if result == nil {
			writer.WriteHeader(404)
			writer.Write([]byte(`{"status":"error", "message": "No such item"}`))
			return
		}
		id, err := strconv.Atoi(result[1])
		if err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte(`{"status":"error", "message": "Invalid item number"}`))
			return
		}
		fmt.Println(id)
		mtdCore.Done(id)
		writer.WriteHeader(200)
		writer.Write([]byte(`{"status":"success"}`))
	})

	server := &http.Server{
		Addr:         ":8000",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServe()
}
