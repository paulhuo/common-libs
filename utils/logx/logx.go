package logx

import (
	"github.com/zeromicro/go-zero/core/logx"
)

// 以下为 logx 所有直接导出函数的再导出，以完全兼容 go-zero 的用法

var (
	// Redefine type aliases for compatibility
	Field        = logx.Field
	WithColor    = logx.WithColor
	WithDuration = logx.WithDuration
	WithFields   = logx.WithFields
	NewLogger    = logx.NewLogger
	NewWriter    = logx.NewWriter
	Reset        = logx.Reset
	SetWriter    = logx.SetWriter
	ErrorStack   = logx.ErrorStack
	Info         = logx.Info
	Infof        = logx.Infof
	Infov        = logx.Infov
	Error        = logx.Error
	Errorf       = logx.Errorf
	Errorv       = logx.Errorv
	Debug        = logx.Debug
	Debugf       = logx.Debugf
	Debugv       = logx.Debugv
	Slow         = logx.Slow
	Slowf        = logx.Slowf
	Slowv        = logx.Slowv
	WithContext  = logx.WithContext
	MustSetup    = logx.MustSetup
	SetLevel     = logx.SetLevel
	Close        = logx.Close

	Severe  = logx.Severe
	Severef = logx.Severef
)

type LogField = logx.LogField
type Logger = logx.Logger
