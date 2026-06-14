package dto

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type RoleWithUsersResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Users     []UserResponse `json:"users"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}
