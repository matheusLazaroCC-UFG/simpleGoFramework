package framework

import (
    "net/http"
)

// Controller é a interface padrão dos controllers
type Controller interface {
    RegisterRoutes(mux *http.ServeMux)
}
