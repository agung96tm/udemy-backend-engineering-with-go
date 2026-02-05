package main

import (
	"context"
	"errors"
	"net/http"
	"socialv3/internal/store"
)

type userKey string

const userCtx userKey = "user"

// getUserHandler godoc
//
//	@Summary		Ambil detail user
//	@Description	Mengambil satu user berdasarkan ID.
//	@Tags			users
//	@Produce		json
//	@Param			userID	path		int			true	"ID user"
//	@Success		200		{object}	store.User	"Detail user"
//	@Failure		400		{object}	object		"ID tidak valid"
//	@Failure		404		{object}	object		"User tidak ditemukan"
//	@Failure		500		{object}	object		"Server error"
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	_ = app.jsonResponse(w, http.StatusOK, user)
}

type FollowUserRequest struct {
	UserID int64 `json:"user_id" validate:"required,min=1"`
}

// followUserHandler godoc
//
//	@Summary		Follow user
//	@Description	User saat ini (dari context) mem-follow user dengan user_id dari body.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int						true	"ID user (follower / yang mem-follow)"
//	@Param			body	body		FollowUserRequest		true	"ID user yang akan di-follow"
//	@Success		200		{object}	map[string]interface{}	"message: user followed successfully"
//	@Failure		400		{object}	object					"Request body tidak valid"
//	@Failure		404		{object}	object					"User tidak ditemukan"
//	@Failure		500		{object}	object					"Server error"
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	var payload FollowUserRequest
	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	err := app.store.Follower.Follow(ctx, followerUser.ID, payload.UserID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	_ = app.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "user followed successfully",
	})
}

// unfollowUserHandler godoc
//
//	@Summary		Unfollow user
//	@Description	User saat ini (dari context) unfollow user dengan user_id dari body.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int						true	"ID user (yang unfollow)"
//	@Param			body	body		FollowUserRequest		true	"ID user yang akan di-unfollow"
//	@Success		200		{object}	map[string]interface{}	"message: user unfollowed successfully"
//	@Failure		400		{object}	object					"Request body tidak valid"
//	@Failure		404		{object}	object					"User tidak ditemukan"
//	@Failure		500		{object}	object					"Server error"
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromContext(r)

	var payload FollowUserRequest
	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	err := app.store.Follower.Unfollow(ctx, unfollowedUser.ID, payload.UserID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	_ = app.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "user unfollowed successfully",
	})
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := app.readID(r, "userID")
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	ctx := r.Context()
	user, ok := ctx.Value(userCtx).(*store.User)
	if !ok {
		return nil
	}
	return user
}
