package handler

import (
	"net/http"
	"time"

	"social-api/internal/auth"
	"social-api/internal/cache"
	"social-api/internal/model"
	"social-api/internal/store"
)

// CreatePost handles the post creation endpoint
func (app *Application) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		app.unauthorizedResponse(w, r)
		return
	}

	// Parse request body
	var input model.PostInput
	err := model.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate input
	err = app.ValidateRequest(input)
	if err != nil {
		validationErrors := app.FormatValidationErrors(err)
		app.validationErrorResponse(w, r, validationErrors)
		return
	}

	// Create post object
	post := &store.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  user.ID,
	}

	// Create post in database
	err = app.PostStore.Create(r.Context(), post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set user for the response
	post.User = user

	// Convert to response
	response := model.PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	// Cache post if enabled
	if app.Cache != nil {
		err = app.Cache.Set(r.Context(), cache.PostKey(post.ID), post, 15*time.Minute)
		if err != nil {
			app.Logger.Printf("Error caching post: %v", err)
		}
	}

	// Send response
	err = model.WriteJSON(w, http.StatusCreated, model.NewResponse(response))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GetPost handles the get post endpoint
func (app *Application) GetPost(w http.ResponseWriter, r *http.Request) {
	// Extract post ID from URL
	id, err := app.GetIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Try to get post from cache
	var post *store.Post
	if app.Cache != nil {
		var cachedPost store.Post
		err = app.Cache.Get(r.Context(), cache.PostKey(id), &cachedPost)
		if err == nil {
			post = &cachedPost
		}
	}

	// Get post from database if not in cache
	if post == nil {
		post, err = app.PostStore.GetByID(r.Context(), id)
		if err != nil {
			app.handleError(w, r, err)
			return
		}

		// Cache post
		if app.Cache != nil {
			err = app.Cache.Set(r.Context(), cache.PostKey(id), post, 15*time.Minute)
			if err != nil {
				app.Logger.Printf("Error caching post: %v", err)
			}
		}
	}

	// Convert to response
	response := model.PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		User: model.UserResponse{
			ID:        post.User.ID,
			Username:  post.User.Username,
			Email:     post.User.Email,
			CreatedAt: post.User.CreatedAt,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	// Send response
	err = model.WriteJSON(w, http.StatusOK, model.NewResponse(response))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// UpdatePost handles the update post endpoint
func (app *Application) UpdatePost(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		app.unauthorizedResponse(w, r)
		return
	}

	// Extract post ID from URL
	id, err := app.GetIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get post from database
	post, err := app.PostStore.GetByID(r.Context(), id)
	if err != nil {
		app.handleError(w, r, err)
		return
	}

	// Check if user is the post owner
	if post.UserID != user.ID {
		app.forbiddenResponse(w, r)
		return
	}

	// Parse request body
	var input model.PostUpdateInput
	err = model.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate input
	err = app.ValidateRequest(input)
	if err != nil {
		validationErrors := app.FormatValidationErrors(err)
		app.validationErrorResponse(w, r, validationErrors)
		return
	}

	// Apply updates if provided
	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}

	// Update post in database
	err = app.PostStore.Update(r.Context(), post)
	if err != nil {
		app.handleError(w, r, err)
		return
	}

	// Invalidate cache
	if app.Cache != nil {
		err = app.Cache.Delete(r.Context(), cache.PostKey(id))
		if err != nil {
			app.Logger.Printf("Error deleting post from cache: %v", err)
		}
	}

	// Convert to response
	response := model.PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	// Send response
	err = model.WriteJSON(w, http.StatusOK, model.NewResponse(response))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// DeletePost handles the delete post endpoint
func (app *Application) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		app.unauthorizedResponse(w, r)
		return
	}

	// Extract post ID from URL
	id, err := app.GetIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get post from database
	post, err := app.PostStore.GetByID(r.Context(), id)
	if err != nil {
		app.handleError(w, r, err)
		return
	}

	// Check if user is the post owner
	if post.UserID != user.ID {
		app.forbiddenResponse(w, r)
		return
	}

	// Delete post from database
	err = app.PostStore.Delete(r.Context(), id)
	if err != nil {
		app.handleError(w, r, err)
		return
	}

	// Invalidate cache
	if app.Cache != nil {
		err = app.Cache.Delete(r.Context(), cache.PostKey(id))
		if err != nil {
			app.Logger.Printf("Error deleting post from cache: %v", err)
		}
	}

	// Send no content response
	w.WriteHeader(http.StatusNoContent)
}

// ListPosts handles the list posts endpoint
func (app *Application) ListPosts(w http.ResponseWriter, r *http.Request) {
	// Get pagination params
	pagination := model.GetPagination(r)

	// Get filter params
	filter := model.PostFilter{}

	// Parse user_id filter if provided
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := app.GetIDParam(r)
		if err == nil {
			filter.UserID = &userID
		}
	}

	// Parse title filter if provided
	if title := r.URL.Query().Get("title"); title != "" {
		filter.Title = &title
	}

	// Parse content filter if provided
	if content := r.URL.Query().Get("content"); content != "" {
		filter.Content = &content
	}

	// Get posts from database
	posts, totalCount, err := app.PostStore.List(r.Context(), pagination, filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Convert to responses
	postResponses := make([]model.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = model.PostResponse{
			ID:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			User: model.UserResponse{
				ID:        post.User.ID,
				Username:  post.User.Username,
				Email:     post.User.Email,
				CreatedAt: post.User.CreatedAt,
			},
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	// Send response with pagination
	err = model.WriteJSON(
		w,
		http.StatusOK,
		model.NewPageResponse(
			postResponses,
			pagination.Page,
			pagination.PageSize,
			totalCount,
		),
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
