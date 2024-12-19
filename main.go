package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
)

type LogData struct {
	Req struct {
		URL        string            `json:"url"`
		QSParams   string            `json:"qs_params"`
		Headers    map[string]string `json:"headers"`
		ReqBodyLen int               `json:"req_body_len"`
	} `json:"req"`
	Rsp struct {
		StatusClass string `json:"status_class"`
		RspBodyLen  int    `json:"rsp_body_len'`
	} `json:"rsp"`
}

func writeLog(logData LogData) {
	file, err := json.MarshalIndent(logData, "", "	")
	if err != nil {
		fmt.Println("Error marshaling JSON: ", err)
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	os.MkdirAll("logs", os.ModePerm)

	filename := fmt.Sprintf("logs/log - %s.json", timestamp)
	if err := os.WriteFile(filename, file, 0644); err != nil {
		fmt.Println("error writing file: ", err)
	}
}

func logging(originFunction func(w http.ResponseWriter, r *http.Request, claims *Claims)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var logData LogData

		// collect Req data
		logData.Req.URL = r.URL.String()
		logData.Req.QSParams = r.URL.RawQuery
		// collect headers
		headers := make(map[string]string)
		for key, value := range r.Header {
			headers[key] = strings.Join(value, ", ")
		}
		logData.Req.Headers = headers
		// find request body len
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			logData.Req.ReqBodyLen = len(bodyBytes)
			r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		}

		// run the original function and record the response
		rec := httptest.NewRecorder()
		claims := &Claims{}
		originFunction(rec, r, claims)

		logData.Rsp.RspBodyLen = rec.Body.Len()
		logData.Rsp.StatusClass = fmt.Sprintf("%dxx", rec.Result().StatusCode/100)

		for key, value := range rec.Header() {
			w.Header()[key] = value
		}

		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)

		// Write and save the json file
		writeLog(logData)
	}
}

func main() {
	http.HandleFunc("/register", logging(Register)) // role: user / admin
	http.HandleFunc("/login", logging(Login))
	http.HandleFunc("/accounts", logging(AccountsHandler))
	http.HandleFunc("/balance", logging(BalanceHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
