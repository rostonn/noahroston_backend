package main

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/login/{provider}", a.loginUser).Methods("POST", "OPTIONS")

	a.Router.HandleFunc("/test", a.Test).Methods("GET")
}
