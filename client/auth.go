package client

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"
)

type AuthInterface struct {
	Client *BboxClient
}

type DeviceTokenResponse struct {
	Device DeviceToken `json:"device"`
}

func (ai *AuthInterface) ObtainBearerToken() error {
	resp, err := ai.Client.Get("/device/token")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// La réponse est un array
	var responses []DeviceTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&responses); err != nil {
		return err
	}

	if len(responses) == 0 {
		return errors.New("no device token in response")
	}

	ai.Client.Bearer = &responses[0].Device
	return nil
}

func (ai *AuthInterface) StartTokenRefresher() error {
	if ai.Client.Bearer == nil {
		return errors.New("can't start before BasicAuth")
	}

	go func() {
		for {
			expiryTime, _ := time.Parse("2006-01-02T00:00:00+0100", ai.Client.Bearer.Expires)
			<-time.After(time.Until(expiryTime) - 1*time.Minute)
			ai.ObtainBearerToken()
		}
	}()

	return nil
}

func (ai *AuthInterface) BasicAuth(password string) error {
	_, err := ai.Client.Post(
		"/login", "application/x-www-form-urlencoded", io.Reader(strings.NewReader("password="+password)),
	)
	if err != nil {
		return err
	}
	ai.ObtainBearerToken()
	return nil
}
