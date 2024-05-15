package hanlder

import (
	"net/http"
	"yadro-go/internal/core/domain"
)

type AuthenticatedHandlerFunc = func(w http.ResponseWriter, req *http.Request, user *domain.User)
