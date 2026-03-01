package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/GlebMoskalev/go-path-backend/internal/utils"
)

type QuizHandler struct {
	quizService *service.QuizService
}

func NewQuizHandler(quizService *service.QuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

func (h *QuizHandler) ListChapters(w http.ResponseWriter, r *http.Request) {
	chapters := h.quizService.ListChapters()
	utils.ResponseWithJSON(w, http.StatusOK, chapters)
}

func (h *QuizHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	chaptersParam := r.URL.Query().Get("chapters")
	limitParam := r.URL.Query().Get("limit")

	var chapterSlugs []string
	if chaptersParam != "" {
		chapterSlugs = strings.Split(chaptersParam, ",")
	}

	limit := 0
	if limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit < 0 {
			utils.ResponseWithError(w, http.StatusBadRequest, "invalid limit")
			return
		}
	}

	questions := h.quizService.GetQuestions(chapterSlugs, limit)
	utils.ResponseWithJSON(w, http.StatusOK, questions)
}

func (h *QuizHandler) CheckAnswer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		QuestionID string `json:"question_id"`
		Answer     int    `json:"answer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.QuestionID == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "question_id is required")
		return
	}

	result, err := h.quizService.CheckAnswer(req.QuestionID, req.Answer)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "question not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, result)
}
