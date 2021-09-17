package repository

import as "github.com/aerospike/aerospike-client-go/v5"

type User struct {
	ID        string
	Timestamp int64
}

func NewUser(ID string, timestamp int64) *User {
	return &User{ID: ID, Timestamp: timestamp}
}

const namespace = "test"

type UserRepository struct {
	client *as.Client
}

func NewUserRepository() (*UserRepository, error) {
	client, err := as.NewClient("127.0.0.1", 3000)
	if err != nil {
		return nil, err
	}
	return &UserRepository{client: client}, nil
}

func (ur *UserRepository) SaveUser(user *User) error {
	key, err := as.NewKey(namespace, "", user.ID)
	if err != nil {
		return err
	}
	return ur.client.PutObject(nil, key, user)
}

func (ur *UserRepository) GetUser(id string) (*User, error) {
	key, err := as.NewKey(namespace, "", id)
	if err != nil {
		return nil, err
	}
	var user User
	err = ur.client.GetObject(nil, key, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur UserRepository) InvokeUDFInScope(begin int64, end int64) error {
	stmt := as.NewStatement(namespace, "")
	err := stmt.SetFilter(as.NewRangeFilter("Timestamp", begin, end))
	if err != nil {
		return err
	}
	return ur.executeUDF(stmt)
}

func (ur UserRepository) InvokeUDF() error {
	stmt := as.NewStatement(namespace, "")
	return ur.executeUDF(stmt)
}

func (ur *UserRepository) executeUDF(stmt *as.Statement) error {
	task, err := ur.client.ExecuteUDF(nil, stmt, "utils", "remove")
	if err != nil {
		return err
	}
	return <-task.OnComplete()
}
