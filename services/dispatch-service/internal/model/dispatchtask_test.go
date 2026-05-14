package model

import (
	"testing"
)

func TestDispatchTask_TableName(t *testing.T) {
	task := DispatchTask{}
	if got := task.TableName(); got != "dispatch_tasks" {
		t.Errorf("TableName() = %s, want dispatch_tasks", got)
	}
}

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		constant string
		expected string
	}{
		{StatusPending, "pending"},
		{StatusActive, "active"},
		{StatusCompleted, "completed"},
		{StatusCancelled, "cancelled"},
	}
	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("status constant = %s, want %s", tt.constant, tt.expected)
		}
	}
}

func TestDispatchTask_Defaults(t *testing.T) {
	task := DispatchTask{}
	if task.Status != "" {
		t.Errorf("default Status should be empty, got %s", task.Status)
	}
}

func TestDispatchTask_GormTags(t *testing.T) {
	task := DispatchTask{ID: 1, VehicleID: 100, Status: StatusActive}
	if task.ID != 1 || task.VehicleID != 100 || task.Status != StatusActive {
		t.Error("DispatchTask struct fields not properly accessible")
	}
}
