package main

import (
	"net/http"
	"socialv3/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Feed user
//	@Description	Mengambil feed post (dari user yang di-follow) dengan paginasi. Query: limit, offset, sort (asc/desc), tags (comma-separated), search, since, until.
//	@Tags			users
//	@Produce		json
//	@Param			limit	query		int					false	"Limit (1-100)"		default(10)
//	@Param			offset	query		int					false	"Offset"			default(0)
//	@Param			sort	query		string				false	"Sort (asc/desc)"	default("desc")
//	@Param			tags	query		string				false	"Filter tags (comma-separated)"
//	@Param			search	query		string				false	"Search"
//	@Param			since	query		string				false	"Filter sejak (datetime)"
//	@Param			until	query		string				false	"Filter sampai (datetime)"
//	@Success		200		{array}		store.PostWithMeta	"Daftar post di feed"
//	@Failure		400		{object}	object				"Query tidak valid"
//	@Failure		500		{object}	object				"Server error"
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "DESC",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	feeds, err := app.store.Posts.GetUserFeed(ctx, int64(198), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	_ = app.jsonResponse(w, http.StatusOK, feeds)
}
