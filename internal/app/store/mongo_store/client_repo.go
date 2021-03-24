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
)

//TODO:Testing client repository
//Client represent the "Clients" collection
type Client struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ClientName string             `bson:"client_name,omitempty"`
	RefTokens  []RefToken         `bson:"ref_tokens,omitempty"`
}

//RefToken represent  attached "RefTokens" document in "Client"
type RefToken struct {
	SessionID primitive.ObjectID `bson:"session_id,omitempty"`
	RefToken  string             `bson:"ref_token,omitempty"`
	ExpIn     primitive.DateTime `bson:"exp_in,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty"`
}

type ClientRepo struct {
	store      store.Store
	clientsCol *mongo.Collection
}

func (c *ClientRepo) fetch(ctx context.Context, query, proj bson.M) (*Client, error) {
	var client *Client
	options := options.FindOne()
	if proj != nil {
		options.SetProjection(proj)
	}
	err := c.clientsCol.FindOne(ctx, query, options).Decode(client)
	return client, err
}

func (c *ClientRepo) FindById(ctx context.Context, id string) (*model.Client, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrInvalidArgument.Newf("Invalid clientId %s", id)
	}

	query := bson.M{
		"_id": oid,
	}

	var client *Client
	client, err = c.fetch(ctx, query, nil)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errors.ErrInvalidArgument.Newf("Invalid clientId %s", id)
		default:
			return nil, errors.NoType.Wrap(err, "")
		}
	}
	return ToClient(client), nil
}

func (c *ClientRepo) FindRefToken(ctx context.Context, clientID, sessionID, refToken string) (*model.ClientRefToken, error) {
	clientObjID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return nil, errors.ErrInvalidArgument.Newf("Invalid client ID %s", clientID)
	}
	sessionObjID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, errors.ErrInvalidArgument.Newf("Invalid session ID %s", sessionID)
	}
	query := bson.M{
		"_id":                   clientObjID,
		"ref_tokens.session_id": sessionObjID,
		"ref_tokens.ref_token":  refToken,
	}
	proj := bson.M{
		"ref_tokens.$": 1,
	}
	client, err := c.fetch(ctx, query, proj)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errors.ErrInvalidArgument.New("Invalid refToken")
		default:
			return nil, errors.NoType.Wrap(err, "")
		}
	}

	rToken := (*client).RefTokens[0]

	return ToClientRefToken(&rToken), nil
}

func (c ClientRepo) CheckRefToken(ctx context.Context, clientID, sessionID, refToken string) (bool, error) {
	clientObjID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		return false, errors.ErrInvalidArgument.Newf("Invalid client ID %s", clientID)
	}
	sessionObjID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return false, errors.ErrInvalidArgument.Newf("Invalid session ID %s", sessionID)
	}
	query := bson.M{
		"_id":                   clientObjID,
		"ref_tokens.session_id": sessionObjID,
		"ref_tokens.ref_token":  refToken,
	}
	proj := bson.M{
		"_id": 1,
	}

	options := options.FindOne()
	if proj != nil {
		options.SetProjection(proj)
	}
	err = c.clientsCol.FindOne(ctx, query, options).Err()
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return false, errors.ErrInvalidArgument.New("Not found refresh token.")
		default:
			return false, errors.NoType.Wrap(err, "")
		}
	}
	return true, nil
}

func (c *ClientRepo) CreateRefToken(ctx context.Context, clientID string, refToken *model.ClientRefToken) error {
	clientObjectID, err := primitive.ObjectIDFromHex(clientID)
	sessionObjectID, err := primitive.ObjectIDFromHex(refToken.SessionID)

	if err != nil {
		return nil
	}
	query := bson.M{
		"_id": clientObjectID,
	}
	pullUpdate := bson.M{
		"$pull": bson.M{
			"ref_tokens": bson.M{
				"session_id": sessionObjectID,
			},
		},
	}

	dbRefToken := ToDbRefToken(refToken)
	_, err = c.clientsCol.UpdateOne(ctx, query, pullUpdate)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return errors.ErrInvalidArgument.New("Invalid clientID")
		default:
			return errors.NoType.New("")
		}
	}
	update := bson.M{
		"$push": bson.M{
			"ref_tokens": dbRefToken,
		},
	}
	res, err := c.clientsCol.UpdateOne(ctx, query, update)

	if err != nil {
		return errors.NoType.Wrap(err, "")
	}
	if res.MatchedCount == 0 {
		return errors.ErrInvalidArgument.Newf("Invalid client ID %s", clientID)
	}

	if res.ModifiedCount == 0 {
		return errors.NoType.New("")
	}

	return nil
}

func (c *ClientRepo) DeleteRefToken(ctx context.Context, clientID, sessionID, refToken string) error {
	clientObjectID, err := primitive.ObjectIDFromHex(clientID)
	//sessionObjectID, err := primitive.ObjectIDFromHex(sessionID)

	if err != nil {
		return nil
	}
	query := bson.M{
		"_id": clientObjectID,
	}
	pullUpdate := bson.M{
		"$pull": bson.M{
			"ref_tokens": bson.M{
				"refToken": refToken,
			},
		},
	}
	res, err := c.clientsCol.UpdateOne(ctx, query, pullUpdate)
	if err != nil {
		return errors.NoType.Wrap(err, "")
	}

	if res.MatchedCount == 0 {
		return errors.ErrInvalidArgument.Newf("Invalid client id %s", clientID)
	}
	if res.ModifiedCount == 0 {
		return errors.ErrInvalidArgument.New("Invalid ref token")
	}
	return nil
}

func ToClientRefToken(dbRefToken *RefToken) *model.ClientRefToken {
	sessionID := dbRefToken.SessionID.Hex()
	expIn := dbRefToken.ExpIn.Time()
	createdAt := dbRefToken.CreatedAt.Time()
	return &model.ClientRefToken{
		SessionID: sessionID,
		RefToken:  dbRefToken.RefToken,
		ExpIn:     expIn,
		CreatedAt: createdAt,
	}
}
func ToDbRefToken(clientRefToken *model.ClientRefToken) *RefToken {
	sessionID, _ := primitive.ObjectIDFromHex(clientRefToken.SessionID)
	expIn := primitive.NewDateTimeFromTime(clientRefToken.ExpIn)
	createdAt := primitive.NewDateTimeFromTime(clientRefToken.CreatedAt)

	return &RefToken{
		SessionID: sessionID,
		RefToken:  clientRefToken.RefToken,
		ExpIn:     expIn,
		CreatedAt: createdAt,
	}
}

func ToClient(dbclient *Client) *model.Client {
	tokens := make([]model.ClientRefToken, 0)
	for _, v := range dbclient.RefTokens {
		refToken := model.ClientRefToken{
			SessionID: v.SessionID.Hex(),
			RefToken:  v.RefToken,
			ExpIn:     v.ExpIn.Time(),
			CreatedAt: v.CreatedAt.Time(),
		}
		tokens = append(tokens, refToken)
	}

	return &model.Client{
		ID:               dbclient.ID.Hex(),
		ClientName:       dbclient.ClientName,
		ClientsRefTokens: tokens,
	}
}
