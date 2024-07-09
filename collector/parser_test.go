package collector

import (
	"news-aggregator/entity/source"
	"testing"
)

func TestGetParserBySourceType(t *testing.T) {
	parsers := InitParsers()

	tests := []struct {
		name        string
		sourceType  source.Type
		expectedErr bool
	}{
		{
			name:        "Test with existing RSS parser",
			sourceType:  source.RSS,
			expectedErr: false,
		},
		{
			name:        "Test with existing JSON parser",
			sourceType:  source.JSON,
			expectedErr: false,
		},
		{
			name:        "Test with existing UsaToday parser",
			sourceType:  source.UsaToday,
			expectedErr: false,
		},
		{
			name:        "Test with non-existent parser",
			sourceType:  "non-existent",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := parsers.GetParserBySourceType(tt.sourceType)
			if (err != nil) != tt.expectedErr {
				t.Errorf("GetParserBySourceType() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if !tt.expectedErr && parser == nil {
				t.Errorf("GetParserBySourceType() parser is nil, but expected a valid parser")
			}
		})
	}
}
