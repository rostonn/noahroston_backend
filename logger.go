package main

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func createLogger(profile string) (*zap.Logger, error) {
	rawJSON := []byte(`{
		"level": "debug",
		"development":true,
		"disableCaller":false,
		"encoding": "console",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "severity",
		  "levelEncoder": "uppercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	// 	s, ok := _levelToCapitalColorString[level]
	// 	if !ok {
	// 		s = _unknownLevelColor.Add(level.CapitalString())
	// 	}
	// 	enc.AppendString("[" + s + "]")
	// }

	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := cfg.Build()

	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return logger, nil

	// logger.Info("logger construction succeeded")
	// logger.Error("Logger error ...")

	// slogger := logger.Sugar()
	// slogger.Infow("Infow() allows tags", "name", "Legolas", "type", 1)
}
