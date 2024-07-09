package sorter

import (
	"news-aggregator/entity/article"
	"testing"
	"time"
)

func TestDateSorter_SortArticle(t *testing.T) {
	articles := []article.Article{
		{Title: "Article 1", Date: time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)},
		{Title: "Article 2", Date: time.Date(2023, 7, 2, 0, 0, 0, 0, time.UTC)},
		{Title: "Article 3", Date: time.Date(2023, 7, 3, 0, 0, 0, 0, time.UTC)},
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
			want:   []string{"Article 1", "Article 2", "Article 3"},
		},
		{
			name:   "descending",
			sortBy: "desc",
			want:   []string{"Article 3", "Article 2", "Article 1"},
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
			sortedArticles, err := dateSorter.SortArticle(articles, tt.sortBy)
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
					t.Fatalf("expected article title %v, got %v", tt.want[i], article.Title)
				}
			}
		})
	}
}
