// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package config

// this is your config struct that you define
// from this you can generate the html page to get the config route
// TODO: custom go generate tool to generate the html template

import "time"

// references
// https://blog.gopheracademy.com/advent-2013/day-03-building-a-twelve-factor-app-in-go/
// variable naming: [PROJECT]_[APP]_[VARNAME]
// [PROJECT]_[APP] is the Hatch prefix

// ZENAUTHConfig the configuration
// TODO: fill in descriptions for each entry
type ZENAUTHConfig struct {
	LogQueries                         bool          `default:"false"`
	LogLevel                           string        `default:"INFO"`
	HashSecret                         string        `required:"true"`
	HashSecretBytes                    []byte        `ignored:"true"`
	APIToken                           string        `required:"true"`
	PasswordResetValidTokenDuration    time.Duration `default:"60m"`
	RandomPasswordLengthBeforeEncoding uint16        `default:"8"`
	JwtUserTokenDuration               time.Duration `default:"8760h"`
	JwtClaimUserID                     string        `default:"userid"`
	JwtClaimUserEmail                  string        `default:"emailaddr"`
	DefaultContentType                 string        `default:"application/json"`
	APITokenHeader                     string        `default:"x-api-token"`
	AuthTokenHeader                    string        `default:"x-authentication-token"`
	UUIDRegex                          string        `default:"[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"`
	Transport                          string        `default:"https"`
	DomainHost                         string        `required:"true"`
	TestDomainHost                     string        `default:"localhost"`
	Port                               uint16        `default:"5000"`
	GRPCPort                           uint16        `default:"5001"`
	MinPasswordLength                  uint16        `default:"8"`

	BcryptCost          uint16 `default:"8"`
	AllowHashDowngrades bool   `default:"false"`

	DeflateCompression    int8          `default:"-1"`
	DrainAndDieTimeout    time.Duration `default:"60s"`
	TransportReadTimeout  time.Duration `default:"60s"`
	TransportWriteTimeout time.Duration `default:"60s"`
	Environment           string        `required:"true"`
	PasswordResetLinkBase string        `ignored:"true"`
	AnalyticsEnabled      bool          `default:"false"`
	MixpanelAPIToken      string        `default:"token"`
	EmailEnabled          bool          `default:"false"`
	MailGunDomain         string        `default:"domain"`
	MailGunFrom           string        `default:"from@email.com"`
	MailGunPublicKey      string        `default:"mailgun"`
	MailGunPrivateKey     string        `default:"mailgun"`
	NewRelicEnabled       bool          `default:"false"`
	NewRelicName          string        `default:"to_be_filled_in"`
	NewRelicPoll          uint16        `default:"60"`
	NewRelicGCEnabled     bool          `default:"false"`
	NewRelicGCPoll        uint16        `default:"60"`
	NewRelicMemEnabled    bool          `default:"false"`
	NewRelicMemPoll       uint16        `default:"60"`
	NewRelicKey           string        `default:"nrkey"`
	RequestIDHeader       string        `default:"X-Request-ID"`

	PostgreSQLHost           string        `default:"localhost"`
	PostgreSQLPort           uint16        `default:"5432"`
	PostgreSQLUsername       string        `default:"postgres"`
	PostgreSQLPassword       string        `required:"false"`
	PostgreSQLDatabase       string        `required:"true"`
	PostgreSQLSSL            *bool         `default:"true"`
	PostgreSQLRetryNumTimes  uint16        `default:"10"`
	PostgreSQLRetrySleepTime time.Duration `default:"30s"`

	// how to override the environment var
	//AccessorServiceFQDN string `envconfig:"ACCESSOR_ENV_DOCKERCLOUD_SERVICE_FQDN"`
	// dependent variable here, ignored, but calculated in computeDependents
	//AccessorURI         string `ignored:"true"`
	TemplatesPath    string `default:"email/templates"`
	AppName          string `default:"ZenAuth"`
	ResetPasswordURL string `required:"true"`
	VerifyEmailURL   string `required:"true"`
}
