package framework

import (
    "log"
    "net/http"
)

type Config struct {
    Port string
}

type App struct {
    mux    *http.ServeMux
    config *Config
}

func NewApp(config *Config) *App {
    return &App{
        mux:    http.NewServeMux(),
        config: config,
    }
}

func (a *App) RegisterController(c Controller) {
    c.RegisterRoutes(a.mux)
}

func (a *App) Start() {
    log.Printf("Servidor iniciado na porta %s...\n", a.config.Port)
    log.Fatal(http.ListenAndServe(":"+a.config.Port, a.mux))
}
