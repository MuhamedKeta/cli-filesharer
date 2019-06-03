package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func fileHandler(file *os.File) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment")
		io.Copy(w, file)

	}
}

func (s Server) routesHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		for _, f := range s.Files {
			fmt.Fprintf(w, "<p><a href=\"http://%s%s/%s\"> http://%s%s/%s </a></p>", s.LocalIP, s.Port, f.Name(), s.LocalIP, s.Port, f.Name())
		}

	}
}

type Server struct {
	Port    string
	Files   []*os.File
	LocalIP string
	Router  *httprouter.Router
}

func (s Server) Close() error {

	var retErr error
	for _, f := range s.Files {
		err := f.Close()
		if err != nil {
			retErr = err
		}
	}

	return retErr
}

func (s *Server) routes() {
	for _, f := range s.Files {
		s.Router.GET("/"+f.Name(), fileHandler(f))
	}

	s.Router.GET("/", s.routesHandler())
}
func (s *Server) AddFile(filePath string) error {

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	s.Files = append(s.Files, f)

	return nil

}

func NewServer(port string, localIP string, filePaths []string) Server {

	srv := Server{
		Port:    port,
		LocalIP: localIP,
		Router:  httprouter.New(),
	}

	for _, f := range filePaths {
		srv.AddFile(f)
	}

	srv.routes()

	return srv
}

func getLocalIP() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	var currentIP string

	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
			}
		}
	}

	return currentIP, nil
}
func main() {

	var port string

	flag.StringVar(&port, "port", ":8008", "port to open conn")
	flag.StringVar(&port, "p", ":8008", "port to open conn")
	flag.Parse()

	var files []string
	for _, f := range flag.Args() {

		files = append(files, f)

	}

	localIP, err := getLocalIP()
	if err != nil {
		log.Fatal(err)
	}

	srv := NewServer(port, localIP, files)
	defer srv.Close()

	fmt.Println("server started on: http://" + localIP + port)

	log.Fatal(http.ListenAndServe(srv.Port, srv.Router))
}
