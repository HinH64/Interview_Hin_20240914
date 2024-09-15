package controllers

import (
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePayment godoc
// @Summary Process a payment
// @Description Process a payment with the given details
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body models.PaymentRequest true "Process payment"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /payments [post]
func CreatePayment(c *gin.Context) {
	var requestBody models.PaymentRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	payment, err := services.ProcessPayment(requestBody)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"payment": payment}
	response.SendResponse(c)
}

// GetPayment godoc
// @Summary Get payment details
// @Description Get details of a specific payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /payments/{id} [get]
func GetPayment(c *gin.Context) {
	idHex := c.Param("id")
	paymentId, _ := primitive.ObjectIDFromHex(idHex)
	payment, err := services.GetPaymentByID(paymentId)
	if err != nil {
		models.SendErrorResponse(c, http.StatusNotFound, "Payment not found")
		return
	}
	models.SendResponseData(c, gin.H{"payment": payment})
}