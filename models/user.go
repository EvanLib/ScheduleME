package models

import (
	"fmt"

	"github.com/EvanLib/hash"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var userPwPerpper = "SOMETHING DUMB AND SCRETE"
var hmac = hash.NewHMAC("something-secrot_:D")

type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"not null;unique_index"`
	Password       string `gorm:"-"`
	HashedPassword string `gorm:"not null"`
	Remember       string `gorm:"-"`
	RememberHash   string `gorm:"not null;unique_index"`
}

type UserService interface {
	ByID(id uint) *User
	ByEmail(email string) *User
	ByRemember(token string) *User
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
	Authenticate(email, password string) *User
}

type UserGorm struct {
	*gorm.DB
}

func (ug *UserGorm) DestructiveReset() {
	ug.DropTable(&User{})
	ug.AutoMigrate(&User{})
	ug.AutoMigrate(&Schedule{})
	ug.Model(&Schedule{}).Related(&User{})
}
func NewUserGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{db}
}

func (ug *UserGorm) ByID(id uint) *User {
	return ug.byQuery(ug.DB.Where("id = ?", id))

}

func (ug *UserGorm) ByEmail(email string) *User {
	return ug.byQuery(ug.DB.Where("email = ?", email))
}

func (ug *UserGorm) ByRemember(token string) *User {
	return ug.byQuery(ug.DB.Where("remember_hash = ?", hmac.String(token)))
}

func (ug *UserGorm) byQuery(query *gorm.DB) *User {
	ret := &User{}
	err := query.First(ret).Error
	switch err {
	case nil:
		return ret
	case gorm.ErrRecordNotFound:
		fmt.Println(err)
		return nil
	default:
		panic(err)
	}
}

func (ug *UserGorm) Create(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password+userPwPerpper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = ""
	if user.RememberHash != "" {
		user.RememberHash = hmac.String(user.Remember)
	}
	return ug.DB.Create(user).Error
}
func (ug *UserGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = hmac.String(user.Remember)
	}
	return ug.DB.Save(user).Error
}
func (ug *UserGorm) Delete(id uint) error {
	user := &User{Model: gorm.Model{ID: id}}
	return ug.DB.Delete(user).Error
}

func (ug *UserGorm) Authenticate(email, password string) *User {
	foundUser := ug.ByEmail(email)
	if foundUser == nil {
		//No User found
		return nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(password+userPwPerpper))
	if err != nil {
		// Invalid password
		return nil
	}
	return foundUser
}
