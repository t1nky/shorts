package router

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"shorts/controllers"
	"shorts/database"
	h "shorts/helper"
	"shorts/models"

	"github.com/gin-gonic/gin"
)

// SetupRouter : Creates default instance of gin and adds routes for it
func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(errorHandler)

	// Routes for authenticated only users
	authorizedV1 := r.Group("v1/", basicAuth())

	// User actions

	// swagger:route GET /me user getCurrentUser
	// Return currently authenticated user's information
	// responses:
	//   400: ResponseError
	//   401: ResponseError
	//   200: UserResponse
	// security:
	//   basic:
	authorizedV1.GET("me", controllers.GetCurrentUser)
	// swagger:route GET /logout user logout
	// Log out current user
	// responses:
	//   401: ResponseError
	// security:
	//   basic:
	authorizedV1.GET("logout", responseUnauthorized)

	// Short links actions

	// swagger:route GET /shorts shortlink getShortlinks
	// Return list of short links created by currently authenticated user
	// responses:
	//   400: ResponseError
	//   401: ResponseError
	//   200: ShortlinksResponse
	// security:
	//   basic:
	authorizedV1.GET("shorts", controllers.GetShortlinks)
	// swagger:route GET /short/{id} shortlink getShortlink
	// Return information about specific short link that was created by currently authenticated user
	// responses:
	//   400: ResponseError
	//   401: ResponseError
	//   200: ShortlinkResponse
	//   404: ResponseError
	// security:
	//   basic:
	authorizedV1.GET("shorts/:id", controllers.GetShortlinkInfo)
	// swagger:route POST /shorts shortlink addShortlink
	// Create a new short link
	// responses:
	//   400: ResponseError
	//   401: ResponseError
	//   201: AddShortResponse
	// security:
	//   basic:
	authorizedV1.POST("shorts", controllers.AddShortlink)
	// swagger:route DELETE /shorts shortlink deleteShortlink
	// Delete specific short link that was created by currently authenticated user
	// responses:
	//   400: ResponseError
	//   401: ResponseError
	//   200: ResponseOK
	//   404: ResponseError
	// security:
	//   basic:
	authorizedV1.DELETE("shorts/:id", controllers.DeleteShortlink)

	publicV1 := r.Group("v1/")

	// swagger:route POST /users user addUser
	// Create a new user
	// responses:
	//   400: ResponseError
	//   201: ResponseOK
	publicV1.POST("users", controllers.AddUser)
	// swagger:route GET /s/{short} shortlink redirectByShortlink
	// Redirect to a full link by a given short link
	// responses:
	//   301: RedirectResponse
	//   400: ResponseError
	publicV1.GET("s/:short", controllers.GetShortlinkRedirect)

	publicV1Stats := publicV1.Group("stats/")

	// swagger:route GET /stats/top stats getShortlinksTop
	// Return top 20 websites that were most often redirected to
	// responses:
	//   400: ResponseError
	//   200: TopDomainsResponse
	publicV1Stats.GET("top", controllers.GetShortlinksTop)
	// swagger:route GET /stats/graph stats getShortlinksGraph
	// Return amount of redirects groupped by day, hour and minute
	// responses:
	//   400: ResponseError
	//   200: ShortlinksGraphResponse
	publicV1Stats.GET("graph", controllers.GetShortlinksGraph)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, h.NewResponseError(h.NewPageNotFoundError()))
	})

	return r
}

// errorHandler : Default error handler
func errorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		fmt.Println(c.Errors)
	}
}

// basicAuth : Check for authentication
func basicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			responseUnauthorized(c)
			return
		}

		authPayload, _ := base64.StdEncoding.DecodeString(auth[1])
		authPair := strings.SplitN(string(authPayload), ":", 2)

		if len(authPair) != 2 {
			responseUnauthorized(c)
			return
		}

		authOK, userID := authenticateUser(authPair[0], authPair[1])
		if authOK {
			c.Set(gin.AuthUserKey, userID)
		} else {
			responseUnauthorized(c)
			return
		}

		c.Next()
	}
}

// authenticateUser: Find user;password pair in the DB
func authenticateUser(username, password string) (bool, uint64) {
	user := models.User{Name: username, Password: password}

	err := database.DB.Where(&user).First(&user)
	if err.Error != nil {
		return false, 0
	}

	return true, user.ID
}

func responseUnauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "Basic")
	c.AbortWithStatusJSON(http.StatusUnauthorized, h.NewResponseError(errors.New("Authentication required")))
}
