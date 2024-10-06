package models

type Message struct {
	Number           int    `json:"number"`
	ChatNumber       int    `json:"chat_number"`
	ApplicationToken string `json:"application_token"`
	Content          string `json:"content"`
}
