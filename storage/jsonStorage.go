package storage

type jsonStorage struct {
	News
	Source
}

// NewStorage returns the new instance of the Storage interface
func NewStorage(newsStorage News, sourceStorage Source) Storage {
	return &jsonStorage{
		News:   newsStorage,
		Source: sourceStorage,
	}
}
