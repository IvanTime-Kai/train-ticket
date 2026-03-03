package model

// ─────────────────────────────────────────
// Request DTOs
// ─────────────────────────────────────────

type RegisterRequest struct {
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"     binding:"required"`
}

type LoginRequest struct {
	Email     string `json:"email"    binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Device    string `json:"device"`
	IPAddress string `json:"ip_address"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp"   binding:"required,len=6"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token"  binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ─────────────────────────────────────────
// Response DTOs
// ─────────────────────────────────────────

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	IsVerified  bool   `json:"is_verified"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}
