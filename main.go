package main

import (
	"fmt"
	"log"

	"github.com/dnahurnyi/proxybot/bot"
	"github.com/dnahurnyi/proxybot/client"
	"github.com/dnahurnyi/proxybot/opts"
	"github.com/dnahurnyi/proxybot/storage/postgres"
	"gopkg.in/go-playground/validator.v9"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config, err := opts.ReadOS()
	if err != nil {
		log.Fatal(fmt.Errorf("read configs: %w", err))
	}
	err = validator.New().Struct(config)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid config: %w", err))
	}

	tgClient, err := client.New(config.AppID, config.AppHash, config.MasterChatID)
	if err != nil {
		fmt.Println("Can't initiate tg client")
		log.Fatal(err)
	}

	repo, err := repositoryPG(config.DB)
	if err != nil {
		fmt.Println("Can't initiate postres repo")
		log.Fatal(err)
	}
	updatesHandler, err := bot.NewUpdatesHandler(tgClient, repo, config.MasterChatID)
	if err != nil {
		fmt.Println("Can't create updates handler")
		log.Fatal(err)
	}

	fmt.Println("Start listener")
	err = tgClient.Start(updatesHandler)
	if err != nil {
		fmt.Println("Listening updates failed")
		log.Fatal(err)
	}
}

func repositoryPG(config opts.DB) (*postgres.Repository, error) {
	dsn := postgresDSN(config.User, config.Password, config.Host, config.Port, config.DBName)
	db, err := gorm.Open(gorm_postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to posgres: %w", err)
	}

	return postgres.New(db)
}

func postgresDSN(user, password, host, port, database string) string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Kiev", host, user, password, database, port)
	return dsn
}
