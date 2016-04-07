package models

import (
	"github.com/codegangsta/inject"
	"github.com/noaway/config"
	"log"
)

var (
	Injector   = inject.New()
	DBChanSend = make(chan []byte, 1024)
	DBChanRece = make(chan []byte, 1024)
	LoginChan  = make(chan []byte, 1024)
	RouteChan  = make(chan []byte, 1024)
)

func Config(group, key string) string {
	conf, err := config.ReadDefault("conf/web.config")
	if err != nil {
		log.Fatalln(err.Error())
		return ""
	}
	value, err := conf.String(group, key)
	if err != nil {
		log.Fatalln(err.Error())
		return ""
	}
	return value
}

func ConfigInt(group, key string) int {
	conf, err := config.ReadDefault("conf/web.config")
	if err != nil {
		log.Fatalln(err.Error())
		return 0
	}
	value, err := conf.Int(group, key)
	if err != nil {
		log.Fatalln(err.Error())
		return 0
	}
	return value
}
