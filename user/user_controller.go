package user

import (
    "encoding/json"
    "net/http"
    "strings"

    "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
)

// UserController expõe endpoints para CRUD de usuários
type UserController struct {
    service UserService
}

// NewUserController cria uma instância do controlador com injeção de serviço
func NewUserController(service UserService) *UserController {
    return &UserController{service: service}
}

// RegisterRoutes registra as rotas no mux
func (uc *UserController) RegisterRoutes(mux *http.ServeMux) {
    // /users e /users/
    mux.HandleFunc("/users", uc.handleUsers)
    mux.HandleFunc("/users/", uc.handleUserByID)
}

// handleUsers lida com GET/POST em /users
func (uc *UserController) handleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        uc.getAllUsers(w, r)
    case http.MethodPost:
        uc.createUser(w, r)
    default:
        http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
    }
}

// getAllUsers chama o service para obter todos os usuários
func (uc *UserController) getAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := uc.service.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// createUser cria um novo usuário
func (uc *UserController) createUser(w http.ResponseWriter, r *http.Request) {
    var newUser User
    if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
        http.Error(w, "Corpo inválido", http.StatusBadRequest)
        return
    }

    created, err := uc.service.Create(newUser)
    if err != nil {
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(created)
}

// handleUserByID lida com GET/PUT/DELETE em /users/{id}
func (uc *UserController) handleUserByID(w http.ResponseWriter, r *http.Request) {
    // Extrai o ID do usuário removendo o prefixo /users/
    id := strings.TrimPrefix(r.URL.Path, "/users/")
    switch r.Method {
    case http.MethodGet:
        uc.getUserByID(w, r, id)
    case http.MethodPut:
        uc.updateUser(w, r, id)
    case http.MethodDelete:
        uc.deleteUser(w, r, id)
    default:
        http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
    }
}

func (uc *UserController) getUserByID(w http.ResponseWriter, r *http.Request, id string) {
    user, err := uc.service.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (uc *UserController) updateUser(w http.ResponseWriter, r *http.Request, id string) {
    var data User
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Corpo inválido", http.StatusBadRequest)
        return
    }
    updated, err := uc.service.Update(id, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updated)
}

func (uc *UserController) deleteUser(w http.ResponseWriter, r *http.Request, id string) {
    if err := uc.service.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    // 204 No Content
    w.WriteHeader(http.StatusNoContent)
}
