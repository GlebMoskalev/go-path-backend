package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}

func TestRecoveryCatchesPanic(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	h := Recovery(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d after panic", w.Code, http.StatusInternalServerError)
	}
}

func TestRecoveryDoesNotAffectNormalRequest(t *testing.T) {
	h := Recovery(okHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d for normal request", w.Code, http.StatusOK)
	}
}

func TestLoggerPassesThroughStatus(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	h := Logger(inner)
	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d — Logger must not change response status", w.Code, http.StatusCreated)
	}
}

func TestLoggerPassesThroughBody(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	h := Logger(inner)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Body.String() != "hello" {
		t.Errorf("body = %q, want %q — Logger must not change response body", w.Body.String(), "hello")
	}
}

func TestChainAppliesMiddleware(t *testing.T) {
	var order []string

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1-before")
			next.ServeHTTP(w, r)
			order = append(order, "m1-after")
		})
	}
	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2-before")
			next.ServeHTTP(w, r)
			order = append(order, "m2-after")
		})
	}

	h := Chain(okHandler(), m1, m2)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if len(order) != 4 {
		t.Fatalf("executed %d middleware hooks, want 4: %v", len(order), order)
	}
	if order[0] != "m1-before" {
		t.Errorf("first hook = %q, want %q — m1 should be outermost", order[0], "m1-before")
	}
}

func TestChainRecoveryAndLogger(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("oh no")
	})

	h := Chain(panicHandler, Logger, Recovery)
	req := httptest.NewRequest(http.MethodGet, "/crash", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500 when panic is caught by Recovery", w.Code)
	}
}
