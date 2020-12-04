package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LoginConfirm struct {
	Login_confirm struct {
		Login_locked  string
		Token        string
		Login_confirm string
	}
}

type SysInfo struct {
	Sysinfo struct {
		Model       string
		Uptime      string
		Hw_version  string
		Serial_num  string
		Date_info   string
		Fw_version  string
		Linerate_us string
		Linerate_ds string
		Lanip       string
		Lanmac      string
		Wangw       string
		Wanmac      string
		Wan_model   string
		Wandns1     string
		Wandns2     string
		Token       string
		Sysinfo     string
	}
}

func getTimeStamp() string{
	t := time.Now()
	return strconv.FormatInt(t.Unix(), 10)
}

func doRequest(method string, url string, header map[string]string) (*http.Response, error) {
	// set the timeout to 3 seconds because if the fastgate restarts, it will not respond
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range header {
		req.Header.Set(key,value)
	}
	return client.Do(req)
}

func main(){
	user := "admin"
	// during my analysis I found that the password is base64 encoded
	pass := base64.URLEncoding.EncodeToString([]byte("admin"))
	defaultHeader :=  map[string]string{
		"Connection": "keep-alive",
		"Pragma": "no-cache",
		"Cache-Control": "no-cache",
		"Accept": "application/json, text/plain, */*",
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36",
		"DNT": "1",
		"Referer": "http://192.168.1.254/",
		"Accept-Language": "en-US,en;q=0.9",
	}


	// get login token
	resp, err := doRequest(
		"GET",
		"http://192.168.1.254/status.cgi?_="+getTimeStamp()+"&cmd=7&nvget=login_confirm",
		defaultHeader,
	)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var loginPage LoginConfirm
	if err := json.Unmarshal(bodyText, &loginPage); err != nil {
		log.Fatal(err)
	}
	log.Printf("Login page token: %s\n", loginPage.Login_confirm.Token)


	// do login and get session cookie
	resp, err = doRequest(
		"GET",
		"http://192.168.1.254/status.cgi?_="+getTimeStamp()+"&cmd=3&nvget=login_confirm&password="+pass+"&remember_me=1&token="+loginPage.Login_confirm.Token+"&username="+user,
		defaultHeader,
	)
	if err != nil {
		log.Fatal(err)
	}
	// get only first part of cookies
	loginCookie := strings.Split(resp.Header.Get("Set-cookie"), ";")[0]
	log.Printf("Session cookie: %s\n",  loginCookie)
	// assign the cookie to the header
	defaultHeader["Cookie"] = loginCookie


	// get system info options token, every page inside the fastgate require a unique token
	resp, err = doRequest(
		"GET",
		"http://192.168.1.254/status.cgi?_="+getTimeStamp()+"&nvget=sysinfo",
		defaultHeader,
	)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var sysInfoPage SysInfo
	if err := json.Unmarshal(bodyText, &sysInfoPage); err != nil {
		log.Fatal(err)
	}
	log.Printf("System info page token: %s\n", sysInfoPage.Sysinfo.Token)

	// the magic (do restart)
	resp, err = doRequest(
		"GET",
		"http://192.168.1.254/status.cgi?_="+getTimeStamp()+"&act=nvset&service=reset&token="+sysInfoPage.Sysinfo.Token,
		defaultHeader,
	)
	if err != nil {
		log.Fatal(err)
	}
}