package flog

import (
	"fmt"
	"github.com/hunterhug/fafacms/core/util/log"
	"os"
	"path/filepath"
	"strings"
)

// Our diy log config
var jsconf = `
{
  "UseShortFile": true,
  "Appenders": {
    "console": {
      "Type": "console"
    },
    "base": {
      "Type": "dailyfile",
      "Target": "%s"
    }
  },
  "Loggers": {
    "baseLogger": {
      "Appenders": [
        "console",
        "base"
      ],
      "Level": "NOTICE"
    },
    "otherLogger": {
      "Appenders": [
        "console"
      ],
      "Level": "NOTICE"
    }
  },
  "Root": {
    "Level": "debug",
    "Appenders": [
      "console"
    ]
  }
}
 `

var Log *log.Logger

// Init the log
func InitLog(logFile string) {
	os.MkdirAll(filepath.Dir(logFile), 0777)
	err := log.Init(fmt.Sprintf(jsconf, logFile))
	if err != nil {
		panic("log error:" + err.Error())
	}

	Log = log.Get("baseLogger")
}

// Set the log level
func SetLogLevel(level string) {
	if num, ok := log.LogLevelMap[strings.ToUpper(level)]; ok {
		Log.SetLevel(num)
	} else {
		panic("no this level")
	}
}
