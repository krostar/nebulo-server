package handler

import (
	"github.com/krostar/nebulo/log"
	"github.com/labstack/echo"
)

// UserInfos is not handled
func UserInfos(c echo.Context) error {
	log.Debugln(c.Get("userPublicKey"))

	// TODO:
	// retourne tout ce que tu peux sur le user :
	// ID           => non
	// B64PublicKey => non
	// FingerPrint  => oui
	// DisplayName  => oui
	// SignUp       => oui
	// LoginFirst   => oui
	// LoginLast    => oui
	return nil
}
