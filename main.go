package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"go.uber.org/zap"
)

var slogger *zap.SugaredLogger
func main() {
	InitLogger()
	defer slogger.Sync()

	slogger.Infof("Serving and listening on port 8080")
	http.HandleFunc("/", nq) 
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func nq(w http.ResponseWriter, r *http.Request) {
  nqresponse, err := http.Get("http://fffff-api/nq")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  defer nqresponse.Body.Close()

  hqresponse, err := http.Get("http://fffff-api/hq")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  defer hqresponse.Body.Close()

  var data interface{}
  if err := json.NewDecoder(nqresponse.Body).Decode(&data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  if err := json.NewDecoder(hqresponse.Body).Decode(&data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  tmpl, err := template.ParseFiles("template.html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  if err := tmpl.Execute(w, data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	slogger = logger.Sugar()
}
