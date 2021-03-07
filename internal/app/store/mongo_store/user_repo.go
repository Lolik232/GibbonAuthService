package mongo_store

import (
	"auth-server/internal/app/model"
	"auth-server/internal/app/store"
	errors "auth-server/pkg/errors/types"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	//UserCommonFields struct contains the common fields
	UserCommonFields struct {
		ID             primitive.ObjectID  `bson:"_id,omitempty"`
		UserName       string              `bson:"username,omitempty"`
		PasswordHash   string              `bson:"password_hash,omitempty"`
		Email          string              `bson:"email,omitempty"`
		EmailConfirmed bool                `bson:"email_confirmed,omitempty"`
		CreatedAt      *primitive.DateTime `bson:"created_at,omitempty"`
		UserInfo       map[string]string   `bson:"user_info,omitempty"`
	}
	//User represents the "Users" collection
	User struct {
		ID             primitive.ObjectID  `bson:"_id,omitempty"`
		UserName       string              `bson:"username,omitempty"`
		PasswordHash   string              `bson:"password_hash,omitempty"`
		Email          string              `bson:"email,omitempty"`
		EmailConfirmed bool                `bson:"email_confirmed,omitempty"`
		CreatedAt      *primitive.DateTime `bson:"created_at,omitempty"`
		UserInfo       map[string]string   `bson:"user_info,omitempty"`
		UserSessions   []UserSession       `bson:"user_sessions,omitempty"`
		ClientRoles    []ClientRole        `bson:"user_roles,omitempty"`
	}
	//ClientRole represent user roles attached document in "User"
	ClientRole struct {
		ClientID primitive.ObjectID `bson:"client_id,omitempty"`
		Roles    []string           `bson:"roles,omitempty"`
	}

	//UserSession represent attached "Sessions" document in "User"
	UserSession struct {
		ID             primitive.ObjectID `bson:"id,omitempty"`
		ClientID       primitive.ObjectID `bson:"client_id,omitempty"`
		Device         string             `bson:"device,omitempty"`
		LastActiveDate primitive.DateTime `bson:"last_active_time,omitempty"`
	}

	//UserClient represents an aggregation result-set for two collections
	UserClient struct {
		ID             primitive.ObjectID  `bson:"_id,omitempty"`
		UserName       string              `bson:"username,omitempty"`
		PasswordHash   string              `bson:"password_hash,omitempty"`
		Email          string              `bson:"email,omitempty"`
		EmailConfirmed bool                `bson:"email_confirmed,omitempty"`
		CreatedAt      *primitive.DateTime `bson:"created_at,omitempty"`
		UserInfo       map[string]string   `bson:"user_info,omitempty"`
		UserSessions   []UserSessionClient `bson:"user_sessions,omitempty"`
		Roles          []UserClientRole    `bson:"roles,omitempty"`
	}

	//UserSessionClient represent attached "Sessions" document in "UserClient"
	UserSessionClient struct {
		SessionID      string             `bson:"session_id,omitempty"`
		Client         Client             `bson:"client,omitempty"`
		Device         string             `bson:"device,omitempty"`
		LastActiveDate primitive.DateTime `bson:"last_active_time,omitempty"`
	}
	//UserClientRole represent attached user roles document in "UserClient"
	UserClientRole struct {
		Client Client   `bson:"client,omitempty"`
		Roles  []string `bson:"roles,omitempty"`
	}
	//
	UserRepo struct {
		store    *Store
		usersCol *mongo.Collection
	}
)

func (u UserRepo) fetch(ctx context.Context, query bson.M, params *store.UserFields) (*UserClient, error) {
	if params == nil {
		params = new(store.UserFields)
	}
	projection := bson.M{}
	projection["_id"] = true
	if params.UserName {
		projection["username"] = params.UserName
	}
	if params.Email {
		projection["email"] = params.Email
	}
	if params.CreatedAt {
		projection["created_at"] = params.CreatedAt
	}
	if params.UserInfo {
		projection["user_info"] = params.UserInfo
	}

	if params.UserRoles {
		projection["user_info"] = params.UserRoles
	}
	if params.UserPasswordHash {
		projection["password_hash"] = params.UserPasswordHash
	}
	var usr *UserClient
	if params.UserSessions == false {
		opt := options.FindOne().SetProjection(projection)
		err := u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				return nil, errors.ErrInvalidArgument.Newf("")
			case mongo.ErrClientDisconnected:
				return nil, errors.ErrDatabaseDown.New("")
			}
		}
	} else {
		proj := bson.D{{"$project", projection}}
		match := bson.D{{"$match", query}}
		pip := bson.D{{"$lookup", bson.D{
			{"from", "clients"},
			{"localField", "user_sessions.client_id"},
			{"foreignField", "_id"},
			{"as", "user_sessions"}}}}
		cur, err := u.usersCol.Aggregate(ctx, mongo.Pipeline{match, pip, proj})
		defer cur.Close(ctx)
		if err != nil {
			switch err {
			case mongo.ErrClientDisconnected:
				return nil, errors.NoType.New("")
			}
		}
		if cur != nil && cur.Next(ctx) {
			err = cur.Decode(&usr)
			if err != nil {
				return nil, errors.NoType.New("")
			}
		} else {
			return nil, errors.ErrInvalidArgument.Newf("")
		}
	}
	return usr, nil
}

func (u UserRepo) FindById(ctx context.Context, id string, params *store.UserFields) (*model.User, error) {
	uid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrInvalidArgument.Newf("Invalid userID: %s", id)
	}
	query := bson.M{"_id": uid}
	usr, err := u.fetch(ctx, query, params)
	if err != nil {
		if errors.GetType(err) == errors.ErrInvalidArgument {
			return nil, errors.ErrInvalidArgument.Newf("Invalid userID %s", id)
		}
		return nil, err
	}

	return ToUserClient(usr), nil
}
func (u UserRepo) FindByName(ctx context.Context, username string, params *store.UserFields) (*model.User, error) {
	query := bson.M{
		"username": username,
	}
	usr, err := u.fetch(ctx, query, params)
	if err != nil {
		return nil, err
	}
	return ToUserClient(usr), nil
}
func (u UserRepo) FindByEmail(ctx context.Context, email string, params *store.UserFields) (*model.User, error) {
	query := bson.M{
		"email": email,
	}
	usr, err := u.fetch(ctx, query, params)
	if err != nil {
		return nil, err
	}
	return ToUserClient(usr), nil
}

func (u UserRepo) Create(ctx context.Context, user *model.User) (string, error) {
	usr := ToDb(user)
	createdTime := primitive.NewDateTimeFromTime(time.Now())
	usr.CreatedAt = &createdTime
	cur, err := u.usersCol.InsertOne(ctx, usr)
	if err != nil {
		switch err {
		case mongo.ErrNilValue, mongo.ErrNoDocuments, mongo.ErrClientDisconnected:
			return "", errors.ErrDatabaseDown.New("")
		}
	}
	return cur.InsertedID.(primitive.ObjectID).Hex(), nil

}
func (u UserRepo) Update(ctx context.Context, userID string, user *model.User) error {
	panic("implement me")
}

func (u UserRepo) FindSessions(ctx context.Context, id string) (*[]model.UserSession, error) {
	panic("implement me")
}
func (u UserRepo) CheckSession(ctx context.Context, id string) error {
	panic("implement me")
}

func (u UserRepo) FindUserClientRoles(ctx context.Context, userID, clientID string) ([]model.UserRole, error) {
	panic("implement me")
}
func (u UserRepo) DeleteById(ctx context.Context, userID string) error {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.ErrInvalidArgument.Newf("Invalid userID %s", userID)
	}
	query := bson.M{
		"_id": ID,
	}
	result, err := u.usersCol.DeleteOne(ctx, query)
	if result != nil && result.DeletedCount == 0 {
		return errors.ErrInvalidArgument.New("Invalid userID.")
	}
	if err != nil {
		return errors.NoType.Newf("")
	}
	return nil
}
func (u UserRepo) DeleteByName(ctx context.Context, username string) error {
	query := bson.M{
		"username": username,
	}
	result, err := u.usersCol.DeleteOne(ctx, query)
	if result != nil && result.DeletedCount == 0 {
		return errors.ErrInvalidArgument.New("Invalid username.")
	}
	if err != nil {
		return errors.NoType.Newf("")
	}
	return nil
}

func (u UserRepo) CheckPassByID(ctx context.Context, userID, passwordHash string) error {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.ErrInvalidArgument.Newf("Invalid userID %s", userID)
	}
	opt := options.FindOne().SetProjection(
		bson.M{
			"_id": 1,
		},
	)
	query := bson.M{
		"_id":           ID,
		"password_hash": passwordHash,
	}
	var user *User
	err = u.usersCol.FindOne(ctx, query, opt).Decode(&user)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.ErrInvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return errors.ErrDatabaseDown.New("")
		}
	}
	if user.ID.Hex() != userID {
		return errors.NoType.New("")
	}
	return nil
}
func (u UserRepo) CheckPassByName(ctx context.Context, username, passwordHash string) error {

	query := bson.M{
		"username":      username,
		"password_hash": passwordHash,
	}
	proj := bson.M{
		"_id":      1,
		"username": 1,
	}
	opt := options.FindOne().SetProjection(proj)
	var usr *User
	err := u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.ErrInvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return errors.ErrDatabaseDown.New("")
		}
	}
	if usr.UserName != username {
		return errors.ErrInvalidPasswordOrUsername.New("")
	}
	return nil
}
func (u UserRepo) CheckPassByEmail(ctx context.Context, email, passwordHash string) error {
	query := bson.M{
		"email":         email,
		"password_hash": passwordHash,
	}
	proj := bson.M{
		"_id":   1,
		"email": 1,
	}
	opt := options.FindOne().SetProjection(proj)
	var usr *User
	err := u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.ErrInvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return errors.ErrDatabaseDown.New("")
		}
	}
	if usr.Email != email {
		return errors.ErrInvalidPasswordOrUsername.New("")
	}
	return nil
}

//Convert "User" database model to "DTO" model without ObjectID's
func ToUserClient(usr *UserClient) *model.User {
	sessions := make([]model.UserSession, 0)
	roles := make([]model.UserRole, 0)
	for _, v := range usr.UserSessions {
		session := model.UserSession{
			SessionID:      v.SessionID,
			ClientName:     v.Client.ClientName,
			Device:         v.Device,
			LastActiveTime: v.LastActiveDate.Time(),
		}
		sessions = append(sessions, session)
	}
	for _, v := range usr.Roles {
		role := model.UserRole{
			ClientName: v.Client.ClientName,
			Roles:      v.Roles,
		}
		roles = append(roles, role)
	}
	var date *time.Time

	if usr.CreatedAt != nil {
		date := &time.Time{}
		*date = usr.CreatedAt.Time()
	}

	return &model.User{
		ID:           usr.ID.Hex(),
		UserName:     usr.UserName,
		Email:        usr.Email,
		UserInfo:     usr.UserInfo,
		UserSessions: sessions,
		CreatedAt:    date,
		Roles:        roles,
	}
}

func ToDb(usr *model.User) *User {
	user := &User{
		UserName:       usr.UserName,
		PasswordHash:   usr.PasswordHash,
		Email:          usr.Email,
		EmailConfirmed: false,
		UserInfo:       usr.UserInfo,
	}
	return user
}
