package config

import (
	"crypto/tls"
	"discord/pkg/valkeys"
	"log"
	"net/url"

	valkey "github.com/valkey-io/valkey-go"
)

var client *valkey.Client

func InitValKey() error {
	if client != nil {
		return nil
	}
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(cfg.Database.Redis.URL)
	if err != nil {
		log.Fatal(err)
	}
	password, _ := u.User.Password()

	lclient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{u.Host},
		Username:    u.User.Username(),
		Password:    password,
		TLSConfig:   &tls.Config{},
	})
	if err != nil {
		panic(err)
	}
	client = &lclient
	return nil
}

func GetValKeyClient() *valkey.Client {
	err := InitValKey()
	if err != nil {
		return nil
	}
	return client
}

var myValkeyClient *valkeys.Valkeys

func GetMyValkeyClient() *valkeys.Valkeys {
	if myValkeyClient == nil {
		myValkeyClient = valkeys.NewValkeysClient(client)
	}
	return myValkeyClient
}
