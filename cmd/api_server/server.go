package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/O-Tempora/Melvad/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type server struct {
	logger   *slog.Logger
	router   *chi.Mux
	pgcon    *pgx.Conn
	rediscon *redis.Client
	service  *service.Service
}

func WithLogger(s *server) *server {
	s.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return s
}
func WithDbConn(s *server, cf *Config) *server {
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cf.Database.User, cf.Database.Password, cf.Database.Host, cf.Database.Port, cf.Database.Name),
	)

	if err = conn.Ping(context.Background()); err != nil {
		s.logger.Error("Log error", slog.String("Error: ", err.Error()))
		os.Exit(1)
	}
	s.pgcon = conn
	return s
}
func WithRedisConn(s *server, cf *Config) *server {
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cf.Redis.Host, cf.Redis.Port),
		Password: cf.Redis.Password,
		DB:       cf.Redis.Db,
	})
	if status := conn.Ping(context.Background()); status.Err() != nil {
		log.Fatal(status.Err().Error())
	}
	s.rediscon = conn
	return s
}

func StartServer(config *Config) error {
	s := &server{}
	s = WithDbConn(WithLogger(s), config)
	s.service = &service.Service{
		Pgconn:   s.pgcon,
		Rediscon: s.rediscon,
	}
	s.Router()
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), s)
}
