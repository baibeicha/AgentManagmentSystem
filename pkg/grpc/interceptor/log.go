package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	log := slog.Default()
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)

		statusCode := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Internal
			}
		}

		logAttrs := []any{
			slog.String("method", info.FullMethod),
			slog.String("code", statusCode.String()),
		}

		if err != nil {
			logAttrs = append(logAttrs, slog.String("err", err.Error()))

			if statusCode == codes.Internal || statusCode == codes.Unavailable {
				log.Error("gRPC request failed", logAttrs...)
			} else {
				log.Debug("gRPC request client error", logAttrs...)
			}
		} else {
			log.Debug("gRPC request completed", logAttrs...)
		}

		return resp, err
	}
}
