package v1

import (
	"context"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestValidateCreate(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)
	f := &FeedList{
		Items: []Feed{
			{Spec: FeedSpec{Name: "feed1"}},
		},
	}
	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(f).Build()

	r := &HotNews{
		Spec: HotNewsSpec{
			DateStart: "",
			DateEnd:   "",
			Keywords:  []string{"keyword1"},
			FeedsName: []string{"feed1"},
		},
	}

	var feedList FeedList
	err := k8sClient.List(context.Background(), &feedList, &client.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list feeds: %v", err)
	}

	_, err = r.ValidateCreate()
	assert.NoError(t, err)
}

func TestValidateUpdate(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)
	f := &FeedList{
		Items: []Feed{
			{Spec: FeedSpec{Name: "feed1"}},
		},
	}
	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(f).Build()

	r := &HotNews{
		Spec: HotNewsSpec{
			DateStart: "",
			DateEnd:   "",
			Keywords:  []string{"keyword1"},
			FeedsName: []string{"feed1"},
		},
	}
	var feedList FeedList
	err := k8sClient.List(context.Background(), &feedList, &client.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list feeds: %v", err)
	}

	_, err = r.ValidateUpdate(r)
	assert.NoError(t, err)
}

func TestValidateHotNews(t *testing.T) {

	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)
	f := &FeedList{
		Items: []Feed{
			{Spec: FeedSpec{Name: "feed1"}},
		},
	}
	k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(f).Build()

	tests := []struct {
		name      string
		hotNews   *HotNews
		expectErr bool
	}{
		{
			name: "valid hotnews",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					DateStart: "2023-01-01",
					DateEnd:   "2023-01-02",
					Keywords:  []string{"keyword1"},
					FeedsName: []string{"feed1"},
				},
			},
			expectErr: false,
		},
		{
			name: "missing keywords",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					DateStart: "2023-01-01",
					DateEnd:   "2023-01-02",
					Keywords:  []string{},
					FeedsName: []string{"feed1"},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid date range",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					DateStart: "2023-01-02",
					DateEnd:   "2023-01-01",
					Keywords:  []string{"keyword1"},
					FeedsName: []string{"feed1"},
				},
			},
			expectErr: true,
		},
	}
	var feedList FeedList
	err := k8sClient.List(context.Background(), &feedList, &client.ListOptions{})
	if err != nil {
		t.Fatalf("failed to list feeds: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.hotNews.validateHotNews()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFeeds(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = AddToScheme(scheme)

	tests := []struct {
		name      string
		hotNews   *HotNews
		feedList  *FeedList
		expectErr bool
		errorMsg  string
	}{
		{
			name: "valid feed",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					FeedsName: []string{"feed1"},
				},
			},
			feedList: &FeedList{
				Items: []Feed{
					{Spec: FeedSpec{Name: "feed1"}},
				},
			},
			expectErr: false,
		},
		{
			name: "feed does not exist",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					FeedsName: []string{"nonexistent_feed"},
				},
			},
			feedList: &FeedList{
				Items: []Feed{
					{Spec: FeedSpec{Name: "feed1"}},
				},
			},
			expectErr: true,
			errorMsg:  "validateFeeds: feed nonexistent_feed does not exist in namespace ",
		},
		{
			name: "one invalid feed",
			hotNews: &HotNews{
				Spec: HotNewsSpec{
					FeedsName: []string{"feed1", "invalid_feed"},
				},
			},
			feedList: &FeedList{
				Items: []Feed{
					{Spec: FeedSpec{Name: "feed1"}},
				},
			},
			expectErr: true,
			errorMsg:  "validateFeeds: feed invalid_feed does not exist in namespace ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient = fake.NewClientBuilder().WithScheme(scheme).WithLists(tt.feedList).Build()

			err := tt.hotNews.validateFeeds()

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
