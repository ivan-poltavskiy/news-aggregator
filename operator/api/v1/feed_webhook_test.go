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
		errorMessage  string
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
			errorMessage:  "name cannot be empty",
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
			errorMessage:  "URL must be a valid URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.feed.ValidateCreate()
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeed_ValidateUpdate(t *testing.T) {
	tests := []struct {
		name          string
		feed          Feed
		expectedError bool
		errorMessage  string
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
			name: "New name too long",
			feed: Feed{
				Spec: FeedSpec{
					Name: "this-name-is-way-too-long-for-feed",
					Url:  "http://valid.url",
				},
			},
			expectedError: true,
			errorMessage:  "name must not exceed 20 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldFeed := &Feed{}
			_, err := tt.feed.ValidateUpdate(oldFeed)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
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
		name          string
		feed          *Feed
		expectErr     bool
		expectedError string
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
					UID:       "test-uid",
				},
			},
			expectErr:     true,
			expectedError: "a Feed with name 'valid-name' already exists in namespace 'default'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkNameUniqueness(test.feed)
			if test.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
