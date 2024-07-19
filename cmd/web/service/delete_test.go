package service_test

import (
	"errors"
	"news-aggregator/cmd/web/service"
	"news-aggregator/storage/mock_aggregator"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSourceByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_aggregator.NewMockStorage(ctrl)

	tests := []struct {
		name       string
		sourceName string
		mockFunc   func()
		expectErr  bool
	}{
		{
			name:       "Success",
			sourceName: "example-source",
			mockFunc: func() {
				mockStorage.EXPECT().DeleteSourceByName("example-source").Return(nil)
			},
			expectErr: false,
		},
		{
			name:       "Failure",
			sourceName: "non-existent-source",
			mockFunc: func() {
				mockStorage.EXPECT().DeleteSourceByName("non-existent-source").Return(errors.New("delete error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := service.DeleteSourceByName(tt.sourceName, mockStorage)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
