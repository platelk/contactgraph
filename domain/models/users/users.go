package users

// ID is a unique identifier assign to each user
type ID uint32

// PhoneNumber define an international phone number
type PhoneNumber string

// User hold the definition of what is a user in the system
type User struct {
	// ID is a unique id across the system which is hidden to external users
	ID          ID          `json:"id"`
	NickName    string      `json:"nick_name"`
	PhoneNumber PhoneNumber `json:"phone_number"`
}
