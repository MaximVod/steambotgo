package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/MaximVod/steambotgo/internal/logger"
)

type AiQueriesAPI struct {
	baseURL string
	client  *http.Client
}

func (f AiQueriesAPI) SearchGamesByUserQuery(ctx context.Context, query string) (string, error) {
	systemPrompt := "Ты — помощник, который исправляет названия видеоигр. Пользователь вводит неточное название. Твоя задача — предложить наиболее вероятное исправленное название из известных видеоигр, даже если уверенность не 100%. Верни ТОЛЬКО одно название. НЕ используй NOT_FOUND, если есть разумное предположение."

	appLogger := logger.New()

	// Очищаем query от лишних пробелов
	query = strings.TrimSpace(query)
	appLogger.Info("AI запрос", "original_query", query)

	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": query,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("не удалось сериализовать запрос: %w", err)
	}

	// Логируем отправляемый запрос для отладки
	appLogger.Info("AI запрос payload", "payload", string(body))

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.artemox.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("не удалось создать запрос: %w", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY не установлен")
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса к AI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный статус от AI API: %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать ответ AI API: %w", err)
	}

	// Логируем ответ для отладки
	appLogger.Info("AI ответ", "response", string(raw))

	// Парсинг ответа в формате chat.completion
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
				Role    string `json:"role"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("не удалось распарсить ответ AI: %w", err)
	}

	if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("AI вернул пустой ответ")
	}

	correctedName := strings.TrimSpace(result.Choices[0].Message.Content)
	appLogger.Info("AI исправил название игры", "original", query, "corrected", correctedName)

	return correctedName, nil
}
