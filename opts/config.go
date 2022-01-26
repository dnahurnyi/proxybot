package opts

type Config struct {
	AppName      string `validate:"required"`
	AppID        int32  `validate:"required"`
	AppHash      string `validate:"required"`
	MasterChatID int64  `validate:"required"`
}
