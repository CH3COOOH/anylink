package handler

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"github.com/bjdgyc/anylink/admin"
)

func LinkHome(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.RemoteAddr)
	// hu, _ := httputil.DumpRequest(r, true)
	// fmt.Println("DumpHome: ", string(hu))

	connection := strings.ToLower(r.Header.Get("Connection"))
	userAgent := strings.ToLower(r.UserAgent())
	if connection == "close" && (strings.Contains(userAgent, "anyconnect") || strings.Contains(userAgent, "openconnect")) {
		w.Header().Set("Connection", "close")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filepath := "./index.htm"
	content ,err :=ioutil.ReadFile(filepath)
	if err !=nil {
		http.Redirect(w, r, "https://henchat.net", 301)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, content)
	}

}

func LinkOtpQr(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	idS := r.FormValue("id")
	jwtToken := r.FormValue("jwt")
	data, err := admin.GetJwtData(jwtToken)
	if err != nil || idS != fmt.Sprint(data["id"]) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	admin.UserOtpQr(w, r)
}
