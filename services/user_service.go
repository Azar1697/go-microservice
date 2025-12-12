package services

import (
	"errors"
	"go-microservice/models"
	"sync"
)

// UserService хранит данные в памяти
type UserService struct {
	users  map[int]models.User // Сама база данных в памяти
	nextID int                 // Счетчик для ID
	mu     sync.RWMutex        // Мьютекс для защиты от гонки данных (Race Condition)
}

// NewUserService создает новый экземпляр сервиса
func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]models.User),
		nextID: 1,
	}
}

// Create добавляет пользователя (Thread-safe)
func (s *UserService) Create(user models.User) models.User {
	s.mu.Lock()         // БЛОКИРУЕМ доступ для всех остальных
	defer s.mu.Unlock() // Разблокируем, когда функция закончится

	user.ID = s.nextID
	s.nextID++
	s.users[user.ID] = user
	return user
}

// GetAll возвращает всех пользователей
func (s *UserService) GetAll() []models.User {
	s.mu.RLock()         // Блокируем ТОЛЬКО на запись (читать можно многим одновременно)
	defer s.mu.RUnlock()

	var allUsers []models.User
	for _, user := range s.users {
		allUsers = append(allUsers, user)
	}
	return allUsers
}

// GetByID ищет пользователя по ID
func (s *UserService) GetByID(id int) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// Update обновляет данные
func (s *UserService) Update(id int, updatedUser models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.users[id]
	if !exists {
		return models.User{}, errors.New("user not found")
	}

	// Сохраняем ID, чтобы он не изменился случайно
	updatedUser.ID = id
	s.users[id] = updatedUser
	return updatedUser, nil
}

// Delete удаляет пользователя
func (s *UserService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
	return nil
}