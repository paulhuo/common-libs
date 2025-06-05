package httpclient

import (
	"context"

	"github.com/google/uuid"
	"gitlab.yx/base-service/common-go/utils/logx"
	"go.opentelemetry.io/otel/trace"
)

func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return generateRequestID()
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return generateRequestID()
}

func logger(ctx context.Context) logx.Logger {
	if ctx == nil {
		return logx.WithContext(context.Background())
	}

	return logx.WithContext(ctx).WithFields(logx.LogField{Key: "action", Value: "httpCurl"})
}

func generateRequestID() string {
	return uuid.New().String()
}

// remove \n from the string
func removeNewline(s []byte) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' {
			result = append(result, s[i])
		}
	}
	return string(result)
}
