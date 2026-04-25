package middleware

import (
	"errors"
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

// =====================  RECOVERY  =====================

func TestRecoveryCatchesStringPanic(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	h := Recovery(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	defer func() {
		if v := recover(); v != nil {
			t.Errorf("panic %q просочился через Recovery", v)
		}
	}()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d after panic", w.Code, http.StatusInternalServerError)
	}
}

func TestRecoveryCatchesErrorPanic(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("error panic"))
	})

	h := Recovery(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	defer func() {
		if v := recover(); v != nil {
			t.Errorf("panic просочился через Recovery: %v", v)
		}
	}()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRecoveryCatchesNilDereference(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p *int
		_ = *p
	})

	h := Recovery(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	defer func() {
		if v := recover(); v != nil {
			t.Errorf("nil-разыменование просочилось через Recovery: %v", v)
		}
	}()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
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
	if w.Body.String() != "ok" {
		t.Errorf("body = %q, want %q (Recovery не должен менять успешный ответ)",
			w.Body.String(), "ok")
	}
}

func TestRecoveryDoesNotEatStatusCode(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	h := Recovery(inner)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d (Recovery не должен менять статус если паники не было)",
			w.Code, http.StatusTeapot)
	}
}

// =====================  LOGGER  =====================

func TestLoggerPassesThroughStatus(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	h := Logger(inner)
	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d — Logger не должен менять статус ответа",
			w.Code, http.StatusCreated)
	}
}

func TestLoggerPassesThroughBody(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	h := Logger(inner)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Body.String() != "hello world" {
		t.Errorf("body = %q, want %q — Logger не должен менять тело ответа",
			w.Body.String(), "hello world")
	}
}

func TestLoggerPassesThroughHeaders(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom", "value-123")
		w.WriteHeader(http.StatusOK)
	})

	h := Logger(inner)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if got := w.Header().Get("X-Custom"); got != "value-123" {
		t.Errorf("X-Custom header = %q, want %q (Logger не должен терять заголовки)",
			got, "value-123")
	}
}

func TestLoggerWithMultipleStatuses(t *testing.T) {
	cases := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent,
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusInternalServerError,
	}
	for _, code := range cases {
		c := code
		inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(c)
		})
		h := Logger(inner)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != c {
			t.Errorf("status = %d, want %d", w.Code, c)
		}
	}
}

// =====================  CHAIN  =====================

func TestChainEmpty(t *testing.T) {
	h := Chain(okHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (Chain без middleware = handler как есть)", w.Code, http.StatusOK)
	}
	if w.Body.String() != "ok" {
		t.Errorf("body = %q, want %q", w.Body.String(), "ok")
	}
}

func TestChainAppliesMiddlewareInOrder(t *testing.T) {
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
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	h := Chain(finalHandler, m1, m2)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	want := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(want) {
		t.Fatalf("выполнено %d хуков, want %d: got=%v want=%v", len(order), len(want), order, want)
	}
	for i, w := range want {
		if order[i] != w {
			t.Errorf("order[%d] = %q, want %q (порядок: %v)", i, order[i], w, order)
		}
	}
}

func TestChainRecoveryAndLogger(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("oh no")
	})

	h := Chain(panicHandler, Logger, Recovery)
	req := httptest.NewRequest(http.MethodGet, "/crash", nil)
	w := httptest.NewRecorder()

	defer func() {
		if v := recover(); v != nil {
			t.Errorf("panic просочился: %v", v)
		}
	}()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500 когда Recovery ловит panic", w.Code)
	}
}

func TestChainLoggerOuterRecoveryInner(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	h := Chain(panicHandler, Recovery, Logger)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	defer func() {
		if v := recover(); v != nil {
			t.Errorf("panic просочился через Recovery: %v", v)
		}
	}()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

// =====================  RESPONSE WRITER WRAPPING =====================

func TestLoggerCapturesActualStatus(t *testing.T) {
	// Косвенная проверка что Logger оборачивает ResponseWriter:
	// если inner-handler вызывает Hijack или другие методы, это не должно ломаться.
	// Также если Logger не оборачивает, статус по умолчанию будет 200,
	// но реальный статус в ответе должен быть актуальным.

	statuses := []int{http.StatusCreated, http.StatusBadRequest, http.StatusInternalServerError}
	for _, want := range statuses {
		w := want
		inner := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(w)
		})
		h := Logger(inner)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != w {
			t.Errorf("Logger потерял реальный статус: got %d, want %d", rec.Code, w)
		}
	}
}
