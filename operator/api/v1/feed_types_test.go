package v1

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetFeedCurrentCondition(t *testing.T) {
	tests := []struct {
		name       string
		conditions []Condition
		expected   Condition
	}{
		{
			name:       "No conditions present",
			conditions: []Condition{},
			expected:   Condition{},
		},
		{
			name: "Single condition present",
			conditions: []Condition{
				{
					Type:            ConditionAdded,
					Success:         true,
					Reason:          "Test",
					Message:         "Feed successfully added",
					LastUpdatedName: "test-user",
					LastUpdateTime:  metav1.Now(),
				},
			},
			expected: Condition{
				Type:            ConditionAdded,
				Success:         true,
				Reason:          "Test",
				Message:         "Feed successfully added",
				LastUpdatedName: "test-user",
			},
		},
		{
			name: "Multiple conditions present",
			conditions: []Condition{
				{
					Type:            ConditionAdded,
					Success:         true,
					Reason:          "Initial add",
					Message:         "Feed added",
					LastUpdatedName: "admin",
					LastUpdateTime:  metav1.Now(),
				},
				{
					Type:            ConditionDeleted,
					Success:         false,
					Reason:          "Not found",
					Message:         "Feed not found",
					LastUpdatedName: "test-user",
					LastUpdateTime:  metav1.Now(),
				},
			},
			expected: Condition{
				Type:            ConditionDeleted,
				Success:         false,
				Reason:          "Not found",
				Message:         "Feed not found",
				LastUpdatedName: "test-user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feedStatus := FeedStatus{
				Conditions: tt.conditions,
			}
			currentCondition := feedStatus.GetCurrentCondition()
			assert.Equal(t, tt.expected.Type, currentCondition.Type)
			assert.Equal(t, tt.expected.Success, currentCondition.Success)
			assert.Equal(t, tt.expected.Reason, currentCondition.Reason)
			assert.Equal(t, tt.expected.Message, currentCondition.Message)
			assert.Equal(t, tt.expected.LastUpdatedName, currentCondition.LastUpdatedName)
		})
	}
}

func TestAddCondition(t *testing.T) {
	feedStatus := FeedStatus{}
	initialTime := metav1.Now()

	condition := Condition{
		Type:            ConditionAdded,
		Success:         true,
		Reason:          "Test Reason",
		Message:         "Feed successfully added",
		LastUpdatedName: "test-user",
		LastUpdateTime:  initialTime,
	}

	feedStatus.SetCondition(condition)
	assert.Equal(t, len(feedStatus.Conditions), 1)

	latestCondition := feedStatus.GetCurrentCondition()

	assert.Equal(t, condition.Type, latestCondition.Type)
	assert.Equal(t, condition.Success, latestCondition.Success)
	assert.Equal(t, condition.Reason, latestCondition.Reason)
	assert.Equal(t, condition.Message, latestCondition.Message)
	assert.Equal(t, condition.LastUpdatedName, latestCondition.LastUpdatedName)
	assert.Equal(t, latestCondition.LastUpdateTime.Time, initialTime.Time)
}
