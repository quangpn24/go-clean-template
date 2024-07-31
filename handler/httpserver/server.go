package httpserver

import (
	"net/http"
	"strings"

	middleware2 "go-clean-template/handler/httpserver/middleware"
	apperror "go-clean-template/pkg/apperror"
	"go-clean-template/pkg/config"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/sentry"
	"go-clean-template/usecase"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Router *echo.Echo
	Config *config.Config
	Logger *zap.SugaredLogger

	TransactionUseCase usecase.ITransactionUseCase
}

func New(options ...Options) (*Server, error) {
	s := Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	for _, fn := range options {
		if err := fn(&s); err != nil {
			return nil, err
		}
	}

	s.RegisterGlobalMiddlewares()

	apiV1 := s.Router.Group("/api/v1")

	s.RegisterHealthCheck(s.Router.Group(""))
	s.RegisterTransactionRoutesV1(apiV1.Group("/transactions"))

	return &s, nil
}

func (s *Server) RegisterGlobalMiddlewares() {
	s.Router.Use(middleware.Recover())
	s.Router.HideBanner = false
	s.Router.HidePort = false
	s.Router.Use(middleware.Secure())
	s.Router.Use(middleware.RequestID())
	s.Router.Use(middleware.Gzip())
	s.Router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	skipPath := []string{
		"/healthz",
		"/api/v1/transactions",
	}

	// Authentication with cognito
	auth := middleware2.NewAuthentication(s.Config.UserPoolID, skipPath, s.Config)
	s.Router.Use(auth.Middleware())

	// CORS
	if s.Config.AllowOrigins != "" {
		aos := strings.Split(s.Config.AllowOrigins, ",")
		s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: aos,
		}))
	}
}

func (s *Server) Start(addr string) error {
	return s.Router.Start(addr)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) RegisterHealthCheck(Router *echo.Group) {
	Router.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK!!!")
	})
}

func (s *Server) handleError(c echo.Context, err error) error {
	s.Logger.Errorw(
		err.Error(),
		zap.String("request_id", s.requestID(c)),
	)

	if e, ok := apperror.ErrorAs(err); ok {
		if e.HTTPCode >= http.StatusInternalServerError {
			sentry.WithContext(c).Error(err)
		}

		return c.JSON(e.HTTPCode, Errs{
			ErrCode: e.Code,
			Message: e.Message,
			RawErr:  e.Raw.Error(),
			Info:    e.Info,
		})
	} else {
		sentry.WithContext(c).Error(err)
		return c.JSON(http.StatusInternalServerError, Errs{
			Message: err.Error(),
		})
	}
}

type Errs struct {
	ErrCode interface{} `json:"err_code,omitempty"`
	Message string      `json:"message,omitempty"`
	RawErr  string      `json:"raw_err,omitempty"`
	Info    interface{} `json:"info,omitempty"`
}

type Success struct {
	Code    interface{} `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *Server) handleSuccess(c echo.Context, httpCode int, data interface{}) error {
	return c.JSON(httpCode, Success{
		Code:    http.StatusText(httpCode),
		Message: "success",
		Data:    data,
	})
}

func (s *Server) requestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
