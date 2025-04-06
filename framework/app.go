package framework

import (
    "log"
    "net/http"
)

// Config define as configurações, como porta
type Config struct {
    Port string
}

// App representa nossa aplicação principal
type App struct {
    mux    *http.ServeMux
    config *Config
}

// NewApp cria uma nova instância da aplicação
func NewApp(config *Config) *App {
    return &App{
        mux:    http.NewServeMux(),
        config: config,
    }
}

// RegisterController registra as rotas de um Controller
func (a *App) RegisterController(c Controller) {
    c.RegisterRoutes(a.mux)
}

// Start inicia o servidor HTTP
func (a *App) Start() {
    log.Printf("Iniciando servidor na porta %s...\n", a.config.Port)
    log.Fatal(http.ListenAndServe(":"+a.config.Port, a.mux))
}
