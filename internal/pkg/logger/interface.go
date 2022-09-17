package logger

type ILogger interface {
	Debug(*Fields)
	Info(*Fields)
	Warn(*Fields)
	Error(*Fields)
	Fatal(*Fields)
	Panic(*Fields)
	Sync()
}
