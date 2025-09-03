package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
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

type BerkeleyClient struct {
	ID         int64
	ServerAddr string
	Port       string
	TimeOffset time.Duration
	isRunning  bool
}

func NewBerkeleyClient(id int64, serverAddr, port string) *BerkeleyClient {
	// Random offset to simulate desync
	offset := time.Duration(rand.Intn(10)-5) * time.Second
	return &BerkeleyClient{
		ID:         id,
		ServerAddr: serverAddr,
		Port:       port,
		TimeOffset: offset,
		isRunning:  true,
	}
}

func (bc *BerkeleyClient) GetCurrentTime() time.Time {
	return time.Now().Add(bc.TimeOffset)
}

func (bc *BerkeleyClient) Adjust(args *AdjustmentArgs, reply *AdjustmentReply) error {
	if args.ClientID != bc.ID {
		return fmt.Errorf("Client ID does not match")
	}

	fmt.Printf("Received server adjustment: %v\n", args.Adjustment)
	bc.TimeOffset += args.Adjustment
	fmt.Printf("New offset: %v, Current Time: %v\n", bc.TimeOffset, bc.GetCurrentTime())

	reply.Success = true
	return nil
}

// Register client on server
func (bc *BerkeleyClient) registerWithServer() error {
	client, err := rpc.Dial("tcp", bc.ServerAddr)
	if err != nil {
		return err
	}
	defer client.Close()

	args := &RegisterArgs{
		ClientID: bc.ID,
		Port:     bc.Port,
	}
	var reply RegisterReply

	err = client.Call("BerkeleyServer.Register", args, &reply)
	if err != nil {
		return err
	}

	if reply.Success {
		fmt.Printf("Client %d registered on server\n", bc.ID)
	}
	return nil
}

func (bc *BerkeleyClient) sendTimeToServer() {
	client, err := rpc.Dial("tcp", bc.ServerAddr)
	if err != nil {
		fmt.Printf("Error connecting on server: %v\n", err)
		return
	}
	defer client.Close()

	args := &TimeInfoArgs{
		ClientID:   bc.ID,
		ClientTime: bc.GetCurrentTime(),
	}
	var reply TimeInfoReply

	err = client.Call("BerkeleyServer.ReportTime", args, &reply)
	if err != nil {
		fmt.Printf("Error sending time info to server: %v\n", err)
		return
	}

	fmt.Printf("Calculated diff: %v\n", reply.TimeDifference)
}

func (bc *BerkeleyClient) startSyncLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for bc.isRunning {
		select {
		case <-ticker.C:
			bc.sendTimeToServer()
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Use: go run client.go [client-id] [server-port] [client-port]")
		fmt.Println("Example: go run client.go 1 8080 8081")
		return
	}

	// Run args
	clientID := int64(0)
	fmt.Sscanf(os.Args[1], "%d", &clientID)
	serverPort := os.Args[2]
	clientPort := os.Args[3]

	serverAddr := "localhost:" + serverPort

	client := NewBerkeleyClient(clientID, serverAddr, clientPort)

	rpc.Register(client)

	listener, err := net.Listen("tcp", ":"+clientPort)
	if err != nil {
		fmt.Printf("Error initializing client: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Client %d listening on port %s\n", clientID, clientPort)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if client.isRunning {
					fmt.Printf("Error accepting the connection: %v\n", err)
				}
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	if err := client.registerWithServer(); err != nil {
		fmt.Printf("Error registering the client: %v\n", err)
		return
	}

	fmt.Printf("Client %d initialized with offset: %v\n", clientID, client.TimeOffset)
	fmt.Printf("Starting time: %v\n", client.GetCurrentTime())

	client.startSyncLoop()

	// This keeps the client running
	select {}
}
