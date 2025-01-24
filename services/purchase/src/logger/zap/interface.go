package loggerZap

import functionCallerInfo "github.com/TimDebug/FitByte/src/logger/helper"

type LoggerInterface interface {
	Info(msg string, function functionCallerInfo.FunctionCaller, data ...interface{})
	Error(msg string, function functionCallerInfo.FunctionCaller, data ...interface{})
	Debug(msg string, function functionCallerInfo.FunctionCaller, data ...interface{})
	Warn(msg string, function functionCallerInfo.FunctionCaller, data ...interface{})
}
