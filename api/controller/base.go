package controller

import "github.com/labstack/echo/v4"

type BaseController interface {
	Register(g *echo.Group)
}
