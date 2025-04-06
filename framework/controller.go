package framework

import "net/http"

// Controller é a interface que todos os controllers devem implementar
type Controller interface {
    RegisterRoutes(mux *http.ServeMux)
}
