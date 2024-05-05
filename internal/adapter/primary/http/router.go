package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
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
	scanner QueryScanner
	updater Updater

	scanTimeout time.Duration
	scanLimit   int
}

type QueryScanner interface {
	Scan(ctx context.Context, query string, useIndex bool) ([]string, error)
}

type Updater interface {
	Update(ctx context.Context) (int, error)
}

func NewRouter(log *slog.Logger, handler *http.ServeMux, scanner QueryScanner, updater Updater, opts ...Option) {
	r := &router{
		log:         log,
		scanner:     scanner,
		updater:     updater,
		scanTimeout: defaultScanTimeout,
		scanLimit:   defaultScanLimit,
	}

	for _, opt := range opts {
		opt(r)
	}

	handler.HandleFunc("POST /update", r.Update)
	handler.HandleFunc("GET /pics", r.Pics)
}

func (r *router) Update(w http.ResponseWriter, req *http.Request) {
	const op = "router.Update"
	log := r.log.With(slog.String("op", op))

	log.Debug("handle update")

	total, err := r.updater.Update(req.Context())
	if err != nil {
		if errors.Is(err, service.ErrUpdateInProgress) {
			responseError(w, http.StatusAccepted, "update in progress")
			return
		}

		log.Error("error updating", logger.Err(err))
		responseError(w, http.StatusInternalServerError, "update failed")
		return
	}

	if err = responseJson(w, &UpdateResponse{Total: total}); err != nil {
		log.Error("failed to response", logger.Err(err))
	}
}

func (r *router) Pics(w http.ResponseWriter, req *http.Request) {
	const op = "router.Pics"
	log := r.log.With(slog.String("op", op))

	log.Debug("handle search")

	if err := req.ParseForm(); err != nil {
		log.Error("failed to parse form", logger.Err(err))
		responseError(w, http.StatusBadRequest, "bad request")
		return
	}
	if !req.Form.Has(formSearch) {
		responseError(w, http.StatusBadRequest, "search param required")
		return
	}

	search := req.FormValue(formSearch)
	ctx, cancel := context.WithTimeout(req.Context(), r.scanTimeout)
	defer cancel()
	res, err := r.scanner.Scan(ctx, search, true)
	if err != nil {
		log.Error("scan error")
		responseError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if ctx.Err() != nil {
		log.Error("scan timeout exceeded")
		responseError(w, http.StatusGatewayTimeout, "scan timeout exceeded")
		return
	}

	if len(res) > r.scanLimit {
		res = res[:r.scanLimit]
	}

	if err = responseJson(w, res); err != nil {
		log.Error("failed to response", logger.Err(err))
	}
}
