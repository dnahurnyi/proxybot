package main

import (
	"fmt"
	"log"

	"github.com/dnahurnyi/proxybot/app"
	"github.com/dnahurnyi/proxybot/client"
	"github.com/dnahurnyi/proxybot/opts"
	"gopkg.in/go-playground/validator.v9"
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

	fmt.Println("Create updates handler")
	updatesHandler, err := app.NewUpdatesHandler(tgClient, config.MasterChatID)
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
