package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	http.HandleFunc("/exec", func(writer http.ResponseWriter, request *http.Request) {
		var err error
		script, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Print("err")
		}
		query := request.URL.Query()
		ip := query.Get("ip")
		port := query.Get("port")
		addr := fmt.Sprintf("%s:%s", ip, port)
		username := query.Get("username")
		password := query.Get("password")
		//err := Exec(Command, writer)
		err = exec(string(script), addr, username, password, writer)
		if err != nil {
			writer.Write([]byte(err.Error()))
		}
	})

	err := http.ListenAndServe(":9080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func exec(script, addr, username, password string, writer http.ResponseWriter) error {
	session, err := connect(username, password, addr)
	if err != nil {
		return err
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		flusher, ok := writer.(http.Flusher)
		if !ok {
			// err 不会发生
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		for {
			read, err := reader.ReadByte()
			if err != nil || err == io.EOF {
				return
			}
			fmt.Fprintf(writer, "%s", string(read))
			flusher.Flush()
		}
	}()
	err = session.Run(script)
	if err != nil {
		writer.Write([]byte(err.Error()))
	}
	wg.Wait()
	return nil
}

func connect(user, password, addr string) (session *ssh.Session, err error) {
	var (
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
	)

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}
