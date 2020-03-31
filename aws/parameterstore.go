package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cultureamp/glamplify/log"
)

type ParameterStore struct {
	session *session.Session
	ssm     *ssm.SSM
	logger  *log.Logger
}

func NewParameterStore(ctx context.Context, profile string) ParameterStore {
	logger := log.New(ctx)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:profile, // eg. "default", or "dev-admin" etc
	}))

	ssm := ssm.New(sess)

	return ParameterStore{
		session: sess,
		ssm: ssm,
		logger:  logger,
	}
}

func (ps ParameterStore) Get(key string) (string, error) {

	result, err := ps.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           &key,
	})
	if err != nil {
		ps.logger.Warn("parameter_store_get_key", log.Fields{
			"key": key,
			"error": err.Error(),
		})
		return "", err
	}

	return *result.Parameter.Value, nil
}
