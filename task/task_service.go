package tasks

import (
    "math"
    "sync"
)

// TaskService define a interface para serviço de tarefas
type TaskService interface {
    FindPrimesInRange(start, end, workerCount int) []int
}

// taskServiceImpl é a implementação concreta do serviço
type taskServiceImpl struct{}

// NewTaskService cria uma instância do serviço de tarefas
func NewTaskService() TaskService {
    return &taskServiceImpl{}
}

// FindPrimesInRange encontra todos os primos entre [start, end] usando workerCount goroutines
func (s *taskServiceImpl) FindPrimesInRange(start, end, workerCount int) []int {
    if start < 2 {
        start = 2
    }
    if workerCount < 1 {
        workerCount = 1
    }

    // Canal de jobs (números para testar)
    jobs := make(chan int, end-start+1)
    // Canal de resultados (enviamos o número se for primo)
    results := make(chan int, end-start+1)

    var wg sync.WaitGroup

    // Inicia os workers
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for n := range jobs {
                if isPrime(n) {
                    results <- n
                }
            }
        }()
    }

    // Preenche o canal jobs com todos os números do intervalo
    go func() {
        for n := start; n <= end; n++ {
            jobs <- n
        }
        close(jobs)
    }()

    // Aguarda os workers terminarem e então fecha o canal de resultados
    go func() {
        wg.Wait()
        close(results)
    }()

    // Coleta todos os primos
    var primes []int
    for prime := range results {
        primes = append(primes, prime)
    }

    return primes
}

// isPrime verifica se um número é primo
func isPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }
    limit := int(math.Sqrt(float64(n)))
    for i := 5; i <= limit; i += 6 {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
    }
    return true
}
