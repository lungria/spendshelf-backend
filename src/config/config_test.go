package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config, err := NewConfig()
	assert.Equal(t, nil, err)

	if httpAddr := os.Getenv("HTTP_ADDR"); httpAddr != "" {
		assert.Equal(t, httpAddr, config.HTTPAddr)
	} else {
		assert.Equal(t, ":8080", config.HTTPAddr)
	}

	if mongoUri := os.Getenv("MONGO_URI"); mongoUri != "" {
		assert.Equal(t, mongoUri, config.MongoURI)
	} else {
		assert.Equal(t, "mongodb://root:toor@localhost:27017", config.MongoURI)
	}

	if dbName := os.Getenv("SPEND_SHELF_DB"); dbName != "" {
		assert.Equal(t, dbName, config.DBName)
	} else {
		assert.Equal(t, "spendShelf", config.DBName)
	}

}
