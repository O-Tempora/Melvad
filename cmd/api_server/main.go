package main

import (
	"flag"
	"log"
	"os"

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
	cf := &Config{}
	file, err := os.OpenFile(defaulConfig, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	dec := yaml.NewDecoder(file)
	if err = dec.Decode(cf); err != nil {
		log.Fatal(err.Error())
	}

	flag.Parse()
	cf.Rc.Host = redis_host
	cf.Rc.Port = redis_port
}
