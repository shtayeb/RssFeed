package handlers

import (
	"net/http"

	"github.com/shtayeb/rssfeed/views"
)

func (cfg *ApiConfig) HandlerNotFoundPage(w http.ResponseWriter, r *http.Request) {
	views.NotFoundPage().Render(r.Context(), w)
}
