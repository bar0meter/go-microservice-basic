package models

import (
	"encoding/json"
	"io"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/utils"
	"github.com/go-playground/validator"
)

// UserModel
type User struct {
	ID                int64  `json:"id" bson:"_id, omitempty"`
	Name              string `json:"name" bson:"name, omitempty"`
	Role              int64  `json:"role" bson:"role, omitempty"`
	Username          string `json:"username" validate:"required" bson:"username, omitempty" `
	Password          string `json:"password" bson:"password, omitempty"`
	ResetPasswordUUID string `json:"resetPassword" bson:"resetPassword, omitempty"`
	VerifyAccountUUID string `json:"verifyUser" bson:"verifyUser, omitempty"`
	Verified          int    `json:"verified" bson:"verified, omitempty"`
}

const (
	Id            = "id"
	ID            = "_id"
	Name          = "name"
	Role          = "role"
	Username      = "username"
	Password      = "password"
	ResetPassword = "resetPassword"
	VerifyUser    = "verifyUser"
	Verified      = "verified"
)

// ParseJSON for unmarshalling/decoding UserObject from payload
func (u *User) ParseJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// Serialize for marshalling UserObject and sending it in the response
func (u *User) Serialize(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// UserRepository interface => All available methods for UserCollection
type UserRepository interface {
	FindByID(ID int64) (*User, error)
	FindByUsername(Username string) (*User, error)
	Save(user *User) error
	ResetPassword(email string) (*utils.TagClaim, string, error)
	ValidateResetPassword(claim *utils.TagClaim) (bool, error)
	VerifyUser(userID int64, email string) (*utils.TagClaim, string, error)
	ValidateVerifyUser(claim *utils.TagClaim) (bool, error)
}
