package main

import (
	"tqif-golang/config"
	"tqif-golang/controller"
	"tqif-golang/middleware"
	"tqif-golang/repository"
	"tqif-golang/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db              *gorm.DB                   = config.SetupDatabaseConnection()
	userRepository  repository.UserRepository  = repository.NewUserRepository(db)
	diaryRepository repository.DiaryRepository = repository.NewDiaryRepository(db)
	jwtService      service.JWTService         = service.NewJWTService()
	userService     service.UserService        = service.NewUserService(userRepository)
	diaryService    service.DiaryService       = service.NewDiaryService(diaryRepository)
	authService     service.AuthService        = service.NewAuthService(userRepository)
	authController  controller.AuthController  = controller.NewAuthController(authService, jwtService)
	userController  controller.UserController  = controller.NewUserController(userService, jwtService)
	diaryController controller.DiaryController = controller.NewDiaryController(diaryService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	diaryRoutes := r.Group("api/diaries", middleware.AuthorizeJWT(jwtService))
	{
		diaryRoutes.GET("/", diaryController.All)
		diaryRoutes.GET("/:id", diaryController.FindByID)
		diaryRoutes.POST("/", diaryController.Insert)
		diaryRoutes.PUT("/:id", diaryController.Update)
		diaryRoutes.DELETE("/:id", diaryController.Delete)
	}

	r.Run()
}
