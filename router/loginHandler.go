package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/gologin/google"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	googleOAuth2 "golang.org/x/oauth2/google"
)

var (
	googleOauth2Config   *oauth2.Config
	facebookOauth2Config *oauth2.Config

	stateConfig = gologin.DebugOnlyCookieConfig
)

func init() {
	googleOauth2Config = &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Host + "/v1/authentication/google/redirect",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	facebookOauth2Config = &oauth2.Config{
		ClientID:     cfg.Facebook.ClientID,
		ClientSecret: cfg.Facebook.ClientSecret,
		RedirectURL:  cfg.Host + "/v1/authentication/facebook/redirect",
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}
}

func providerLogin(c *gin.Context) {
	provider := c.Param("provider")
	if strings.EqualFold(provider, "google") {
		g := googleLoginHandler()
		g(c)
	} else if strings.EqualFold(provider, "facebook") {
		f := facebookLoginHandler()
		f(c)
	}
	return
}

func providerRedirect(c *gin.Context) {
	provider := c.Param("provider")
	if strings.EqualFold(provider, "google") {
		g := googleRedirectHandler()
		g(c)
	} else if strings.EqualFold(provider, "facebook") {
		f := facebookRedirectHandler()
		f(c)
	}
	return
}

func googleLoginHandler() gin.HandlerFunc {
	return gin.WrapH(google.StateHandler(stateConfig, google.LoginHandler(googleOauth2Config, nil)))
}

func googleRedirectHandler() gin.HandlerFunc {
	return gin.WrapH(google.StateHandler(stateConfig, google.CallbackHandler(googleOauth2Config, googleSetJWT(), nil)))
}

func googleSetJWT() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data map[string]string
		data["email"] = googleUser.Email
		data["gender"] = googleUser.Gender
		data["picture"] = googleUser.Picture
		data["name"] = googleUser.Name

		renderTokenPage(w, data)
	})
}

func facebookLoginHandler() gin.HandlerFunc {
	return gin.WrapH(facebook.StateHandler(stateConfig, facebook.LoginHandler(facebookOauth2Config, nil)))
}

func facebookRedirectHandler() gin.HandlerFunc {
	return gin.WrapH(facebook.StateHandler(stateConfig, facebook.CallbackHandler(facebookOauth2Config, facebookSetJWT(), nil)))
}

func facebookSetJWT() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var data map[string]string
		data["id"] = facebookUser.ID
		data["email"] = facebookUser.Email
		data["name"] = facebookUser.Name

		renderTokenPage(w, data)
	})
}

func renderTokenPage(w http.ResponseWriter, data map[string]string) {
	jwtoken := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()

	for key, value := range data {
		claims[key] = value
	}
	jwtoken.Claims = claims
	tokenString, err := jwtoken.SignedString([]byte(secretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		return
	}
	response := token{tokenString}
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := page{
		Token: string(json),
	}
	renderTemplate(w, "authenticated", &page)
}

type token struct {
	Token string `json:"token"`
}

type page struct {
	Token string
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	t, err := template.ParseFiles("cmd/movie/" + tmpl + ".html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, p)
}
