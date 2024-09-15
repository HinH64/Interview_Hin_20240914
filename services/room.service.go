package services

import (
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateRoom(request models.RoomRequest) (*db.Room, error) {
	room := &db.Room{
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
	}

	err := mgm.Coll(room).Create(room)
	if err != nil {
		return nil, errors.New("cannot create new room")
	}

	return room, nil
}

func GetRooms(page int, limit int) ([]db.Room, error) {
	var rooms []db.Room

	if page < 0 {
		page = 0
	}

	skip := page * limit

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Room{}).SimpleFind(&rooms, bson.M{}, findOptions)

	if err != nil {
		return nil, fmt.Errorf("cannot find rooms: %w", err)
	}

	return rooms, nil
}

func GetRoom(roomID primitive.ObjectID) (*db.Room, error) {
	room := &db.Room{}
	err := mgm.Coll(room).FindByID(roomID, room)
	if err != nil {
		return nil, errors.New("cannot find room")
	}
	return room, nil
}

func UpdateRoom(roomID primitive.ObjectID, request *models.RoomRequest) error {
	room := &db.Room{}
	err := mgm.Coll(room).FindByID(roomID, room)
	if err != nil {
		return errors.New("cannot find room")
	}

	room.Name = request.Name
	room.Description = request.Description
	room.Status = request.Status

	err = mgm.Coll(room).Update(room)
	if err != nil {
		return errors.New("cannot update room")
	}

	return nil
}

func DeleteRoom(roomID primitive.ObjectID) error {
	room := &db.Room{}
	err := mgm.Coll(room).FindByID(roomID, room)
	if err != nil {
		return errors.New("cannot find room")
	}

	err = mgm.Coll(room).Delete(room)
	if err != nil {
		return errors.New("cannot delete room")
	}

	return nil
}