package services

import (
	"Interview_Hin_20240914/models"
	db "Interview_Hin_20240914/models/db"
	"errors"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateReservation(request models.ReservationRequest) (*db.Reservation, error) {
	roomID, err := primitive.ObjectIDFromHex(request.RoomID)
	if err != nil {
		return nil, errors.New("invalid room ID")
	}

	playerIDs := make([]primitive.ObjectID, len(request.PlayerIDs))
	for i, playerID := range request.PlayerIDs {
		playerIDs[i], err = primitive.ObjectIDFromHex(playerID)
		if err != nil {
			return nil, errors.New("invalid player ID")
		}
	}

	// Check if room exists
	room := &db.Room{}
	err = mgm.Coll(room).FindByID(roomID, room)
	if err != nil {
		return nil, errors.New("room not found")
	}

	// Check if players exist
	for _, playerID := range playerIDs {
		player := &db.Player{}
		err = mgm.Coll(player).FindByID(playerID, player)
		if err != nil {
			return nil, errors.New("player not found")
		}
	}

	// Check if the room is available at the requested time
	conflictingReservation := &db.Reservation{}
	err = mgm.Coll(conflictingReservation).First(bson.M{
		"roomId":          roomID,
		"bookingDatetime": request.BookingDatetime,
	}, conflictingReservation)

	if err == nil {
		return nil, errors.New("room is already reserved at this time")
	}

	bookingDatetime, err := time.Parse("2006/01/02 15:04", request.BookingDatetime)
	if err != nil {
		return nil, errors.New("invalid date format for bookingDatetime")
	}

	reservation := &db.Reservation{
		RoomID:          roomID,
		PlayerIDs:       playerIDs,
		BookingDatetime: bookingDatetime,
	}

	err = mgm.Coll(reservation).Create(reservation)
	if err != nil {
		return nil, errors.New("cannot create new reservation")
	}

	return reservation, nil
}

func GetReservations(roomID string, date *time.Time, page, limit int) ([]db.GetReservation, error) {
	var reservations []db.GetReservation

	filter := bson.M{}
	if roomID != "" {
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			return nil, errors.New("invalid room ID")
		}
		filter["roomId"] = roomObjID
	}
	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		filter["bookingDatetime"] = bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		}
	}

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Reservation{}).SimpleFind(&reservations, filter, findOptions)
	if err != nil {
		return nil, errors.New("error fetching reservations")
	}

	// Fetch room names and player names
	for i := range reservations {
		// Get room name
		room := &db.Room{}
		err := mgm.Coll(room).FindByID(reservations[i].RoomID, room)
		if err == nil {
			reservations[i].RoomName = room.Name
		}

		// Get player names
		reservations[i].PlayerNames = make([]string, len(reservations[i].PlayerIDs))
		for j, playerID := range reservations[i].PlayerIDs {
			player := &db.Player{}
			err := mgm.Coll(player).FindByID(playerID, player)
			if err == nil {
				reservations[i].PlayerNames[j] = player.Name
			} else {
				reservations[i].PlayerNames[j] = "Unknown Player"
			}
		}
	}

	return reservations, nil
}