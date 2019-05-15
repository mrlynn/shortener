package storage

import "github.com/Dimitriy14/shortener/models"

type Storage interface {
	GetURL(code string) (string, error)
	SaveUrl(url string) (string, error)
	GetInfo() ([]models.Shortener, error)
}

var storage Storage

func SetStorage(s Storage) {
	storage = s
}

func SaveUrl(url string) (string, error) {
	return storage.SaveUrl(url)
}

func GetURL(code string) (string, error) {
	return storage.GetURL(code)
}

func GetInfo() ([]models.Shortener, error) {
	return storage.GetInfo()
}
