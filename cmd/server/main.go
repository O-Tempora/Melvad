package main

import (
	"flag"
	"log"
	"os"

	apiserver "github.com/O-Tempora/Melvad/cmd/api_server"
	"github.com/O-Tempora/Melvad/config"
	"gopkg.in/yaml.v3"
)

var (
	redis_host string
	redis_port int
)

const (
	defaulConfig = "config/default.yaml"
)

func init() {
	flag.StringVar(&redis_host, "host", "localhost", "Redis host")
	flag.IntVar(&redis_port, "port", 5554, "Redis port")
}
func main() {
	cf := &config.Config{}
	file, err := os.OpenFile(defaulConfig, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	dec := yaml.NewDecoder(file)
	if err = dec.Decode(cf); err != nil {
		log.Fatal(err.Error())
	}

	flag.Parse()
	cf.Redis.Host = redis_host
	cf.Redis.Port = redis_port

	if err = apiserver.StartServer(cf); err != nil {
		log.Fatal(err.Error())
	}
}
