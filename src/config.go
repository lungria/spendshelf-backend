package main

import (
	"os"
)

type Config struct {
	HTTTAddr string
	MongoURI string
	DBName   string
}

func NewConfig() *Config {
	c := Config{}
	defaultAddr := ":80"
	defaultDBName := "SpendShelf"
	defaultMongoURI := "mongodb://root:toor@localhost:27017"

	envAddr := os.Getenv("WEB_HOOK_ADDR")
	envDBName := os.Getenv("SPEND_SHELF_DB")
	envMongoURI := os.Getenv("MONGO_URI")

	if envAddr != "" {
		c.HTTTAddr = envAddr
	} else {
		c.HTTTAddr = defaultAddr
	}

	if envDBName != "" {
		c.DBName = envDBName
	} else {
		c.DBName = defaultDBName

	}

	if envMongoURI != "" {
		c.MongoURI = envMongoURI
	} else {
		c.MongoURI = defaultMongoURI
	}

	return &c
}
