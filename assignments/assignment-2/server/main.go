package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Structs for client-server communication
type RegisterArgs struct {
	ClientID int64
	Port     string
}

type RegisterReply struct {
	Success bool
}

type TimeInfoArgs struct {
	ClientID   int64
	ClientTime time.Time
}

type TimeInfoReply struct {
	TimeDifference time.Duration
}

type AdjustmentArgs struct {
	ClientID   int64
	Adjustment time.Duration
}

type AdjustmentReply struct {
	Success bool
}

type BerkeleyServer struct {
	mu        sync.Mutex
	clients   map[int64]string // clientID: port
	timeDiffs map[int64]time.Duration
}

func NewBerkeleyServer() *BerkeleyServer {
	return &BerkeleyServer{
		clients:   make(map[int64]string),
		timeDiffs: make(map[int64]time.Duration),
	}
}

// Register a client
func (bs *BerkeleyServer) Register(args *RegisterArgs, reply *RegisterReply) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if _, exists := bs.clients[args.ClientID]; exists {
		fmt.Println("Failed to register a client, this ID is already taken")
		reply.Success = false
		return nil
	}

	bs.clients[args.ClientID] = args.Port
	fmt.Printf("Client %d registered on port %s\n", args.ClientID, args.Port)
	reply.Success = true
	return nil
}

// Receives time info from client
func (bs *BerkeleyServer) ReportTime(args *TimeInfoArgs, reply *TimeInfoReply) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	serverTime := time.Now()
	diff := serverTime.Sub(args.ClientTime)
	bs.timeDiffs[args.ClientID] = diff

	fmt.Printf("Client %d - Diff: %v\n", args.ClientID, diff)
	reply.TimeDifference = diff
	return nil
}

// Client confirms the adjustment
func (bs *BerkeleyServer) ConfirmAdjustment(args *AdjustmentArgs, reply *AdjustmentReply) error {
	fmt.Printf("Client %d confirmed the adjustment: %v\n", args.ClientID, args.Adjustment)
	reply.Success = true
	return nil
}

func (bs *BerkeleyServer) sendAdjustmentToClient(clientID int64, adjustment time.Duration) error {
	bs.mu.Lock()
	port, exists := bs.clients[clientID]
	bs.mu.Unlock()

	if !exists {
		return fmt.Errorf("Client %d not found", clientID)
	}

	client, err := rpc.Dial("tcp", "localhost:"+port)
	if err != nil {
		return err
	}
	defer client.Close()

	var reply AdjustmentReply
	args := &AdjustmentArgs{
		ClientID:   clientID,
		Adjustment: adjustment,
	}

	err = client.Call("BerkeleyClient.Adjust", args, &reply)
	if err != nil {
		return err
	}

	if reply.Success {
		fmt.Printf("Adjustment sent to client %d: %v\n", clientID, adjustment)
	}
	return nil
}

func (bs *BerkeleyServer) runBerkeleyAlgorithm() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		bs.mu.Lock()

		if len(bs.clients) == 0 {
			fmt.Println("No client registered...")
			bs.mu.Unlock()
			continue
		}

		fmt.Println("\n=== Executing Berkerley Algorithm ===")

		// Collect time diffs
		timeDiffs := make(map[int64]time.Duration)
		timeDiffs[0] = 0 // server

		for clientID, diff := range bs.timeDiffs {
			timeDiffs[clientID] = diff
			fmt.Printf("Client: %d - Diff: %v\n", clientID, diff)
		}

		var sum time.Duration
		count := 0

		for id, diff := range timeDiffs {
			if id == 0 { // server
				continue
			}
			sum += diff
			count++
		}

		mean := sum / time.Duration(count)
		fmt.Printf("Calculated mean: %v\n", mean)

		// Send adjustmnets to clients
		for clientID, diff := range timeDiffs {
			if clientID == 0 { // server
				continue
			}

			adjustment := mean + 1/diff
			fmt.Printf("Sending time adjustment to client %d: %v\n", clientID, adjustment)

			bs.mu.Unlock()
			if err := bs.sendAdjustmentToClient(clientID, adjustment); err != nil {
				fmt.Printf("Error sending adjustmen to client %d: %v\n", clientID, err)
			}
			bs.mu.Lock()
		}

		fmt.Println("=== Sync done ===")
		bs.mu.Unlock()
	}
}

func main() {
	server := NewBerkeleyServer()

	rpc.Register(server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error initializing Berkerley Server: %v\n", err)
	}
	defer listener.Close()

	fmt.Println("Berkerley server listening on Port 8080")
	fmt.Println("Waiting for client connections...")

	// go routine to accept new clients
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Error accepting the connection: %v\n", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	server.runBerkeleyAlgorithm()
}
