package config

import (
	"encoding/json"
	"github.com/hunterhug/fafacms/core/util/kv"
	"github.com/hunterhug/fafacms/core/util/mail"
	"github.com/hunterhug/fafacms/core/util/oss"
	"github.com/hunterhug/fafacms/core/util/rdb"
)

var (
	//  Global config!
	FaFaConfig *Config
)

type Config struct {
	DefaultConfig MyConfig                   // default config
	OssConfig     oss.Key                    // oss like aws s3
	DbConfig      rdb.MyDbConfig             // mysql config
	SessionConfig kv.MyRedisConf             // redis config for user session
	MailConfig    mail.Sender `json:"Email"` // email config
}

// Some especial my config
type MyConfig struct {
	WebPort       string
	LogPath       string
	StoragePath   string
	LogDebug      bool
	StorageOss    bool
	CloseRegister bool
}

// Let the config struct to json file, just for test
func JsonOutConfig(config Config) (string, error) {
	raw, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	back := string(raw)
	return back, nil
}
