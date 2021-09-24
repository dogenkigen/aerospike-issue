package repository

import (
	as "github.com/aerospike/aerospike-client-go/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestUserRepository_SaveUser(t *testing.T) {
	repository, err := NewUserRepository()
	require.NoError(t, err)
	err = cleanState(repository.client)
	require.NoError(t, err)

	user := NewUser(uuid.NewString(), time.Now().UTC().Unix())

	err = repository.SaveUser(user)
	require.NoError(t, err)

	getUser, err := repository.GetUser(user.ID)
	require.NoError(t, err)

	assert.Equal(t, user, getUser)
}

func TestUserRepository_InvokeUDFInScope(t *testing.T) {
	repository, err := NewUserRepository()
	require.NoError(t, err)
	err = cleanState(repository.client)
	require.NoError(t, err)

	ts := time.Now().UTC().Unix()
	user := NewUser(uuid.NewString(), ts)

	err = repository.SaveUser(user)
	require.NoError(t, err)

	err = repository.InvokeUDFInScope(ts-100, ts+100)
	require.NoError(t, err)

	getUser, err := repository.GetUser(user.ID)
	assert.Nil(t, getUser)
	assert.True(t, strings.Contains(err.Error(), "Key not found"))
}

func TestUserRepository_InvokeUDF(t *testing.T) {
	repository, err := NewUserRepository()
	require.NoError(t, err)
	err = cleanState(repository.client)
	require.NoError(t, err)

	ts := time.Now().UTC().Unix()
	user := NewUser(uuid.NewString(), ts)

	err = repository.SaveUser(user)
	require.NoError(t, err)

	err = repository.InvokeUDF()
	require.NoError(t, err)

	getUser, err := repository.GetUser(user.ID)
	assert.Nil(t, getUser)
	assert.True(t, strings.Contains(err.Error(), "Key not found"))
}

func BenchmarkUserRepository_InvokeUDFInScope(b *testing.B) {
	repository, err := NewUserRepository()
	if err != nil {
		b.Fatal(err)
	}
	err = cleanState(repository.client)
	if err != nil {
		b.Fatal(err)
	}
	ts := time.Now().UTC().Unix()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = repository.SaveUser(NewUser(uuid.NewString(), ts))
		if err != nil {
			b.Fatal(err)
		}
		err = repository.InvokeUDFInScope(ts-100, ts+100)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserRepository_InvokeUDF(b *testing.B) {
	repository, err := NewUserRepository()
	if err != nil {
		b.Fatal(err)
	}
	err = cleanState(repository.client)
	if err != nil {
		b.Fatal(err)
	}
	ts := time.Now().UTC().Unix()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = repository.SaveUser(NewUser(uuid.NewString(), ts))
		if err != nil {
			b.Fatal(err)
		}
		err = repository.InvokeUDF()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func cleanState(client *as.Client) error {
	utcNow := time.Now().UTC()
	return client.Truncate(nil, namespace, "", &utcNow)
}
