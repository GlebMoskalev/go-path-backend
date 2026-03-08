package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

var (
	ErrTaskNotPassed = errors.New("task not passed")
)

type AIService struct {
	log              *zap.Logger
	client           *openai.Client
	serviceTask      *TaskService
	modelPassedTests string
	systemPromptTask string
	userPromptTask   string
	maxTokensTask    int
	temperature      float32
	topP             float32
}

func NewAIService(
	log *zap.Logger,
	serviceTask *TaskService,
	apiKey, apiUrl, modelPassedTests, systemPromptTask, userPromptTask string,
	maxTokensTask int,
	temperature, topP float32,
) *AIService {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = apiUrl
	client := openai.NewClientWithConfig(config)
	ai := &AIService{
		log:              log,
		serviceTask:      serviceTask,
		client:           client,
		modelPassedTests: modelPassedTests,
		systemPromptTask: systemPromptTask,
		userPromptTask:   userPromptTask,
		maxTokensTask:    maxTokensTask,
		temperature:      temperature,
		topP:             topP,
	}
	return ai
}

func (s *AIService) AnalyzePassedCodeTask(ctx context.Context, chapterSlug, taskSlug, code string, userID uuid.UUID) (string, error) {
	task, err := s.serviceTask.GetTask(ctx, chapterSlug, taskSlug, &userID)
	if err != nil {
		s.log.Error("failed get task", zap.Error(err), zap.String("chapterSlug", chapterSlug), zap.String("taskSlug", taskSlug))
		return "", err
	}

	if task.Solved == nil || !*task.Solved {
		s.log.Warn("task not passed")
		return "", ErrTaskNotPassed
	}

	replacer := strings.NewReplacer(
		"{{title}}", task.Title,
		"{{description}}", task.Description,
		"{{code}}", code,
	)
	userContent := replacer.Replace(s.userPromptTask)

	chatResp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.modelPassedTests,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.systemPromptTask},
			{Role: "user", Content: userContent},
		},
		MaxCompletionTokens: s.maxTokensTask,
		Temperature:         s.temperature,
		TopP:                s.topP,
	})

	if err != nil {
		s.log.Error("failed create chat", zap.Error(err))
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		s.log.Error("empty choices from LLM")
		return "", errors.New("empty response from AI")
	}
	return chatResp.Choices[0].Message.Content, nil
}
