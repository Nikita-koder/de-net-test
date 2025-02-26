package handlers

import (
	"de-net/libs/4_common/safe_go"
)

// Канал для работы с запросами
var handlersRequests = make(chan interface{})

// Обработчик очереди задач
// FIXME: тут стоит задуматься об универсиализации, код повторяется
func Worker() {
	for query := range handlersRequests {
		switch q := query.(type) {
		case LoginQuery:
			safe_go.SafeGo(q.r.logger, func() {
				token, status := login(q.r.logger, q.r.Login, q.r.Password)

				q.responseChan <- LoginResponse{
					token:  token,
					status: status,
				}
			})

		}
	}
}
