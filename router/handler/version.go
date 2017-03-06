package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

type build struct {
	BuildVersion string `json:"build_version"`
	BuildTime    string `json:"build_time"`
}

var (
	// BuildTime is filled by the main package
	BuildTime = ""
	// BuildVersion is filled by the main package
	BuildVersion = ""
)

// Version handle the route /version.
// Return the build time and tag version or the revision if tag doesnt exist
/**
 * @api {get} /version Version of the API
 * @apiDescription GIT revision or TAG of the API
 * @apiName Version
 * @apiGroup Other
 *
 * @apiExample {curl} Usage example
 *		$>curl -X GET --cacert ca.crt -v "https://api.nebulo.io/version/"
 *
 * @apiSuccess (Success) {nothing} 200 OK
 * @apiSuccessExample {json} Success example
 *		HTTP/1.1 200 "OK"
 *		{
 *			"build_version": "0.1.0",
 *			"build_time": "2017-02-16-0624 UTC",
 *		}
 */
func Version(c echo.Context) error {
	b := build{
		BuildVersion: BuildVersion,
		BuildTime:    BuildTime,
	}
	return c.JSON(http.StatusOK, b)
}
