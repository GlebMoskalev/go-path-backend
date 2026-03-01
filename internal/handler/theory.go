package handler

import (
	"errors"
	"net/http"

	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/GlebMoskalev/go-path-backend/internal/utils"
	"github.com/go-chi/chi/v5"
)

type TheoryHandler struct {
	theoryService *service.TheoryService
}

func NewTheoryHandler(theoryService *service.TheoryService) *TheoryHandler {
	return &TheoryHandler{theoryService: theoryService}
}

func (h *TheoryHandler) ListChapter(w http.ResponseWriter, r *http.Request) {
	chapters := h.theoryService.ListChapters()
	utils.ResponseWithJSON(w, http.StatusOK, chapters)
}

func (h *TheoryHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "chapterSlug")

	chapter, err := h.theoryService.GetChapter(slug)
	if err != nil {
		if errors.Is(err, service.ErrChapterNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "chapter not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, chapter)
}

func (h *TheoryHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
	chapterSlug := chi.URLParam(r, "chapterSlug")
	lessonSlug := chi.URLParam(r, "lessonSlug")

	lesson, err := h.theoryService.GetLesson(chapterSlug, lessonSlug)
	if err != nil {
		if errors.Is(err, service.ErrLessonNotFound) || errors.Is(err, service.ErrChapterNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "lesson not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, lesson)
}
