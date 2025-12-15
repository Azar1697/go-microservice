package utils

import (
	"fmt"
	"time"
)

var LogChannel = make(chan string, 100)

// StartLogger запускает "вечный цикл", который слушает канал
func StartLogger() {
	for msg := range LogChannel {
		fmt.Print(msg)
	}
}

// Log отправляет сообщение в канал (не блокируя основной поток)
func LogUserAction(action string, userID int) {
	msg := fmt.Sprintf("[AUDIT] Time: %s | Action: %s | UserID: %d\n",
		time.Now().Format("2006-01-02 15:04:05"), action, userID)
	
	// Кидаем в канал. select нужен, чтобы не зависнуть, если канал переполнен
	select {
	case LogChannel <- msg:
		
	default:
		
		fmt.Printf("Error: Log channel full, dropped log for user %d\n", userID)
	}
}