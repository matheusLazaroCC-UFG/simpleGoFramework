package framework

import "net/http"

// Controller é a interface para todos os controladores
type Controller interface {
    RegisterRoutes(mux *http.ServeMux)
}
