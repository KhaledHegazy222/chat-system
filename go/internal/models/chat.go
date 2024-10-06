package models

type Chat struct {
	Number           int    `json:"number"`
	ApplicationToken string `json:"application_token"`
	Title            string `json:"title"`
}
