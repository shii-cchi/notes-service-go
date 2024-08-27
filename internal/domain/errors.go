package domain

const (
	ErrUndefinedEnvParam = "parameter is undefined"
	ErrParsingAccessTTL  = "error parsing access ttl"
	ErrParsingRefreshTTL = "error parsing refresh ttl"
)

const (
	ErrGettingRefreshTokenFromCookie = "error getting refresh token from cookie"
	ErrCreatingRefreshToken          = "error creating refresh token"
	ErrCreatingAccessToken           = "error creating access token"
	ErrSavingRefreshToken            = "error saving refresh token to db"
	ErrGettingRefreshTokenFromDB     = "error getting refresh token by user id from db"
	ErrInvalidRefreshToken           = "invalid refresh token"
	ErrInvalidAccessToken            = "invalid access token"
	ErrRefreshTokenUndefined         = "refresh token is undefined"
	ErrAccessTokenUndefined          = "access token is undefined"
)

const (
	ErrParsingUserCredentialsInput = "error parsing user credentials"
	ErrInvalidUserCredentialsInput = "invalid user credentials input(login and password must be at least 6 characters long and can't be empty)"
	ErrCheckingUserExist           = "error checking user exist"
	ErrUserAlreadyExists           = "user with this login already exists"
	ErrHashingPassword             = "error hashing password"
	ErrCreatingUser                = "error creating user"
	ErrGettingPassword             = "error getting password by login from db"
	ErrUserNotFound                = "user does not exist"
	ErrWrongPassword               = "wrongPassword"
	ErrWrongCredentials            = "error wrong credentials(login or password)"
	ErrLogin                       = "login error"
	ErrLogout                      = "logout error"
	ErrRefresh                     = "refresh error"
)

const (
	ErrParsingNoteInput       = "error parsing note input"
	ErrInvalidNoteInput       = "invalid note input(both 'name' and 'content' fields are required and can't be empty)"
	ErrParsingID              = "error parsing id"
	ErrGettingNotes           = "error getting notes"
	ErrCreatingNote           = "error creating note"
	ErrCheckingSpellingErrors = "error checking spelling errors"
	ErrSpellingText           = "error spelling text"
)
