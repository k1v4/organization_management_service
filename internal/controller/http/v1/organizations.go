package v1

//type containerRoutes struct {
//	t usecase.IArticleService
//	l logger.Logger
//}
//
//func newArticleRoutes(handler *echo.Group, t usecase.IArticleService, l logger.Logger) {
//	r := &containerRoutes{t, l}
//
//	// GET /api/v1/articles/{id}
//	handler.GET("/articles/:id", r.GetArticle)
//
//	// POST /api/v1/articles
//	handler.POST("/articles", r.PostArticle)
//
//	// DELETE /api/v1/articles/{id}
//	handler.DELETE("/articles/:id", r.DeleteArticle)
//
//	// GET /api/v1/articles?limit=5&offset=0
//	handler.GET("/articles", r.ListArticles)
//
//	// GET /api/v1/user_articles
//	handler.GET("/user_articles", r.GetArticlesByUser)
//}
