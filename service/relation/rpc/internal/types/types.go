package types

type AddIsFollowMsg struct {
	FromId int64 ` json:"fromId,omitempty"`
	ToId   int64 ` json:"toId,omitempty"`
}
