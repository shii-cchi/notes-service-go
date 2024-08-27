package domain

const (
	ErrUndefinedEnvParam = "parameter is undefined"
	ErrParsingAccessTTL  = "error parsing access ttl"
	ErrParsingRefreshTTL = "error parsing refresh ttl"
	ErrLoadingConfig     = "error loading config"
	ErrConnectingToDb    = "error connecting to db"
)

const (
	ErrParsingUserCredentials        = "error parsing user credentials"
	ErrInvalidUserCredentials        = "invalid user credentials(login and password must be at least 6 characters long and can't be empty)"
	ErrCookieNotFound                = "cookie with name refresh_token not found"
	ErrGettingRefreshTokenFromCookie = "error getting refresh token from cookie"
	ErrLogin                         = "login error"
)

const (
	ErrParsingNoteInput = "error parsing note input"
	ErrInvalidNoteInput = "invalid note input(both 'name' and 'content' fields are required and can't be empty)"
)

const (
	ErrCheckingUserExist         = "error checking user exist"
	ErrUserAlreadyExists         = "user with this login already exists"
	ErrHashingPassword           = "error hashing password"
	ErrCreatingRefreshToken      = "error creating refresh token"
	ErrCreatingUser              = "error creating user"
	ErrCreatingAccessToken       = "error creating access token"
	ErrSavingRefreshToken        = "error saving refresh token to db"
	ErrGettingPassword           = "error getting password by login from db"
	ErrUserNotFound              = "user does not exist"
	ErrWrongCredentials          = "error wrong credentials(login or password)"
	ErrLogout                    = "logout error"
	ErrGettingRefreshTokenFromDB = "error getting refresh token by user id from db"
	ErrInvalidRefreshToken       = "invalid refresh token"
	ErrRefresh                   = "refresh error"
)

const (
	ErrInvalidAccessToken     = "invalid access token"
	ErrParsingID              = "error parsing id"
	ErrGettingNotes           = "error getting notes"
	ErrCreatingNote           = "error creating note"
	ErrCheckingSpellingErrors = "error checking spelling errors"
	ErrSpellingText           = "error spelling text"
)
