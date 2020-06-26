package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
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

func (s Service) ReportMeterReading(reading int) error {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://fgmwebdiszpp.nkmenergia.hu/sap/bc/webdynpro/sap/zusz2_wd_meroallas_rogz",
		nil,
	)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Mobile Safari/537.36")

	resp, err := s.client.Do(req)
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

	form := doc.Find("#sap\\.client\\.SsrClient\\.form")
	action, ok := form.Attr("action")
	if !ok {
		return fmt.Errorf("can't find sap.client.SsrClient.form's action")
	}
	sapwdsecureid, ok := form.Find("#sap-wd-secure-id").Attr("value")
	if !ok {
		return fmt.Errorf("can't find sap-wd-secure-id")
	}

	data := url.Values{}
	data.Set("sap-charset", "utf-8")
	data.Set("sap-wd-secure-id", sapwdsecureid)
	data.Set("_stateful_", "X")
	data.Set("SAPEVENTQUEUE", "ComboBox_Change~E002Id~E004WD8A~E005Value~E004"+strconv.Itoa(reading)+"~E003~E002ResponseData~E004delta~E005EnqueueCardinality~E004single~E005Delay~E004full~E003~E002~E003~E001Button_Press~E002Id~E004WD94~E003~E002ResponseData~E004delta~E005ClientAction~E004submit~E003~E002~E003")

	req, err = http.NewRequest(
		http.MethodPost,
		"https://fgmwebdiszpp.nkmenergia.hu/"+action,
		strings.NewReader(data.Encode()),
	)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Mobile Safari/537.36")

	resp, err = s.client.Do(req)
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
	err = service.ReportMeterReading(1036)
	if err != nil {
		fmt.Println(err)
		return
	}
}
