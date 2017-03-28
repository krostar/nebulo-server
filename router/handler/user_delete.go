package handler

import (
	"fmt"
	"net/http"

	"github.com/krostar/nebulo/router/httperror"
	"github.com/labstack/echo"

	"github.com/krostar/nebulo/tools/cert"
	up "github.com/krostar/nebulo/user/provider"
)

// UserDelete handle the route DELETE /user/.
// Delete all the users informations and revoke certificate
/**
 * @api {delete} /user Delete user profile
 * @apiDescription Delete the user profile, wiping every data about the user.
 * @apiName User - delete profile
 * @apiGroup User
 *
 * @apiExample {curl} Usage example
 *		$>curl -X DELETE -v --cert bob.crt --key bob.key "https://api.nebulo.io/user/"
 *
 * @apiSuccess (Success) {nothing} 202 Accepted
 * @apiSuccessExample {json} Success example
 *		HTTP/1.1 202 "Accepted"
 *
 * @apiError (Errors 4XX) {json} 401 Unauthorized: missing client certificate
 * @apiError (Errors 4XX) {json} 404 Not found: user not found
 * @apiError (Errors 5XX) {json} 500 Internal server error: server failed to handle the request
 */
func UserDelete(c echo.Context) (err error) {
	u, err := GetLoggedUser(c.Get("user"))
	if err != nil {
		return httperror.UserNotFound()
	}

	if err = cert.Revoke(c.Request().TLS.PeerCertificates[0]); err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to revoke certificate: %v", err))
	}
	// cp.P.LeaveAll(user, delete personal copy = true)
	// fp.P.DeleteAll(user)
	if err = up.P.Delete(u); err != nil {
		return httperror.HTTPInternalServerError(fmt.Errorf("unable to delete user profile: %v", err))
	}

	return c.NoContent(http.StatusAccepted)
}
