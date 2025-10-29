package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CORSInterceptor handles CORS for gRPC-Web
func CORSInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Add CORS headers to response
		md := metadata.Pairs(
			"Access-Control-Allow-Origin", "*",
			"Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers", "Content-Type, Authorization",
			"Access-Control-Max-Age", "86400",
		)
		
		grpc.SetHeader(ctx, md)

		return handler(ctx, req)
	}
}
