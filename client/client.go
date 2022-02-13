package client

import (
	"fmt"
	"path/filepath"

	"github.com/zelenin/go-tdlib/client"
)

type Client struct {
	tgClient     *client.Client
	masterChatID int64
}

func New(appID int32, apiHash string, masterChatID int64) (*Client, error) {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  appID,
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	logVerbosity := client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})

	tdlibClient, err := client.NewClient(authorizer, logVerbosity)
	if err != nil {
		return nil, fmt.Errorf("create tg client: %w", err)
	}

	return &Client{
		tgClient:     tdlibClient,
		masterChatID: masterChatID,
	}, nil
}
