package enums

type LogAction string

const (
	LogActionRegister        LogAction = "register"
	LogActionLogin           LogAction = "login"
	LogActionLogout          LogAction = "logout"
	LogActionEnterRoom       LogAction = "enter_room"
	LogActionExitRoom        LogAction = "exit_room"
	LogActionJoinChallenge   LogAction = "join_challenge"
	LogActionChallengeResult LogAction = "challenge_result"
)

func (s LogAction) IsValid() bool {
	switch s {
	case LogActionRegister, LogActionLogin, LogActionLogout, LogActionEnterRoom, LogActionExitRoom, LogActionJoinChallenge, LogActionChallengeResult:
		return true
	}
	return false
}