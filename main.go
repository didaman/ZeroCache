package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"zerocache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *zerocache.Group {
	newgroup := zerocache.NewGroup("scores", 2<<10, zerocache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[DB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	return newgroup
}

func startCacheServer(addr string, addrs []string, zero *zerocache.Group) {
	peers := zerocache.NewHTTPPool(addr)
	peers.Set(addrs...)
	zero.RegisterPeers(peers)
	log.Println("zerocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))

}

func startAPIServer(apiAddr string, zero *zerocache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := zero.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Zerocache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	zero := createGroup()
	if api {
		go startAPIServer(apiAddr, zero)
	}
	startCacheServer(addrMap[port], []string(addrs), zero)
}
