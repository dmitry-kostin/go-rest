package models

type Ping struct {
	Pong    string `json:"pong"`
	Version string `json:"version"`
	Env     string `json:"env"`
}
