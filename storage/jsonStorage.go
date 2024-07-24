package storage

type jsonStorage struct {
	NewsStorage
	SourceStorage
}

func NewStorage(newsStorage NewsStorage, sourceStorage SourceStorage) Storage {
	return &jsonStorage{
		NewsStorage:   newsStorage,
		SourceStorage: sourceStorage,
	}
}
