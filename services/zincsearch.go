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

var base_url = ""
var usr = ""
var password = ""

func Init() {
	if base_url != "" {
		return
	}
	base_url = os.Getenv("ZINCSEARCH_URL")
	usr = os.Getenv("USER_ZINCSEARCH")
	password = os.Getenv("PASSWORD_ZINCSEARCH")
}

func Search(bdy *models.SearchBody) models.Response {

	if bdy.MaxResults <= 0 {
		return models.Response{
			Status:  400,
			Message: "MaxResults debe ser mayor a 0",
			Result:  nil,
		}
	}
	if bdy.From < 0 {
		return models.Response{
			Status:  400,
			Message: "From no puede ser negativo",
			Result:  nil,
		}
	}

	b := models.Body{
		SearchType: bdy.TypeSearch,
		From:       bdy.From,
		MaxResults: bdy.MaxResults,
		Query: models.QueryS{
			Term: bdy.Term,
		},
		Sort: []any{"Date:desc"},
	}
	jsonData, err := json.Marshal(b)
	if err != nil {
		return models.Response{
			Status:  500,
			Message: "Error al serializar la petición: " + err.Error(),
			Result:  nil,
		}
	}
	resp := SendRequest("POST", "/maildir/_search", jsonData)
	if resp == nil {
		return models.Response{
			Status:  500,
			Message: "Error al enviar la petición a ZincSearch",
			Result:  nil,
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Response{
			Status:  500,
			Message: "Error al leer la respuesta de ZincSearch: " + err.Error(),
			Result:  nil,
		}
	}
	bodyReader := io.NopCloser(bytes.NewBuffer(body))
	var result map[string]interface{}
	err = json.NewDecoder(bodyReader).Decode(&result)
	if err != nil {
		return models.Response{
			Status:  500,
			Message: "Error al decodificar la respuesta de ZincSearch: " + err.Error(),
			Result:  nil,
		}
	}
	return models.Response{
		Status:  200,
		Message: "Búsqueda exitosa",
		Result:  result,
	}
}

func GetIndexByName(name string) []string {
	resp := SendRequest("POST", "index_name?name="+name, nil)
	if resp == nil {
		log.Fatal("Error al obtener índice")
		return nil
	}
	defer resp.Body.Close()
	var result []string
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func BulkIndex(jsonData []byte) models.Response {

	resp := SendRequest("POST", "/_bulkv2", jsonData)
	if resp == nil {
		return models.Response{
			Status:  500,
			Message: "Error al enviar la petición de indexación a ZincSearch",
			Result:  nil,
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Response{
			Status:  500,
			Message: "Error al leer la respuesta de ZincSearch: " + err.Error(),
			Result:  nil,
		}
	}
	bodyReader := io.NopCloser(bytes.NewBuffer(body))
	var result map[string]interface{}
	err = json.NewDecoder(bodyReader).Decode(&result)
	if err != nil {
		return models.Response{
			Status:  500,
			Message: "Error al decodificar la respuesta de ZincSearch: " + err.Error(),
			Result:  nil,
		}
	}

	return models.Response{
		Status:  200,
		Message: "Indexación exitosa",
		Result:  result,
	}

}

func SendRequest(action_type string, endpoint string, jsonData []byte) *http.Response {
	Init()
	req, err := http.NewRequest(action_type, base_url+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(string(jsonData))
		log.Fatal("Error creating request:", err)
		return nil
	}

	req.SetBasicAuth(usr, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err.Error())
		return nil
	}
	return resp
}
