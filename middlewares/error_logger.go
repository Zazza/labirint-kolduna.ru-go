package middlewares

import (
	"bytes"
	"gamebook-backend/config"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ResponseWriterWrapper оборачивает gin.ResponseWriter для перехвата тела ответа
type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// ErrorLoggerMiddleware логирует все запросы с не-200 статусом
func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Читаем тело запроса
		requestBody, err := io.ReadAll(ctx.Request.Body)
		if err == nil && len(requestBody) > 0 {
			// Восстанавливаем тело для дальнейшей обработки
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Оборачиваем ResponseWriter для перехвата тела ответа
		w := &ResponseWriterWrapper{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
		}
		ctx.Writer = w

		// Вызываем следующий handler
		ctx.Next()

		// Получаем статус код ответа
		statusCode := ctx.Writer.Status()

		// Логируем только ошибки (не 2xx)
		if statusCode < 200 || statusCode >= 300 {
			// Получаем информацию о запросе
			method := ctx.Request.Method
			path := ctx.Request.URL.Path
			query := ctx.Request.URL.RawQuery
			clientIP := ctx.ClientIP()
			userAgent := ctx.Request.UserAgent()

			// Получаем текст ошибки из тела ответа
			errorBody := strings.TrimSpace(w.Body.String())

			// Формируем сообщение об ошибке
			errorMsg := "HTTP Error: " + method + " " + path
			if query != "" {
				errorMsg += "?" + query
			}
			errorMsg += " | Status: " + strconv.Itoa(statusCode)
			errorMsg += " | IP: " + clientIP
			errorMsg += " | UA: " + userAgent

			// Добавляем текст ошибки
			if errorBody != "" {
				if len(errorBody) > 500 {
					errorBody = errorBody[:500] + "..."
				}
				errorMsg += " | Response Error: " + errorBody
			}

			// Логируем ошибку
			config.LogError(errorMsg)
		}
	}
}
