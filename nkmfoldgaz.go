package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewService() Service {
	cookieJar, _ := cookiejar.New(nil)

	return Service{client: &http.Client{
		Jar: cookieJar,
	}}
}

type Service struct {
	client *http.Client
}

func (s Service) Login(username string, password string) error {
	resp, err := s.client.Get("https://fgmwebdiszpp.nkmenergia.hu/sap/bc/webdynpro/sap/zusz2_wd_meroallas_rogz")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	xsrf, ok := doc.Find(`[name="sap-login-XSRF"]`).Attr("value")
	if !ok {
		return fmt.Errorf("can't find form with sap-login-XSRF")
	}

	data := url.Values{}
	data.Set("sap-system-login-oninputprocessing", "onLogin")
	data.Set("sap-urlscheme", "")
	data.Set("sap-system-login", "onLogin")
	data.Set("sap-system-login-basic_auth", "")
	data.Set("sap-client", "100")
	data.Set("sap-language", "EN")
	data.Set("sap-accessibility", "")
	data.Set("sap-login-XSRF", xsrf)
	data.Set("sap-system-login-cookie_disabled", "")
	data.Set("sap-alias", username)
	data.Set("sap-password", password)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://fgmwebdiszpp.nkmenergia.hu/sap/bc/webdynpro/sap/zusz2_wd_meroallas_rogz",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Mobile Safari/537.36")

	fmt.Printf("%#v\n", req)
	fmt.Println("")
	u, _ := url.ParseRequestURI("https://fgmwebdiszpp.nkmenergia.hu")
	for _, c := range s.client.Jar.Cookies(u) {
		fmt.Printf("%#v\n", c)
	}
	fmt.Println("")

	resp, err = s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	fmt.Println(bodyString)

	for _, c := range resp.Cookies() {
		fmt.Printf("%#v\n", c)
	}
	fmt.Println("")
	for _, c := range s.client.Jar.Cookies(u) {
		fmt.Printf("%#v\n", c)
	}

	return nil
}

func main() {
	username := os.Getenv("NKM_USERNAME")
	password := os.Getenv("NKM_PASSWORD")
	service := NewService()
	err := service.Login(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = service.Login(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
}
