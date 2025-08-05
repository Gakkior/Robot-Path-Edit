// Package logger 鎻愪緵缁熶竴鐨勬棩蹇楃鐞嗗姛鑳?
//
// 璁捐鍙傝€冿細
// - Kubernetes鐨勭粨鏋勫寲鏃ュ織绯荤粺
// - Grafana鐨勬棩蹇楃鐞?
// - Prometheus鐨勬棩蹇楄鑼?
//
// 鐗圭偣锛?
// 1. 缁撴瀯鍖栨棩蹇楋細鏀寔JSON鍜屾枃鏈牸寮?
// 2. 鏃ュ織杞浆锛氭敮鎸佹枃浠跺ぇ灏忓拰鏃堕棿杞浆
// 3. 涓婁笅鏂囨劅鐭ワ細鏀寔閾捐矾杩借釜
package logger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"robot-path-editor/internal/config"
)

// Init 鍒濆鍖栨棩蹇楃郴缁?
// 鍙傝€僈ubernetes鐨勬棩蹇楀垵濮嬪寲娴佺▼
func Init(cfg config.LoggerConfig) {
	// 璁剧疆鏃ュ織绾у埆
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.WithError(err).Warn("瑙ｆ瀽鏃ュ織绾у埆澶辫触锛屼娇鐢ㄩ粯璁ょ骇鍒?info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 璁剧疆鏃ュ織鏍煎紡
	switch strings.ToLower(cfg.Format) {
	case "json":
		// JSON鏍煎紡 - 閫傚悎鐢熶骇鐜鍜屾棩蹇楁敹闆嗙郴缁?
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	default:
		// 鏂囨湰鏍煎紡 - 閫傚悎寮€鍙戠幆澧?
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// 璁剧疆杈撳嚭鐩爣
	switch strings.ToLower(cfg.Output) {
	case "file":
		if cfg.File == "" {
			logrus.Warn("鏃ュ織鏂囦欢璺緞涓虹┖锛屼娇鐢ㄩ粯璁よ矾寰?)
			cfg.File = "./logs/app.log"
		}

		// 纭繚鏃ュ織鐩綍瀛樺湪
		logDir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logrus.WithError(err).Error("鍒涘缓鏃ュ織鐩綍澶辫触")
			return
		}

		// 鎵撳紑鏃ュ織鏂囦欢
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.WithError(err).Error("鎵撳紑鏃ュ織鏂囦欢澶辫触")
			return
		}

		logrus.SetOutput(file)
	default:
		// 榛樿杈撳嚭鍒版爣鍑嗚緭鍑?
		logrus.SetOutput(os.Stdout)
	}

	// 璁剧疆璋冪敤鑰呬俊鎭姤鍛?
	logrus.SetReportCaller(true)

	logrus.WithFields(logrus.Fields{
		"level":  cfg.Level,
		"format": cfg.Format,
		"output": cfg.Output,
	}).Info("鏃ュ織绯荤粺鍒濆鍖栧畬鎴?)
}
