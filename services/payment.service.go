package services

import (
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"errors"
	"math/rand"
	"time"

	"Interview_Hin_20240914/enums"

	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessPayment(request models.PaymentRequest) (*db.Payment, error) {
	playerID, err := primitive.ObjectIDFromHex(request.PlayerID)
	if err != nil {
		return nil, errors.New("invalid player ID")
	}

	payment := &db.Payment{
		PlayerID:      playerID,
		PaymentMethod: request.PaymentMethod,
		Amount:        request.Amount,
		Status:        "pending",
		TransactionId: "",
	}

	var transactionId string

	switch request.PaymentMethod {
	case enums.PaymentMethodCreditCard:
		transactionId, err = processCreditCardPayment()
	case enums.PaymentMethodBankTransfer:
		transactionId, err = processBankTransferPayment()
	case enums.PaymentMethodThirdParty:
		transactionId, err = processThirdPartyPayment()
	case enums.PaymentMethodBlockchain:
		transactionId, err = processBlockchainPayment()
	default:
		return nil, errors.New("invalid payment method")
	}

	if err != nil {
		payment.Status = "failed"
		payment.Details += " Error: " + err.Error()
	} else {
		payment.Status = "success"
		payment.Details += " Payment success: " + time.Now().Format("2006-01-02 15:04:05")
		payment.TransactionId = transactionId

		// Get current player info
		player, err := GetPlayerByID(playerID)
		if err != nil {
			return nil, errors.New("failed to get player information")
		}

		// Update player's balance
		updatedBalance := player.Balance + request.Amount
		updateRequest := &models.PlayerRequest{
			Name:           player.Name,
			LevelID:        player.LevelID.Hex(),
			Balance:        updatedBalance,
			ChallengeCount: player.ChallengeCount,
		}

		err = UpdatePlayer(playerID, updateRequest)
		if err != nil {
			return nil, errors.New("failed to update player's balance")
		}
	}

	err = mgm.Coll(payment).Create(payment)
	if err != nil {
		return nil, errors.New("cannot create new payment record")
	}

	return payment, nil
}

func GetPaymentByID(id primitive.ObjectID) (*db.Payment, error) {
	payment := &db.Payment{}
	
	err := mgm.Coll(payment).FindByID(id, payment)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	return payment, nil
}

func processCreditCardPayment() (string, error) {
	// Simulate credit card payment gateway call with 50% chance of success or failure
	if rand.Intn(2) == 0 {
		return "", errors.New("creditCard payment failed")
	}
	return uuid.New().String(), nil
}

func processBankTransferPayment() (string, error) {
	// Simulate bank transfer service call
	if rand.Intn(2) == 0 {
		return "", errors.New("bankTransfer payment failed")
	}
	return uuid.New().String(), nil
}

func processThirdPartyPayment() (string, error) {
	// Simulate third-party payment platform call
	if rand.Intn(2) == 0 {
		return "", errors.New("thirdParty payment failed")
	}
	return uuid.New().String(), nil
}

func processBlockchainPayment() (string, error) {
	// Simulate blockchain payment gateway call
	if rand.Intn(2) == 0 {
		return "", errors.New("blockchain payment failed")
	}
	return uuid.New().String(), nil
}