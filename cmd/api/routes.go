package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func (app *application) Router() *echo.Echo {

	e := echo.New()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	//middlewares
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// Set the timeouts
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 30 * time.Second
	e.Server.IdleTimeout = time.Minute

	e.GET("/login", app.loginHandler)
	e.GET("/callback", app.callbackHandler)
	e.GET("/user", app.userHandler)
	e.GET("/logout", app.logoutHandler)
	e.RouteNotFound("/*", app.routeNotFoundHandler)

	return e
}

func (app *application) loginHandler(c echo.Context) error {

	state, err := generateRandomState()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// save the state inside the session
	sess, _ := session.Get("session", c)

	sess.Values["auth-state"] = state
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	c.Redirect(http.StatusTemporaryRedirect, app.auth.AuthCodeURL(state))
	return nil
}

func (app *application) callbackHandler(c echo.Context) error {
	app.logger.Info("callbackHandler")
	app.logger.Info("Getting the cookie",
		zap.String("############# ", c.QueryParam("code")),
		zap.String("############# auth-session ", c.QueryParam("state")),
	)

	sess, _ := session.Get("session", c)
	if c.QueryParam("state") != sess.Values["auth-state"] {
		return c.String(http.StatusBadRequest, "Invalid state paramater")
	}
	// exchange an authorization code for a token
	token, err := app.auth.Exchange(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		return c.String(http.StatusUnauthorized, "Failed to exchange authoization code for a token")
	}

	idToken, err := app.auth.VerifyIDToken(c.Request().Context(), token)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	sess.Values["access_token"] = token.AccessToken
	sess.Values["profile"] = profile

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.Redirect(http.StatusTemporaryRedirect, "/user")
	return nil
}

func (app *application) userHandler(c echo.Context) error {
	app.logger.Info("userHandler")
	return c.String(http.StatusOK, "userHandler")
}

func (app *application) logoutHandler(c echo.Context) error {
	app.logger.Info("logoutHandler")

	logoutUrl, err := url.Parse("http://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	scheme := "http"

	returnTo, err := url.Parse(scheme + "://" + c.Request().Host)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	sess, _ := session.Get("session", c)

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	app.logger.Info("logoutUrl" + logoutUrl.String())
	c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	// clear cookies
	sess.Values["access_token"] = nil
	sess.Values["profile"] = nil
	return nil
}

func (app *application) routeNotFoundHandler(c echo.Context) error {
	app.logger.Info("routeNotFoundHandler")
	return c.String(http.StatusNotFound, "routeNotFoundHandler")
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}