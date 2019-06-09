package oauth

type oauthError struct {
	Error   error
	Message string
	Code    int
}
