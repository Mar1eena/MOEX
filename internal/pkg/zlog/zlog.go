package zlog

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type ZLogger interface {
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)

	Infof(msg string, args ...any)
	Fatalf(msg string, args ...any)
	Errorf(msg string, args ...any)
}

type zlog struct {
	zl zerolog.Logger
}

func New() *zlog {
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &zlog{zl: zl}
}

func (z *zlog) Trace(msg string, args ...any) {
	z.zl.Trace().Msg(msg)
}

func (z *zlog) Debug(msg string, args ...any) {
	z.zl.Debug().Msg(msg)
}

func (z *zlog) Info(msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	z.zl.Info().Msg(msg)
}

func (z *zlog) Warn(msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	z.zl.Warn().Msg(msg)
}

func (z *zlog) Error(msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	z.zl.Error().Msg(msg)
}

func (z *zlog) Fatal(msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	z.zl.Fatal().Msg(msg)
}

func (z *zlog) Panic(msg string, args ...any) {
	z.zl.Panic().Msg(msg)
}

func (z *zlog) Fatalf(msg string, args ...any) {
	z.zl.Fatal().Msg(msg)
}

func (z *zlog) Infof(msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	z.zl.Info().Msg(msg)
}

func (z *zlog) Errorf(msg string, args ...any) {
	z.zl.Error().Msg(msg)
}
