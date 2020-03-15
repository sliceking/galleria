package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sliceking/galleria/hash"
	"github.com/sliceking/galleria/rand"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID is return when an invalid ID is passed to a method like delete
	ErrInvalidID = errors.New("models: ID was invalid")
	// ErrInvalidPassword is returned when an invalid password is used to auth
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "IamAsuperSecretString"
const hmacSecretKey = "secret-hmac-key"

// User represents a user in our application
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserDB is used to interact with a users db
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

//UserService is a set of methods used to manipulate and work with user model
type UserService interface {
	// Authenticate will verify the provided email and password are correct
	// if they are correct the user corrosponding to that email will be returned
	// otherwise youll receive an error.
	Authenticate(email, password string) (*User, error)
	UserDB
}

// NewUserService is used at the start of the application to open a connection
// to the database
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := NewUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		UserDB: ug,
		hmac:   hmac,
	}

	return &userService{
		UserDB: uv,
	}, nil
}

var _ UserService = &userService{}

// UserService is used to communicate with the db for user data
type userService struct {
	UserDB
}

// Authenticate can be used to authenticate a user with the provided
// email and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash), []byte(password+userPwPepper),
	)
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}

	return nil
}

var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// Create will create a user in the database and fill the ID, CreatedAt,
// UpdatedAt and DeletedAt fields
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user, uv.bcryptPassword, uv.setRememberIfUnset, uv.hmacRemember)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will has a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user, uv.bcryptPassword, uv.hmacRemember)
	if err != nil {
		return err
	}

	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}

	return uv.UserDB.Update(user)
}

// Delete will soft delete a user in the database
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	return uv.UserDB.Delete(id)
}

// ByRemember will has the remember token and then call the ByRemember on the
// UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}

	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

// bcryptPassword will hash a users password with a predefined pepper
// (userPwPepper) and bcrypt if the password field is not empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	// zero the password out for safety
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}

	user.RememberHash = uv.hmac.Hash(user.Remember)

	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token

	return nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

// NewUserGorm accepts a postgres connection string and returns a new instance
// of the userGorm Type
func NewUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &userGorm{
		db: db,
	}, nil
}

// Create will create a user in the database and fill the ID, CreatedAt,
// UpdatedAt and DeletedAt fields
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will update the provided user with all the data
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will soft delete a user in the database
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}

	return ug.db.Delete(&user).Error
}

// ByID will find a user in the database or return an error
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)

	return &user, err
}

// ByEmail will find a user in the database by email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)

	return &user, err
}

// ByRemember looks up a given user by the given remember token and
// returns that user. This method expects the remember token to already be
// hashed. Errors are the same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Close will close the UserService database conenction
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}

	return ug.AutoMigrate()
}

// AutoMigrate will attempt to migrate the db automatically if there is a schema
// change, doesn't always work
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	return nil
}

// first will query using the provided gorm.db and fetch the first record
// and place it into dst, if dst is not a pointer it will not update it
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	return err
}
