package handler

import (
	"errors"
	"net/http"

	"github.com/GlebMoskalev/go-path-backend/internal/middleware"
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
	userID := middleware.OptionalUserID(r.Context())
	chapters := h.theoryService.ListChapters(r.Context(), userID)
	utils.ResponseWithJSON(w, http.StatusOK, chapters)
}

func (h *TheoryHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "chapterSlug")
	userID := middleware.OptionalUserID(r.Context())

	chapter, err := h.theoryService.GetChapter(r.Context(), slug, userID)
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
	userID := middleware.OptionalUserID(r.Context())

	lesson, err := h.theoryService.GetLesson(r.Context(), chapterSlug, lessonSlug, userID)
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

func (h *TheoryHandler) MarkLessonCompleted(w http.ResponseWriter, r *http.Request) {
	chapterSlug := chi.URLParam(r, "chapterSlug")
	lessonSlug := chi.URLParam(r, "lessonSlug")
	userID, ok := middleware.UserIDFromContext(r.Context())

	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.theoryService.MarkLessonCompleted(r.Context(), userID, chapterSlug, lessonSlug)
	if err != nil {
		if errors.Is(err, service.ErrChapterNotFound) || errors.Is(err, service.ErrLessonNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "lesson not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "lesson marked as completed"})
}
