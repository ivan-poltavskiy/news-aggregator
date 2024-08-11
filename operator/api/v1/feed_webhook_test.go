package v1

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestFeed_ValidateCreate(t *testing.T) {
	newScheme := runtime.NewScheme()
	_ = AddToScheme(newScheme)
	k8sClient = fake.NewClientBuilder().WithScheme(newScheme).Build()
	tests := []struct {
		name          string
		feed          Feed
		expectedError bool
	}{
		{
			name: "Valid feed",
			feed: Feed{
				Spec: FeedSpec{
					Name: "valid-name",
					Url:  "http://valid.url",
				},
			},
			expectedError: false,
		},
		{
			name: "Empty name",
			feed: Feed{
				Spec: FeedSpec{
					Name: "",
					Url:  "http://valid.url",
				},
			},
			expectedError: true,
		},
		{
			name: "Invalid URL",
			feed: Feed{
				Spec: FeedSpec{
					Name: "valid-name",
					Url:  "invalid-url",
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.feed.ValidateCreate()
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}

func TestFeed_ValidateUpdate(t *testing.T) {
	tests := []struct {
		name          string
		feed          Feed
		expectedError bool
	}{
		{
			name: "Valid feed",
			feed: Feed{
				Spec: FeedSpec{
					Name: "valid-name",
					Url:  "http://valid.url",
				},
			},
			expectedError: false,
		},
		{
			name: "NewName too long",
			feed: Feed{
				Spec: FeedSpec{
					Name: "this-name-is-way-too-long",
					Url:  "http://valid.url",
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFeed := &Feed{}
			_, err := tt.feed.ValidateUpdate(oldFeed)
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}

func TestFeed_ValidateDelete(t *testing.T) {
	tests := []struct {
		name          string
		feed          Feed
		expectedError bool
	}{
		{
			name: "Valid feed for delete",
			feed: Feed{
				Spec: FeedSpec{
					Name: "valid-name",
					Url:  "http://valid.url",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.feed.ValidateDelete()
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}

func TestCheckNameUnique(t *testing.T) {
	newScheme := runtime.NewScheme()
	_ = AddToScheme(newScheme)

	existingFeed := &Feed{
		Spec: FeedSpec{
			Name: "valid-name",
			Url:  "http://valid.url",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			UID:       "valid-uid",
		},
	}
	existingFeedsList := &FeedList{
		Items: []Feed{*existingFeed},
	}
	k8sClient = fake.NewClientBuilder().WithScheme(newScheme).WithLists(existingFeedsList).Build()

	tests := []struct {
		name      string
		feed      *Feed
		expectErr bool
	}{
		{
			name: "unique name",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "unique name",
					Url:  "http://valid.url",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: existingFeed.Namespace,
					UID:       "valid-uid",
				},
			},
			expectErr: false,
		},

		{
			name: "not unique name",
			feed: &Feed{
				Spec: FeedSpec{
					Name: "valid-name",
					Url:  "http://valid.url",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: existingFeed.Namespace,
					UID:       "test",
				},
			},
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkNameUniqueness(test.feed)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
