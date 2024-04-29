package model

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken []byte `json:"refreshToken"`
}
type TokenModel struct {
	UserGUID     string `json:"userGUID"`
	AccessToken  string `json:"accessToken"`
	RefreshToken []byte `json:"refreshToken"`
}

type Token struct {
	RefreshToken string `json:"refreshToken"`
}
