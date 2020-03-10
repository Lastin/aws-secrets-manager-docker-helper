package secretsmanager

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/docker/docker-credential-helpers/credentials"
)

type AWSSecretHelper struct {
}
var notImplemented = errors.New("not implemented")

// ensure ECRHelper adheres to the credentials.Helper interface
var _ credentials.Helper = (*AWSSecretHelper)(nil)

type DockerCredentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (AWSSecretHelper) Add(creds *credentials.Credentials) error {
	// This does not seem to get called
	return notImplemented
}

func (AWSSecretHelper) Delete(serverURL string) error {
	// This does not seem to get called
	return notImplemented
}

func (self AWSSecretHelper) Get(serverURL string) (username string, password string, err error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sess.Config.Region = aws.String("eu-west-2")
	secManager := secretsmanager.New(sess)
	out, err := secManager.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("docker-hub"),
		VersionId:    nil,
		VersionStage: nil,
	})
	if err != nil {
		return
	}
	dockerCredentials := DockerCredentials{}
	err = json.Unmarshal([]byte(*out.SecretString), &dockerCredentials)
	if err != nil {
		return
	}
	return dockerCredentials.Username, dockerCredentials.Password, nil
}

func (self AWSSecretHelper) List() (map[string]string, error) {
	return map[string]string{}, notImplemented
}