package mongo_store

import (
	"auth-server/internal/app/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO:Implement client repository.
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
	store      *Store
	clientsCol *mongo.Collection
}

func (c ClientRepo) FindById(ctx context.Context, id string) (*model.Client, error) {
	panic("implement me")
}

func (c ClientRepo) FindRefToken(ctx context.Context, clientID, sessionID, refToken string) (*model.ClientRefToken, error) {
	panic("implement me")
}

func (c ClientRepo) CheckRefToken(ctx context.Context, clientID, sessionID, refToken string) error {
	panic("implement me")
}
