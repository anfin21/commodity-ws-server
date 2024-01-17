package logger

func DebugOutputLogger(
	logger *Logger,
) func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	return func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Infof(
			"%s\t %s\t --> %s\t (%d handlers)",
			httpMethod, absolutePath, handlerName, nuHandlers,
		)
	}
}
