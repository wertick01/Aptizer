package logger

type Config struct {
	// Уровень логов
	// Возможные значения: debug, info, warn, error, panic, fatal
	Level string

	// Конфиг zap для подмены стандартного.
	// Если заполнен, то Level игнорируется
	CustomRawZapCfg string

	// Версия приложения в формате major.minor.patch
	Version string

	// Наименование объекта, на котором развёрнуто приложение (сервер, под и т.д.)
	Host string

	// Наименование микросервиса, приложения  и т.д.
	Context string
}
