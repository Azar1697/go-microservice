package handlers

import (
	"encoding/json"
	"net/http"
	
	"go-microservice/services"
)

type IntegrationHandler struct {
	service *services.IntegrationService
}

func NewIntegrationHandler(service *services.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{service: service}
}

// UploadFile - POST /api/upload
func (h *IntegrationHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Ограничиваем размер загрузки (например, 10 МБ)
	r.ParseMultipartForm(10 << 20)

	// Получаем файл из формы (ключ "file")
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Вызываем сервис для загрузки в MinIO
	err = h.service.UploadFile(header.Filename, header.Size, file, header.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "Error uploading to MinIO", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully", "filename": header.Filename})
}