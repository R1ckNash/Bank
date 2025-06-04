package slog_helper

import "log/slog"

// there is a tradeoff in naming packages like `sl`(slog)
//some people prefer short names, some people prefer long names
//in lib directory will be defined some useful function
//analog Utils in Java

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
