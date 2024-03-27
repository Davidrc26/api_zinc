package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../environment/environment.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	mux := Routes()
	server := NewServer(mux)
	server.Run()
	/* query := `{
	        "search_type": "match",
	        "query":
	        {
	            "term": "Sorry for the delay"
	        },
	        "from": 0,
	        "max_results": 20,
	        "_source": ["Body"]
	    }`
	req, err := http.NewRequest("POST", "http://localhost:4080/api/shackleton-s/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "#22171Drc#")
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
	/* fmt.Println(string(body))

	var result map[string]interface{}
	err = json.NewDecoder(bodyReader).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result["hits"]) */
}
