package controllers

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/enums"
	"Interview_Hin_20240914/middlewares/validators"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"Interview_Hin_20240914/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
)

func setupTestPaymentRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/payments", validators.CreatePaymentValidator(), controllers.CreatePayment)
	r.GET("/payments/:id", controllers.GetPayment)
	return r
}

func TestCreatePayment(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test level
	testLevel := &db.Level{Name: "Test Level"}
	err := mgm.Coll(testLevel).Create(testLevel)
	assert.NoError(t, err)
	
	// Create test player
	testPlayer := &db.Player{
		Name:           "Test Player",
		Balance:        1000,
		ChallengeCount: 0,
		LevelID:        testLevel.ID,
	}
	err = mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Setup router
	router := setupTestPaymentRouter()

	// Test case
	paymentRequest := models.PaymentRequest{
		PlayerID:      testPlayer.ID.Hex(),
		PaymentMethod: enums.PaymentMethodCreditCard,
		Amount:        100,
	}
	jsonValue, _ := json.Marshal(paymentRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify payment creation
	paymentData, ok := response.Data["payment"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testPlayer.ID.Hex(), paymentData["playerId"])
	assert.Equal(t, string(enums.PaymentMethodCreditCard), paymentData["paymentMethod"])
	assert.Equal(t, float64(100), paymentData["amount"])

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Level{}).Delete(testLevel)
	assert.NoError(t, err)
	payment := &db.Payment{}
	err = mgm.Coll(payment).FindByID(paymentData["id"], payment)
	assert.NoError(t, err)
	err = mgm.Coll(payment).Delete(payment)
	assert.NoError(t, err)
}

func TestGetPayment(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test level
	testLevel := &db.Level{Name: "Test Level"}
	err := mgm.Coll(testLevel).Create(testLevel)
	assert.NoError(t, err)
	// Create test player
	testPlayer := &db.Player{
		Name:           "Test Player",
		Balance:        1000,
		ChallengeCount: 0,
		LevelID:        testLevel.ID,
	}
	err = mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Create test payment
	testPayment := &db.Payment{
		PlayerID:      testPlayer.ID,
		PaymentMethod: enums.PaymentMethodCreditCard,
		Amount:        100,
		Status:        "completed",
		TransactionId: "test-transaction-id",
	}
	err = mgm.Coll(testPayment).Create(testPayment)
	assert.NoError(t, err)

	// Setup router
	router := setupTestPaymentRouter()

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/payments/"+testPayment.ID.Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify payment retrieval
	paymentData, ok := response.Data["payment"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testPayment.ID.Hex(), paymentData["id"])
	assert.Equal(t, testPlayer.ID.Hex(), paymentData["playerId"])
	assert.Equal(t, string(enums.PaymentMethodCreditCard), paymentData["paymentMethod"])
	assert.Equal(t, float64(100), paymentData["amount"])
	assert.Equal(t, "completed", paymentData["status"])
	assert.Equal(t, "test-transaction-id", paymentData["transactionId"])

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Level{}).Delete(testLevel)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Payment{}).Delete(testPayment)
	assert.NoError(t, err)
}