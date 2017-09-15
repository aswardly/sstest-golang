//Package user provides the business domain model definitions of User
package user

import (
	"fmt"
	"sync"

	"github.com/go-errors/errors"
	"golang.org/x/crypto/bcrypt"
)

//User is business domain model definition of user
type User struct {
	id       string
	password []byte //password hash
	name     string
	address  string
	status   string
	mu       sync.Mutex
}

//StatusActive is const for 'active' user status
const StatusActive string = "A"

//StatusInactive is const for 'inactive' user status
const StatusInactive string = "I"

//StatusSuspended is const for 'suspended' user status
const StatusSuspended string = "S"

//statusSlice is a map of known status code and its label pairs
var statusSlice = map[string]string{
	StatusActive:    "Active",
	StatusInactive:  "Inactive",
	StatusSuspended: "Suspended",
}

//New creates a new user model struct, initializes it's properties and returns a reference to it
func New(id, name, address string) (*User, *errors.Error) {
	//default password for new user is "changeme"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("Failed creating user: %v", err.Error()), 0)
	}
	return &User{
		id,
		hashedPassword,
		name,
		address,
		StatusInactive,
		*new(sync.Mutex),
	}, nil
}

//ID is a getter function for returning a user's id
func (u *User) ID() string {
	return u.id
}

//Password is getter function for returning a user's password hash
func (u *User) Password() []byte {
	return u.password
}

//Name is a getter function for returning a user's name
func (u *User) Name() string {
	return u.name
}

//Address is a getter function for returning a user's stock
func (u *User) Address() string {
	return u.address
}

//Status is a getter function for returning a user's status
func (u *User) Status() string {
	return u.status
}

//SetID is a setter function for setting a user's id
func (u *User) SetID(id string) *User {
	u.id = id
	return u
}

//SetPassword is a setter function for setting user's password from a plaintext string value
func (u *User) SetPassword(password string) (*User, *errors.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(fmt.Errorf("Failed setting password: %v", err.Error()), 0)
	}
	u.password = hashedPassword
	return u, nil
}

//SetPasswordHash is a setter function for setting user's password from a password hash
func (u *User) SetPasswordHash(hash []byte) (*User, *errors.Error) {
	//check whether given hash is a valid bcrypt hash (possible by trying to check bcrypt hash cost)
	if _, err := bcrypt.Cost(hash); err != nil {
		return nil, errors.Wrap(err, 1)
	}

	u.password = hash
	return u, nil
}

//SetName is a setter function for setting a user's name
func (u *User) SetName(name string) *User {
	u.name = name
	return u
}

//SetAddress is a setter function for setting a user's address
func (u *User) SetAddress(address string) *User {
	u.address = address
	return u
}

//SetStatus is a setter function for setting a user's status
func (u *User) SetStatus(status string) (*User, *errors.Error) {
	if _, ok := statusSlice[status]; false == ok {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set unknown status type: %v", status), 0)
	}
	u.status = status
	return u, nil
}

//Business logic methods

//Activate is a function for activating user
func (u *User) Activate() {
	u.SetStatus(StatusActive)
}

//Deactivate is a function for deactivating user
func (u *User) Deactivate() {
	u.SetStatus(StatusInactive)
}

//Suspend is a function for suspending user
func (u *User) Suspend() {
	u.SetStatus(StatusSuspended)
}

//ValidatePassword is a function for validating whether a given string is the user's password
//returns true if given password matches, else return false and an error
func (u *User) ValidatePassword(password string) (bool, *errors.Error) {
	if err := bcrypt.CompareHashAndPassword(u.password, []byte(password)); err != nil {
		return false, errors.Wrap(err, 1)
	}
	return true, nil
}

//CanOrder is a function for inquiring whether an user can place order or not
//returns true if user can place order, else return false and an error
func (u *User) CanOrder() (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if StatusActive != u.status {
		return false, errors.Wrap(fmt.Errorf("User status is not active"), 0)
	}
	return true, nil
}
