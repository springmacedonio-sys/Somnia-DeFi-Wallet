package models

type AuthLoginRequest struct {
	AuthProvider   string `json:"auth_provider"`
	AuthExternalID string `json:"auth_external_id"`
}

type AuthLoginResponse struct {
	PrivateUserInfo PrivateUserInfo `json:"private_user_info"`
}

type AuthRegisterRequest struct {
	WalletName      string `json:"wallet_name"`
	AuthProvider    string `json:"auth_provider"`
	AuthExternalID  string `json:"auth_external_id"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
}

type AuthRegisterResponse struct {
	PrivateUserInfo PrivateUserInfo `json:"private_user_info"`
}
