package middleware

import (
	"context"
	"strings"

	"discord/internal/auth/util"
	"discord/internal/common/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor validates JWT tokens
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if isPublicEndpoint(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		userID, err := util.ValidateJWT(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		ctx = context.WithValue(ctx, "user_id", userID)

		return handler(ctx, req)
	}
}

func StreamAuthInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if isPublicEndpoint(info.FullMethod) {
			return handler(srv, ss)
		}

		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		userID, err := util.ValidateJWT(token)
		if err != nil {
			return status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		newCtx := context.WithValue(ss.Context(), "user_id", userID)

		wrapped := &WrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		return handler(srv, wrapped)
	}
}

// isPublicEndpoint checks if endpoint requires authentication
func isPublicEndpoint(method string) bool {
	publicEndpoints := []string{
		"/service.auth.AuthService/Register",
		"/service.auth.AuthService/Login",
		"/service.auth.AuthService/ForgotPassword",
		"/service.auth.AuthService/ResetPassword",
	}

	for _, endpoint := range publicEndpoints {
		if method == endpoint {
			return true
		}
	}
	return false
}

// WrappedServerStream wraps grpc.ServerStream with a custom context
type WrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *WrappedServerStream) Context() context.Context {
	return w.ctx
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (int32, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return 0, errors.ErrUnauthorized
	}
	return userID, nil
}
