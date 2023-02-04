package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"go.uber.org/zap"
)

type AllData struct {
  NQData interface{}
  HQData interface{}
}

var slogger *zap.SugaredLogger
func main() {
	InitLogger()
	defer slogger.Sync()

	slogger.Infof("Serving and listening on port 8080")
	http.HandleFunc("/", root) 
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	slogger.Infof("Retrieving nq shopping list from api")
  nqresponse, err := http.Get("http://fffff-api/nq")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  defer nqresponse.Body.Close()

	slogger.Infof("Retrieving hq shopping list from api")
  hqresponse, err := http.Get("http://fffff-api/hq")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  defer hqresponse.Body.Close()

  var allData AllData
	slogger.Infof("Decoding nq data from json")
  if err := json.NewDecoder(nqresponse.Body).Decode(&allData.NQData); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  slogger.Infof("After decoding nq data: %s", allData.NQData)
	slogger.Infof("Decoding hq data from json")
  if err := json.NewDecoder(hqresponse.Body).Decode(&allData.HQData); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  slogger.Infof("After decoding hq data: %s", allData.HQData)

  slogger.Infof("Retrieving html template")
  tmpl, err := template.ParseFiles("template.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  slogger.Infof("Rendering html template")
  if err := tmpl.Execute(w, allData); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	slogger = logger.Sugar()
}
