package pkg

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type user struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Biography string    `json:"biography"`
}

type Application struct {
	mu   sync.Mutex
	Data map[uuid.UUID]user
}

func NewDB() *Application {
	return &Application{
		Data: make(map[uuid.UUID]user),
	}
}

func (a *Application) FindAll() []user {
	a.mu.Lock()
	defer a.mu.Unlock()
	users := []user{}

	for i, u := range a.Data {
		users = append(users, user{
			ID:        i,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Biography: u.Biography,
		})
	}

	return users
}

func (a *Application) Insert(firstName, lastName, bio string) (user, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	id := uuid.New()

	a.Data[id] = user{
		FirstName: firstName,
		LastName:  lastName,
		Biography: bio,
		ID:        id,
	}

	u, ok := a.Data[id]
	if !ok {
		return user{}, errors.New("There was an error while saving the user to the database")
	}

	return u, nil
}

func (a *Application) FindById(id string) *user {
	a.mu.Lock()
	defer a.mu.Unlock()
	err := uuid.Validate(id)
	if err != nil {
		return nil
	}

	formatId := uuid.MustParse(id)

	user, ok := a.Data[formatId]
	if !ok {
		return nil
	}

	return &user
}

func (a *Application) Update(id, firstName, lastName, bio string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	u := a.FindById(id)
	if u == nil {
		return errors.New("The user with the specified ID does not exist")
	}

	if firstName == "" || bio == "" {
		return errors.New("Please provide name and bio for the user")
	}

	if lastName == "" {
		lastName = u.LastName
	}

	a.Data[u.ID] = user{
		FirstName: firstName,
		LastName:  lastName,
		Biography: bio,
		ID:        u.ID,
	}

	return nil
}

func (a *Application) Delete(id string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	user := a.FindById(id)
	if user == nil {
		return errors.New("The user with the specified ID does not exist")
	}
	delete(a.Data, user.ID)
	return nil
}
