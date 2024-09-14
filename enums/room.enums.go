package enums

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusUnavailable RoomStatus = "unavailable"
	RoomStatusMaintenance RoomStatus = "maintenance"
)

func (s RoomStatus) IsValid() bool {
	switch s {
	case RoomStatusAvailable, RoomStatusUnavailable, RoomStatusMaintenance:
		return true
	}
	return false
}