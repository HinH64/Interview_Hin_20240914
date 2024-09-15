package services

import (
	"Interview_Hin_20240914/enums"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"Interview_Hin_20240914/services"
	"testing"

	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProcessPayment(t *testing.T) {
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

	// Test cases
	testCases := []struct {
		name           string
		paymentRequest models.PaymentRequest
		expectError    bool
	}{
		{
			name: "Valid Credit Card Payment",
			paymentRequest: models.PaymentRequest{
				PlayerID:      testPlayer.ID.Hex(),
				PaymentMethod: enums.PaymentMethodCreditCard,
				Amount:        100,
			},
			expectError: false,
		},
		{
			name: "Valid Bank Transfer Payment",
			paymentRequest: models.PaymentRequest{
				PlayerID:      testPlayer.ID.Hex(),
				PaymentMethod: enums.PaymentMethodBankTransfer,
				Amount:        200,
			},
			expectError: false,
		},
		{
			name: "Invalid Payment Method",
			paymentRequest: models.PaymentRequest{
				PlayerID:      testPlayer.ID.Hex(),
				PaymentMethod: "InvalidMethod",
				Amount:        300,
			},
			expectError: true,
		},
		{
			name: "Invalid Player ID",
			paymentRequest: models.PaymentRequest{
				PlayerID:      "invalidID",
				PaymentMethod: enums.PaymentMethodCreditCard,
				Amount:        400,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payment, err := services.ProcessPayment(tc.paymentRequest)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, payment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payment)
				assert.Equal(t, tc.paymentRequest.PlayerID, payment.PlayerID.Hex())
				assert.Equal(t, tc.paymentRequest.PaymentMethod, payment.PaymentMethod)
				assert.Equal(t, tc.paymentRequest.Amount, payment.Amount)


				// Clean up the created payment
				err = mgm.Coll(&db.Payment{}).Delete(payment)
				assert.NoError(t, err)
			}
		})
	}

	// Clean up the test player
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(testLevel).Delete(testLevel)
	assert.NoError(t, err)
}

func TestGetPaymentByID(t *testing.T) {
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

	// Test case: Get existing payment
	payment, err := services.GetPaymentByID(testPayment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, testPayment.ID, payment.ID)
	assert.Equal(t, testPlayer.ID, payment.PlayerID)
	assert.Equal(t, testPayment.PaymentMethod, payment.PaymentMethod)
	assert.Equal(t, testPayment.Amount, payment.Amount)
	assert.Equal(t, testPayment.Status, payment.Status)
	assert.Equal(t, testPayment.TransactionId, payment.TransactionId)

	// Test case: Get non-existing payment
	nonExistingID := primitive.NewObjectID()
	payment, err = services.GetPaymentByID(nonExistingID)
	assert.Error(t, err)
	assert.Nil(t, payment)

	// Clean up
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Payment{}).Delete(testPayment)
	assert.NoError(t, err)
	err = mgm.Coll(testLevel).Delete(testLevel)
	assert.NoError(t, err)
}