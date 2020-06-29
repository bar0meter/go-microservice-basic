package db

import (
	"context"
	log "github.com/frost060/go-microservice-basic/rest-api-mongo/logging"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/models"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepo model having user collection and logger instances
type UserRepo struct {
	Instance *mongo.Collection
}

// NewUserRepo creates a new user repo
func NewUserRepo(db *mongo.Database) *UserRepo {
	instance := db.Collection(USER)
	return &UserRepo{instance}
}

// FindByID returns a user from db given user id
func (u *UserRepo) FindByID(id int64) (*models.User, error) {
	var user models.User

	err := u.Instance.FindOne(context.TODO(),
		bson.M{models.ID: id}).Decode(&user)
	if err != nil {
		log.Error("FindByID, error occurred while querying DB", "error", err)
		return nil, err
	}

	return &user, nil
}

// FindByUsername returns a user from db given username
func (u *UserRepo) FindByUsername(username string) (*models.User, error) {
	var user models.User

	err := u.Instance.FindOne(context.TODO(),
		bson.M{models.Username: username}).Decode(&user)
	if err != nil {
		log.Error("FindByUsername, error occurred while querying DB", "error", err)
		return nil, err
	}

	return &user, nil
}

// Save saves a user to the db
func (u *UserRepo) Save(user *models.User) error {
	user.ID = utils.GenerateID()
	opts := options.Update().SetUpsert(true)
	findQuery := bson.M{models.Username: user.Username}

	updateResult, err := u.Instance.UpdateOne(context.TODO(), findQuery, &user, opts)
	if err != nil {
		log.Error("Error while saving document to users db", "error", err)
		return err
	}

	log.Info("Inserted a single user document: ", updateResult.UpsertedID)
	return nil
}

// ResetPassword is used for resetting user password.
// Password expires in 1 day and can be used only once.
// Note that we are not returning error here if we dont find the user for security reasons.
// Always send status OK, saying will send an email if email belongs to an user.
func (u *UserRepo) ResetPassword(email string) (*utils.TagClaim, string, error) {
	var user models.User
	findQuery := bson.M{models.Username: email}

	err := u.Instance.FindOne(context.TODO(), findQuery).Decode(&user)
	if err != nil {
		log.Error("RestPassword, error occurred while querying DB", "error", err)
		return nil, "", nil
	}

	uuidToken, err := uuid.NewRandom()
	if err != nil {
		log.Error("ResetPassword, error occurred while setting reset password uuid token", "error", err)
		return nil, "", err
	}

	opts := options.Update().SetUpsert(false)
	updateQuery := bson.D{{
		"$set", bson.M{models.ResetPassword: uuidToken.String()},
	}}
	_, err = u.Instance.UpdateOne(context.TODO(), findQuery, updateQuery, opts)
	if err != nil {
		log.Error("ResetPassword, error occurred while setting reset password uuid token", "error", err)
		return nil, "", err
	}

	return utils.GenerateResetPasswordTag(email, uuidToken.String()), user.Name, nil
}

// ValidateResetPassword validates claim in the url by checking if Identity and uuid combination exists in the user db.
// If it exists then we just remove the saved uuid (As tag is one time use only).
func (u *UserRepo) ValidateResetPassword(claim *utils.TagClaim) (bool, error) {
	var user models.User
	findQuery := bson.M{models.Username: claim.Identity, models.ResetPassword: claim.ID}
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
	updateQuery := bson.D{{
		"$unset", bson.M{models.ResetPassword: 1},
	}}

	err := u.Instance.FindOneAndUpdate(context.TODO(), findQuery, updateQuery, opts).Decode(&user)
	if err != nil {
		return false, err
	}

	return user.ResetPasswordUUID == "", nil
}

func (u *UserRepo) VerifyUser(userID int64, email string) (*utils.TagClaim, string, error) {
	var user models.User

	findQuery := bson.M{models.ID: userID}
	if email != "" {
		findQuery[models.Username] = email
	}

	uuidToken, err := uuid.NewRandom()
	if err != nil {
		log.Error("VerifyUser, error occurred while setting verify email uuid token", "error", err)
		return nil, "", err
	}

	// Set options for FindOneAndUpdate
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
	updateQuery := bson.D{{
		"$set", bson.M{models.VerifyUser: uuidToken.String()},
	}}

	err = u.Instance.FindOneAndUpdate(context.TODO(), findQuery, updateQuery, opts).Decode(&user)
	if err != nil {
		return nil, "", err
	}

	return utils.GenerateVerifyEmailTag(user.Username, uuidToken.String()), user.Name, nil
}

// ValidateVerifyUser validates the verify link user clicks on and sets verified flag in db.
func (u *UserRepo) ValidateVerifyUser(claim *utils.TagClaim) (bool, error) {
	var user models.User
	findQuery := bson.M{models.Username: claim.Identity, models.VerifyUser: claim.ID}
	opts := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
	updateQuery := bson.D{{
		"$unset", bson.M{models.VerifyUser: 1},
	}, {
		"$set", bson.M{models.Verified: 1},
	}}

	err := u.Instance.FindOneAndUpdate(context.TODO(), findQuery, updateQuery, opts).Decode(&user)
	if err != nil {
		return false, err
	}

	return user.VerifyAccountUUID == "" && user.Verified == 1, nil
}
