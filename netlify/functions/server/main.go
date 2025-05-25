package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rebelopsio/macro-calc/internal/handlers"
	"github.com/rebelopsio/macro-calc/internal/templates"
	"github.com/a-h/templ"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	// Static files - served directly by Netlify
	
	// Routes
	e.GET("/", func(c echo.Context) error {
		component := templates.Index()
		return renderTempl(c, component)
	})

	e.POST("/calculate", handlers.Calculate)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	echoLambda = echoadapter.New(e)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}

func renderTempl(c echo.Context, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return component.Render(c.Request().Context(), c.Response())
}