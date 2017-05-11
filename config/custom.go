// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package config

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"

	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/routes"
)

// GetURL returns a full url based off of the config
// and passed in path
func (c *ZENAUTHConfig) GetURL(path string) string {
	var buffer bytes.Buffer
	buffer.WriteString(c.Transport)
	buffer.WriteString("://")
	buffer.WriteString(c.DomainHost)
	if string(c.DomainHost[len(c.DomainHost)-1]) != ":" {
		buffer.WriteString(":")
	}
	if c.Port == 0 {
		// assume 80
		buffer.WriteString("80")
	} else {
		buffer.WriteString(strconv.FormatUint(uint64(c.Port), 10))
	}
	buffer.WriteString(path)
	return buffer.String()
}

// this is where custom stuff for your config can go

// ComputeDependents computes variables that depend on previous config vars
// and it performs sanity checks on convienent defaults
func (c *ZENAUTHConfig) ComputeDependents() error {

	c.PasswordResetLinkBase = c.GetURL(routes.V1 + routes.ResourceUsers + routes.ResourceResetPassword + "?token=")
	// check secret value if in production
	if c.Environment == constants.EnvironmentProduction {
		if reg, err := regexp.Compile(`^([a-zA-Z_]{1}[a-zA-Z0-9_]{31})$`); err != nil {
			return err
		} else if !reg.MatchString(c.HashSecret) {
			return errors.New("HashSecret contains invalid characters or is not the right length")
		}
	}
	// computed once here so we don't have to keep converting to bytes
	c.HashSecretBytes = []byte(c.HashSecret)

	// if you specified things, check that they are not the defaults
	if c.AnalyticsEnabled && c.MixpanelAPIToken == "token" {
		return errors.New("if Mixpanel is enabled you need a proper token")
	}

	// ensure valid new relic configuration
	if c.NewRelicEnabled {
		if len(c.NewRelicKey) == 0 {
			return errors.New("if New Relic is enabled, you need a proper key")
		}
	}

	// ensure we have a valid email setting
	if c.EmailEnabled {
		if c.MailGunDomain == "domain" {
			return errors.New("if Email is enabled, you need a proper mailgun domain")
		}
		if c.MailGunPublicKey == "mailgun" {
			return errors.New("if Email is enabled, you need a proper mailgun public key")
		}
		if c.MailGunPrivateKey == "mailgun" {
			return errors.New("if Email is enabled, you need a proper mailgun private key")
		}
	}

	// *********calculate your custom dependent variable(s) here***********
	//c.AccessorURI = "http://" + c.AccessorServiceFQDN + ":" + c.AccessorPort + routes.V1

	return nil
}
