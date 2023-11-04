package service

import (
	"context"

	"github.com/O-Tempora/Melvad/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Pgconn   *pgx.Conn
	Rediscon *redis.Client
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
		returning id`,
		user.Name, user.Age,
	)
	if err = res.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}
