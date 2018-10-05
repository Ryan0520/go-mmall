package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_NOT_EXIST_USERNAME                  = 10001
	ERROR_CHECK_USERNAME                      = 10002
	ERROR_QUERY_USER_BY_USERNAME_AND_PASSWORD = 10003
	ERROR_PASSWORD_NOT_CORRECT                = 10004
	ERROR_EXIST_USERNAME                      = 10005
	ERROR_EXIST_PHONE                         = 10006
	ERROR_CHECK_PHONE                         = 10007
	ERROR_EXIST_EMAIL                         = 10008
	ERROR_CHECK_EMAIL                         = 10009
	ERROR_SAVE_USER                           = 10010
	ERROR_QUERY_QUESTION                      = 10011
	ERROR_NOT_EXIST_QUESTION                  = 10012
	ERROR_CHECK_QUESTION_ANSWER               = 10013
	ERROR_QUESTION_ANSWER_NOT_CORRECT         = 10014
	ERROR_USER_NOT_LOGIN                      = 10015
	ERROR_QUERY_USER                          = 10016
	ERROR_USER_NOT_EXIST                      = 10017
	ERROR_RESET_PASSWORD                      = 10018
	ERROR_UPDATE_USER                         = 10019
	ERROR_FORGET_RESET_PASSWORD_TOKEN         = 10020
	ERROR_NOT_ADMIN                           = 10021
	ERROR_SAVE_CATEGORY                       = 10022
	ERROR_CHECK_CATEGORY_BY_PARENT_ID         = 10023
	ERROR_NOT_EXIST_CATEGORY                  = 10024
	ERROR_GET_CATEGORY                        = 10025
	ERROR_UPDATE_CATEGORY                     = 10026
)
