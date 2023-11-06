package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/O-Tempora/Melvad/config"
	"github.com/O-Tempora/Melvad/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var serv Service

func TestMain(m *testing.M) {
	cf := &config.Config{}
	file, err := os.OpenFile("../../config/test_config.yaml", os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	dec := yaml.NewDecoder(file)
	if err = dec.Decode(cf); err != nil {
		log.Fatal(err.Error())
	}
	pgconn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cf.Database.User, cf.Database.Password, cf.Database.Host, cf.Database.Port, cf.Database.Name),
	)
	if err = pgconn.Ping(context.Background()); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	serv.Pgconn = pgconn

	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cf.Redis.Host, cf.Redis.Port),
		Password: cf.Redis.Password,
		DB:       cf.Redis.Db,
	})
	if status := conn.Ping(context.Background()); status.Err() != nil {
		log.Fatal(status.Err().Error())
	}
	serv.Rediscon = conn
}

func TestIncr(t *testing.T) {
	serv.Rediscon.IncrBy(context.Background(), "Test", 45)
	tests := []*struct {
		err      error
		Expected int64
		Actual   int64
		req      models.RedisIncrRequest
	}{
		{
			Expected: 55,
			req: models.RedisIncrRequest{
				Key:   "Test",
				Value: 10,
			},
		},
		{
			Expected: 55,
			req: models.RedisIncrRequest{
				Key:   "Test",
				Value: 0,
			},
		},
		{
			Expected: 17,
			req: models.RedisIncrRequest{
				Key:   "Bruh",
				Value: 17,
			},
		},
		{
			Expected: 20,
			req: models.RedisIncrRequest{
				Key:   "Bruh",
				Value: 3,
			},
		},
	}
	for _, tt := range tests {
		tt.Actual, tt.err = serv.RedisIncrease(&tt.req)
		if assert.NoError(t, tt.err) {
			assert.Equal(t, tt.Expected, tt.Actual)
		}
	}
}

func TestHmac(t *testing.T) {
	tests := []*struct {
		err      error
		Actual   string
		Expected string
		req      models.Hmac512Request
	}{
		{
			Expected: "70cf5c654a3335e493c263498b849b1aa425012914f8b5e77f4b7b7408ad207db9758f7c431887aa8f4885097e3bc032ee78238157c2ad43e900b69c60aee71e",
			req: models.Hmac512Request{
				Text: "",
				Key:  "1",
			},
		},
		{
			Expected: "9ba1f63365a6caf66e46348f43cdef956015bea997adeb06e69007ee3ff517df10fc5eb860da3d43b82c2a040c931119d2dfc6d08e253742293a868cc2d82015",
			req: models.Hmac512Request{
				Text: "test",
				Key:  "test",
			},
		},
		{
			Expected: "148ee5a9970f8de96d8b02bbd796fc6ca43657f1493c4a1ec6a17968c7e1573e7928bedbd15b49093b450fb82850b7cc16849459a3934be7d7f9f9e6ace4fe09",
			req: models.Hmac512Request{
				Text: "hmac",
				Key:  "",
			},
		},
	}
	for _, tt := range tests {
		tt.Actual, tt.err = serv.SignHmac512(&tt.req)
		if assert.NoError(t, tt.err) {
			assert.Equal(t, tt.Expected, tt.Actual)
		}
	}
}

func TestPgInsert(t *testing.T) {
	tests := []*struct {
		err      error
		Actual   int
		Expected int
		req      models.PgRequest
	}{
		{
			Expected: 1,
			req: models.PgRequest{
				Name: "Alex",
				Age:  21,
			},
		},
		{
			Expected: 2,
			req: models.PgRequest{
				Name: "Dima",
				Age:  4,
			},
		},
		{
			Expected: 0,
			err:      errors.New("ERROR: new row for relation \"users\" violates check constraint \"users_age_check\" (SQLSTATE 23514)"),
			req: models.PgRequest{
				Name: "Kolya",
				Age:  -1,
			},
		},
	}
	for _, tt := range tests {
		err := tt.err
		tt.Actual, tt.err = serv.InsertUser(&tt.req)
		if err != nil {
			assert.Error(t, tt.err)
		} else {
			assert.Equal(t, tt.Expected, tt.Actual)
		}
	}
}
