package controllers

import "github.com/kalyaniandhare/fullstack/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Users routes
	s.Router.HandleFunc("/log", middlewares.SetMiddlewareJSON(s.CreateLogConfig)).Methods("POST")
	s.Router.HandleFunc("/log", middlewares.SetMiddlewareJSON(s.GetLogsConfig)).Methods("GET")
	s.Router.HandleFunc("/log/{id}", middlewares.SetMiddlewareJSON(s.GetLogDetail)).Methods("GET")

	//Posts routes
	s.Router.HandleFunc("/log-config", middlewares.SetMiddlewareJSON(s.CreateLog)).Methods("POST")
	s.Router.HandleFunc("/log-config", middlewares.SetMiddlewareJSON(s.GetAllLogs)).Methods("GET")
	s.Router.HandleFunc("/log-config/{id}", middlewares.SetMiddlewareJSON(s.GetLog)).Methods("GET")
}
