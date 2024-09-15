package services

import (
	"Interview_Hin_20240914/enums"
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"Interview_Hin_20240914/services"
	"testing"

	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig()
	services.InitMongoDBTest()
	
	// Test case
	request := models.RoomRequest{
		Name:        "Test Room",
		Description: "Test Description",
		Status:      enums.RoomStatusAvailable,
	}

	room, err := services.CreateRoom(request)
	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, request.Name, room.Name)
	assert.Equal(t, request.Description, room.Description)
	assert.Equal(t, request.Status, room.Status)

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(room)
	assert.NoError(t, err)
}

func TestGetRooms(t *testing.T) {
	// Setup test database connection
	services.LoadConfig()
	services.InitMongoDBTest()

	// Create test rooms
	room1 := &db.Room{Name: "Unit Test Room 1", Status: enums.RoomStatusAvailable}
	room2 := &db.Room{Name: "Unit Test Room 2", Status: enums.RoomStatusUnavailable}
	err := mgm.Coll(room1).Create(room1)
	assert.NoError(t, err)
	err = mgm.Coll(room2).Create(room2)
	assert.NoError(t, err)

	// Test case
	rooms, err := services.GetRooms(0, 10)
	assert.NoError(t, err)
	
	// Check if the created rooms are in the result
	foundRoom1 := false
	foundRoom2 := false
	for _, room := range rooms {
		if room.ID == room1.ID {
			foundRoom1 = true
		}
		if room.ID == room2.ID {
			foundRoom2 = true
		}
	}
	assert.True(t, foundRoom1, "Room 1 not found in the result")
	assert.True(t, foundRoom2, "Room 2 not found in the result")

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(room1)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Room{}).Delete(room2)
	assert.NoError(t, err)
}

func TestGetRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig()
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Test case
	room, err := services.GetRoom(testRoom.ID)
	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, testRoom.ID, room.ID)
	assert.Equal(t, testRoom.Name, room.Name)

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(testRoom)
	assert.NoError(t, err)
}

func TestUpdateRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig()
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Test case
	updateRequest := &models.RoomRequest{
		Name:   "Updated Room",
		Status: enums.RoomStatusUnavailable,
	}
	err = services.UpdateRoom(testRoom.ID, updateRequest)
	assert.NoError(t, err)

	// Verify update
	updatedRoom, err := services.GetRoom(testRoom.ID)
	assert.NoError(t, err)
	assert.Equal(t, updateRequest.Name, updatedRoom.Name)
	assert.Equal(t, updateRequest.Status, updatedRoom.Status)

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(testRoom)
	assert.NoError(t, err)
}

func TestDeleteRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig()
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Test case
	err = services.DeleteRoom(testRoom.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = services.GetRoom(testRoom.ID)
	assert.Error(t, err)
}