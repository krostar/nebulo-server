package handler

import (
	"net/http"

	"github.com/krostar/nebulo/router/httperror"
	"github.com/labstack/echo"
)

// UserInfos handle the route GET /user/.
// Return the whole profile of the logged user
/**
 * @api {get} /user Get profile infos
 * @apiDescription Return all the informations are stored about the users
 * @apiName User - Get profile infos
 * @apiGroup User
 *
 * @apiExample {curl} Usage example
 *		$>curl -X GET -v --cert bob.crt --key bob.key "https://api.nebulo.io/user/"
 *
 * @apiSuccess (Success) {nothing} 200 OK
 * @apiSuccessExample {json} Success example
 *		HTTP/1.1 200 "OK"
 *		{
 *			"key_fingerprint": "SHA256:r3nwvMos/pMuSuDLmWt0owQVUViqNw6Tn0mCZ0FLbUs",
 *			"display_name": "bob",
 *			"signup": "2017-03-11T15:47:54.153661099-08:00",
 *			"login_first": "2017-03-11T15:47:59.316160305-08:00",
 *			"login_last": "2017-03-11T15:48:12.89865226-08:00"
 *		}
 *
 * @apiError (Errors 4XX) {json} 401 Unauthorized: missing client certificate
 * @apiError (Errors 4XX) {json} 404 Not found: user not found
 * @apiError (Errors 5XX) {json} 500 Internal server error: server failed to handle the request
 */
func UserInfos(c echo.Context) (err error) {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	return c.JSONPretty(http.StatusOK, u, "    ")
}
