/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var (
	getAction map[string]interface{}
)

func PublicCfgVersion(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("[Success] config version  === " + config.Version))
}

func getRouteName(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}

func ValidatePostData(rw http.ResponseWriter, req *http.Request) ([]byte, bool) {
	header := req.Header.Get("Content-Type")
	if header != "application/json" {
		rsp, _ := json.Marshal(&ResponseData{2, "request must follow application/json"})
		httprsp(rw, rsp)
		return nil, false
	}
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	return body, true
}

func httprsp(rw http.ResponseWriter, rsp []byte) {
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rsp)
}

func actionMap(cloud string, service string, action string) (string, bool) {
	d, err := ioutil.ReadFile("ActionMap.json")
	if err != nil {
		raven.CaptureError(err, nil)
	}
	json.Unmarshal(d, &getAction)

	getCloud, ok := getAction[action].(map[string]interface{})
	if ok {
		getService, ok := getCloud[cloud].(map[string]interface{})
		if ok {
			value, ok := getService[service].(string)
			if ok {
				return value, true
			}
			return "", false
		}
		return "", false
	}
	return "", false
}

type ResponseData struct {
	Code int8        `json:"code"`
	Data interface{} `json:"data"`
}

func ErrorMessage(c int8) string {
	msg := ""
	switch c {
	// 1 - 10 Receive an incorrect request
	case 1:
		msg = "[Failed] AAA Failed"
	case 2:
		msg = "[Failed] Params missing or incorrect"
	case 3:
		msg = "[Failed] Request Timeout"
	case 4:
		msg = "[Failed] Request Forbidden"
	case 5:
		msg = "[Failed] Invalid Public Token"
	case 10:
		msg = "[Failed] Do not hack the system"
	// 100 - 120 System Error
	// 120 - 127 Middleware Error
	case 125:
		msg = "[Error] MQ crash"
	case 126:
		msg = "[Error] DB crash"
	case 127:
		msg = "[Error] RPC crash"
	default:
		msg = "[Error] Unknown problem"
	}
	return msg
}

func ResponseErrmsg(c int8) *ResponseData {
	return &ResponseData{c, ErrorMessage(c)}
}
