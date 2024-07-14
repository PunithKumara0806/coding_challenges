package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type ServerStatus struct {
	Address  string
	UpStatus bool
}

type ServerConn struct {
	mu            sync.Mutex
	servers       []ServerStatus
	active_server []int
}

func (sc *ServerConn) healthCheck() {
	sc.mu.Lock()
	fmt.Println("started healthCheck")
	defer sc.mu.Unlock()
	sc.active_server = make([]int, 0)
	for i, v := range sc.servers {
		_, err := http.Head(v.Address)
		if err != nil {
			sc.servers[i].UpStatus = false
		} else {
			sc.servers[i].UpStatus = true
			sc.active_server = append(sc.active_server, i)
		}
	}
	fmt.Println(sc.active_server)
	fmt.Println("finished healthCheck")
}

func (sc *ServerConn) getNext() ServerStatus {
	//evenly distribution
	i := rand.Int() % len(sc.active_server)
	return sc.servers[sc.active_server[i]]
}

var sc ServerConn

func createFull_addrs(addrs []string, ports map[int][]string) []ServerStatus {
	serverStatus := make([]ServerStatus, 0)
	for i, addr := range addrs {
		for _, port := range ports[i] {
			serverStatus = append(serverStatus,
				ServerStatus{
					Address:  addr + ":" + port,
					UpStatus: false,
				})
		}
	}
	return serverStatus
}

func GetResponse(rw http.ResponseWriter, r *http.Request, depth int) {
	v := sc.getNext()
	if !v.UpStatus {
		v = sc.getNext()
		if !v.UpStatus {
			rw.WriteHeader(500)
			return
		}
	}
	res, err := http.Get(v.Address + r.URL.Path)
	if err != nil {
		fmt.Println(err)
		//limit the max number of retries
		if depth <= 3 {
			GetResponse(rw, r, depth+1)
		}
		return
	}
	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	rw.Write(b)
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	if res.StatusCode >= 500 {
		fmt.Println("res : ", res)
		v.UpStatus = false
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage error:\n loadBalancer N => N servers ( 8081 - 808N )")
		return
	}
	addrs := []string{
		"http://localhost",
	}
	ports := make(map[int][]string)
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Usage error:\n loadBalancer N => N servers ( 8081 - 808N )")
		return
	}
	for i := 1; i <= n; i++ {
		ports[0] = append(ports[0], "808"+strconv.Itoa(i))
	}
	sc.servers = createFull_addrs(addrs, ports)
	sc.healthCheck()
	go func() {
		ch := time.NewTicker(time.Duration(5) * time.Second)
		for _ = range ch.C {
			sc.healthCheck()

		}
	}()
	http.HandleFunc("/",
		func(rw http.ResponseWriter, r *http.Request) {
			GetResponse(rw, r, 1)
		})
	http.ListenAndServe(":8080", nil)
}
