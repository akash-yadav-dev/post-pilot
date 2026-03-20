package social

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ConnectAccount(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}

	var req ConnectAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.svc.ConnectAccount(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect social account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *Handler) ListAccounts(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}

	accounts, err := h.svc.ListAccounts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list social accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid social account id"})
		return
	}

	if err := h.svc.DeleteAccount(c.Request.Context(), userID, accountID); err != nil {
		if IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "social account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete social account"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) PublishPost(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}

	postID, err := uuid.Parse(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	resp, err := h.svc.PublishPostNow(c.Request.Context(), userID, postID)
	if err != nil {
		if IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to publish post: %v", err)})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func userIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("auth_user_id")
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
