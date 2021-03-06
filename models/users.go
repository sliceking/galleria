package models

import (
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sliceking/galleria/hash"
	"github.com/sliceking/galleria/rand"
	"golang.org/x/crypto/bcrypt"
)

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
func NewUserService(db *gorm.DB, pepper, hmacKey string) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac, pepper)

	return &userService{
		UserDB: uv,
		pepper: pepper,
	}
}

var _ UserService = &userService{}

// UserService is used to communicate with the db for user data
type userService struct {
	UserDB
	pepper string
}

// Authenticate can be used to authenticate a user with the provided
// email and password.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash), []byte(password+us.pepper),
	)
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrPasswordIncorrect
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
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
	pepper     string
}

func newUserValidator(udb UserDB, hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
		pepper:     pepper,
	}
}

// Create will create a user in the database and fill the ID, CreatedAt,
// UpdatedAt and DeletedAt fields
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable,
	)

	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will has a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable,
	)

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
	var user User
	user.ID = id

	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

// ByEmail will normalize the email address before calling ByEmail on the
// UserDB field.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}

	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
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

	pwBytes := []byte(user.Password + uv.pepper)
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

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}

	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}

	if n < 32 {
		return ErrRememberTooShort
	}

	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}

	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}

		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}

	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}

	return nil
}

func (uv *userValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		// Email address is not taken.
		return nil
	}
	if err != nil {
		return err
	}

	// We found a user with this email address  If the found user has the same
	// ID as this user, it is an update and this is the same user.
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}

	return nil
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

// first will query using the provided gorm.db and fetch the first record
// and place it into dst, if dst is not a pointer it will not update it
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	return err
}
