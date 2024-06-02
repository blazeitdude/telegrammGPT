package gptClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"telegrammGPT/pkg/botLogger"
	"telegrammGPT/pkg/historyCache"
)

type GptConfiguration struct {
	ApiURL    string `yaml:"ApiURL"`
	ApiKey    string `yaml:"ApiKey"`
	Model     string `yaml:"model"`
	MaxTokens string `yaml:"max_Tokens"`
}

type GptResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Text         string    `json:"text"`
	Index        int       `json:"index"`
	Logprobs     *Logprobs `json:"logprobs,omitempty"`
	FinishReason string    `json:"finish_reason"`
}

type Logprobs struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float64            `json:"token_logprobs"`
	TopLogprobs   []map[string]float64 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GptClient struct {
	Cache  *historyCache.Cache
	Client *http.Client
	conf   GptConfiguration
}

type GPTRequest struct {
	Model     string                 `json:"model"`
	Messages  []historyCache.Message `json:"messages"`
	MaxTokens int                    `json:"max_tokens"`
}

func InitGpt(configuration GptConfiguration) GptClient {
	var gptClient GptClient
	gptClient.Cache = historyCache.NewCache()
	gptClient.conf = configuration
	gptClient.Client = &http.Client{}
	return gptClient
}

func (c *GptClient) SendMessage(content string, userID string) (string, error) {
	logger := botLogger.GetLogger()
	message := historyCache.Message{
		Role:    "user",
		Content: content,
	}

	userCache := c.Cache.GetUserCache(userID)

	messagesWithHistory := append(userCache.Messages, message)

	requestBody := GPTRequest{
		Model:     c.conf.Model,
		Messages:  messagesWithHistory,
		MaxTokens: 150,
	}

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

	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Debug("Non-200 response: %s\n", resp)
		return "", nil
	}

	var gptResponse GptResponse
	err = json.Unmarshal(byteResp, &gptResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return "", nil
	}

	fmt.Println("Response from OpenAI:")
	for _, choice := range gptResponse.Choices {
		logger.Logger.Debug(choice.Text)
	}

	return gptResponse.Choices[0].Text, nil
}
