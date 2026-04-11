package service

import (
	"context"
	"errors"
	"strings"

	"github.com/GlebMoskalev/go-path-backend/internal/config"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

var (
	ErrTaskNotPassed = errors.New("task not passed")
)

type AIService struct {
	log                 *zap.Logger
	client              *openai.Client
	serviceTask         *TaskService
	serviceProject      *ProjectService
	modelPassedTests    string
	systemPromptTask    string
	systemPromptProject string
	systemPromptError   string
	userPromptTask      string
	userPromptProject   string
	userPromptError     string
	maxTokensTask       int
	maxTokensProject    int
	temperature         float32
	topP                float32
}

func NewAIService(
	log *zap.Logger,
	serviceTask *TaskService,
	serviceProject *ProjectService,
	aiCfg config.AIConfig,
) *AIService {
	openaiCfg := openai.DefaultConfig(aiCfg.ApiKey)
	openaiCfg.BaseURL = aiCfg.ApiUrl
	client := openai.NewClientWithConfig(openaiCfg)
	return &AIService{
		log:                 log,
		serviceTask:         serviceTask,
		serviceProject:      serviceProject,
		client:              client,
		modelPassedTests:    aiCfg.ModelPassedTests,
		systemPromptTask:    aiCfg.SystemPromptTask,
		userPromptTask:      aiCfg.UserPromptTask,
		systemPromptProject: aiCfg.SystemPromptProject,
		userPromptProject:   aiCfg.UserPromptProject,
		systemPromptError:   aiCfg.SystemPromptError,
		userPromptError:     aiCfg.UserPromptError,
		maxTokensTask:       aiCfg.MaxTokensTask,
		maxTokensProject:    aiCfg.MaxTokensProject,
		temperature:         aiCfg.Temperature,
		topP:                aiCfg.TopP,
	}
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
		ChatTemplateKwargs: map[string]any{
			"enable_thinking": false,
		},
	})

	if err != nil {
		s.log.Error("failed create chat", zap.Error(err))
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		s.log.Error("empty choices from LLM")
		return "", errors.New("empty response from AI")
	}

	content := chatResp.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		s.log.Warn("LLM returned empty content",
			zap.String("model", s.modelPassedTests),
			zap.String("finishReason", string(chatResp.Choices[0].FinishReason)),
			zap.Int("promptTokens", chatResp.Usage.PromptTokens),
			zap.Int("completionTokens", chatResp.Usage.CompletionTokens),
		)
		return "", errors.New("AI returned empty recommendation")
	}

	return content, nil
}

func (s *AIService) AnalyzePassedCodeProject(ctx context.Context, projectSlug, stepSlug, code string, userID uuid.UUID) (string, error) {
	step, err := s.serviceProject.GetStep(ctx, projectSlug, stepSlug, &userID)
	if err != nil {
		s.log.Error("failed get project step",
			zap.Error(err),
			zap.String("projectSlug", projectSlug),
			zap.String("stepSlug", stepSlug),
		)
		return "", err
	}

	if step.Solved == nil || !*step.Solved {
		s.log.Warn("project step not passed")
		return "", ErrTaskNotPassed
	}

	replacer := strings.NewReplacer(
		"{{title}}", step.Title,
		"{{description}}", step.Description,
		"{{code}}", code,
	)
	userContent := replacer.Replace(s.userPromptProject)

	chatResp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.modelPassedTests,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.systemPromptProject},
			{Role: "user", Content: userContent},
		},
		MaxCompletionTokens: s.maxTokensProject,
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

	content := chatResp.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		s.log.Warn("LLM returned empty content",
			zap.String("model", s.modelPassedTests),
			zap.String("finishReason", string(chatResp.Choices[0].FinishReason)),
			zap.Int("promptTokens", chatResp.Usage.PromptTokens),
			zap.Int("completionTokens", chatResp.Usage.CompletionTokens),
		)
		return "", errors.New("AI returned empty recommendation")
	}

	return content, nil
}

func (s *AIService) AnalyzeErrorTask(ctx context.Context, chapterSlug, taskSlug, code, errorOutput string, userID uuid.UUID) (string, error) {
	task, err := s.serviceTask.GetTask(ctx, chapterSlug, taskSlug, &userID)
	if err != nil {
		s.log.Error("failed get task", zap.Error(err), zap.String("chapterSlug", chapterSlug), zap.String("taskSlug", taskSlug))
		return "", err
	}

	replacer := strings.NewReplacer(
		"{{title}}", task.Title,
		"{{description}}", task.Description,
		"{{code}}", code,
		"{{error}}", errorOutput,
	)
	userContent := replacer.Replace(s.userPromptError)

	s.log.Debug("analyzing error",
		zap.String("systemPrompt", s.systemPromptError),
		zap.Int("userContentLen", len(userContent)),
		zap.String("model", s.modelPassedTests),
	)

	chatResp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.modelPassedTests,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.systemPromptError},
			{Role: "user", Content: userContent},
		},
		MaxCompletionTokens: s.maxTokensTask,
		Temperature:         s.temperature,
		TopP:                s.topP,
		ChatTemplateKwargs: map[string]any{
			"enable_thinking": false,
		},
	})

	if err != nil {
		s.log.Error("failed create chat for error analysis", zap.Error(err))
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		s.log.Error("empty choices from LLM for error analysis")
		return "", errors.New("empty response from AI")
	}

	content := chatResp.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		s.log.Warn("LLM returned empty content for error analysis",
			zap.String("model", s.modelPassedTests),
			zap.String("finishReason", string(chatResp.Choices[0].FinishReason)),
		)
		return "", errors.New("AI returned empty error analysis")
	}

	return content, nil
}

func (s *AIService) AnalyzeErrorProject(ctx context.Context, projectSlug, stepSlug, code, errorOutput string, userID uuid.UUID) (string, error) {
	step, err := s.serviceProject.GetStep(ctx, projectSlug, stepSlug, &userID)
	if err != nil {
		s.log.Error("failed get project step",
			zap.Error(err),
			zap.String("projectSlug", projectSlug),
			zap.String("stepSlug", stepSlug),
		)
		return "", err
	}

	replacer := strings.NewReplacer(
		"{{title}}", step.Title,
		"{{description}}", step.Description,
		"{{code}}", code,
		"{{error}}", errorOutput,
	)
	userContent := replacer.Replace(s.userPromptError)

	s.log.Debug("analyzing project error",
		zap.String("systemPrompt", s.systemPromptError),
		zap.Int("userContentLen", len(userContent)),
		zap.String("model", s.modelPassedTests),
	)

	chatResp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.modelPassedTests,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.systemPromptError},
			{Role: "user", Content: userContent},
		},
		MaxCompletionTokens: s.maxTokensProject,
		Temperature:         s.temperature,
		TopP:                s.topP,
		ChatTemplateKwargs: map[string]any{
			"enable_thinking": false,
		},
	})

	if err != nil {
		s.log.Error("failed create chat for error analysis", zap.Error(err))
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		s.log.Error("empty choices from LLM for error analysis")
		return "", errors.New("empty response from AI")
	}

	content := chatResp.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		s.log.Warn("LLM returned empty content for error analysis",
			zap.String("model", s.modelPassedTests),
			zap.String("finishReason", string(chatResp.Choices[0].FinishReason)),
		)
		return "", errors.New("AI returned empty error analysis")
	}

	return content, nil
}
