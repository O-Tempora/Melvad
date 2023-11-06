package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"

	"github.com/O-Tempora/Melvad/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Pgconn   *pgx.Conn
	Rediscon *redis.Client
}

type IService interface {
	InsertUser(user *models.PgRequest) (int, error)
	RedisIncrease(m *models.RedisIncrRequest) (int64, error)
	SignHmac512(m *models.Hmac512Request) (string, error)
}

func (s *Service) InsertUser(user *models.PgRequest) (int, error) {
	var id int
	_, err := s.Pgconn.Exec(context.Background(),
		`create table if not exists users (
			id serial not null primary key,
			name text not null,
			age int not null check (age >= 0)
		)`)
	if err != nil {
		return -1, err
	}

	res := s.Pgconn.QueryRow(context.Background(),
		`insert into users
		(name, age) 
		values ($1, $2)
		on conflict do nothing
		returning id`,
		user.Name, user.Age,
	)
	if err = res.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (s *Service) RedisIncrease(m *models.RedisIncrRequest) (int64, error) {
	res := s.Rediscon.IncrBy(context.Background(), m.Key, m.Value)
	return res.Val(), res.Err()
}

func (s *Service) SignHmac512(m *models.Hmac512Request) (string, error) {
	hash := hmac.New(sha512.New, []byte(m.Key))
	_, err := hash.Write([]byte(m.Text))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
