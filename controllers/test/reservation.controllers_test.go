package controllers

import (
	"Interview_Hin_20240914/controllers"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestReservationsRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/reservations", controllers.ListReservations)
	r.POST("/reservations", controllers.CreateReservation)
	return r
}

func TestListReservations(t *testing.T) {
	// Setup test database connection
	setupTestDB()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: "available"}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Create test player
	testPlayer := &db.Player{Name: "Test Player"}
	err = mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Create test reservations
	now := time.Now()
	reservation1 := &db.Reservation{
		RoomID:          testRoom.ID,
		PlayerIDs:       []primitive.ObjectID{testPlayer.ID},
		BookingDatetime: now,
	}
	err = mgm.Coll(reservation1).Create(reservation1)
	assert.NoError(t, err)

	reservation2 := &db.Reservation{
		RoomID:          testRoom.ID,
		PlayerIDs:       []primitive.ObjectID{testPlayer.ID},
		BookingDatetime: now.Add(24 * time.Hour),
	}
	err = mgm.Coll(reservation2).Create(reservation2)
	assert.NoError(t, err)

	// Setup router
	router := setupTestReservationsRouter()

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reservations?room_id="+testRoom.ID.Hex()+"&date="+now.Format("2006-01-02"), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	reservations, ok := response.Data["reservations"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(reservations))

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(testRoom)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Reservation{}).Delete(reservation1)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Reservation{}).Delete(reservation2)
	assert.NoError(t, err)
}

func TestCreateReservation(t *testing.T) {
	// Setup test database connection
	setupTestDB()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: "available"}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Create test player
	testPlayer := &db.Player{Name: "Test Player"}
	err = mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Setup router
	router := setupTestReservationsRouter()

	// Test case
	reservationRequest := models.ReservationRequest{
		RoomID:          testRoom.ID.Hex(),
		PlayerIDs:       []string{testPlayer.ID.Hex()},
		BookingDatetime: time.Now().Add(24 * time.Hour).Format("2006/01/02 15:04"),
	}
	jsonValue, _ := json.Marshal(reservationRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify reservation creation
	reservationData, ok := response.Data["reservation"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testRoom.ID.Hex(), reservationData["roomId"])
	assert.Equal(t, 1, len(reservationData["playerIds"].([]interface{})))
	assert.Equal(t, testPlayer.ID.Hex(), reservationData["playerIds"].([]interface{})[0])

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(testRoom)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	reservationID, _ := primitive.ObjectIDFromHex(reservationData["id"].(string))
	_, err = mgm.Coll(&db.Reservation{}).DeleteOne(context.TODO(), bson.M{"id": reservationID})
	assert.NoError(t, err)
}