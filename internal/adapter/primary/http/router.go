package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"
	"yadro-go/internal/adapter/primary"
	"yadro-go/internal/adapter/primary/http/middleware"
	"yadro-go/internal/adapter/primary/http/protocol"
	"yadro-go/internal/core/domain"
	"yadro-go/internal/core/service"
	"yadro-go/pkg/logger"
)

const (
	formSearch = "search"

	defaultScanTimeout = 1 * time.Minute
	defaultScanLimit   = 10
)

type router struct {
	log     *slog.Logger
	scanner primary.QueryScanner
	updater primary.Updater
	auth    primary.Auth

	scanTimeout time.Duration
	scanLimit   int
}

func ApplyRouter(
	log *slog.Logger,
	handler *http.ServeMux,
	scanner primary.QueryScanner,
	updater primary.Updater,
	auth primary.Auth,
	authMiddleware *middleware.AuthMiddleware,
	rpsMiddleware *middleware.RpsLimitMiddleware,
	concurrencyMiddleware *middleware.ConcurrencyLimitMiddleware,
	opts ...Option,
) {
	r := &router{
		log:         log,
		scanner:     scanner,
		updater:     updater,
		auth:        auth,
		scanTimeout: defaultScanTimeout,
		scanLimit:   defaultScanLimit,
	}

	for _, opt := range opts {
		opt(r)
	}

	handler.HandleFunc("POST /login", r.Login)
	handler.HandleFunc("POST /update", authMiddleware.WithAuth(domain.ROLE_ADMIN, r.Update))
	handler.HandleFunc("GET /pics", concurrencyMiddleware.WithConcurrencyLimit(
		authMiddleware.WithAuth(
			domain.ROLE_USER,
			rpsMiddleware.WithRpsLimit(r.Pics))),
	)
}

func (r *router) Update(w http.ResponseWriter, req *http.Request, user *domain.User) {
	const op = "router.Update"
	log := r.log.With(slog.String("op", op), slog.String("uname", user.Username))

	log.Debug("handle update")

	total, err := r.updater.Update(req.Context())
	if err != nil {
		if errors.Is(err, service.ErrUpdateInProgress) {
			protocol.ResponseError(w, http.StatusAccepted, "update in progress")
			return
		}

		log.Error("error updating", logger.Err(err))
		protocol.ResponseError(w, http.StatusInternalServerError, "update failed")
		return
	}

	if err = protocol.ResponseJson(w, &protocol.UpdateResponse{Total: total}); err != nil {
		log.Error("failed to response", logger.Err(err))
	}
}

func (r *router) Pics(w http.ResponseWriter, req *http.Request, user *domain.User) {
	const op = "router.Pics"
	log := r.log.With(slog.String("op", op), slog.String("uname", user.Username))

	log.Debug("handle search")

	if err := req.ParseForm(); err != nil {
		log.Error("failed to parse form", logger.Err(err))
		protocol.ResponseError(w, http.StatusBadRequest, "bad request")
		return
	}
	if !req.Form.Has(formSearch) {
		protocol.ResponseError(w, http.StatusBadRequest, "search param required")
		return
	}

	search := req.FormValue(formSearch)
	ctx, cancel := context.WithTimeout(req.Context(), r.scanTimeout)
	defer cancel()
	res, err := r.scanner.Scan(ctx, search, true)
	if err != nil {
		log.Error("scan error")
		protocol.ResponseError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if ctx.Err() != nil {
		log.Error("scan timeout exceeded")
		protocol.ResponseError(w, http.StatusGatewayTimeout, "scan timeout exceeded")
		return
	}

	if len(res) > r.scanLimit {
		res = res[:r.scanLimit]
	}

	if err = protocol.ResponseJson(w, res); err != nil {
		log.Error("failed to response", logger.Err(err))
	}
}

func (r *router) Login(w http.ResponseWriter, req *http.Request) {
	const op = "router.Login"
	log := r.log.With(slog.String("op", op))

	log.Debug("handle login")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error("failed to read body", logger.Err(err))
		protocol.ResponseError(w, http.StatusBadRequest, "bad request")
		return
	}

	var loginRequest protocol.LoginRequest
	if err = json.Unmarshal(body, &loginRequest); err != nil {
		log.Error("failed to unmarshal login request", logger.Err(err))
		protocol.ResponseError(w, http.StatusBadRequest, "bad request")
		return
	}

	if len(loginRequest.Username) == 0 {
		protocol.ResponseError(w, http.StatusBadRequest, "username is mandatory")
		return
	}

	token, err := r.auth.Login(req.Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			protocol.ResponseError(w, http.StatusUnauthorized, "wrong credentials")
			return
		}

		log.Error("failed to login", logger.Err(err))
		protocol.ResponseError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err = protocol.ResponseJson(w, protocol.LoginResponse{Token: token}); err != nil {
		log.Error("failed to response", logger.Err(err))
		return
	}
}
