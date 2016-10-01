package middleware

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/hkparker/Wave/models"
)

//
// Ensure that a request is authenticated
//
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.Request.URL.Path
		if !PublicEndpoint(endpoint) {
			session_cookie, err := c.Request.Cookie("wave_session")

			if err != nil {
				c.Redirect(302, "/login")
				c.Abort()
				log.WithFields(log.Fields{
					"at":     "middleware.Authentication",
					"reason": "missing wave_session cookie",
					"error":  err.Error(),
				}).Info("redirecting unauthenticated request")
				return
			}

			var user models.User
			if session, err := models.SessionFromID(session_cookie.Value); err == nil {
				if user, err = session.User(); err != nil {
					c.Redirect(302, "/login")
					c.Abort()
					log.WithFields(log.Fields{
						"at":     "middleware.Authentication",
						"reason": "could not find user for session",
					}).Info("redirecting unauthenticated request")
					return
				}
			} else {
				c.Redirect(302, "/login")
				c.Abort()
				log.WithFields(log.Fields{
					"at":     "middleware.Authentication",
					"reason": "wave_session header does not exist in session record",
				}).Info("redirecting unauthenticated request")
				return
			}

			if AdminProtected(endpoint) && !user.Admin {
				c.JSON(401, gin.H{"error": "permission denied"})
				c.Abort()
				log.WithFields(log.Fields{
					"at":       "middleware.Authentication",
					"reason":   "user is not administrator",
					"user_id":  user.ID,
					"endpoint": endpoint,
				}).Warn("blocking unauthenticated request")
				return
			}
		}
	}
}

//
// Given an endpoint, return if the endpoint is accessible without authentication.
//
func PublicEndpoint(url string) bool {
	switch url {
	case "/login":
		return true
	}
	return false
}

//
// Given an endpoint, return if the endpoint can only be accessed by users
// with the admin role.
//
func AdminProtected(url string) bool {
	switch url {
	case "/users/create":
		return true
	}
	return false
}
