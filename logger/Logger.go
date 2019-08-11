package logger

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"github.com/vortgo/ma-parser/utils"
	"gopkg.in/sohlich/elogrus.v7"
	"os"
)

type Data map[string]interface{}

type Context struct {
	BandId     int
	AlbumId    int
	SongId     int
	Collection string
}

type Logger struct {
	logrus  *logrus.Logger
	context Context
	data    map[string]interface{}
}

func New() *Logger {
	log := logrus.New()
	client, err := elastic.NewClient(elastic.SetHttpClient(utils.CustomHttpClient), elastic.SetSniff(false), elastic.SetHealthcheck(false), elastic.SetURL(os.Getenv("ELASTIC_URL")))
	if err != nil {
		log.Panic(err)
	}
	hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "parser-log")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)
	logger := Logger{logrus: log}
	return &logger
}

func (logger *Logger) SetContext(context Context) *Logger {
	logger.context = context
	return logger
}

func (logger *Logger) SetData(data map[string]interface{}) *Logger {
	logger.data = data
	return logger
}

func (logger *Logger) reset() {
	logger.context = Context{}
	logger.data = map[string]interface{}{}
}

func (logger *Logger) prepare() *logrus.Entry {
	data, _ := json.Marshal(logger.data)

	return logger.logrus.WithFields(logrus.Fields{
		"band_id":    logger.context.BandId,
		"album_id":   logger.context.AlbumId,
		"song_id":    logger.context.SongId,
		"collection": logger.context.Collection,
		"data":       string(data),
	})
}

func (logger *Logger) Trace(args ...interface{}) {
	entry := logger.prepare()
	entry.Trace(args)
	logger.reset()
}

func (logger *Logger) Debug(args ...interface{}) {
	entry := logger.prepare()
	entry.Debug(args)
	logger.reset()
}

func (logger *Logger) Print(args ...interface{}) {
	entry := logger.prepare()
	entry.Print(args)
	logger.reset()
}

func (logger *Logger) Info(args ...interface{}) {
	entry := logger.prepare()
	entry.Info(args)
	logger.reset()
}

func (logger *Logger) Warn(args ...interface{}) {
	entry := logger.prepare()
	entry.Warn(args)
	logger.reset()
}

func (logger *Logger) Warning(args ...interface{}) {
	entry := logger.prepare()
	entry.Warning(args)
	logger.reset()
}

func (logger *Logger) Error(args ...interface{}) {
	entry := logger.prepare()
	entry.Error(args)
	logger.reset()
}

func (logger *Logger) Fatal(args ...interface{}) {
	entry := logger.prepare()
	entry.Fatal(args)
	logger.reset()
	entry.Logger.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	entry := logger.prepare()
	entry.Panic(args)
	logger.reset()
	panic(fmt.Sprint(args...))
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Tracef(format, args...)
	logger.reset()
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Debugf(format, args...)
	logger.reset()
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Infof(format, args...)
	logger.reset()
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Printf(format, args...)
	logger.reset()
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Warnf(format, args...)
	logger.reset()
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Warningf(format, args...)
	logger.reset()
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Errorf(format, args...)
	logger.reset()
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Fatalf(format, args...)
	logger.reset()
	entry.Logger.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	entry := logger.prepare()
	entry.Panicf(format, args...)
	logger.reset()
}

func (logger *Logger) Traceln(args ...interface{}) {
	entry := logger.prepare()
	entry.Traceln(args...)
	logger.reset()
}

func (logger *Logger) Debugln(args ...interface{}) {
	entry := logger.prepare()
	entry.Debugln(args...)
	logger.reset()
}

func (logger *Logger) Infoln(args ...interface{}) {
	entry := logger.prepare()
	entry.Infoln(args...)
	logger.reset()
}

func (logger *Logger) Println(args ...interface{}) {
	entry := logger.prepare()
	entry.Println(args...)
	logger.reset()
}

func (logger *Logger) Warnln(args ...interface{}) {
	entry := logger.prepare()
	entry.Warnln(args...)
	logger.reset()
}

func (logger *Logger) Warningln(args ...interface{}) {
	entry := logger.prepare()
	entry.Warningln(args...)
	logger.reset()
}

func (logger *Logger) Errorln(args ...interface{}) {
	entry := logger.prepare()
	entry.Errorln(args...)
	logger.reset()
}

func (logger *Logger) Fatalln(args ...interface{}) {
	entry := logger.prepare()
	entry.Fatalln(args...)
	logger.reset()
	entry.Logger.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	entry := logger.prepare()
	entry.Panicln(args...)
	logger.reset()
}
