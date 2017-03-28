package handler

import (
	"fmt"
	"net/http"

	"github.com/krostar/nebulo/router/httperror"
	"github.com/labstack/echo"

	cvalidator "github.com/krostar/nebulo/tools/validator"
	up "github.com/krostar/nebulo/user/provider"
	validator "gopkg.in/validator.v2"
)

// UserEdit handle the route PUT /user/.
// Perform a profile modification and return the whole profile of the logged user
/**
 * @api {put} /user Update user profile
 * @apiDescription Perform a profile modification and return the whole profile of the logged user.
 * You can't revoke and issue a new certificate from this call. The only updatable field is: "display_name".
 * @apiName User - update profile infos
 * @apiGroup User
 *
 * @apiExample {curl} Usage example
 *		$>curl -X PUT -v --cert bob.crt --key bob.key "https://api.nebulo.io/user/" --data "{\"display_name\": \"bob\"}"
 *
 * @apiSuccess (Success) {nothing} 200 OK
 * @apiSuccessExample {json} Success example
 *		HTTP/1.1 200 "OK"
 *		{
 *			"key_fingerprint": "SHA256:r3nwvMos/pMuSuDLmWt0owQVUViqNw6Tn0mCZ0FLbUs",
 *			"display_name": "Bob",
 *			"signup": "2017-03-30T05:21:45.787568658-07:00",
 *			"login_first": "2017-03-30T05:22:33.543332295-07:00",
 *			"login_last": "2017-03-30T05:34:29.527686645-07:00"
 *		}
 *
 * @apiError (Errors 4XX) {json} 400 Bad request: bad json input
 * @apiError (Errors 4XX) {json} 401 Unauthorized: missing client certificate
 * @apiError (Errors 4XX) {json} 404 Not found: user not found
 * @apiError (Errors 5XX) {json} 500 Internal server error: server failed to handle the request
 */
func UserEdit(c echo.Context) (err error) {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	if err = c.Bind(u); err != nil {
		return httperror.HTTPBadRequestError(fmt.Errorf("unable to bind json to user: %v", err))
	}

	if err = validator.WithTag("validator-update").Validate(u); err != nil {
		return cvalidator.HTTPErrors(err)
	}

	if err = up.P.Update(u, map[string]interface{}{
		"display_name": u.DisplayName,
	}); err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to save user: %v", err))
	}

	return c.JSONPretty(http.StatusOK, u, "    ")
}
