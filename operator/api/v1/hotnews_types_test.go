package v1

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestSetCondition(t *testing.T) {
	tests := []struct {
		name               string
		initialConditions  []Condition
		newCondition       Condition
		expectedConditions []Condition
	}{
		{
			name:              "Add new condition to empty status",
			initialConditions: []Condition{},
			newCondition: Condition{
				Type:            ConditionAdded,
				Success:         true,
				Reason:          "",
				Message:         "",
				LastUpdatedName: "test",
			},
			expectedConditions: []Condition{
				{
					Type:            ConditionAdded,
					Success:         true,
					Reason:          "",
					Message:         "",
					LastUpdatedName: "test",
				},
			},
		},
		{
			name: "Add new condition to non-empty status",
			initialConditions: []Condition{
				{
					Type:            ConditionAdded,
					Success:         true,
					Reason:          "",
					Message:         "",
					LastUpdatedName: "test",
				},
			},
			newCondition: Condition{
				Type:            ConditionAdded,
				Success:         true,
				Reason:          "",
				Message:         "",
				LastUpdatedName: "new-name",
			},
			expectedConditions: []Condition{
				{
					Type:            ConditionAdded,
					Success:         true,
					Reason:          "",
					Message:         "",
					LastUpdatedName: "new-name",
				},
			},
		},
		{
			name: "Update existing condition",
			initialConditions: []Condition{
				{
					Type:            ConditionAdded,
					Success:         true,
					Reason:          "",
					Message:         "",
					LastUpdatedName: "test",
				},
			},
			newCondition: Condition{
				Type:            ConditionAdded,
				Success:         false,
				Reason:          "Failed to add",
				Message:         "Feed addition failed",
				LastUpdatedName: "test",
			},
			expectedConditions: []Condition{
				{
					Type:            ConditionAdded,
					Success:         false,
					Reason:          "Failed to add",
					Message:         "Feed addition failed",
					LastUpdatedName: "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hotNewsStatus := &HotNewsStatus{
				Conditions: tt.initialConditions,
			}

			hotNewsStatus.SetCondition(tt.newCondition)

			assert.Equal(t, tt.expectedConditions, hotNewsStatus.Conditions)
		})
	}
}

func TestGetCurrentCondition(t *testing.T) {
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
			hotNewsStatus := HotNewsStatus{
				Conditions: tt.conditions,
			}
			currentCondition := hotNewsStatus.GetCurrentCondition()
			assert.Equal(t, tt.expected.Type, currentCondition.Type)
			assert.Equal(t, tt.expected.Success, currentCondition.Success)
			assert.Equal(t, tt.expected.Reason, currentCondition.Reason)
			assert.Equal(t, tt.expected.Message, currentCondition.Message)
			assert.Equal(t, tt.expected.LastUpdatedName, currentCondition.LastUpdatedName)
		})
	}
}
