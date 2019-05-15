package storage

type Storage interface {
	GetURL(code string) (string, error)
	SaveUrl(url string) (string, error)
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
