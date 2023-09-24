package config

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
)

var (
	Cfg *ini.File
	err error
)

func init() {
	// 读取配置信息
	Cfg, err = ini.Load("./settings.ini")
	if err != nil {
		fmt.Println("加载配置文件失败！", err)
		panic("读取配置信息失败")
	}
}

func InitLogger() {
	// 1.encoder
	encoder := zap.NewProductionEncoderConfig()
	encoder.TimeKey = "STRFTIME"                                    // 修改时间Key的名称
	encoder.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime) // 修改时间格式化样式
	encoder.LevelKey = "LEVEL"                                      // 修改level的名称
	encoder.EncodeLevel = zapcore.CapitalLevelEncoder               // level大写

	// 2.writesyncer
	logFile, _ := os.OpenFile("./logs/zaplog.logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// 即想写入到文件又想写入到终端时：（同时把日志打印在多个文件中内）
	// zapcore.AddSync(logFile)
	// zapcore.AddSync(os.Stdout)
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(logFile), zapcore.AddSync(os.Stdout))
	// 如果要求把报错级别的日志单独存放在一个日志文件中，普通级别的日志存在另一个日志文件中
	errLogFile, _ := os.OpenFile("./logs/zaplog.err.logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	errWriteSyncer := zapcore.AddSync(errLogFile)

	// 3.loglevel
	// zapcore.DebugLevel\InfoLevel\WarnLevel\ErrorLevel\PanicLevel\FataLevel
	// 从配置文件读取或从命令行获取了当前项目的所处的环境
	level, err := zapcore.ParseLevel(Cfg.Section("app").Key("level").String())
	if err != nil {
		level = zapcore.InfoLevel
	}

	// zapcore.LevelEnabler.Enabled()
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		multiWriteSyncer,
		level,
	)
	// 如果要求把报错级别的日志单独存放在一个日志文件中，普通级别的日志存在另一个日志文件中
	errZapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		errWriteSyncer,
		zap.ErrorLevel,
	)

	// 如果要求把报错级别的日志单独存放在一个日志文件中，普通级别的日志存在另一个日志文件中
	cusLogger := zap.New(zapcore.NewTee(zapCore, errZapCore), zap.AddCaller())

	zap.ReplaceGlobals(cusLogger)
}
