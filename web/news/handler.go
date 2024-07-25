package news

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/client"
	"news-aggregator/storage"
)

type HandlerForNews struct {
	service *Service
}

// NewNewsHandler returns the new instance of the news handler
func NewNewsHandler(storage storage.Storage) *HandlerForNews {
	return &HandlerForNews{
		service: NewService(storage),
	}
}

// FetchNewsHandler handles requests for fetching news.
func (h *HandlerForNews) FetchNewsHandler(w http.ResponseWriter, client client.Client) {

	news, err := client.FetchNews()
	if err != nil {
		logrus.Error("Failed to fetch news ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client.Print(news)
}
