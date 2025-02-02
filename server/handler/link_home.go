package handler

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"

	"github.com/bjdgyc/anylink/admin"
	"github.com/bjdgyc/anylink/dbdata"
)

func LinkHome(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.RemoteAddr)
	// hu, _ := httputil.DumpRequest(r, true)
	// fmt.Println("DumpHome: ", string(hu))
	w.Header().Set("Server", "AnyLinkOpenSource")
	connection := strings.ToLower(r.Header.Get("Connection"))
	userAgent := strings.ToLower(r.UserAgent())
	if connection == "close" && (strings.Contains(userAgent, "anyconnect") || strings.Contains(userAgent, "openconnect")) {
		w.Header().Set("Connection", "close")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	index := &dbdata.SettingOther{}
	if err := dbdata.SettingGet(index); err != nil {
		return
	}
	if index.Homeindex == "" {
		// index.Homeindex = "AnyLink 是一个企业级远程办公 SSL VPN 软件，可以支持多人同时在线使用。"
		filepath := "/var/www/anylink/welcome.html"
		content ,err :=ioutil.ReadFile(filepath)
		if err !=nil {
			http.Redirect(w, r, "https://henchat.net", 301)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, content)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, index.Homeindex)
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
