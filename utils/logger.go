package utils

import (
	"fmt"
	"time"
)

// Канал, в который мы будем кидать сообщения.
// Буфер 100 значит, что мы можем закинуть 100 сообщений, даже если их еще не успели напечатать.
var LogChannel = make(chan string, 100)

// StartLogger запускает "вечный цикл", который слушает канал
func StartLogger() {
	for msg := range LogChannel {
		// Печатаем сообщение, которое пришло из канала
		// fmt.Print(msg)
		_ = msg
	}
}

// Log отправляет сообщение в канал (не блокируя основной поток)
func LogUserAction(action string, userID int) {
	// Формируем строку
	msg := fmt.Sprintf("[AUDIT] Time: %s | Action: %s | UserID: %d\n",
		time.Now().Format("2006-01-02 15:04:05"), action, userID)
	
	// Кидаем в канал. select нужен, чтобы не зависнуть, если канал переполнен
	select {
	case LogChannel <- msg:
		// Успешно отправили
	default:
		// Канал полон, сообщение теряется (чтобы не тормозить сервер)
		fmt.Printf("Error: Log channel full, dropped log for user %d\n", userID)
	}
}