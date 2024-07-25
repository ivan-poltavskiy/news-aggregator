package news

import (
	"github.com/golang/mock/gomock"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/storage/mock_aggregator"
	"testing"
)

func TestSaveNews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_aggregator.NewMockStorage(ctrl)

	sourceEntity := source.Source{Name: "TestSource"}
	parsedNews := []news.News{
		{Title: news.Title("New Title 1")},
		{Title: news.Title("New Title 2")},
	}
	existingNews := []news.News{
		{Title: news.Title("Old Title")},
	}

	mockStorage.EXPECT().GetNewsBySourceName(sourceEntity.Name, mockStorage).Return(existingNews, nil)
	mockStorage.EXPECT().SaveNews(sourceEntity, append(existingNews, parsedNews...)).Return(sourceEntity, nil)

	service := Service{
		storage: mockStorage,
	}

	updatedSource, err := service.SaveNews(sourceEntity, parsedNews)
	if err != nil {
		t.Errorf("SaveNews() error = %v", err)
		return
	}

	if updatedSource != sourceEntity {
		t.Errorf("SaveNews() = %v, want %v", updatedSource, sourceEntity)
	}
}
