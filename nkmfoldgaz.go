package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
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

	data := url.Values{}
	data.Set("sap-alias", username)
	data.Set("sap-password", password)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://fgmwebdiszpp.nkmenergia.hu/sap/bc/webdynpro/sap/zusz2_wd_bejelentkezes",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Mobile Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
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
}
