package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-backend/internal/core/domain"
	"github.com/ericolvr/sec-backend/internal/core/services"
	"github.com/ericolvr/sec-backend/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var createDTO dto.UserRequest

	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	user := &domain.User{
		Name:     createDTO.Name,
		Mobile:   createDTO.Mobile,
		Password: createDTO.Password,
		UserType: createDTO.UserType,
		Status:   true,
	}

	if err := h.userService.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar usuário",
			"details": err.Error(),
		})
		return
	}

	responseDTO := dto.ToUserResponse(user)
	c.JSON(http.StatusCreated, responseDTO)
}

func (h *UserHandler) List(c *gin.Context) {

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parâmetro limit inválido",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parâmetro offset inválido",
		})
		return
	}

	users, err := h.userService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro interno do servidor",
		})
		return
	}

	// Convert domain users to response DTOs
	var responseUsers []dto.UserResponse
	for _, user := range users {
		responseUsers = append(responseUsers, dto.ToUserResponse(user))
	}

	c.JSON(http.StatusOK, responseUsers)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	user, err := h.userService.GetByID(c.Request.Context(), idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuário não encontrado",
		})
		return
	}

	responseDTO := dto.ToUserResponse(user)

	c.JSON(http.StatusOK, responseDTO)
}

func (h *UserHandler) GetByMobile(c *gin.Context) {
	mobile := c.Param("mobile")

	user, err := h.userService.GetByMobile(c.Request.Context(), mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuário não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := dto.ToUserResponse(user)

	c.JSON(http.StatusOK, responseDTO)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var updateDTO dto.UserUpdate
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		log.Printf("Error binding JSON in user update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	updateDTO.ID = idInt

	// Get current user to preserve status
	currentUser, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuário não encontrado",
		})
		return
	}

	user := &domain.User{
		ID:       updateDTO.ID,
		Name:     updateDTO.Name,
		Mobile:   updateDTO.Mobile,
		UserType: updateDTO.UserType,
		Status:   currentUser.Status, // Preserve current status
	}

	if err := h.userService.Update(c.Request.Context(), user); err != nil {
		if err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedUser, err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuário deletado com sucesso",
		"data":    dto.ToUserResponse(deletedUser),
	})
}
