package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func hello(rw http.ResponseWriter, r *http.Request, n int) {
	// fmt.Printf("UserAgent : %s called\n", r.UserAgent())
	// fmt.Printf("URL : %s called\n", r.URL)
	rw.Write([]byte(fmt.Sprintf("%d server %s\n", n, r.URL.Path)))
}

func main() {
	var wg sync.WaitGroup
	if len(os.Args) != 2 {
		fmt.Printf("Give %d, takes 1 argument\n", len(os.Args))
		fmt.Println("Usage : server N\n N => number of servers to start")
		return
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Usage : server N\n N => number of servers to start")
		return
	}
	wg.Add(n)
	for ; n > 0; n-- {
		i := n
		go func() {
			mux := http.NewServeMux()
			s := http.Server{Addr: ":808" + strconv.Itoa(i), Handler: mux}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
				hello(rw, r, i)
			})
			mux.HandleFunc("/exit", func(rw http.ResponseWriter, r *http.Request) {
				s.Shutdown(ctx)
				wg.Done()
				fmt.Printf("exited %d\n", i)
			})
			fmt.Println("listening in port : 808" + strconv.Itoa(i))
			go func() {
				s.ListenAndServe()
			}()
		}()
	}
	wg.Wait()
}
