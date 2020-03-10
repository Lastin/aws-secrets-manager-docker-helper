package main

import (
	"aws-secrets-manager-docker-credentials-helper/secretsmanager"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/docker/docker-credential-helpers/credentials"
	"io"
	"os"
)

func main() {
	helper := secretsmanager.AWSSecretHelper{}
	username, password, err := helper.Get("")
	if err != nil {
		fmt.Println("failed to retrieve credentials", err)
	} else {
		fmt.Printf("Username: %s. Password: %s", username, password)
	}
	Serve()
}

func Serve() {
	var err error
	if len(os.Args) != 2 {
		err = fmt.Errorf("Usage: %s <store|get|erase|list|version>", os.Args[0])
	}

	if err == nil {
		err = HandleCommand(os.Args[1], os.Stdin, os.Stdout)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
}

func HandleCommand(key string, reader io.Reader, writer io.Writer) error {
	helper := secretsmanager.AWSSecretHelper{}
	switch key {
	case "store":
		return helper.Add(&credentials.Credentials{})
	case "get":
		return getCredentials(helper, writer)
	case "erase":
		return helper.Delete("")
	case "list":
		return nil
	}
	return fmt.Errorf("Unknown credential action `%s`", key)
}

type Credentials struct {
	ServerURL string
	Username  string
	Secret    string
}

func getCredentials(helper secretsmanager.AWSSecretHelper, writer io.Writer) error {
	buffer := new(bytes.Buffer)
	buffer.Reset()
	username, secret, _ := helper.Get("")
	resp := Credentials{
		ServerURL: "",
		Username:  username,
		Secret:    secret,
	}
	if err := json.NewEncoder(buffer).Encode(resp); err != nil {
		return err
	}
	_, _ = fmt.Fprint(writer, buffer.String())
	return nil
}