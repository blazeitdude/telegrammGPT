package gptClient

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
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

func (c *GptClient) SendMessage(message string) (string, error) {
	logger := botLogger.GetLogger()
	messages := []Message{
		{
			Role:    "user",
			Content: message,
		},
	}

	requestBody := RequestBody{
		Model:       c.conf.Model,
		Messages:    messages,
		Temperature: 0.7,
	}

	// Преобразуем тело запроса в JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return "", nil
	}

	if err != nil {
		logger.Logger.Debug("Failed to Marshal Request to ChatGPT")
		return "", nil
	}

	req, err := http.NewRequest("POST", c.conf.ApiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Logger.Debug("Failed to configure request to Gpt:", err)
		return "", nil
	}
	req.Header.Set("Authorization", "Bearer "+c.conf.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-200 response: %s\n", body)
		return "", nil
	}

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return "", nil
	}

	fmt.Println("Response from OpenAI:")
	for _, choice := range responseBody.Choices {
		fmt.Println(choice.Text)
	}
	return responseBody.Choices[0].Text, nil
}
