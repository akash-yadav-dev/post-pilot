package handler

import (
	"errors"
	"fmt"
	"net/http"

	"post-pilot/apps/api/internal/posts/model"
	"post-pilot/apps/api/internal/posts/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for the posts domain.
type Handler struct {
	svc service.PostService
}

func NewHandler(svc service.PostService) *Handler {
	return &Handler{svc: svc}
}

// CreatePost POST /api/v1/posts
func (h *Handler) CreatePost(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.svc.CreatePost(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPost GET /api/v1/posts/:id
func (h *Handler) GetPost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, err := h.svc.GetPost(c.Request.Context(), id)
	if errors.Is(err, service.ErrPostNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// ListPosts GET /api/v1/posts
func (h *Handler) ListPosts(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}

	posts, err := h.svc.ListUserPosts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// UpdatePost PATCH /api/v1/posts/:id
func (h *Handler) UpdatePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.svc.UpdatePost(c.Request.Context(), id, req)
	if errors.Is(err, service.ErrPostNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost DELETE /api/v1/posts/:id
func (h *Handler) DeletePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.svc.DeletePost(c.Request.Context(), id); errors.Is(err, service.ErrPostNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	c.Status(http.StatusNoContent)
}

// userIDFromContext extracts the authenticated user's UUID from the Gin context (set by auth middleware).
func userIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}
	id, err := uuid.Parse(fmt.Sprintf("%v", raw))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}
	return id, true
}
