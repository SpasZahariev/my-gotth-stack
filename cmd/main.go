package main

import (
	"goth/internal/config"
	database "goth/internal/store/db"
	"goth/internal/store/dbstore"
	"goth/internal/templates"
	"os"

	m "goth/internal/middleware"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

// This function will render teh templ component Into
// a gin context's Response Writer
func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

/*
* Set to production at build time
* used to determine what assets to load
 */
var Environment = "development"

func init() {
	os.Setenv("env", Environment)
}

func main() {
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := gin.Default()

	cfg := config.MustLoadConfig()

	db := database.MustOpen(cfg.DatabaseName)
	// passwordhash := passwordhash.NewHPasswordHash()

	// userStore := dbstore.NewUserStore(
	// 	dbstore.NewUserStoreParams{
	// 		DB:           db,
	// 		PasswordHash: passwordhash,
	// 	},
	// )

	sessionStore := dbstore.NewSessionStore(
		dbstore.NewSessionStoreParams{
			DB: db,
		},
	)

	// fileServer := http.FileServer(http.Dir("./static"))
	// r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Static("/static", "./static")

	authMiddleware := m.NewAuthMiddleware(sessionStore, cfg.SessionCookieName)

	// Apply middleware
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		m.TextHTMLMiddleware(),
		m.CSPMiddleware(),
		authMiddleware.AddUserToContext(),
	)

	r.GET("/", func(c *gin.Context) {
		render(c, 200, templates.About())
	})

	// r.NoRoute(gin.WrapF(handlers.NewNotFoundHandler().ServeHTTP))

	// r.GET("/", gin.WrapF(handlers.NewHomeHandler().ServeHTTP))
	// r.GET("/about", gin.WrapF(handlers.NewAboutHandler().ServeHTTP))
	// r.GET("/register", gin.WrapF(handlers.NewGetRegisterHandler().ServeHTTP))
	// r.POST("/register", gin.WrapF(handlers.NewPostRegisterHandler(handlers.PostRegisterHandlerParams{
	// 	UserStore: userStore,
	// }).ServeHTTP))
	// r.GET("/login", gin.WrapF(handlers.NewGetLoginHandler().ServeHTTP))
	// r.POST("/login", gin.WrapF(handlers.NewPostLoginHandler(handlers.PostLoginHandlerParams{
	// 	UserStore:         userStore,
	// 	SessionStore:      sessionStore,
	// 	PasswordHash:      passwordhash,
	// 	SessionCookieName: cfg.SessionCookieName,
	// }).ServeHTTP))
	// r.POST("/logout", gin.WrapF(handlers.NewPostLogoutHandler(handlers.PostLogoutHandlerParams{
	// 	SessionCookieName: cfg.SessionCookieName,
	// }).ServeHTTP))

	// r.Run("localhost:4000")
	r.Run()

	// killSig := make(chan os.Signal, 1)

	// signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	// srv := &http.Server{
	// 	Addr:    cfg.Port,
	// 	Handler: r,
	// }

	// go func() {
	// 	err := srv.ListenAndServe()
	//
	// 	if errors.Is(err, http.ErrServerClosed) {
	// 		logger.Info("Server shutdown complete")
	// 	} else if err != nil {
	// 		logger.Error("Server error", slog.Any("err", err))
	// 		os.Exit(1)
	// 	}
	// }()
	//
	// logger.Info("Server started", slog.String("port", cfg.Port), slog.String("env", Environment))
	// <-killSig
	//
	// logger.Info("Shutting down server")
	//
	// // Create a context with a timeout for shutdown
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	//
	// // Attempt to gracefully shut down the server
	// if err := srv.Shutdown(ctx); err != nil {
	// 	logger.Error("Server shutdown failed", slog.Any("err", err))
	// 	os.Exit(1)
	// }
	//
	// logger.Info("Server shutdown complete")
}
