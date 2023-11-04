package main

type Config struct {
	Host     string      `yaml:"host"`
	Port     int         `yaml:"port"`
	Redis    redisConfig `yaml:"redis_config"`
	Database dbConfig    `yaml:"db_config"`
}
type redisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`
}
type dbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
}
