package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ParameterStore struct {
	session *session.Session
	ssm     *ssm.SSM
}

func NewParameterStore(profile string) ParameterStore {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:profile, // eg. "default", or "dev-admin" etc
	}))

	ssm := ssm.New(sess)

	return ParameterStore{
		session: sess,
		ssm: ssm,
	}
}

func (ps ParameterStore) Get(key string) (string, error) {

	result, err := ps.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           &key,
	})
	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}
