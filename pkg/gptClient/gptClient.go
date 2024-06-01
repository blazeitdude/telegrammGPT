package gptClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"telegrammGPT/pkg/botLogger"
)

type GptConfiguration struct {
	ApiURL    string `yaml:"ApiURL"`
	ApiKey    string `yaml:"ApiKey"`
	Model     string `yaml:"model"`
	MaxTokens string `yaml:"max_Tokens"`
}

type GptResponse struct {
	ResponseHeader string
	ResponseBody   string
}

type GptClient struct {
	Client *http.Client
	conf   GptConfiguration
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

func InitGpt(configuration GptConfiguration) GptClient {
	var gptClient GptClient
	gptClient.conf = configuration
	gptClient.Client = &http.Client{}
	return gptClient
}

func (c *GptClient) SendMessage(message string) (GptResponse, error) {
	logger := botLogger.GetLogger()
	var response GptResponse
	messages := []Message{
		{
			Role:    "user",
			Content: "Say this is a test!",
		},
	}

	requestBody := RequestBody{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		Temperature: 0.7,
	}

	// Преобразуем тело запроса в JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return GptResponse{}, nil
	}

	if err != nil {
		logger.Logger.Debug("Failed to Marshal Request to ChatGPT")
		return GptResponse{}, nil
	}

	req, err := http.NewRequest("POST", c.conf.ApiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Logger.Debug("Failed to configure request to Gpt:", err)
		return GptResponse{}, nil
	}
	req.Header.Set("Authorization", "Bearer "+c.conf.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		logger.Logger.Debug("Failed to send Message to GPT", err)
		return GptResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Logger.Debug("Failed to close response body ReaderCloser??", err)
			return
		}
	}(resp.Body)

	readedBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Debug("Failed to Read Response from GPT", err)
		return GptResponse{}, err
	}
	logger.Logger.Debugf("header from response: %v", resp.Header)
	response.ResponseBody = string(readedBody)
	return response, nil
}
