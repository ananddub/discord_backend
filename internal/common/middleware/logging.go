package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logs all gRPC requests
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Log request
		log.Printf("[gRPC] --> %s", info.FullMethod)

		// Call handler
		resp, err := handler(ctx, req)

		// Log response
		duration := time.Since(start)
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		log.Printf("[gRPC] <-- %s [%s] %v", info.FullMethod, statusCode, duration)

		return resp, err
	}
}

// StreamLoggingInterceptor logs streaming gRPC requests
func StreamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		log.Printf("[gRPC Stream] --> %s", info.FullMethod)

		err := handler(srv, ss)

		duration := time.Since(start)
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		log.Printf("[gRPC Stream] <-- %s [%s] %v", info.FullMethod, statusCode, duration)

		return err
	}
}
