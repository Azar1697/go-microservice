package services

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IntegrationService struct {
	minioClient *minio.Client
	bucketName  string
}

// NewIntegrationService создает подключение к MinIO
func NewIntegrationService(endpoint, accessKeyID, secretAccessKey, bucketName string) *IntegrationService {
	// Инициализация клиента
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false, 
	})
	if err != nil {
		log.Fatalln("Failed to connect to MinIO:", err)
	}

	// Проверяем, существует ли бакет (папка), если нет — создаем
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Println("Error checking bucket:", err)
	} else if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Println("Failed to create bucket:", err)
		} else {
			log.Println("Successfully created bucket:", bucketName)
		}
	}

	return &IntegrationService{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

// UploadFile загружает файл в хранилище
func (s *IntegrationService) UploadFile(filename string, fileSize int64, reader io.Reader, contentType string) error {
	ctx := context.Background()

	// Загружаем поток данных
	info, err := s.minioClient.PutObject(ctx, s.bucketName, filename, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded %s of size %d\n", filename, info.Size)
	return nil
}