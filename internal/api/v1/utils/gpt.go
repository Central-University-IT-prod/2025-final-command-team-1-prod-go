package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"example.com/m/internal/api/v1/infrastructure/logger"
	"example.com/m/internal/config"
)

type TokenResponse struct {
	IAMToken  string `json:"iamToken"`
	ExpiresAt string `json:"expiresAt"`
}

func UpdateYandexIAMToken() error {
	url := "https://iam.api.cloud.yandex.net/iam/v1/tokens"

	payload := map[string]string{
		"yandexPassportOauthToken": config.Config.YandexOAuthToken,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return err
	}
	if tokenResp.IAMToken == "" {
		config.IAMToken = "empty"
		logger.Logger.Fatal("IAM TOEKN WASNT UPDATED")
	} else {
		config.IAMToken = tokenResp.IAMToken
	}
	return nil
}
