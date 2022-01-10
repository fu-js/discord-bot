package log

import "go.uber.org/zap"

var Zap *zap.SugaredLogger

func init() {
	log, _ := zap.NewDevelopment()
	Zap = log.Sugar()
}
