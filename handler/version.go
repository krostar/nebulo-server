package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// Version handle the route /version.
// Return the TAG version or the revision if TAG doesnt exist
/**
 * @api {get} /version Version of the API
 * @apiDescription GIT revision or TAG of the API
 * @apiName Version
 * @apiGroup Other
 *
 * @apiExample {curl} Usage example
 *		$>curl -X GET "http://127.0.0.1:17241/version"
 *
 * @apiSuccess (Success) {nothing} 200 OK
 * @apiSuccessExample {json} Success example
 *		HTTP/1.1 200 "OK"
 *		{
 *		}
 */
func Version(c echo.Context) error {

	var version string
	version = ""

	return c.JSON(http.StatusOK, version)
}
