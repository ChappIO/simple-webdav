package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/net/webdav"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"simple-webdav/pkg/dav"
	"simple-webdav/pkg/htpasswd"
	"syscall"
)

var version = "development"
var printVersion = false
var printHelp = false
var homeFolder = ""
var dataRoot = ""
var htpasswordFile = ""
var serverPort = 0

func init() {
	flag.BoolVar(&printVersion, "version", false, "Show the version")
	flag.BoolVar(&printHelp, "help", false, "Show this page")
	flag.IntVar(&serverPort, "serverPort", 5780, "The port on which the http server listens")
	flag.StringVar(&homeFolder, "homeDir", "", "The path to the root directory for data storage. (default $HOME/.simple-webdav)")
	flag.StringVar(&htpasswordFile, "htpasswd", ".htpasswd", "The path to the .htpasswd file. This path is relative to the homeDir")
	flag.StringVar(&dataRoot, "dataDir", "data", "The path to the root folder for data files. This path is relative to the homeDir")
	flag.Parse()
	if printHelp {
		flag.Usage()
		os.Exit(0)
	}
	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if homeFolder == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("homeDir not found, please provide it using the --homeDir option")
		}
		homeFolder = path.Clean(path.Join(home, ".simple-webdav"))
	}
	htpasswordFile = path.Clean(path.Join(homeFolder, htpasswordFile))
	dataRoot = path.Clean(path.Join(homeFolder, dataRoot))
}

func main() {
	if err := os.MkdirAll(dataRoot, 0770); err != nil {
		log.Fatalf("failed to create data directory: %s", err)
	}
	htpwd, err := htpasswd.LoadHtPasswordFile(htpasswordFile)
	if err != nil {
		log.Fatalf("failed to read htpassword file: %s", err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Fatalf("failed to start the http server: %s", err)
	}
	log.Printf("server started on %s", listen.Addr())
	router := http.NewServeMux()
	mainWebdav(htpwd, router)

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("404: %s", request.URL)
		http.Error(writer, "Not Found", 404)
	})
	server := http.Server{
		Handler: router,
	}
	go func() {
		if err := server.Serve(listen); err != http.ErrServerClosed {
			log.Fatalf("http server crashed: %s", err)
		}
	}()

	// wait for kill signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("shutting down...")
	server.Close()
	listen.Close()
	log.Println("done")
}

func mainWebdav(htpwd *htpasswd.File, router *http.ServeMux) {
	router.HandleFunc("/.well-known/webdav", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/webdav/", 301)
	})
	fs := &dav.UserScopedFileSystem{
		RootDir:    dataRoot,
		SubDir:     "Files",
		FileSystem: webdav.Dir(dataRoot),
	}
	router.Handle("/webdav/", htpwd.BasicAuth(
		&webdav.Handler{
			Prefix:     "/webdav/",
			FileSystem: fs,
			LockSystem: webdav.NewMemLS(),
		},
	))
	ctx := context.TODO()
	for _, username := range htpwd.Users() {
		fs.Mkdir(context.WithValue(ctx, "Username", username), "", 0770)
	}
}
