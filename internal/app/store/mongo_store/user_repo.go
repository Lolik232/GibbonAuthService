package mongo_store

import (
	errors "auth-server/internal/app"
	"auth-server/internal/app/model"
	"auth-server/internal/app/store"
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
		UserCommonFields
		UserSessions []UserSession `bson:"user_sessions,omitempty"`
		ClientRoles  []ClientRole  `bson:"user_roles,omitempty"`
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
		UserCommonFields
		UserSessions []UserSessionClient `bson:"user_sessions,omitempty"`
		Roles        []UserClientRole    `bson:"roles,omitempty"`
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
		usersCol *mongo.Collection
	}
)

func (u UserRepo) FindById(ctx context.Context, id string, params *store.UserFields) (*model.User, error) {
	projection := bson.M{}
	projection["username"] = params.UserName
	projection["email"] = params.Email
	projection["created_at"] = params.CreatedAt
	projection["user_info"] = params.UserInfo
	uid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.InvalidArgument.Newf("Invalid userID: %s", id)
	}
	query := bson.M{"_id": uid}
	var usr *UserClient
	if params.UserSessions == false {
		opt := options.FindOne().SetProjection(projection)
		err = u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				return nil, errors.InvalidArgument.Newf("Invalid userID: %s", id)
			case mongo.ErrClientDisconnected:
				return nil, errors.InternalServerError.New("")
			}
		}
	} else {
		proj := bson.D{{"$project", projection}}
		match := bson.D{{"$match", bson.D{{"_id", uid}}}}
		pip := bson.D{{"$lookup", bson.D{{"from", "clients"}, {"localField", "user_sessions.client_id"}, {"foreignField", "_id"}, {"as", "user_sessions"}}}}
		cur, err := u.usersCol.Aggregate(ctx, mongo.Pipeline{match, pip, proj})
		defer cur.Close(ctx)
		if err != nil {
			switch err {
			case mongo.ErrClientDisconnected:
				return nil, errors.InternalServerError.New("")
			}
		}
		if cur != nil && cur.Next(ctx) {
			err = cur.Decode(&usr)
			if err != nil {
				return nil, errors.InternalServerError.New("")
			}
		} else {
			return nil, errors.InvalidArgument.Newf("Invalid userID %s", id)
		}
	}
	return ToUserClient(usr), nil
}
func (u UserRepo) Create(ctx context.Context, user *model.User) (string, error) {
	usr := ToDb(user)
	*usr.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	cur, err := u.usersCol.InsertOne(ctx, usr)
	if err != nil {
		switch err {
		case mongo.ErrNilValue, mongo.ErrNoDocuments, mongo.ErrClientDisconnected:
			return "", errors.InternalServerError.New("")
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

func (u UserRepo) Authenticate(ctx context.Context, login, passwordHash string) (*model.User, error) {
	proj := bson.M{
		"username": 1,
		"email":    1,
	}
	opt := options.FindOne().SetProjection(proj)
	query := bson.M{
		"$or": bson.M{
			"username": login,
			"email":    login,
		},
		"password_hash": passwordHash,
	}
	var usr *UserClient
	err := u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errors.InvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return nil, errors.InternalServerError.New("")
		}
	}
	return ToUserClient(usr), nil
}
func (u UserRepo) AuthenticateByID(ctx context.Context, userID, passwordHash string) (*model.User, error) {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.InvalidArgument.Newf("Invalid userID %s", userID)
	}
	proj := bson.M{
		"username": 1,
		"email":    1,
	}
	opt := options.FindOne().SetProjection(proj)
	query := bson.M{
		"_id":           ID,
		"password_hash": passwordHash,
	}
	var usr *UserClient
	err = u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errors.InvalidPassword.New("")
		case mongo.ErrClientDisconnected:
			return nil, errors.InternalServerError.New("")
		}
	}
	return ToUserClient(usr), nil
}
func (u UserRepo) AuthenticateByEmail(ctx context.Context, email, passwordHash string) (*model.User, error) {
	proj := bson.M{
		"username": 1,
		"email":    1,
	}
	opt := options.FindOne().SetProjection(proj)
	query := bson.M{
		"email":         email,
		"password_hash": passwordHash,
	}
	var usr *UserClient
	err := u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errors.InvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return nil, errors.InternalServerError.New("")
		}
	}
	return ToUserClient(usr), nil
}
func (u UserRepo) AuthenticateByName(ctx context.Context, username, passwordHash string) (*model.User, error) {
	proj := bson.M{
		"username": 1,
		"email":    1,
	}
	opt := options.FindOne().SetProjection(proj)
	query := bson.M{
		"username":      username,
		"password_hash": passwordHash,
	}
	var usr *UserClient
	err := u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errors.InvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return nil, errors.InternalServerError.New("")
		}
	}
	return ToUserClient(usr), nil
}

func (u UserRepo) CheckPassword(ctx context.Context, login, passwordHash string) error {
	proj := bson.M{
		"username": 1,
		"email":    1,
	}
	opt := options.FindOne().SetProjection(proj)
	query := bson.M{
		"$or": bson.M{
			"username": login,
			"email":    login,
		},
		"password_hash": passwordHash,
	}
	var usr *UserClient
	err := u.usersCol.FindOne(ctx, query, opt).Decode(&usr)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.InvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return errors.InternalServerError.New("")
		}
	}
	return nil
}
func (u UserRepo) CheckPasswordById(ctx context.Context, userID, passwordHash string) error {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.InvalidArgument.Newf("Invalid userID %s", userID)
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
			return errors.InvalidPassword.New("")
		case mongo.ErrClientDisconnected:
			return errors.InternalServerError.New("")
		}
	}
	if user.ID.Hex() != userID {
		return errors.InvalidArgument.New("")
	}
	return nil
}
func (u UserRepo) CheckPasswordByName(ctx context.Context, username, passwordHash string) error {

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
			return errors.InvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return errors.InternalServerError.New("")
		}
	}
	if usr.UserName != username {
		return errors.InvalidPasswordOrUsername.New("")
	}
	return nil
}
func (u UserRepo) CheckPasswordByEmail(ctx context.Context, email, passwordHash string) error {
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
			return errors.InvalidPasswordOrUsername.New("")
		case mongo.ErrClientDisconnected:
			return errors.InternalServerError.New("")
		}
	}
	if usr.Email != email {
		return errors.InvalidPasswordOrUsername.New("")
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
		UserCommonFields: UserCommonFields{
			UserName:       usr.UserName,
			PasswordHash:   usr.PasswordHash,
			Email:          usr.Email,
			EmailConfirmed: false,
			UserInfo:       usr.UserInfo,
		},
	}
	return user
}
