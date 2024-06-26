package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Password string
	Id       string
}

func dbConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil
	}
	return db
}

// user routes ============================

func CreateUser(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.String(http.StatusBadRequest, "email or password is empty")
	} else {
		db := dbConnect()

		if db == nil {
			return c.String(http.StatusInternalServerError, "error")
		}

		if db.Where("Email = ?", email).Find(&User{}).RowsAffected != 0 {
			return c.String(http.StatusBadRequest, "email already exists")
		}

		result := db.Create(&User{
			Email:    email,
			Password: password,
			Id:       uuid.NewString(),
		})

		if result.Error != nil {
			return c.String(http.StatusInternalServerError, "error")
		} else {
			return c.String(http.StatusOK, "created")
		}
	}
}

func FindUser(c echo.Context) error {
	if db := dbConnect(); db != nil {
		var users []User
		result := db.Find(&users)
		if result.Error != nil {
			return c.String(http.StatusInternalServerError, "error")
		} else {
			for _, user := range users {
				fmt.Println(user.Email, user.Password, user.Id)
			}
		}

	}
	return c.String(http.StatusOK, "hello")
}

func UpdateUser(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func DeleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

// login ========================================
func Login(c echo.Context) error {
	db := dbConnect()
	if db == nil {
		return c.String(http.StatusInternalServerError, "error")
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.String(http.StatusBadRequest, "email or password is empty")
	}

	var user User
	result := db.Where("Email = ? AND Password = ?", email, password).Find(&user)
	if result.RowsAffected != 0 {
		// make jwt access token ========================
		claims := jwt.MapClaims{
			"user_id": user.Id,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		accessToken, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{"token": accessToken})
	} else {
		return c.JSON(http.StatusOK, map[string]string{"error": "invalid email or password"})
	}
}

