package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message Types
const (
	REQUEST = "REQUEST"
	GRANT   = "GRANT"
	RELEASE = "RELEASE"
)

type Message struct {
	Type      string
	ProcessID int
	Timestamp time.Time
}

type Process struct {
	ID             int
	Cluster        *Cluster
	MessageChannel chan Message
	StopChannel    chan bool
}

type Coordinator struct {
	Process *Process
	Queue   []int
}

type Cluster struct {
	Processes   map[int]*Process
	Coordinator *Coordinator
	Resource    sync.Mutex
	StopChannel chan bool
}

func generateRandomID(existingIDs map[int]bool) int {
	for {
		id := rand.Intn(100000) + 1
		if !existingIDs[id] {
			return id
		}
	}
}

func newCluster() *Cluster {
	return &Cluster{
		Processes:   make(map[int]*Process),
		StopChannel: make(chan bool),
	}
}

func (cluster *Cluster) getExistingIDs() map[int]bool {
	existingIDs := make(map[int]bool)
	for id := range cluster.Processes {
		existingIDs[id] = true
	}
	return existingIDs
}

func (cluster *Cluster) newProcess() *Process {
	//cluster.Resource.Lock()
	//defer cluster.Resource.Unlock()

	id := generateRandomID(cluster.getExistingIDs())

	process := &Process{
		ID:             id,
		Cluster:        cluster,
		MessageChannel: make(chan Message, 100),
		StopChannel:    make(chan bool),
	}

	cluster.Processes[id] = process

	fmt.Printf("\nProcess PID: %d / Created", id)
	if len(cluster.Processes) == 1 {
		coordinator := &Coordinator{
			Process: process,
		}
		cluster.Coordinator = coordinator
		fmt.Printf("\nProcess PID: %d / Is the Coordinator", id)
	}
	process.Start()
	return process
}

func (cluster *Cluster) killCoordinator() {
	//cluster.Resource.Lock()
	//defer cluster.Resource.Unlock()

	if cluster.Coordinator == nil {
		return
	}

	fmt.Printf("\nProcess PID: %d / Coordinator killed", cluster.Coordinator.Process.ID)
	close(cluster.Coordinator.Process.StopChannel)
	delete(cluster.Processes, cluster.Coordinator.Process.ID)
	// Limpar a Queue
	cluster.Coordinator.Queue = []int{}

	if len(cluster.Processes) <= 0 {
		return
	}

	// Eleger novo coordenador aleatoriamente
	processIDs := make([]int, 0, len(cluster.Processes))
	for id := range cluster.Processes {
		processIDs = append(processIDs, id)
	}

	randomIndex := rand.Intn(len(processIDs) - 1)
	newCoordinatorId := processIDs[randomIndex]

	cluster.Coordinator.Process = cluster.Processes[newCoordinatorId]
	fmt.Printf("\nProcess PID: %d / Is the new Coordinator", newCoordinatorId)
}

// Inicia o cluster
func (cluster *Cluster) Start() {
	// Cria o primeiro processo (que será o 1o Coordenador)
	cluster.newProcess()

	// Go routine para criar novos processos a cada 40s
	go func() {
		ticker := time.NewTicker(40 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cluster.newProcess()
			case <-cluster.StopChannel:
				return
			}
		}
	}()

	// Go routine para matar o coordenador a cada 1min
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cluster.killCoordinator()
			case <-cluster.StopChannel:
				return
			}
		}
	}()
}

// Para o cluster
func (cluster *Cluster) Stop() {
	close(cluster.StopChannel)
}

// Inicia um processo
func (process *Process) Start() {
	go process.run()
}

func (process *Process) run() {
	// Setando o timer para pegar um recurso a cada 10-25s
	resourceTicker := time.NewTimer(time.Duration(10+rand.Intn(16)) * time.Second)
	defer resourceTicker.Stop()

	for {
		select {
		case <-resourceTicker.C:
			// Tenta acessar o recurso
			process.requestResource()
			// Reseta o ticker para o próximo intervalo
			resourceTicker.Reset(time.Duration(10+rand.Intn(16)) * time.Second)
		case msg := <-process.MessageChannel:
			process.handleMessage(msg)
		case <-process.StopChannel:
			return
		}
	}
}

func (process *Process) requestResource() {
	if process.Cluster.Coordinator == nil {
		fmt.Printf("\nProcess PID: %d / No Coordinator found", process.ID)
		return
	}

	fmt.Printf("\nProcess PID: %d / Requested resource", process.ID)
	process.Cluster.Coordinator.Process.MessageChannel <- Message{
		Type:      REQUEST,
		ProcessID: process.ID,
		Timestamp: time.Now(),
	}
}

func (process *Process) accessResource() {
	process.Cluster.Resource.Lock()
	duration := time.Duration(5+rand.Intn(11)) * time.Second
	fmt.Printf("\nProcess PID: %d / Process accessing the resource for %s", process.ID, duration)
	time.Sleep(duration)
	process.Cluster.Resource.Unlock()
	fmt.Printf("\nProcess PID: %d / Process releasing the resource", process.ID)

	process.Cluster.Coordinator.Process.MessageChannel <- Message{
		Type:      RELEASE,
		ProcessID: process.ID,
		Timestamp: time.Now(),
	}
}

func (process *Process) releaseResource() {
	if len(process.Cluster.Coordinator.Queue) > 1 {
		process.Cluster.Coordinator.Queue = process.Cluster.Coordinator.Queue[1:]
		nextProcessId := process.Cluster.Coordinator.Queue[0]

		if nextProc, exists := process.Cluster.Processes[nextProcessId]; exists {
			nextProc.MessageChannel <- Message{
				Type:      GRANT,
				ProcessID: process.ID,
				Timestamp: time.Now(),
			}
			fmt.Printf("\nProcess PID: %d / Coordinator granting access to PID: %d", process.ID, nextProcessId)
		}
	}
}

func (process *Process) handleMessage(msg Message) {
	switch msg.Type {
	case REQUEST:
		if process.ID == process.Cluster.Coordinator.Process.ID {
			//process.Cluster.Resource.Lock()
			process.Cluster.Coordinator.Queue = append(process.Cluster.Coordinator.Queue, msg.ProcessID)
			fmt.Printf("\nProcess PID: %d / Coordinator added PID: %d to the queue", process.ID, msg.ProcessID)
			//process.Cluster.Resource.Unlock()

			if len(process.Cluster.Coordinator.Queue) == 1 {
				process.Cluster.Processes[msg.ProcessID].MessageChannel <- Message{
					Type:      GRANT,
					ProcessID: process.ID,
					Timestamp: time.Now(),
				}
				fmt.Printf("\nProcess PID: %d / Coordinator granting access to PID: %d", process.ID, msg.ProcessID)
			}
		}

	case GRANT:
		fmt.Printf("\nProcess PID: %d / Process received permition to access the resource", process.ID)
		process.accessResource()

	case RELEASE:
		if process.ID == process.Cluster.Coordinator.Process.ID {
			process.releaseResource()
		}
	}
}

func main() {
	fmt.Println("Initializing Centralized Mutex System")
	fmt.Println("================================================")

	system := newCluster()

	system.Start()

	system.newProcess()
	system.newProcess()
	system.newProcess()

	// Executa por 5 minutos para demonstração
	time.Sleep(5 * time.Minute)
	system.Stop()

	fmt.Println("System finished")
}
