package router

import (
	"FeedCraft/internal/admin"
	"FeedCraft/internal/controller"
	"FeedCraft/internal/craft"
	"FeedCraft/internal/middleware"
	"FeedCraft/internal/recipe"
	"FeedCraft/internal/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"log"
	"net/http"
)

func RegisterRouters(router *gin.Engine) {
	envClient := util.GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
	}
	siteBaseUrl := envClient.GetString("SITE_BASE_URL")
	router.LoadHTMLFiles("web/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"SiteBaseUrl": siteBaseUrl,
		})
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	corsMiddleware := cors.New(corsConfig)
	//corsMiddleware := cors.Default()
	router.Use(corsMiddleware)

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/login", admin.LoginAuth)
	}

	craftRouters := router.Group("/craft")
	{
		craftRouters.GET("/:craft-name", craft.Entry)
	}
	recipeRoutes := router.Group("/recipe")
	{
		recipeRoutes.GET("/:id", recipe.CustomRecipe)
	}

	// admin api
	adminApi := router.Group("/api/admin")
	adminApi.Use(middleware.JwtAuthMiddleware(), corsMiddleware)
	{
		adminApi.GET("/admin-login-test", adminLoginTest)
		adminApi.POST("/craft-debug/advertorial", craft.DebugCheckIfAdvertorial)
		adminApi.POST("/craft-debug/common-llm-call-test", admin.LLMDebug)

		adminApi.POST("/recipes", controller.CreateCustomRecipe)
		adminApi.GET("/recipes", controller.ListCustomRecipe)
		adminApi.GET("/recipes/:id", controller.GetCustomRecipe)
		adminApi.PUT("/recipes/:id", controller.UpdateCustomRecipe)
		adminApi.DELETE("/recipes/:id", controller.DeleteCustomRecipe)

		adminApi.GET("/craft-flows", controller.ListCraftFlows)
		adminApi.POST("/craft-flows", controller.CreateCraftFlow)
		adminApi.GET("/craft-flows/:name", controller.GetCraftFlow)
		adminApi.PUT("/craft-flows/:name", controller.UpdateCraftFlow)
		adminApi.DELETE("/craft-flows/:name", controller.DeleteCraftFlow)

	}

}
func adminLoginTest(c *gin.Context) {
	ret := map[string]string{
		"success": "ok",
	}
	c.JSON(http.StatusOK, ret)
}
