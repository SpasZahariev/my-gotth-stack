package middleware

import (
	"context"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"goth/internal/store"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type key string

var NonceKey key = "nonces"

type Nonces struct {
	Htmx            string
	ResponseTargets string
	Tw              string
	Clerk           string
	ClerkSignIn     string
	HtmxCSSHash     string
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware() gin.HandlerFunc {
	// To use the same nonces in all responses, move the Nonces
	// struct creation to here, outside the handler.

	return func(c *gin.Context) {
		// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)
		// Create a new Nonces struct for every request when here.
		// move to outside the handler to use the same nonces in all responses
		nonceSet := Nonces{
			Htmx:            generateRandomString(16),
			ResponseTargets: generateRandomString(16),
			Tw:              generateRandomString(16),
			Clerk:           generateRandomString(16),
			ClerkSignIn:     generateRandomString(16),
			HtmxCSSHash:     "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg=",
		}

		// set nonces in context
		// ctx := context.WithValue(r.Context(), NonceKey, nonceSet)
		c.Set(string(NonceKey), nonceSet)
		// insert the nonces into the content security policy header
		cspHeader := fmt.Sprintf(
			"default-src 'self'; "+
				"script-src 'nonce-%s' 'nonce-%s' 'nonce-%s' 'nonce-%s' https://modern-colt-57.clerk.accounts.dev https://img.clerk.com; "+
				// "style-src 'nonce-%s' '%s' 'unsafe-inline'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"connect-src 'self' https://modern-colt-57.clerk.accounts.dev; worker-src 'self' blob:; "+
				"img-src 'self' https://img.clerk.com;",
			nonceSet.Htmx,
			nonceSet.ResponseTargets,
			nonceSet.Clerk,
			nonceSet.ClerkSignIn,
		// nonceSet.Tw,
		// nonceSet.HtmxCSSHash
		)
		c.Header("Content-Security-Policy", cspHeader)

		c.Next()
		// next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func TextHTMLMiddleware() gin.HandlerFunc {
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)
	return func(c *gin.Context) {
		// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Next()
		// next.ServeHTTP(w, r)
	}
}

// get the Nonce from the context, it is a struct called Nonces,
// so we can get the nonce we need by the key, i.e. HtmxNonce
func GetNonces(c *gin.Context) Nonces {
	// nonceSet := ctx.Value(NonceKey)
	nonceSet, exists := c.Get(string(NonceKey))
	if !exists {
		log.Fatal("error getting nonce set - is nil")
	}

	nonces, ok := nonceSet.(Nonces)

	if !ok {
		log.Fatal("error getting nonce set - not ok")
	}

	return nonces
}

func GetHtmxNonce(c *gin.Context) string {
	nonceSet := GetNonces(c)

	return nonceSet.Htmx
}

func GetResponseTargetsNonce(c *gin.Context) string {
	nonceSet := GetNonces(c)
	return nonceSet.ResponseTargets
}

func GetTwNonce(c *gin.Context) string {
	nonceSet := GetNonces(c)
	return nonceSet.Tw
}

func GetClerkSignInNonce(c *gin.Context) string {
	nonceSet := GetNonces(c)
	return nonceSet.ClerkSignIn
}

func GetClerkNonce(c *gin.Context) string {
	nonceSet := GetNonces(c)
	return nonceSet.Clerk
}

type AuthMiddleware struct {
	sessionStore      store.SessionStore
	sessionCookieName string
}

func NewAuthMiddleware(sessionStore store.SessionStore, sessionCookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore:      sessionStore,
		sessionCookieName: sessionCookieName,
	}
}

type UserContextKey string

var UserKey UserContextKey = "user"

func (m *AuthMiddleware) AddUserToContext() gin.HandlerFunc {

	return func(c *gin.Context) {
		sessionCookie, err := c.Cookie(m.sessionCookieName)

		if err != nil {
			fmt.Println("error getting session cookie", err)
			c.Next()
			return
		}

		decodedValue, err := b64.StdEncoding.DecodeString(sessionCookie)

		if err != nil {
			c.Next()
			return
		}

		splitValue := strings.Split(string(decodedValue), ":")

		if len(splitValue) != 2 {
			c.Next()
			return
		}

		sessionID := splitValue[0]
		userID := splitValue[1]

		fmt.Println("sessionID", sessionID)
		fmt.Println("userID", userID)

		user, err := m.sessionStore.GetUserFromSession(sessionID, userID)

		if err != nil {
			c.Next()
			return
		}

		c.Set(string(UserKey), user)

		c.Next()
	}
}

func GetUser(ctx context.Context) *store.User {
	user := ctx.Value(UserKey)
	if user == nil {
		return nil
	}

	return user.(*store.User)
}

// func GetUser(c *gin.Context) *store.User {
// 	user, exists := c.Get(string(UserKey))
// 	if !exists {
// 		return nil
// 	}
//
// 	return user.(*store.User)
// }
