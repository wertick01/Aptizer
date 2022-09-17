package logger

type Fields struct {
	*Detail
	UserID  string
	Code    string
	Message string
}
