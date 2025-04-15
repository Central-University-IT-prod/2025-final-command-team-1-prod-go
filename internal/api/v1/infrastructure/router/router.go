package router

import (
	docs "example.com/m/docs"
	"example.com/m/internal/api/v1/adapters/controllers"
	"example.com/m/internal/api/v1/infrastructure/middlewares"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const prefix string = "/api/v1"

type Router struct {
	e   *gin.Engine
	am  middlewares.AuthMiddleware
	adm middlewares.AdminMiddleware
}

func NewRouter(e *gin.Engine, am *middlewares.AuthMiddleware, adm *middlewares.AdminMiddleware) *Router {
	return &Router{
		e:   e,
		am:  *am,
		adm: *adm,
	}
}

func (r *Router) BindUserRoutes(uc *controllers.UserController) {
	r.e.POST(prefix+"/users", uc.CreateUser)
	r.e.GET(prefix+"/users/:username", r.am.Authenticate(), uc.GetUserByUsername)
	r.e.GET(prefix+"/users/me", r.am.Authenticate(), uc.GetUserProfile)
	r.e.PATCH(prefix+"/users/me", r.am.Authenticate(), uc.UpdateUserProfile)
	r.e.POST(prefix+"/users/bind_token", r.am.Authenticate(), uc.BindPushToken)
}

func (r *Router) BindAuthRoutes(ac *controllers.AuthController) {
	r.e.POST(prefix+"/auth", ac.AuthorizeUser)
	r.e.PATCH(prefix+"/auth/changePassword", r.am.Authenticate(), ac.ChangePassword)
}

func (r *Router) BindSwaggerRoutes() {
	docs.SwaggerInfo.BasePath = prefix
	docs.SwaggerInfo.Version = "v1"

	ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8000/api/v1/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))
	r.e.GET(prefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (r *Router) BindMetricsRoutes(m *controllers.MetricController) {
	r.e.GET(prefix+"/metrics", m.GetMetrics())
}

func (r *Router) BindPostRoutes(pc *controllers.PostController) {
	r.e.POST(prefix+"/posts", r.am.Authenticate(), pc.CreatePost)
	r.e.GET(prefix+"/posts/:id", r.am.Authenticate(), pc.GetPost)
	r.e.GET(prefix+"/posts/my", r.am.Authenticate(), pc.GetMyPosts)
	r.e.PUT(prefix+"/posts/:id/favorites", r.am.Authenticate(), pc.AddFavorite)
	r.e.DELETE(prefix+"/posts/:id/favorites", r.am.Authenticate(), pc.DeleteFavorite)
	r.e.GET(prefix+"/posts/available", r.am.Authenticate(), pc.GetAllAvailablePosts)
	r.e.GET(prefix+"/posts/favorites", r.am.Authenticate(), pc.GetAllFavourites)
	r.e.GET(prefix+"/posts/search", r.am.Authenticate(), pc.SearchByTitleOrAuthorOrGenre)
	r.e.GET(prefix+"/posts/booked", r.am.Authenticate(), pc.GetAllMyBooked)
	r.e.POST(prefix+"/posts/:id/image", r.am.Authenticate(), pc.AddImage)
}

func (r *Router) BindBookingRoutes(bc *controllers.BookingController) {
	r.e.POST(prefix+"/posts/:id/booking", r.am.Authenticate(), bc.BookBook)
	r.e.DELETE(prefix+"/posts/:id/booking", r.am.Authenticate(), bc.DeleteBooking)
	r.e.PUT(prefix+"/posts/:id/mark-taken", r.am.Authenticate(), bc.MarkAsTaken)
}

func (r *Router) BindReviewRoutes(rc *controllers.ReviewController) {
	r.e.POST(prefix+"/reviews", r.am.Authenticate(), rc.CreateReview)
	r.e.GET(prefix+"/users/:username/reviews", r.am.Authenticate(), rc.GetReviewsForUser)
}

func (r *Router) BindPlaceRoutes(pc *controllers.PlaceController) {
	r.e.POST(prefix+"/places", r.am.Authenticate(), r.adm.CheckAdminStatus(), pc.CreatePlace)
	r.e.GET(prefix+"/places", pc.GetPlaces)
	r.e.DELETE(prefix+"/places/:id", r.am.Authenticate(), r.adm.CheckAdminStatus(), pc.DeletePlace)
}

func (r *Router) BindChatBotRoutes(cc *controllers.ChatBotController) {
	r.e.POST(prefix+"/chat", r.am.Authenticate(), cc.SendMessage)
	r.e.GET(prefix+"/chat", r.am.Authenticate(), cc.GetChat)
}