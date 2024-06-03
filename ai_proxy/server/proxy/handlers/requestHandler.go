// handlers/requestHandler.go
package handlers

import (
	"ai-proxy/server/proxy/config"
	"ai-proxy/server/proxy/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"context"
	"net/http"
	"sync/atomic"

	openai "github.com/sashabaranov/go-openai"
)

var currentNode uint32

func ProxyRequest(c *gin.Context) {
	apiKey := c.GetHeader("API-Key")
	if apiKey == "" {
		c.JSON(400, gin.H{"error": "API key is required"})
		return
	}

	apiKeyInfo, err := models.GetAPIKeyInfo(apiKey)
	if err != nil {
		config.Logger.Error("Failed to get API key info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API key info"})
		return
	}

	node := apiKeyInfo.Nodes[atomic.AddUint32(&currentNode, 1)%uint32(len(apiKeyInfo.Nodes))]

	client := openai.NewClient(apiKey)
	request := openai.CompletionRequest{
		Prompt:    c.PostForm("prompt"),
		MaxTokens: 100,
	}

	response, err := client.CreateCompletion(context.Background(), request)
	if err != nil {
		config.Logger.Error("Failed to get response from OpenAI", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from OpenAI"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response.Choices[0].Text, "node": node})
}
