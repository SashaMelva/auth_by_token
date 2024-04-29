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

type RefreshToken struct {
	UserGUID     string `json:"userGUID"`
	RefreshToken []byte `json:"refreshToken"`
}
