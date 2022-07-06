package router

import (
	"context"
	"fmt"
	"mts/task-service/internal/grpc/validateservice"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

func (rs Router) Validate(w http.ResponseWriter, r *http.Request) (string, error) {
	ref, err := r.Cookie("refreshToken")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error refreshToken cookie "))
		return "", err
	}

	grpcAddr := fmt.Sprintf("%s:%d", rs.Cfg.Grpc.Host, rs.Cfg.Grpc.Port)
	rs.Log.Info().Timestamp().Msg("Send validate")
	cwt, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		rs.Log.Fatal().Err(err)
		return "", err
	}
	defer conn.Close()

	os := validateservice.NewValidateClient(conn)

	newOrder := &validateservice.ValidateRequest{
		Token: ref.Value,
	}

	or, err := os.Validate(cwt, newOrder)
	if err != nil {
		rs.Log.Error().Err(err)
		return "", err
	}
	if len(or.Login) == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Login not auth"))
		return "", err
	}

	rs.Log.Warn().Timestamp().Msgf("Login: %s is validate", or.Login)
	return or.Login, nil
}
