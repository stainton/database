package common

// Error code for http response
const (
	SUCCESS = 200

	// About database
	DB_CONNECT      = 10000
	DB_QUERY_ERROR  = 10001
	DB_RESULT_SCAN  = 10002
	DB_INSERT_ERROR = 10003
	DB_UPDATE_ERROR = 10004

	// About http
	INVALID_PARAMS = 11000
	INTERNAL       = 11001
)
