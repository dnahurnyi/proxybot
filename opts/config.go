package opts

type Config struct {
	AppName      string `validate:"required"`
	AppID        int32  `validate:"required"`
	AppHash      string `validate:"required"`
	MasterChatID int64  `validate:"required"`
	DB           DB     `validate:"required"`
}

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSL      bool
}
