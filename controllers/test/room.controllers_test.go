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

func setupTestRoomRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/rooms", controllers.ListRooms)
	r.POST("/rooms", validators.CreateRoomValidator(), controllers.CreateRoom)
	r.GET("/rooms/:id", validators.PathIdValidator(), controllers.GetRoom)
	r.PUT("/rooms/:id", validators.PathIdValidator(), validators.UpdateRoomValidator(), controllers.UpdateRoom)
	r.DELETE("/rooms/:id", validators.PathIdValidator(), controllers.DeleteRoom)
	return r
}

func TestListRooms(t *testing.T) {
	// Setup test database connection
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test rooms
	room1 := &db.Room{Name: "Room 1", Status: enums.RoomStatusAvailable}
	room2 := &db.Room{Name: "Room 2", Status: enums.RoomStatusUnavailable}
	err := mgm.Coll(room1).Create(room1)
	assert.NoError(t, err)
	err = mgm.Coll(room2).Create(room2)
	assert.NoError(t, err)

	// Setup router
	router := setupTestRoomRouter()

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(room1)
	assert.NoError(t, err)
	err = mgm.Coll(&db.Room{}).Delete(room2)
	assert.NoError(t, err)
}

func TestCreateRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Setup router
	router := setupTestRoomRouter()

	// Test case
	roomRequest := models.RoomRequest{
		Name:        "Test Room",
		Description: "Test Description",
		Status:      enums.RoomStatusAvailable,
	}
	jsonValue, _ := json.Marshal(roomRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	// Clean up
	roomData, ok := response.Data["room"].(map[string]interface{})
	if !ok {
		t.Errorf("response.Data['room'] is not a map[string]interface{}")
		return
	}
	roomID, ok := roomData["id"].(string)
	if !ok {
		t.Errorf("roomData['id'] is not a string")
		return
	}
	room := &db.Room{}
	err = mgm.Coll(room).FindByID(roomID, room)
	assert.NoError(t, err)
	err = mgm.Coll(room).Delete(room)
	assert.NoError(t, err)
}

func TestGetRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Setup router
	router := setupTestRoomRouter()

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms/"+testRoom.ID.Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(testRoom)
	assert.NoError(t, err)
}

func TestUpdateRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Setup router
	router := setupTestRoomRouter()

	// Test case
	updateRequest := models.RoomRequest{
		Name:        "Updated Room",
		Description: "Updated Description",
		Status:      enums.RoomStatusUnavailable,
	}
	jsonValue, _ := json.Marshal(updateRequest)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/rooms/"+testRoom.ID.Hex(), bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify update
	updatedRoom := &db.Room{}
	err = mgm.Coll(updatedRoom).FindByID(testRoom.ID, updatedRoom)
	assert.NoError(t, err)
	assert.Equal(t, updateRequest.Name, updatedRoom.Name)
	assert.Equal(t, updateRequest.Description, updatedRoom.Description)
	assert.Equal(t, updateRequest.Status, updatedRoom.Status)

	// Clean up
	err = mgm.Coll(&db.Room{}).Delete(testRoom)
	assert.NoError(t, err)
}

func TestDeleteRoom(t *testing.T) {
	// Setup test database connection
	services.LoadConfig("../../.env")
	services.InitMongoDBTest()

	// Create test room
	testRoom := &db.Room{Name: "Test Room", Status: enums.RoomStatusAvailable}
	err := mgm.Coll(testRoom).Create(testRoom)
	assert.NoError(t, err)

	// Setup router
	router := setupTestRoomRouter()

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/rooms/"+testRoom.ID.Hex(), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify deletion
	deletedRoom := &db.Room{}
	err = mgm.Coll(deletedRoom).FindByID(testRoom.ID, deletedRoom)
	assert.Error(t, err)
}