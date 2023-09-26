package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/jordic/goics"

	"github.com/araddon/dateparse"
	"github.com/manifoldco/promptui"
)

func main() {

	prompt := promptui.Prompt{
		Label: "Username: ",
	}

	usernameResult, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	prompt = promptui.Prompt{
		Label: "Password: ",
	}

	passwordResult, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	client := &http.Client{
		Jar: jar,
		// disable redirection (mainly for the login part)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	log.Print("Logging in Kordis (MyGES API)...")
	account, err := loginKordis(usernameResult, passwordResult, client)
	if err != nil {
		log.Fatal(err)
	}

	profile, err := fetchMe(account.AccessToken, client)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Login successful. Welcome %s %s (UID: %d).", profile.Result.Firstname, profile.Result.Name, profile.Result.UID)

	prompt = promptui.Prompt{
		Label: "Start Date (mm/dd/yyyy): ",
		Validate: func(input string) error {
			_, err := dateparse.ParseAny(input)
			if err != nil {
				return errors.New("invalid date")
			}
			return nil
		},
	}

	// startResult, err := prompt.Run()
	startResult := "08/28/2023"

	// if err != nil {
	// 	log.Fatalf("Prompt failed %v\n", err)
	// }

	// prompt = promptui.Prompt{
	// 	Label: "End Date (mm/dd/yyyy): ",
	// 	Validate: func(input string) error {
	// 		_, err := dateparse.ParseAny(input)
	// 		if err != nil {
	// 			return errors.New("invalid date")
	// 		}
	// 		return nil
	// 	},
	// }

	// endResult, err := prompt.Run()
	endResult := "07/27/2024"

	// if err != nil {
	// 	log.Fatalf("prompt failed %v\n", err)
	// }

	startDate, err := dateparse.ParseAny(startResult)
	if err != nil {
		log.Fatalf("parsing start date failed %v\n", err)
	}

	endDate, err := dateparse.ParseAny(endResult)
	if err != nil {
		log.Fatalf("parsing end date failed %v\n", err)
	}

	agendaResult, err := fetchAgenda(account.AccessToken, fmt.Sprintf("%d", startDate.UnixMilli()), fmt.Sprintf("%d", endDate.UnixMilli()), client)
	if err != nil {
		log.Fatalf("error while fetching agenda %v", err)
	}

	for _, result := range agendaResult.Result {
		log.Printf("%s en %s - %s avec %s ", result.Type, result.Modality, result.Name, result.Discipline.Teacher)
	}

	cal := addEvent(agendaResult.Result)

	file, err := os.OpenFile("calendar.ics", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := goics.NewICalEncode(file)

	cal.Write(encoder)
}

// loginKordis is used to login on the Kordis platform, giving us an access token
func loginKordis(username, password string, client *http.Client) (*Account, error) {
	formatedString := fmt.Sprintf("%s:%s", username, password)
	// encode username:password into base64 string
	encodedToken := b64.StdEncoding.EncodeToString([]byte(formatedString))

	req, err := http.NewRequest("GET", "https://authentication.kordis.fr/oauth/authorize?response_type=token&client_id=skolae-app", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"94\", \"Google Chrome\";v=\"94\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://authentication.kordis.fr")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "skolae-app-ios/3.5.0 (com.reseauges.skolae.app; build:26; iOS 15.0.1) Alamofire/4.9.1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://authentication.kordis.fr/login?service=https%3A%2F%2Fmyges.fr%2Fj_spring_cas_security_check")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")
	req.Header.Set("Authorization", "Basic "+encodedToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// todo: add more status code error
	if resp.StatusCode == 401 {
		return nil, errors.New("invalid login")
	}

	// retrieve location header
	locationHeader := resp.Header.Get("Location")

	// remove skolae protocol url
	locationHeader = strings.TrimPrefix(locationHeader, "comreseaugesskolae:/oauth2redirect#")

	splitedHeader := strings.Split(locationHeader, "&")

	accountCreds := Account{}
	for _, param := range splitedHeader {
		paramSplited := strings.Split(param, "=")
		switch paramSplited[0] {
		case "access_token":
			accountCreds.AccessToken = paramSplited[1]
		case "token_type":
			accountCreds.TokenType = paramSplited[1]
		case "expires_in":
			accountCreds.ExpiresIn = paramSplited[1]
		case "scope":
			accountCreds.Scope = paramSplited[1]
		}
	}

	return &accountCreds, err
}

// fetchMe is used to retrieve all the profiles information (optional, just to add more info and proximity with the user)
func fetchMe(accessToken string, client *http.Client) (*KordisResponseProfile, error) {
	req, err := http.NewRequest("GET", "https://api.kordis.fr/me/profile", nil)
	if err != nil {
		return nil, err
	}

	req.Host = "api.kordis.fr"
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", "skolae-app-ios/3.5.0 (com.reseauges.skolae.app; build:26; iOS 15.0.1) Alamofire/4.9.1")
	req.Header.Set("Accept-Language", "fr-FR;q=1.0, en-FR;q=0.9, en-GB;q=0.8, ar-FR;q=0.7, ja-FR;q=0.6, el-FR;q=0.5")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profileJSON KordisResponseProfile
	err = json.Unmarshal(body, &profileJSON)
	if err != nil {
		return nil, err
	}

	return &profileJSON, nil
}

// fetchAgenda is used to fetch the agenda, start & end date are epoch in milliseconds
func fetchAgenda(accessToken, startTimestamp, endTimestamp string, client *http.Client) (*KordisResponseAgenda, error) {

	req, err := http.NewRequest("GET", "https://api.kordis.fr/me/agenda?start="+startTimestamp+"&end="+endTimestamp, nil)
	if err != nil {
		return nil, err
	}
	req.Host = "api.kordis.fr"
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("User-Agent", "skolae-app-ios/3.5.0 (com.reseauges.skolae.app; build:26; iOS 15.0.1) Alamofire/4.9.1")
	req.Header.Set("Accept-Language", "fr-FR;q=1.0, en-FR;q=0.9, en-GB;q=0.8, ar-FR;q=0.7, ja-FR;q=0.6, el-FR;q=0.5")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// todo: add more status code error
	if resp.StatusCode == 401 {
		return nil, errors.New("invalid login")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var agendaJSON KordisResponseAgenda
	err = json.Unmarshal(body, &agendaJSON)
	if err != nil {
		return nil, err
	}

	return &agendaJSON, err
}

// addEvent function is used to add event into the calendar
func addEvent(courses []KordisResultAgenda) goics.Componenter {

	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")

	for _, course := range courses {
		startDate := time.Unix(course.StartDate/1000, 0)
		endDate := time.Unix(course.EndDate/1000, 0)

		s := goics.NewComponent()
		s.SetType("VEVENT")

		k, v := goics.FormatDateTimeField("DTSTART", startDate)
		v = strings.TrimSuffix(v, "\n") + "Z"
		s.AddProperty(k, v)

		k, v = goics.FormatDateTimeField("DTEND", endDate)
		v = strings.TrimSuffix(v, "\n") + "Z"
		s.AddProperty(k, v)

		s.AddProperty("SUMMARY", fmt.Sprintf("%s - %s", course.Type, course.Name))

		if course.Rooms != nil {
			s.AddProperty("DESCRIPTION", fmt.Sprintf("Avec %s en %s (Campus: %s)", course.Teacher, course.Rooms[0].Name, course.Rooms[0].Campus))
		} else if course.Modality == "Distanciel" {
			s.AddProperty("DESCRIPTION", fmt.Sprintf("Avec %s en DISTANCIEL (Check Teams/MyGES)", course.Teacher))
		} else {
			s.AddProperty("DESCRIPTION", fmt.Sprintf("Avec %s - Pas de salle précisé", course.Teacher))
		}

		c.AddComponent(s)
	}

	return c
}
