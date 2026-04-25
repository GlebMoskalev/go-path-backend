package middleware

// Реализуйте middleware для HTTP-сервера.
//
// Импорты:
//   import (
//       "log"
//       "net/http"
//       "time"
//   )
//
// 1. Обёртка для перехвата статуса ответа:
//    type responseWriter struct {
//        http.ResponseWriter
//        status int
//    }
//    func (rw *responseWriter) WriteHeader(code int) — сохраняет code в status
//
// 2. func Logger(next http.Handler) http.Handler
//    — оберни ResponseWriter в responseWriter (status по умолчанию 200)
//    — запомни start := time.Now()
//    — вызови next.ServeHTTP(rw, r)
//    — залогируй: log.Printf("%s %s → %d %s", r.Method, r.URL.Path, rw.status, time.Since(start))
//
// 3. func Recovery(next http.Handler) http.Handler
//    — используй defer func() { if v := recover(); v != nil { ... } }()
//    — при панике: log.Printf("panic: %v", v) и WriteHeader(500)
//
// 4. func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler
//    — применяй middleware справа налево:
//      for i := len(middlewares) - 1; i >= 0; i-- { h = middlewares[i](h) }
//    — верни h
