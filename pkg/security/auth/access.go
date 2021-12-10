package auth

// Access is the model of authorized access.
type Access struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// NewAccess create a new Access model.
func NewAccess(refreshToken string, accessToken *Token) (*Access, error) {
	access := Access{
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
	}

	var err error
	access.AccessToken, err = accessToken.SignedString()
	if err != nil {
		return nil, err
	}

	return &access, nil
}
