package router

import (
	"github.com/EmilyOng/tusk-manager/backend/constants"
	"github.com/EmilyOng/tusk-manager/backend/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() (router *gin.Engine) {
	router = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{constants.FrontendLocalHostUrl, constants.FrontendProductionUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Use(handlers.SetAuthUser)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/signup", handlers.SignUp)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/", handlers.IsAuthenticated)
		}
		guard := api.Group("/", handlers.AuthGuard)
		{
			states := guard.Group("/states")
			{
				states.POST("/", handlers.CreateState)
				states.PUT("/", handlers.UpdateState)
				states.DELETE("/:state_id", handlers.DeleteState)
			}
			boards := guard.Group("/boards")
			{
				boards.GET("/", handlers.GetUserBoards)
				boards.PUT("/", handlers.UpdateBoard)
				boards.DELETE("/:board_id", handlers.DeleteBoard)
				boards.GET("/:board_id", handlers.GetBoard)
				boards.POST("/", handlers.CreateBoard)
				boards.GET("/:board_id/tasks", handlers.GetBoardTasks)
				boards.GET("/:board_id/tags", handlers.GetBoardTags)
				boards.GET("/:board_id/states", handlers.GetBoardStates)
				boards.GET("/:board_id/members", handlers.GetBoardMemberProfiles)
			}
			tasks := guard.Group("/tasks")
			{
				tasks.POST("/", handlers.CreateTask)
				tasks.PUT("/", handlers.UpdateTask)
				tasks.DELETE("/:task_id", handlers.DeleteTask)
			}
			tags := guard.Group("/tags")
			{
				tags.POST("/", handlers.CreateTag)
				tags.DELETE("/:tag_id", handlers.DeleteTag)
				tags.PUT("/", handlers.UpdateTag)
			}
			members := guard.Group("/members")
			{
				members.POST("/", handlers.CreateMember)
				members.PUT("/", handlers.UpdateMember)
				members.DELETE("/:member_id", handlers.DeleteMember)
			}
		}
	}
	return
}
