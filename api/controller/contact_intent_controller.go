package controller

import (
	"net/http"
	"strconv"

	"github.com/gabrielfmcoelho/platform-core/bootstrap"
	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
	"github.com/gin-gonic/gin"
)

type ContactIntentController struct {
	ContactIntentUsecase domain.ContactIntentUsecase
	Env                  *bootstrap.Env
}

// @Summary Create a new contact intent
// @Description Submit a contact request from the website form (agende demonstração)
// @Tags ContactIntent
// @ID createContactIntent
// @Accept json
// @Produce json
// @Param contactIntent body domain.CreateContactIntent true "Contact intent object"
// @Success 201 {object} domain.SuccessResponse "Contact intent created successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request - Invalid input"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /contact-intent [post]
func (cic *ContactIntentController) CreateContactIntent(c *gin.Context) {
	var contactIntent domain.CreateContactIntent
	if err := c.ShouldBindJSON(&contactIntent); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid input: " + err.Error(),
		})
		return
	}

	err := cic.ContactIntentUsecase.Create(c, &contactIntent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to create contact intent: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, parser.ToSuccessResponse(gin.H{
		"message": "Contact intent submitted successfully. We will get back to you soon!",
	}))
}

// @Summary Get all contact intents
// @Description Get all contact intents from the database (admin only)
// @Tags ContactIntent
// @ID fetchContactIntents
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.SuccessResponse{data=[]domain.PublicContactIntent} "List of contact intents"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /contact-intents [get]
func (cic *ContactIntentController) FetchContactIntents(c *gin.Context) {
	contactIntents, err := cic.ContactIntentUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch contact intents: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(contactIntents))
}

// @Summary Update contact intent status
// @Description Update the status of a contact intent (admin only)
// @Tags ContactIntent
// @ID updateContactIntentStatus
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Contact Intent ID"
// @Param status body domain.UpdateContactIntentStatus true "Status update object"
// @Success 200 {object} domain.SuccessResponse "Status updated successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request - Invalid input"
// @Failure 404 {object} domain.ErrorResponse "Not Found - Contact intent not found"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /contact-intent/{id}/status [patch]
func (cic *ContactIntentController) UpdateContactIntentStatus(c *gin.Context) {
	// Parse ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid contact intent ID",
		})
		return
	}

	// Parse request body
	var statusUpdate domain.UpdateContactIntentStatus
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid input: " + err.Error(),
		})
		return
	}

	// Update status
	err = cic.ContactIntentUsecase.UpdateStatus(c, uint(id), statusUpdate.Status)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, domain.ErrorResponse{
				Message: "Contact intent not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to update status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(gin.H{
		"message": "Status updated successfully",
	}))
}
