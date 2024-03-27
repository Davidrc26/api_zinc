package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Davidrc26/api_zinc.git/models"
)

type Body struct {
	SearchType string `json:"search_type"`
	Query      QueryS `json:"query"`
	MaxResults int    `json:"max_results"`
}

type QueryS struct {
	Term string `json:"term"`
}

const BASE_URL = "http://localhost:4080/api/"

func Search(bdy *models.SearchBody) map[string]interface{} {
	usr := os.Getenv("USER_ZINCSEARCH")
	password := os.Getenv("PASSWORD_ZINCSEARCH")
	b := Body{
		SearchType: bdy.TypeSearch,
		MaxResults: 20,
		Query: QueryS{
			Term: bdy.Term,
		},
	}
	jsonData, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", BASE_URL+"maildir/_search", bytes.NewReader(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(usr, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyReader := io.NopCloser(bytes.NewBuffer(body))
	/* fmt.Println(string(body)) */

	var result map[string]interface{}
	err = json.NewDecoder(bodyReader).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetIndexByName(name string) []string {
	usr := os.Getenv("USER_ZINCSEARCH")
	password := os.Getenv("PASSWORD_ZINCSEARCH")

	req, err := http.NewRequest("GET", BASE_URL+"index_name?name="+name, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(usr, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)

	var result []string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
	/* fmt.Println(string(body)) */

}
