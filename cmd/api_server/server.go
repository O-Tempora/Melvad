package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
)

type server struct {
	logger   *slog.Logger
	router   *chi.Router
	pgcon    *pgx.Conn
	rediscon *redis.Client
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
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cf.Db.User, cf.Db.Password, cf.Db.Host, cf.Db.Port, cf.Db.Name),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
	s.pgcon = conn
	return s
}
func WithRedisConn(s *server, cf *Config) *server {
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cf.Rc.Host, cf.Rc.Port),
		Password: cf.Rc.Password,
		DB:       cf.Rc.Db,
	})
	if status := conn.Ping(context.Background()); status.Err() != nil {
		log.Fatal(status.Err().Error())
	}
	s.rediscon = conn
	return s
}
