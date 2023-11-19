package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandleLXL(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	log.Printf("\"/terabee/lxl\" endpoint called with method %s:\n%s", req.Method, req.Body)
}

// a default endpoint to confirm receipt of a http-post
func HandleDefaultEndpoint(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	log.Printf("\"/\" root endpoint called with method %s:\n%s", req.Method, req.Body)
}

func main() {
	// ctx := context.Background()
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	panic(err)
	// }
	// defer cli.Close()

	// // pull image from remote repo
	// reader, err := cli.ImagePull(ctx, "docker.io/library/debian:bookworm-slim", types.ImagePullOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// defer reader.Close()
	// io.Copy(os.Stdout, reader)

	// resp, err := cli.ContainerCreate(ctx, &container.Config{
	// 	Image: "debian:bookworm-slim",
	// 	Cmd:   []string{"echo", "hello world"},
	// 	Tty:   false,
	// }, nil, nil, nil, "")
	// if err != nil {
	// 	panic(err)
	// }

	// if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
	// 	panic(err)
	// }

	// statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// case <-statusCh:
	// }

	// out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	// if err != nil {
	// 	panic(err)
	// }

	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// set logging flags
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// define handlers
	http.HandleFunc("/terabee/lxl", HandleLXL)
	http.HandleFunc("/", HandleDefaultEndpoint)

	// start the server
	default_port := 8080
	log.Println("== Starting HTTP Server ==")
	log.Printf("== Listening on port %d ==", default_port)
	http.ListenAndServe(fmt.Sprintf(":%d", default_port), nil)
}
