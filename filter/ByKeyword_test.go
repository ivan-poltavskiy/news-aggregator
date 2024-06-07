package filter

import (
	"NewsAggregator/entity/article"
	"reflect"
	"testing"
)

func TestByKeyword_Filter(t *testing.T) {
	type fields struct {
		Keywords []string
	}
	type args struct {
		articles []article.Article
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []article.Article
	}{
		{name: "Find by one keyword in the title",
			fields: fields{
				Keywords: []string{
					"ukr"},
			},
			args: args{
				articles: []article.Article{
					{Title: "Article 1"},
					{Title: "ukranian"},
					{Title: "Ukraine"},
					{Title: "ukr"},
					{Title: "someWord"},
					{Title: "someukrWord"},
				},
			},
			want: []article.Article{
				{Title: "ukranian"},
				{Title: "Ukraine"},
				{Title: "ukr"},
				{Title: "someukrWord"},
			},
		},

		{name: "Find by one keyword in the description",
			fields: fields{
				Keywords: []string{
					"ukr"},
			},
			args: args{
				articles: []article.Article{
					{Description: "Article 1"},
					{Description: "ukranian"},
					{Description: "Ukraine"},
					{Description: "ukr"},
					{Description: "someWord"},
					{Description: "someukrWord"},
				},
			},
			want: []article.Article{
				{Description: "ukranian"},
				{Description: "Ukraine"},
				{Description: "ukr"},
				{Description: "someukrWord"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := ByKeyword{
				Keywords: tt.fields.Keywords,
			}
			if got := f.Filter(tt.args.articles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actual result: = %v Expexted: %v", got, tt.want)
			}
		})
	}
}
