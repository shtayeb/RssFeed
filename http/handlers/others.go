package handlers

import (
	"net/http"

	"github.com/shtayeb/rssfeed/views"
)

func HandlerNotFoundPage(w http.ResponseWriter, r *http.Request) {
	views.NotFoundPage().Render(r.Context(), w)
}
