package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/verath/owbot-bot/owbot"
	"os"
	"os/signal"
	"time"
)

func main() {
	var (
		botId   string
		token   string
		logFile string
		dbFile  string
		debug   bool
	)
	flag.StringVar(&botId, "id", "", "The Bot ID of the bot")
	flag.StringVar(&token, "token", "", "The secret token for the bot")
	flag.StringVar(&logFile, "logfile", "", "A path to a file for writing logs (default is stdout)")
	flag.StringVar(&dbFile, "dbfile", "", "A path to a file to be used for bolt database")
	flag.BoolVar(&debug, "debug", false, "Set to true to log debug messages")
	flag.Parse()

	// TODO: This is not a great solution for required config...
	if botId == "" || token == "" {
		println("Both Bot ID and Bot Token is required")
		os.Exit(-1)
	}

	// Create a logrus instance (logger)
	logger := logrus.New()
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module":   "main",
				"filename": logFile,
				"error":    err,
			}).Fatal("Could not open file for logging")
		}
		logger.Formatter = &logrus.TextFormatter{ForceColors: true}
		logger.Out = f
	}

	// If debug is true, log debug messages
	if debug {
		logger.Level = logrus.DebugLevel
	}

	// Create a Bolt instance (database)
	var db *bolt.DB
	var err error
	if dbFile != "" {
		// Create the Bolt db
		db, err = bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 5 * time.Second})
		if err != nil {
			logger.WithFields(logrus.Fields{
				"module":   "main",
				"filename": dbFile,
				"error":    err,
			}).Fatal("Could not open file as bolt database")
		}
		defer func() {
			if err := db.Close(); err != nil {
				logger.WithField("error", err).Fatal("Could not close bolt db")
			}
		}()
	}

	bot, err := owbot.NewBot(logger, db, botId, token)
	if err != nil {
		logger.WithFields(logrus.Fields{"module": "main", "error": err}).Error("Could not creating bot")
		return
	}

	if err := bot.Start(); err != nil {
		logger.WithFields(logrus.Fields{"module": "main", "error": err}).Error("Could not start bot")
		return
	}

	// Run until asked to quit
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, os.Kill)
	<-interruptChan

	if err := bot.Stop(); err != nil {
		logger.WithFields(logrus.Fields{"module": "main", "error": err}).Error("Could not stop bot")
		return
	}
}
