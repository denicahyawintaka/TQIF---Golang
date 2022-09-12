package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"tqif-golang/dto"
	"tqif-golang/entity"
	"tqif-golang/helper"
	"tqif-golang/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// DiaryController is a ...
type DiaryController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type diaryController struct {
	diaryService service.DiaryService
	jwtService   service.JWTService
}

// NewDiaryController create a new instances of DiaryController
func NewDiaryController(diaryServ service.DiaryService, jwtServ service.JWTService) DiaryController {
	return &diaryController{
		diaryService: diaryServ,
		jwtService:   jwtServ,
	}
}

func (c *diaryController) All(context *gin.Context) {
	var diaries []entity.Diary = c.diaryService.All()
	res := helper.BuildResponse(true, "OK", diaries)
	context.JSON(http.StatusOK, res)
}

func (c *diaryController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var diary entity.Diary = c.diaryService.FindByID(id)
	if (diary == entity.Diary{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK", diary)
		context.JSON(http.StatusOK, res)
	}
}

func (c *diaryController) Insert(context *gin.Context) {
	var diaryCreateDTO dto.DiaryCreateDTO
	errDTO := context.ShouldBindJSON(&diaryCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			diaryCreateDTO.UserID = convertedUserID
		}
		result := c.diaryService.Insert(diaryCreateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *diaryController) Update(context *gin.Context) {
	var diaryUpdateDTO dto.DiaryUpdateDTO
	errDTO := context.ShouldBindJSON(&diaryUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.diaryService.IsAllowedToEdit(userID, diaryUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			diaryUpdateDTO.UserID = id
		}
		result := c.diaryService.Update(diaryUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *diaryController) Delete(context *gin.Context) {
	var diary entity.Diary
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	diary.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.diaryService.IsAllowedToEdit(userID, diary.ID) {
		c.diaryService.Delete(diary)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *diaryController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
