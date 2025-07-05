package cav

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"

// logger can be overridden by WithLogger option.
// If not set, it defaults to xlog.New.
// It is used for logging messages in the client.
var xlogger = xlog.New
