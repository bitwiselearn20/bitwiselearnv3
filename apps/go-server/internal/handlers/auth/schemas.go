package auth

// Request bodies ported from schemas/auth.py (CamelModel -> camelCase JSON).

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type forgotPasswordRequest struct {
	Email string  `json:"email"`
	Role  *string `json:"role"`
}

type verifyOtpRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type resetPasswordRequest struct {
	NewPassword string  `json:"newPassword"`
	Role        *string `json:"role"`
}
