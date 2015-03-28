package log

import golog "log"

const (
	DEBUG_LEVEL int = 0
	INFO_LEVEL  int = 1
	WARN_LEVEL  int = 2
	FATAL_LEVEL int = 3
)

var logLevel = WARN_LEVEL

func SetLevel(level int) {
	logLevel = level
}

func Debug(messages ...interface{}) {
	if logLevel <= DEBUG_LEVEL {
		golog.SetPrefix("[DEBUG]")
		golog.Println(messages)
	}
}

func Info(messages ...interface{}) {
	if logLevel <= INFO_LEVEL {
		golog.SetPrefix("[INFO]")
		golog.Println(messages)
	}
}

func Warn(messages ...interface{}) {
	if logLevel <= WARN_LEVEL {
		golog.SetPrefix("[WARN]")
		golog.Println(messages)
	}
}

func Fatal(messages ...interface{}) {
	if logLevel <= FATAL_LEVEL {
		golog.SetPrefix("[FATAL]")
		golog.Fatalln(messages)
	}
}
