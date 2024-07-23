package service

import (
	"github.com/golang/mock/gomock"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	mock_newsStorage "news-aggregator/storage/news/mock_aggregator"
	"news-aggregator/storage/source/mock_aggregator"
	"testing"
)

func TestSaveNews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNewsStorage := mock_newsStorage.NewMockNewsStorage(ctrl)
	mockSourceStorage := mock_aggregator.NewMockStorage(ctrl)

	sourceEntity := source.Source{Name: "TestSource"}
	parsedNews := []news.News{
		{Title: news.Title("New Title 1")},
		{Title: news.Title("New Title 2")},
	}
	existingNews := []news.News{
		{Title: news.Title("Old Title")},
	}

	mockNewsStorage.EXPECT().GetNewsBySourceName(sourceEntity.Name, mockSourceStorage).Return(existingNews, nil)
	mockNewsStorage.EXPECT().SaveNews(sourceEntity, append(existingNews, parsedNews...)).Return(sourceEntity, nil)

	updatedSource, err := SaveNews(sourceEntity, mockNewsStorage, mockSourceStorage, parsedNews)
	if err != nil {
		t.Errorf("SaveNews() error = %v", err)
		return
	}

	if updatedSource != sourceEntity {
		t.Errorf("SaveNews() = %v, want %v", updatedSource, sourceEntity)
	}
}
