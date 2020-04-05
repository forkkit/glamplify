package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cultureamp/glamplify/cache"
	"time"
)

type ParameterStore struct {
	session *session.Session
	ssm     *ssm.SSM
	cache  *cache.Cache
}

func NewParameterStore(profile string) *ParameterStore {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile, // eg. "default", or "dev-admin" etc
	}))

	ssm := ssm.New(sess)
	c := cache.New()

	return &ParameterStore{
		session: sess,
		ssm:     ssm,
		cache:   c,
	}
}

func (ps ParameterStore) Get(key string) (string, error) {

	if x, found := ps.cache.Get(key); found {
		if val, ok := x.(string); ok {
			return val, nil
		}
	}

	// This makes a network call - can be slow...
	result, err := ps.ssm.GetParameter(&ssm.GetParameterInput{
		Name: &key,
	})
	if err != nil {
		return "", err
	}
	val :=  *result.Parameter.Value

	// cache this for a minute, in case multiple calls request the same key in a short duration
	ps.cache.Set(key, val, 1 * time.Minute)
	return val, nil
}
