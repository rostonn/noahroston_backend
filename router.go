package main

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/login/{provider}", a.loginUser).Methods("POST", "OPTIONS")

	a.Router.HandleFunc("/test", a.Test).Methods("GET")

	a.Router.HandleFunc("/oauth/check_token", a.checkToken).Methods("POST", "OPTIONS")
}
