package logger

import "github.com/fatih/color"

type levelFormatter struct {
	ColorType string
	Type      string
}

var defaultFormatters = map[Level]*levelFormatter{
	DebugLevel: {
		Type:      "[DBUG]",
		ColorType: color.HiBlueString("%s", "[DBUG]"),
	},
	InfoLevel: {
		Type:      "[INFO]",
		ColorType: color.GreenString("[INFO]"),
	},
	WarnLevel: {
		Type:      "[WARN]",
		ColorType: color.HiYellowString("[WARN]"),
	},
	ErrorLevel: {
		Type:      "[ERRO]",
		ColorType: color.HiRedString("[ERRO]"),
	},
}
