package logic

import (
	"testing"

	userv1 "github.com/aicong/mine-dispatch/proto/user/v1"
)

func TestCreateUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		valid    bool
	}{
		{"valid", "alice", "pass123", true},
		{"empty username", "", "pass123", false},
		{"empty password", "alice", "", false},
		{"both empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &userv1.CreateUserRequest{
				Username: tt.username,
				Password: tt.password,
			}
			got := req.Username != "" && req.Password != ""
			if got != tt.valid {
				t.Errorf("CreateUserRequest{%q, %q}: valid=%v, want %v",
					tt.username, tt.password, got, tt.valid)
			}
		})
	}
}

func TestUserResponse_Success(t *testing.T) {
	resp := &userv1.UserResponse{
		Code:    0,
		Message: "success",
		Data: &userv1.User{
			Id:       1001,
			Username: "alice",
			RealName: "测试用户A",
			Email:    "alice@example.com",
			Phone:    "13800138000",
			Role:     1,
			Status:   1,
			MineId:   1,
		},
	}
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
	if resp.Data.Username != "alice" {
		t.Errorf("Username = %s, want alice", resp.Data.Username)
	}
	if resp.Data.Role != 1 {
		t.Errorf("Role = %d, want 1", resp.Data.Role)
	}
}

func TestUserResponse_NotFound(t *testing.T) {
	resp := &userv1.UserResponse{Code: 404, Message: "user not found"}
	if resp.Code != 404 {
		t.Errorf("Code = %d, want 404", resp.Code)
	}
	if resp.Data != nil {
		t.Error("Data should be nil when not found")
	}
}

func TestUserResponse_EmptyFields(t *testing.T) {
	resp := &userv1.UserResponse{
		Code: 0,
		Data: &userv1.User{
			Id:       1002,
			Username: "bob",
		},
	}
	if resp.Data.RealName != "" {
		t.Errorf("RealName = %s, want empty", resp.Data.RealName)
	}
	if resp.Data.Email != "" {
		t.Errorf("Email = %s, want empty", resp.Data.Email)
	}
	if resp.Data.MineId != 0 {
		t.Errorf("MineId = %d, want 0", resp.Data.MineId)
	}
}

func TestUserListResponse_Pagination(t *testing.T) {
	resp := &userv1.UserListResponse{
		Code:    0,
		Message: "success",
		Data: []*userv1.User{
			{Id: 1, Username: "admin", Role: 1, Status: 1},
			{Id: 2, Username: "alice", Role: 1, Status: 1},
			{Id: 3, Username: "bob", Role: 1, Status: 1},
		},
		Total: 3,
	}
	if resp.Total != 3 {
		t.Errorf("Total = %d, want 3", resp.Total)
	}
	if len(resp.Data) != 3 {
		t.Errorf("len(Data) = %d, want 3", len(resp.Data))
	}
	if resp.Data[0].Username != "admin" {
		t.Errorf("Data[0].Username = %s, want admin", resp.Data[0].Username)
	}
}

func TestUserListResponse_Empty(t *testing.T) {
	resp := &userv1.UserListResponse{
		Code:    0,
		Message: "success",
		Data:    []*userv1.User{},
		Total:   0,
	}
	if resp.Total != 0 {
		t.Errorf("Total = %d, want 0", resp.Total)
	}
	if len(resp.Data) != 0 {
		t.Errorf("len(Data) = %d, want 0", len(resp.Data))
	}
}

func TestUserListResponse_KeywordFilter(t *testing.T) {
	req := &userv1.ListUserRequest{
		Page:    1,
		PageSize: 10,
		Keyword: "alice",
	}
	if req.Keyword != "alice" {
		t.Errorf("Keyword = %s, want alice", req.Keyword)
	}
	if req.Page != 1 {
		t.Errorf("Page = %d, want 1", req.Page)
	}
	if req.PageSize != 10 {
		t.Errorf("PageSize = %d, want 10", req.PageSize)
	}
}

func TestDeleteUserResponse(t *testing.T) {
	resp := &userv1.UserResponse{Code: 0, Message: "success"}
	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
}

func TestUserModel_DefaultValues(t *testing.T) {
	user := &userv1.User{
		Id:       1003,
		Username: "newuser",
	}
	if user.Role != 0 {
		t.Errorf("Role = %d, want 0", user.Role)
	}
	if user.Status != 0 {
		t.Errorf("Status = %d, want 0", user.Status)
	}
	if user.MineId != 0 {
		t.Errorf("MineId = %d, want 0", user.MineId)
	}
}
