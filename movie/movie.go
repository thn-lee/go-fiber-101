package movie

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thn-lee/go-fiber-101/database"
	"github.com/thn-lee/go-fiber-101/model"
	"gorm.io/gorm"
)

type Movie struct {
	Id        uint
	Title     string
	Genre     string
	Rating    json.Number
	Studio    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func GetMovies(ctx *fiber.Ctx) (err error) {

	var resForm model.ResponseForm
	var errResArr []model.ResponseError
	var movies []Movie

	db := database.DBConn
	err = db.Find(&movies).Error

	if err != nil {
		log.Printf("Func GetMovies err: %+v", err)

		ctx.Status(http.StatusInternalServerError)

		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  "Func GetMovies",
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}

		errResArr = append(errResArr, errResObj)
		return ctx.JSON(map[string]interface{}{"erroes": errResArr})
	}

	resForm = model.ResponseForm{
		Success: bool(err == nil),
		Result:  fiber.Map{"movies": movies},
	}
	return ctx.JSON(resForm)
}

func GetMovieById(ctx *fiber.Ctx) (err error) {
	var resForm model.ResponseForm
	var errResArr []model.ResponseError

	id := ctx.Params("id")

	intId, err := strconv.Atoi(id)

	if err != nil {
		log.Printf("GetMovieById err: %+v", err)
		ctx.Status(http.StatusBadRequest)
		errResObj := model.ResponseError{
			Code:    http.StatusBadRequest,
			Source:  "GetMovieById",
			Title:   http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	var movie Movie
	db := database.DBConn

	err = db.Where(Movie{Id: uint(intId)}).Find(&movie).Error
	if err != nil {
		log.Printf("GetMovie err: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  http.StatusInternalServerError,
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}
	resForm = model.ResponseForm{
		Success: bool(err == nil),
		Result:  fiber.Map{"Movie": movie},
	}
	return ctx.JSON(resForm)
}

func CreateMovie(ctx *fiber.Ctx) (err error) {
	var resForm model.ResponseForm
	var errResArr []model.ResponseError
	movie := new(Movie)

	if err := ctx.BodyParser(movie); err != nil {
		log.Printf("CreateMovie err: %+v", err)
		ctx.Status(http.StatusUnprocessableEntity)
		errResObj := model.ResponseError{
			Code:    http.StatusUnprocessableEntity,
			Source:  "CreateMovie",
			Title:   http.StatusText(http.StatusUnprocessableEntity),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	db := database.DBConn

	dbTx := db.Begin()

	defer dbTx.Rollback()

	if err := dbTx.Create(&movie).Error; err != nil {
		ctx.Status(http.StatusInternalServerError)
		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  "CreateMovies",
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"erros": errResArr})
	}

	err = dbTx.Commit().Error
	if err != nil {
		log.Printf("CreateMovie err: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  "Func GetMovies",
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	resForm = model.ResponseForm{
		Success: bool(err == nil),
		Result:  fiber.Map{"movie": movie},
	}
	return ctx.JSON(resForm)
}

func DeleteMovie(ctx *fiber.Ctx) (err error) {
	var resForm model.ResponseForm
	var errResArr []model.ResponseError

	id := ctx.Params("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("DeleteMovie err: %+v", err)
		ctx.Status(http.StatusBadRequest)
		errResObj := model.ResponseError{
			Code:    http.StatusBadRequest,
			Source:  "DeleteMovie",
			Title:   http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"erros": errResArr})
	}

	var movie Movie
	db := database.DBConn

	dbTx := db.Begin()
	defer dbTx.Rollback()

	err = dbTx.Where(Movie{Id: uint(intId)}).First(&movie).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(404).SendString("Movie not found")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("DeleteMovie err: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  "DeleteMovie",
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	err = dbTx.Where(Movie{Id: uint(intId)}).Delete(&movie).Error
	if err != nil {
		log.Printf("DeleteMovie err: %+v", err)
		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  "DeleteMovie",
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	err = dbTx.Commit().Error
	if err != nil {
		log.Printf("DeleteMovie err: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  "DeleteMovie",
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	resForm = model.ResponseForm{
		Success: bool(err == nil),
		Result:  fiber.Map{"movie": movie},
	}
	return ctx.JSON(resForm)
}

func UpdateMovie(ctx *fiber.Ctx) (err error) {
	var resForm model.ResponseForm
	var errResArr []model.ResponseError
	movie := new(Movie)

	// get path(id)
	id := ctx.Params("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("UpdateMovie err: %+v", err)
		ctx.Status(http.StatusBadRequest)
		errResObj := model.ResponseError{
			Code:    http.StatusBadRequest,
			Source:  "UpdateMovie",
			Title:   http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	db := database.DBConn
	dbTx := db.Begin()
	defer dbTx.Rollback()

	err = dbTx.Where(Movie{Id: uint(intId)}).First(&movie).Error
	// check if record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(404).SendString("Movie not found")
	}

	// parse body into movie model
	if err := ctx.BodyParser(movie); err != nil {
		log.Printf("UpdateMovie err: %+v", err)
		ctx.Status(http.StatusUnprocessableEntity)
		errResObj := model.ResponseError{
			Code:    http.StatusUnprocessableEntity,
			Source:  "UpdateMovie",
			Title:   http.StatusText(http.StatusUnprocessableEntity),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	if err := dbTx.Where(Movie{Id: uint(intId)}).Updates(&movie).Error; err != nil {
		log.Printf("UpdateBook err: %+v", err)
		ctx.Status(http.StatusInternalServerError)
		errResObj := model.ResponseError{
			Code:    http.StatusInternalServerError,
			Source:  "UpdateBook",
			Title:   http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		}
		errResArr = append(errResArr, errResObj)
		return ctx.JSON(fiber.Map{"errors": errResArr})
	}

	resForm = model.ResponseForm{
		Success: bool(err == nil),
		Result:  fiber.Map{"movie": movie},
	}

	return ctx.JSON(resForm)

}
