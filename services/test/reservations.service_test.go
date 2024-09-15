package services

import (
	"Interview_Hin_20240914/enums"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"Interview_Hin_20240914/services"
	"testing"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func TestCreateReservation(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Create test player
	testPlayer := &db.Player{Name: "Test Player"}
	err = mgm.Coll(testPlayer).Create(testPlayer)
	assert.NoError(t, err)

	// Test case
	bookingTime := time.Now().Add(24 * time.Hour)
	request := models.ReservationRequest{
		RoomID:          testRoom.ID.Hex(),
		PlayerIDs:       []string{testPlayer.ID.Hex()},
		BookingDatetime: bookingTime.Format("2006/01/02 15:04"),
	}

	reservation, err := services.CreateReservation(request)
	assert.NoError(t, err)
	assert.NotNil(t, reservation)
	assert.Equal(t, testRoom.ID, reservation.RoomID)
	assert.Equal(t, []primitive.ObjectID{testPlayer.ID}, reservation.PlayerIDs)
	assert.Equal(t, bookingTime.Format("2006-01-02 15:04"), reservation.BookingDatetime.Format("2006-01-02 15:04"))

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(testRoom)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Player{}).Delete(testPlayer)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Reservation{}).Delete(reservation)
	assert.NoError(t, err)
}

func TestGetReservations(t *testing.T) {
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
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

	// Test case
	reservations, err := services.GetReservations(testRoom.ID.Hex(), &now, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, reservations, 1)
	assert.Equal(t, reservation1.ID, reservations[0].ID)

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