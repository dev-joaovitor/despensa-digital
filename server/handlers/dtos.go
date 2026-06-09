package handlers

// auth
type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type SendRecoveryCodeDTO struct {
	Email string `json:"email"`
}

type VerifyRecoveryCodeDTO struct {
	Email string `json:"email"`
	Code string `json:"code"`
}

type ChangePasswordDTO struct {
	NewPassword string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

// users
type CreateUserDTO struct {
	FullName string `json:"full_name"`
	Email string `json:"email"`
	Password string `json:"password"`

	InvitationCode *string `json:"invitation_code"`
	HouseholdName *string `json:"household_name"`
}

type UpdateUserDTO struct {
	FullName string `json:"full_name"`
	Email string `json:"email"`

	Code string `json:"code"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}
