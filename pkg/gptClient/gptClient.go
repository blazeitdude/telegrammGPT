package gptClient

import (
	"bytes"
	"encoding/json"
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

func InitGpt(configuration GptConfiguration) GptClient {
	var gptClient GptClient
	gptClient.conf = configuration
	gptClient.Client = &http.Client{}
	return gptClient
}

func (c *GptClient) SendMessage(message string) (GptResponse, error) {
	logger := botLogger.GetLogger()
	var response GptResponse
	var messagesArray = []map[string]string{}
	messagesArray = append(messagesArray, map[string]string{})
	messagesArray[0]["role"] = "user"
	messagesArray = append(messagesArray, map[string]string{})
	messagesArray[1]["content"] = message

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       c.conf.Model,
		"temperature": 0.7,
		"max_tokens":  c.conf.MaxTokens,
		"messages":    messagesArray,
	})

	if err != nil {
		logger.Logger.Debug("Failed to Marshal Request to ChatGPT")
		return GptResponse{}, nil
	}

	req, err := http.NewRequest("POST", c.conf.ApiURL, bytes.NewBuffer(requestBody))
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
