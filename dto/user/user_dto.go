package userDto

type NewUserDto struct {
	Name     string `json:"name" form:"name" validate:"required,min=3,max=60"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=60"`
	Role     string `json:"role" form:"role" validate:"required"`
	TenantID string `json:"tenant_id" form:"tenant_id" validate:"required"`
}

type ChangePasswordDto struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required,min=8,max=60"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required,min=8,max=60"`
}

type UserLogin struct {
	Email    string `json:"email" form:"email" validate:"required,email,min=5,max=60"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=60"`
}
