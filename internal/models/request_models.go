package models

type PgRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RedisIncrRequest struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type Hmac512Request struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}
