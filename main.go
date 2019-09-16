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

func fileHandler(filePath string) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		file, err := os.Open(filePath)
		if err != nil {

		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment")
		io.Copy(w, file)

	}
}

func (s Server) routesHandler() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		fmt.Fprintf(w, "<h1>Hosted Files</h1>")
		for _, f := range s.Files {
			fmt.Fprintf(w, "<p><a href=\"http://%s%s/%s\"> http://%s%s/%s </a></p>", s.LocalIP, s.Port, f, s.LocalIP, s.Port, f)
		}

		fmt.Fprintf(w, `<br/><h1>Upload File</h1><br/>
		<form
			enctype="multipart/form-data"
			action="upload"
			method="post">
			<input type="file" name="myFile" />
			<input type="submit" value="upload" />
		  </form>
		</body>
	  </html>`)

	}
}

type Server struct {
	Port    string
	Files   []string
	LocalIP string
	Router  *httprouter.Router
}

func (s *Server) routes() {

	for _, f := range s.Files {
		s.Router.GET("/"+f, fileHandler(f))
	}

	s.Router.GET("/", s.routesHandler())

	s.Router.POST("/upload", uploadFile)

}

func NewServer(port string, localIP string, filePaths []string) Server {

	srv := Server{
		Port:    port,
		LocalIP: localIP,
		Router:  httprouter.New(),
		Files:   filePaths,
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

	fmt.Println("server started on: http://" + localIP + port)

	log.Fatal(http.ListenAndServe(srv.Port, srv.Router))
}
