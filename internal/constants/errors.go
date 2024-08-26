package constants

const (
	ErrUndefinedEnvParam = "parameter is undefined"
	ErrParsingAccessTTL  = "error parsing access ttl"
	ErrParsingRefreshTTL = "error parsing refresh ttl"
	ErrLoadingConfig     = "error loading config"
	ErrConnectingToDb    = "error connecting to db"
)

const (
	ErrParsingUserCredentials = "error parsing user credentials"
	ErrInvalidUserCredentials = "invalid user credentials"
	ErrCookieNotFound         = "cookie with name refresh_token not found"
	ErrGettingRefreshToken    = "error getting refresh token from cookie"
	ErrAccessTokenNotFound    = "access token not found"
)
