package google

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchEmailGmail(client *http.Client, url string) (any, error) {
	resProfile, _ := client.Get(url)
	body, _ := io.ReadAll(resProfile.Body)
	log.Printf("fetchEmailGmail.body: %v", string(body))
	type profileT struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
		Name          string `json:"name"`
		FirstName     string `json:"given_name"`
		LastName      string `json:"family_name"`
	}
	var pp profileT
	json.Unmarshal(body, &pp)
	log.Printf("%v: %v", "", string(body))
	// log.Printf("profile %+v\n\n\n ", pp)
	if pp.Email == "" {
		log.Printf("%v: %v", "", string(body))
		log.Printf("profile %+v", pp)
		return nil, fmt.Errorf("unable to find email in profile")
	}
	_ = map[string]string{
		"provider":   "google",
		"id":         pp.Id,
		"first_name": pp.FirstName,
		"last_name":  pp.LastName,
		"name":       pp.Name,
		"picture":    pp.Picture,

		// "verified_email"

	}
	// return &AuthProfile{Email: pp.Email, ID: pp.Id, Picture: pp.Picture, Attr: attr}, nil
	return nil, nil
}
