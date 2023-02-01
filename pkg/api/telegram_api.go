package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	baseUrl    = "https://api.telegram.org"
	getUpdates = "getUpdates"
)

func urlAPIMethod(token string, method string) string {
	return fmt.Sprintf("%s/bot%s/%s", baseUrl, token, method)
}

type TelegramAPIInterface interface {
}

type Poller struct {
	token          string
	offset         int
	limit          int
	timeout        int
	allowedUpdates []string
}

func NewPoller(token string, timeout int) *Poller {

	return &Poller{
		token:   token,
		timeout: timeout,
	}
}

type TelegramAPI struct {
	token     string
	Poller    *Poller
	channelID string
}

func NewTelegramAPI(token string, channelID string) *TelegramAPI {
	return &TelegramAPI{
		token:     token,
		Poller:    NewPoller(token, 30),
		channelID: channelID,
	}
}

func (poll *Poller) generateUpdateUrl(updateID int) string {

	return fmt.Sprintf("%s?offset=%d?timeout=%d", urlAPIMethod(poll.token, getUpdates), updateID, poll.timeout)
}

func (poll *Poller) GetUpdates() ([]Update, error) {
	resp, err := http.Get(poll.generateUpdateUrl(poll.offset))

	body, err := ioutil.ReadAll(resp.Body)

	var data UpdateResponse
	err = json.Unmarshal(body, &data)

	if len(data.Result) > 0 {
		poll.offset = data.Result[len(data.Result)-1].UpdateID + 1
	}

	//fmt.Println(poll.offset, data)

	return data.Result, err
}

func (poll *Poller) ListenUpdates() chan Update {
	c := make(chan Update)

	go func() {
		for {
			updates, err := poll.GetUpdates()

			if err != nil {
				log.Fatalln(err)
			}

			for _, upd := range updates {

				c <- upd
			}

		}
	}()

	return c
}
