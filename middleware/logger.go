package middleware

import (
	"fmt"
	"net"
	"net/http"
	
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// LoggerConfig is a struct that holds the configuration for the Logger middleware.
type LoggerConfig struct {
	IncludeRemoteAddr bool
	IncludeRemoteIP   bool
	IncludeUserAgent  bool
	IncludeReferer    bool
	IncludeRequestID  bool
	IncludeEndpoint   bool
	IncludeMethod     bool
	IncludeProto      bool
	IncludeVersion    bool
	IncludeHost       bool
}

// DefaultLoggerConfig returns a LoggerConfig with default values.
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		IncludeRemoteIP:  true,
		IncludeUserAgent: true,
		IncludeRequestID: true,
		IncludeEndpoint:  true,
		IncludeMethod:    true,
		// Skip less useful ones by default
		IncludeRemoteAddr: false,
		IncludeReferer:    false,
		IncludeProto:      false,
		IncludeVersion:    false,
		IncludeHost:       false,
	}
}

// LoggerWithConfig is a middleware that loads a logger into the request's ctx,
// along with the data specified in the LoggerConfig parameter.
func LoggerWithConfig(logger zerolog.Logger, config LoggerConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logContext := logger.With()
			
			if config.IncludeRemoteAddr {
				logContext = logContext.Str("addr", r.RemoteAddr)
			}
			
			if config.IncludeRemoteIP {
				if r.RemoteAddr != "" {
					ip, _, err := net.SplitHostPort(r.RemoteAddr)
					if err == nil {
						logContext = logContext.Str("remote_ip", ip)
					}
				}
			}
			
			if config.IncludeUserAgent {
				logContext = logContext.Str("user_agent", r.UserAgent())
			}
			
			if config.IncludeReferer {
				logContext = logContext.Str("referer", r.Referer())
			}
			
			if config.IncludeRequestID {
				requestID := r.Header.Get("Request-Id")
				if requestID == "" {
					id, ok := hlog.IDFromRequest(r)
					if !ok {
						id = xid.New()
					}
					// Always ensure the ID is in the context for downstream handlers
					ctx := hlog.CtxWithID(r.Context(), id)
					r = r.WithContext(ctx)
					requestID = id.String()
				}
				logContext = logContext.Str("req_id", requestID)
			}
			
			if config.IncludeEndpoint {
				logContext = logContext.Str("endpoint", r.URL.String())
			}
			
			if config.IncludeMethod {
				logContext = logContext.Str("method", r.Method)
			}
			
			if config.IncludeProto {
				logContext = logContext.Str("proto", r.Proto)
			}
			
			if config.IncludeVersion {
				version := fmt.Sprintf("%d.%d", r.ProtoMajor, r.ProtoMinor)
				logContext = logContext.Str("version", version)
			}
			
			if config.IncludeHost {
				logContext = logContext.Str("host", r.Host)
			}
			
			enrichedLogger := logContext.Logger()
			ctx := enrichedLogger.WithContext(r.Context())
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logger is a middleware that loads a logger into the request's ctx, along
// with some useful data:
// 	- RemoteIP
// 	- UserAgent
// 	- RequestID
// 	- Method
//  - URL
//
// When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white.
//
// N.B.: it uses rs/zerolog logger
func Logger(logger zerolog.Logger) Middleware {
	return LoggerWithConfig(logger, DefaultLoggerConfig())
}
