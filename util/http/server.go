package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func ReadBody(req *http.Request) []byte {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic("error parsing body of request")
	}
	return body
}

func HandleFlow(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)
	log.Printf("\"/terabee/flow\" endpoint called with method %s:\n%s", req.Method, string(body))
	RunContainer("docker.io/library/terabee", string(body))
}

// a default endpoint to confirm receipt of a http-post
func HandleDefaultEndpoint(w http.ResponseWriter, req *http.Request) {
	// for testing just print the body.
	body := ReadBody(req)
	log.Printf("\"/\" root endpoint called with method %s:\n%s", req.Method, string(body))
}

func RunContainer(img, msg string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// // pull image from remote repo
	// reader, err := cli.ImagePull(ctx, "docker.io/library/debian:bookworm-slim", types.ImagePullOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// defer reader.Close()
	// io.Copy(os.Stdout, reader)
	dest := "http://chaosbox.princeton.edu/frost/v1.1/Observations"
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: img,
		Cmd: []string{"/bin/sh", "-c",
			fmt.Sprintf("convert/convert_flow %s 1 | post/post %s admin admin", msg, dest)},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

func main() {
	// set logging flags
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// define handlers
	http.HandleFunc("/terabee/flow", HandleFlow)
	http.HandleFunc("/", HandleDefaultEndpoint)

	// start the server
	default_port := 8080
	log.Println("== Starting HTTP Server ==")
	log.Printf("== Listening on port %d ==", default_port)
	http.ListenAndServe(fmt.Sprintf(":%d", default_port), nil)
}
