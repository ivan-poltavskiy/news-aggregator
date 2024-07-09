package sorter

import (
	"news-aggregator/entity/news"
	"testing"
	"time"
)

func TestDateSorter_SortArticle(t *testing.T) {
	articles := []news.News{
		{Title: "News 1", Date: time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)},
		{Title: "News 2", Date: time.Date(2023, 7, 2, 0, 0, 0, 0, time.UTC)},
		{Title: "News 3", Date: time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)},
	}

	dateSorter := DateSorter{}

	tests := []struct {
		name      string
		sortBy    string
		want      []string
		expectErr bool
	}{
		{
			name:   "ascending",
			sortBy: "asc",
			want:   []string{"News 1", "News 2", "News 3"},
		},
		{
			name:   "descending",
			sortBy: "desc",
			want:   []string{"News 3", "News 2", "News 1"},
		},
		{
			name:      "invalid parameter",
			sortBy:    "invalid",
			want:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortedArticles, err := dateSorter.SortNews(articles, tt.sortBy)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if err != nil && tt.expectErr {
				expectedErr := "wrong sorting parameter: " + tt.sortBy
				if err.Error() != expectedErr {
					t.Fatalf("expected error %v, got %v", expectedErr, err)
				}
				return
			}

			for i, article := range sortedArticles {
				if string(article.Title) != tt.want[i] {
					t.Fatalf("expected news title %v, got %v", tt.want[i], article.Title)
				}
			}
		})
	}
}
