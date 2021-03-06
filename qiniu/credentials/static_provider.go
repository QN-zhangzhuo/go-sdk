package credentials

import (
	"github.com/QN-zhangzhuo/go-sdk/qiniu/qerr"
)

// StaticProviderName provides a name of Static provider
const StaticProviderName = "StaticProvider"

var (
	// ErrStaticCredentialsEmpty is emitted when static credentials are empty.
	ErrStaticCredentialsEmpty = qerr.New("EmptyStaticCreds", "static credentials are empty", nil)
)

// A StaticProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticProvider struct {
	Value
}

// NewStaticCredentials returns a pointer to a new Credentials object
// wrapping a static credentials value provider.
func NewStaticCredentials(id, secret string) *Credentials {
	return NewCredentials(&StaticProvider{Value: Value{
		AccessKey: id,
		SecretKey: []byte(secret),
	}})
}

// NewStaticCredentialsFromCreds returns a pointer to a new Credentials object
// wrapping the static credentials value provide. Same as NewStaticCredentials
// but takes the creds Value instead of individual fields
func NewStaticCredentialsFromCreds(creds Value) *Credentials {
	return NewCredentials(&StaticProvider{Value: creds})
}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s *StaticProvider) Retrieve() (Value, error) {
	if s.AccessKey == "" || string(s.SecretKey) == "" {
		return Value{ProviderName: StaticProviderName}, ErrStaticCredentialsEmpty
	}

	if len(s.Value.ProviderName) == 0 {
		s.Value.ProviderName = StaticProviderName
	}
	return s.Value, nil
}
