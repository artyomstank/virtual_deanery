package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger представляет собой обёртку над zap.Logger для структурированного
// логирования с постоянными полями и уровнями.
type Logger struct {
	logger *zap.Logger
}

// New создаёт новый экземпляр Logger с заданным уровнем логирования.
// Параметр level принимает значения "debug", "info", "warn", "error".
// При неизвестном уровне по умолчанию используется "info".
// В режиме "debug" вывод производится в читаемом текстовом формате через консоль,
// в остальных случаях — в JSON на stdout. Каждое сообщение снабжается меткой
// времени в формате RFC3339.
func New(level string) *Logger {
	parsedLevel := parseLevel(level)

	var cfg zap.Config
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	// Единый формат времени для всех режимов
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(parsedLevel)
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := cfg.Build()
	if err != nil {
		// Резервный вариант — без дополнительных зависимостей
		panic("failed to build logger: " + err.Error())
	}

	return &Logger{logger: logger}
}

// parseLevel преобразует строку уровня в zapcore.Level.
// При ошибке возвращает InfoLevel.
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug записывает отладочное сообщение с дополнительными полями.
func (l *Logger) Debug(msg string, fields map[string]any) {
	l.logger.Debug(msg, mapToFields(fields)...)
}

// Info записывает информационное сообщение с дополнительными полями.
func (l *Logger) Info(msg string, fields map[string]any) {
	l.logger.Info(msg, mapToFields(fields)...)
}

// Warn записывает предупреждение с дополнительными полями.
func (l *Logger) Warn(msg string, fields map[string]any) {
	l.logger.Warn(msg, mapToFields(fields)...)
}

// Error записывает сообщение об ошибке. Если err не равен nil, в лог добавляется
// поле "error" со значением err.Error(). Дополнительные поля из fields
// также добавляются при их наличии.
func (l *Logger) Error(msg string, err error, fields map[string]any) {
	zapFields := mapToFields(fields)
	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}
	l.logger.Error(msg, zapFields...)
}

// With возвращает новый Logger с добавленным постоянным полем key=value.
// Метод удобен для привязки контекста, например:
//
//	logger.With("service", "user").Info("запрос обработан", nil)
func (l *Logger) With(key string, value any) *Logger {
	return &Logger{logger: l.logger.With(zap.Any(key, value))}
}

// mapToFields преобразует map в срез zap.Field.
// Если fields == nil, возвращает пустой срез, чтобы избежать паники.
func mapToFields(fields map[string]any) []zap.Field {
	if fields == nil {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
