package filter

import (
	"news-aggregator/entity/news"
	"reflect"
	"testing"
)

func TestByKeyword_Filter(t *testing.T) {
	type fields struct {
		Keywords []string
	}
	type args struct {
		articles []news.News
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []news.News
	}{
		{name: "Find by one keyword in the title",
			fields: fields{
				Keywords: []string{
					"ukr"},
			},
			args: args{
				articles: []news.News{
					{Title: "News 1"},
					{Title: "ukranian"},
					{Title: "Ukraine"},
					{Title: "ukr"},
					{Title: "someWord"},
					{Title: "someukrWord"},
				},
			},
			want: []news.News{
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
				articles: []news.News{
					{Description: "News 1"},
					{Description: "ukranian"},
					{Description: "Ukraine"},
					{Description: "ukr"},
					{Description: "someWord"},
					{Description: "someukrWord"},
				},
			},
			want: []news.News{
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
