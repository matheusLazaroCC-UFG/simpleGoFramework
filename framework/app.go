package framework

import (
    "log"
    "net/http"
)

// Config define configurações básicas (porta, etc.)
type Config struct {
    Port string
}

// App representa a aplicação principal
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

// RegisterController registra as rotas de um controller no mux
func (a *App) RegisterController(c Controller) {
    c.RegisterRoutes(a.mux)
}

// Start inicia o servidor HTTP na porta configurada
func (a *App) Start() {
    log.Printf("Servidor iniciado na porta %s...\n", a.config.Port)
    log.Fatal(http.ListenAndServe(":"+a.config.Port, a.mux))
}
