package lg

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
}
