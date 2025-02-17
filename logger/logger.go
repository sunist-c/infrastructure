package logger

import (
	"bytes"
	"encoding/json"
)

type Logger interface {
	Trace(field Fields)
	Debug(field Fields)
	Info(field Fields)
	Notice(field Fields)
	Warn(field Fields)
	Error(field Fields)
	Panic(field Fields)
	Log(level Level, field Fields)
}

type logger struct {
	level  Level
	writer FileWriter
}

func (l logger) Trace(field Fields) {
	l.Log(LevelTrace, field)
}

func (l logger) Debug(field Fields) {
	l.Log(LevelDebug, field)
}

func (l logger) Info(field Fields) {
	l.Log(LevelInfo, field)
}

func (l logger) Notice(field Fields) {
	l.Log(LevelNotice, field)
}

func (l logger) Warn(field Fields) {
	l.Log(LevelWarn, field)
}

func (l logger) Error(field Fields) {
	l.Log(LevelError, field)
}

func (l logger) Panic(field Fields) {
	l.Log(LevelPanic, field)
}

func (l logger) Log(level Level, field Fields) {
	if l.level.shouldLog(level) {
		bg := newBackground(level)
		go writeLog(l.writer, bg, field)
	}
}

func NewLogger(minLevel Level) Logger {
	return &logger{level: minLevel, writer: DefaultFileWriter}
}

func writeLog(writer FileWriter, bg background, field Fields) {
	logEntry := &entry{}
	field.export(logEntry)
	bg.export(logEntry)

	buffer := bytes.Buffer{}
	_ = json.NewEncoder(&buffer).Encode(logEntry)

	writer.WriteLog(buffer.Bytes())
}
