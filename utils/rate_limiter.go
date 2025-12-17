package utils

import (
	"net/http"
	"golang.org/x/time/rate"
)

// СОЗДАЕМ ЛИМИТЕР ЗДЕСЬ (Глобально). 
// Теперь он точно один на всё приложение и не пересоздается.
var globalLimiter = rate.NewLimiter(1000, 50) 

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// Добавим отладочный принт (уберешь перед сдачей)
		// Если в логах docker этого нет — значит middleware вообще не подключена!
		// fmt.Println("Checking limit...") 

		if !globalLimiter.Allow() {
			// fmt.Println("Limit exceeded!") // Чтобы видеть отказы в логах
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}