package main

import (
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

type middleware func(http.Handler) http.Handler

func createStack(xs ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (app *application) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)

		logData := fmt.Sprintf(
			"%d %s %s %v",
			wrapped.statusCode,
			r.Method,
			r.URL.Path,
			time.Since(start))

		app.logger.Info(logData)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mtx     sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mtx.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mtx.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.LimiterEnabled {
			ip := realip.FromRequest(r)

			mtx.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.config.LimiterRPS),
						app.config.LimiterBurst,
					),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mtx.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mtx.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// metricsResponseWriter wraps an existing http.ResponseWriter and also contains
// a field for recording the response status code, and a boolean flag to
// indicate whether the response headers have already been written.
type metricsResponseWriter struct {
	wrapped       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

// newMetricsResponseWriter returns a new metricsResponseWriter instance which
// wraps a given http.ResponseWriter and has a status code of 200 (which is the
// status code that Go will send in a HTTP response by default).
func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}

func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)
	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.wrapped.Write(b)
}

func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}

func (app *application) metrics(next http.Handler) http.Handler {
	// Initialize the new expvar variables when the middleware chain is first
	// built.
	var (
		totalRequestsReceived           = expvar.NewInt("total_requests_received")
		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
		totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")
	)

	// The following code will be run for every request.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the time that we started to process the request.
		start := time.Now()

		// Use the Add() method to increment the number of requests received by 1.
		totalRequestsReceived.Add(1)

		// Create a new metricsResponseWriter, which wraps the original
		// http.ResponseWriter value that the metrics middleware received.
		mw := newMetricsResponseWriter(w)

		// Call the next handler in the chain using the new metricsResponseWriter as
		// the http.ResponseWriter value.
		next.ServeHTTP(mw, r)

		// On the way back up the middleware chain, increment the number of
		// responses sent by 1.
		totalResponsesSent.Add(1)

		// Calculate the number of microseconds since we began to process the
		// request, then increment the total processing time by this amount.
		// At this point, the response status code should be stored in the
		// mw.statusCode field. Note that the expvar map is string-keyed, so we need
		// to use the strconv.Itoa() function to convert the status code (which is
		// an integer) to a string. Then we use the Add() method on our new
		// totalResponsesSentByStatus map to increment the count for the given
		// status code by 1.
		totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)

		// Calculate the number of microseconds since we began to process the
		// request, then increment the total processing time by this amount.
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")
		if origin != "" {
			for _, trustedOrigin := range app.config.CorsTrustedOrigins {
				if origin == trustedOrigin {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// Treat it as a preflight request.
					if r.Method == http.MethodOptions &&
						r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response headers.
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						// w.Header().Set("Content-Type", "text/event-stream")
						// w.Header().Set("Cache-Control", "no-cache")
						// w.Header().Set("Connection", "keep-alive")
						// Return from the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
