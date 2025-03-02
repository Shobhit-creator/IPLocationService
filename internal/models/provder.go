package models

type Provider struct {
	Name       string `json:"name"`
	Url        string `json:"url"`
	Token      string `json:"token"`
	DailyLimit int    `json:"dailyLimit"`
}
