package storage

import "github.com/Dimitriy14/shortener/models"

type Storage interface {
	GetURL(code string) (string, error)
	SaveUrl(url string) (string, error)
	GetInfo() ([]models.Shortener, error)
}

var (
	repository Storage
)

func SetStorage(s Storage) {
	repository = s
}

func SaveUrl(url string) (string, error) {
	return repository.SaveUrl(url)
}

func GetURL(code string) (string, error) {
	return repository.GetURL(code)
}

func GetInfo() ([]models.Shortener, error) {
	return repository.GetInfo()
}
