package constants

type PrivacyType int

const (
	Private PrivacyType = iota
	Public
	Friend
)

type PostType int

const (
	PersonalPost PostType = iota
	GroupPost
)
