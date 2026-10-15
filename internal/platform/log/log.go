package log

import (
	logrus "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/internal/platform/config"
)

func InitLog() {
	cfg := config.Config()
	if level, err := logrus.ParseLevel(cfg.Log.Level); err == nil {
		logrus.SetLevel(level)
	}
}
