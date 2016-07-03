package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"
)

const (
	API_BASE_URL = "https://discordapp.com/api"
)

type Gateway struct {
	Url string `json:"url"`
}

func (s *Session) GetGateway() (*Gateway, error) {
	resp, err := http.Get(API_BASE_URL + "/gateway")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	gateway := Gateway{}
	err = json.NewDecoder(resp.Body).Decode(&gateway)
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}

// https://discordapp.com/developers/docs/resources/channel#create-message
func (s *Session) CreateMessage(channelId string, content string) error {
	url := API_BASE_URL + "/channels/" + channelId + "/messages"
	params := struct {
		Content string `json:"content"`
	}{content}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("authorization", s.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Non-200 response")
	}

	return nil
}
