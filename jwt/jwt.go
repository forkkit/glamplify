package jwt

type Payload struct {
	Customer      string // uuid
	RealUser      string // uuid
	EffectiveUser string // uid
}

