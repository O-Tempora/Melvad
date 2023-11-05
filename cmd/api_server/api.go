package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/O-Tempora/Melvad/internal/models"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}, err error) {
	w.WriteHeader(code)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		s.logger.LogAttrs(context.Background(), slog.LevelError, "Response with error:",
			slog.String("URL", r.URL.Path),
			slog.String("Method", r.Method),
			slog.Int("HTTP Code", code),
			slog.String("HTTP Status", http.StatusText(code)),
			slog.String("Error", err.Error()),
		)
		return
	}

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, "Response:",
		slog.String("URL", r.URL.Path),
		slog.String("Method", r.Method),
		slog.Int("HTTP Code", code),
		slog.String("HTTP Status", http.StatusText(code)),
	)
}

func (s *server) Router() {
	s.router = chi.NewMux()
	s.router.Post("/sign/hmacsha512", s.handleSignHmac)
	s.router.Post("/redis/incr", s.handleRedisIncr)
	s.router.Post("/postgres/users", s.handlePgInsert)
}

func (s *server) handlePgInsert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	model := &models.PgRequest{}
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		s.respond(w, r, 400, nil, err)
		return
	}
	id, err := s.service.InsertUser(model)
	if err != nil {
		s.respond(w, r, 500, nil, err)
		return
	}
	s.respond(w, r, 200, struct {
		Id int `json:"id"`
	}{
		Id: id,
	}, nil)
}
func (s *server) handleSignHmac(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	model := &models.Hmac512Request{}
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		s.respond(w, r, 400, nil, err)
		return
	}
	hx, err := s.service.SignHmac512(model)
	if err != nil {
		s.respond(w, r, 500, nil, err)
		return
	}
	s.respond(w, r, 200, hx, nil)
}
func (s *server) handleRedisIncr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	model := &models.RedisIncrRequest{}
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		s.respond(w, r, 400, nil, err)
		return
	}
	i, err := s.service.RedisIncrease(model)
	if err != nil {
		s.respond(w, r, 500, nil, err)
		return
	}
	s.respond(w, r, 200, struct {
		Value int64 `json:"value"`
	}{
		Value: i,
	}, nil)
}
