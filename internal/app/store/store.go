//Package store represent interfaces for repositories and repositories storage
package store

//Store is repositories storage
//Methods returns the repositories
type Store interface {
	User() UserRepository
	Client() ClientRepository
}
