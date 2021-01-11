package mongo_store

import (
	"auth-server/internal/app/store"
	"go.mongodb.org/mongo-driver/mongo"
)

//Store is mongoDB database storage
type Store struct {
	db               *mongo.Database
	userRepository   *UserRepo
	clientRepository *ClientRepo
}

//User returns the "Users" repository
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepo{
		usersCol: s.db.Collection("col"),
	}
	return s.userRepository
}

//Client returns the "Clients" repository
func (s *Store) Client() store.ClientRepository {
	if s.clientRepository != nil {
		return s.clientRepository
	}
	s.clientRepository = &ClientRepo{
		clientsCol: s.db.Collection("col"),
	}
	return s.clientRepository
}
