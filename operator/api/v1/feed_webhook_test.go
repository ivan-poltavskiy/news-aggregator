package v1

import (
	"testing"
)

func TestFeed_ValidateCreate(t *testing.T) {
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
			name: "Name too long",
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

//func TestFeed_validateFeed(t *testing.T) {
//
//	tests := []struct {
//		name          string
//		feed          Feed
//		expectedError bool
//	}{
//		{
//			name: "Valid feed",
//			feed: Feed{
//				Spec: FeedSpec{
//					Name: "valid-name",
//					Url:  "http://valid.url",
//				},
//			},
//			expectedError: false,
//		},
//		{
//			name: "Empty name",
//			feed: Feed{
//				Spec: FeedSpec{
//					Name: "",
//					Url:  "http://valid.url",
//				},
//			},
//			expectedError: true,
//		},
//		{
//			name: "Invalid URL",
//			feed: Feed{
//				Spec: FeedSpec{
//					Name: "valid-name",
//					Url:  "invalid-url",
//				},
//			},
//			expectedError: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := tt.feed.validateFeed()
//			if (err != nil) != tt.expectedError {
//				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
//			}
//		})
//	}
//}
