package oauth

type OauthError struct {
	Error   error
	Message string
	Code    int
}
