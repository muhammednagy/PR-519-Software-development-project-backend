package controllers

import (
	"github.com/gogearbox/gearbox"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/middlewares"
)

func (server *Server) initializeRoutes(gb gearbox.Gearbox) []*gearbox.Route {

	return []*gearbox.Route{
		// Home Route
		gb.Get("/", server.Home),

		// Login Route
		gb.Post("/login",server.Login),
		gb.Options("/login",server.LoginOptions),

		//Users routes
		gb.Post("/users", server.CreateUser),
		gb.Options("/users", server.CreateUserOptions),
		gb.Put("/users/{id}", middlewares.SetMiddlewareAuthentication, server.UpdateUser),
		gb.Options("/users/{id}", middlewares.SetMiddlewareAuthentication, server.UpdateUserOptions),
		gb.Delete("/users/{id}", middlewares.SetMiddlewareAuthentication, server.DeleteUser),
	}
}
