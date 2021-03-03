package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Runner struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type Runners []Runner

func deleteRunner(apiKey string, ID int64) error {
	u := fmt.Sprintf("https://gitlab.com/api/v4/runners/%d", ID)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)

	req.Header.Add("PRIVATE-TOKEN", apiKey)
	c := http.Client{}
	_, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func getRunners(apiKey string) (Runners, error) {
	runners := Runners{}

	req, _ := http.NewRequest(http.MethodGet, "https://gitlab.com/api/v4/runners?scope=offline", nil)

	req.Header.Add("PRIVATE-TOKEN", apiKey)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	b, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &runners)
	if err != nil {
		log.Fatal(err)
	}

	return runners, err
}

func main() {
	apiKey := os.Getenv("GITLAB_PRIVATE_TOKEN")

	runners, err := getRunners(apiKey)
	if err != nil {
		return
	}

	for _, r := range runners {
		fmt.Println(r.ID)
		deleteRunner(apiKey, r.ID)
	}
}
