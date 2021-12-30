package mymain

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugaredLogger *zap.SugaredLogger

// 初始化日誌
func initializeSugaredLogger() {

	// 设置一些基本日志格式
	encoder := zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			MessageKey:     `msg`,
			LevelKey:       `level`,
			TimeKey:        `time`,
			NameKey:        `logger`,
			CallerKey:      `file`,
			StacktraceKey:  `stacktrace`,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	)

	sugaredLogger =
		zap.New(
			zapcore.NewTee(
				getZapCoreUnderlyingCore(encoder, zapcore.InfoLevel),
				getZapCoreUnderlyingCore(encoder, zapcore.ErrorLevel),
			),
			zap.AddCaller(),
		).Sugar() // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数
}

// 取得zap core底下的核心
func getZapCoreUnderlyingCore(zapcoreEncoder zapcore.Encoder, zapcoreLevel zapcore.Level) zapcore.Core {
	return zapcore.NewCore(
		zapcoreEncoder,
		zapcore.AddSync(
			getWriter(
				getFilePathString(
					zapcoreLevel.String(),
				),
			),
		),
		zap.LevelEnablerFunc(
			func(level zapcore.Level) bool {
				return level >= zapcoreLevel
			},
		),
	)
}

// 取得檔案路徑字串
func getFilePathString(fileNameString string) string {
	return fmt.Sprintf(
		`%s/%s.%s`,
		`./logs/error`,
		fileNameString,
		`log`,
	)
}

// 取得日誌文件的io.Writer抽象
func getWriter(filePathString string) io.Writer {

	if hook, err :=
		rotatelogs.New(
			regexp.
				MustCompile(`(.+)(\.[^\.]+)$`).
				ReplaceAllString(
					filePathString,
					`$1-%Y%m%d$2`,
				), // 生成rotatelogs的Logger 实际生成的文件名 {level}-YYmmdd.log
			rotatelogs.WithMaxAge(time.Hour*24*7),     // 保存7天内的日志
			rotatelogs.WithRotationTime(time.Hour*24), // 每1天(整日)分割一次日志
		); err != nil {
		log.Panic(err)
		return nil
	} else {
		return hook
	}

}
