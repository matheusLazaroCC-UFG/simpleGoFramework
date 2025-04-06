package user

import (
    "fmt"
    "sync"
)

// User representa o modelo de dados do usuário
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// UserOperationType enumera os tipos de operação possíveis
type UserOperationType int

const (
    OpCreate UserOperationType = iota
    OpUpdate
    OpDelete
    OpGetByID
    OpGetAll
)

// userOperation descreve uma requisição de operação enviada ao serviço
type userOperation struct {
    opType  UserOperationType
    user    User
    userID  string
    replyCh chan userOperationResult
}

// userOperationResult descreve o resultado devolvido para cada operação
type userOperationResult struct {
    user  User
    users []User
    err   error
}

// UserService define as operações disponíveis para o Controller
type UserService interface {
    GetAll() ([]User, error)
    GetByID(id string) (User, error)
    Create(user User) (User, error)
    Update(id string, user User) (User, error)
    Delete(id string) error
}

// userServiceImpl mantém o estado de usuários e gerencia o canal de operações
type userServiceImpl struct {
    mu         sync.Mutex
    users      map[string]User
    operations chan userOperation
    wg         sync.WaitGroup
}

// NewUserService cria um UserService com um pool de workers
func NewUserService(workerCount int) UserService {
    s := &userServiceImpl{
        users:      make(map[string]User),
        operations: make(chan userOperation),
    }

    // Inicia os workers
    for i := 0; i < workerCount; i++ {
        s.wg.Add(1)
        go s.worker()
    }

    return s
}

// worker processa as operações recebidas no channel
func (s *userServiceImpl) worker() {
    defer s.wg.Done()

    for op := range s.operations {
        switch op.opType {
        case OpCreate:
            createdUser, err := s.createUser(op.user)
            op.replyCh <- userOperationResult{user: createdUser, err: err}

        case OpUpdate:
            updatedUser, err := s.updateUser(op.userID, op.user)
            op.replyCh <- userOperationResult{user: updatedUser, err: err}

        case OpDelete:
            err := s.deleteUser(op.userID)
            op.replyCh <- userOperationResult{err: err}

        case OpGetByID:
            usr, err := s.getByID(op.userID)
            op.replyCh <- userOperationResult{user: usr, err: err}

        case OpGetAll:
            usrs, err := s.getAll()
            op.replyCh <- userOperationResult{users: usrs, err: err}
        }
    }
}

// Métodos privados que manipulam o mapa de usuários, protegidos por Mutex

func (s *userServiceImpl) getAll() ([]User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    result := make([]User, 0, len(s.users))
    for _, u := range s.users {
        result = append(result, u)
    }
    return result, nil
}

func (s *userServiceImpl) getByID(id string) (User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    u, exists := s.users[id]
    if !exists {
        return User{}, fmt.Errorf("usuário com ID %s não encontrado", id)
    }
    return u, nil
}

func (s *userServiceImpl) createUser(u User) (User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.users[u.ID]; exists {
        return User{}, fmt.Errorf("usuário com ID %s já existe", u.ID)
    }
    s.users[u.ID] = u
    return u, nil
}

func (s *userServiceImpl) updateUser(id string, u User) (User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.users[id]; !exists {
        return User{}, fmt.Errorf("usuário com ID %s não encontrado", id)
    }
    u.ID = id
    s.users[id] = u
    return u, nil
}

func (s *userServiceImpl) deleteUser(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.users[id]; !exists {
        return fmt.Errorf("usuário com ID %s não encontrado", id)
    }
    delete(s.users, id)
    return nil
}

// Métodos públicos que postam operações no canal e aguardam a resposta

func (s *userServiceImpl) GetAll() ([]User, error) {
    replyCh := make(chan userOperationResult)
    s.operations <- userOperation{
        opType:  OpGetAll,
        replyCh: replyCh,
    }
    result := <-replyCh
    return result.users, result.err
}

func (s *userServiceImpl) GetByID(id string) (User, error) {
    replyCh := make(chan userOperationResult)
    s.operations <- userOperation{
        opType:  OpGetByID,
        userID:  id,
        replyCh: replyCh,
    }
    result := <-replyCh
    return result.user, result.err
}

func (s *userServiceImpl) Create(u User) (User, error) {
    replyCh := make(chan userOperationResult)
    s.operations <- userOperation{
        opType:  OpCreate,
        user:    u,
        replyCh: replyCh,
    }
    result := <-replyCh
    return result.user, result.err
}

func (s *userServiceImpl) Update(id string, u User) (User, error) {
    replyCh := make(chan userOperationResult)
    s.operations <- userOperation{
        opType:  OpUpdate,
        userID:  id,
        user:    u,
        replyCh: replyCh,
    }
    result := <-replyCh
    return result.user, result.err
}

func (s *userServiceImpl) Delete(id string) error {
    replyCh := make(chan userOperationResult)
    s.operations <- userOperation{
        opType:  OpDelete,
        userID:  id,
        replyCh: replyCh,
    }
    result := <-replyCh
    return result.err
}
