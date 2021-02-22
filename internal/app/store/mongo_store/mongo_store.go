package mongo_store

import (
	st "auth-server/internal/app/store"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UsersCollection   = "users"
	ClientsCollection = "client"
)

//Store is a mongoDB database storage
type Store struct {
	db               *mongo.Database
	userRepository   *UserRepo
	clientRepository *ClientRepo
}

func NewStore(db *mongo.Database) *Store {
	return &Store{db: db}
}

//User returns the "Users" repository
func (s *Store) User() st.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepo{
		store:    s,
		usersCol: s.db.Collection(UsersCollection),
	}
	return s.userRepository
}

//Client returns the "Clients" repository
func (s *Store) Client() st.ClientRepository {
	if s.clientRepository != nil {
		return s.clientRepository
	}
	s.clientRepository = &ClientRepo{
		store:      s,
		clientsCol: s.db.Collection(ClientsCollection),
	}

	return s.clientRepository
}
