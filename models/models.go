package models

import "time"

type SearchBody struct {
	TypeSearch string `json:"type_search"`
	Term       string `json:"term"`
	From       int    `json:"from"`
	MaxResults int    `json:"max_results"`
}

type Body struct {
	SearchType string `json:"search_type"`
	Query      QueryS `json:"query"`
	From       int    `json:"from"`
	MaxResults int    `json:"max_results"`
}

type QueryS struct {
	Term string `json:"term"`
}

type Response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Result  map[string]interface{} `json:"result,omitempty"`
}

type Data struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}

type Email struct {
	ID                        int       `json:"ID"`
	Message_ID                string    `json:"Message-ID"`
	Date                      time.Time `json:"Date"`
	From                      string    `json:"from"`
	To                        string    `json:"to"`
	Subject                   string    `json:"subject"`
	Mime_Version              string    `json:"Mime-Version"`
	Content_Type              string    `json:"Content-Type"`
	Content_Transfer_Encoding string    `json:"Content-Transfer-Encoding"`
	X_From                    string    `json:"X-From"`
	X_To                      string    `json:"X-To"`
	X_cc                      string    `json:"X-cc"`
	X_bcc                     string    `json:"X-bcc"`
	X_Folder                  string    `json:"X-Folder"`
	X_Origin                  string    `json:"X-Origin"`
	X_FileName                string    `json:"X-FileName"`
	Cc                        string    `json:"Cc"`
	Body                      string    `json:"Body"`
}
