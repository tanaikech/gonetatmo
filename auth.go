// Package main (auth.go) :
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/tanaikech/gonetatmo/netatmo"
	"github.com/urfave/cli"
)

const (
	cfgFile    = "gonetatmo.cfg"
	scope      = "read_station"
	cfgpathenv = "GONETATMO_CFG_PATH"
)

// para : Initial parameters
type para struct {
	pstart  time.Time
	WorkDir string
}

// tokens : Structure for retrieving tokens
type tokens struct {
	Accesstoken  string   `json:"access_token"`
	Refreshtoken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	ExpiresIn    int64    `json:"expires_in"`
	ExpireIn     int64    `json:"expire_in"`
	EndTime      int64    `json:"end_time"`
	EndTimeDate  string   `json:"end_time_date"`
}

// configFile : Structure for config file
type configFile struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Mail         string `json:"-"`
	Pass         string `json:"-"`
	*tokens
	GoogleApiKey string `json:"google_api_key"`
}

// materials : Materials for this application
type materials struct {
	*para
	*tokens
	*configFile
}

// makecfgfile :
func (m *materials) makecfgfile() {
	file, _ := json.MarshalIndent(m.configFile, "", "\t")
	ioutil.WriteFile(filepath.Join(m.para.WorkDir, cfgFile), file, 0777)
	fmt.Printf("Updated '%s' at %s. \n", cfgFile, m.para.WorkDir)
}

// getTokens : Retrieve tokens.
func (m *materials) getTokens(body []byte) {
	json.Unmarshal(body, &m.tokens)
	m.configFile.tokens = m.tokens
	m.tokens.EndTime = m.para.pstart.Unix() + m.tokens.ExpiresIn
	m.tokens.EndTimeDate = time.Unix(m.tokens.EndTime, 0).In(time.Local).Format("20060102_15:04:05_MST")
	m.makecfgfile()
}

// getAccesstokenByRefreshtoken : Retrieve access token by existing refresh token.
func (m *materials) getAccesstokenByRefreshtoken() error {
	tokenparams := url.Values{}
	tokenparams.Set("grant_type", "refresh_token")
	tokenparams.Set("refresh_token", m.configFile.tokens.Refreshtoken)
	tokenparams.Set("client_id", m.configFile.ClientId)
	tokenparams.Set("client_secret", m.configFile.ClientSecret)
	body, err := netatmo.GetTokens(tokenparams)
	if err != nil {
		return err
	}
	m.getTokens(body)
	return nil
}

// getNewRefreshtoken : Retrieve new refresh token.
func (m *materials) getNewRefreshtoken() error {
	tokenparams := url.Values{}
	tokenparams.Set("grant_type", "password")
	tokenparams.Set("client_id", m.configFile.ClientId)
	tokenparams.Set("client_secret", m.configFile.ClientSecret)
	tokenparams.Set("username", m.configFile.Mail)
	tokenparams.Set("password", m.configFile.Pass)
	tokenparams.Set("scope", scope)
	body, err := netatmo.GetTokens(tokenparams)
	if err != nil {
		return err
	}
	m.getTokens(body)
	return nil
}

// chkParamsForTokens : Check parameters for retrieving tokens.
func (m *materials) chkParamsForTokens(c *cli.Context) bool {
	if c.String("googleapikey") != "" {
		m.configFile.GoogleApiKey = c.String("googleapikey")
	}
	if c.String("clientid") != "" && c.String("clientsecret") != "" && c.String("email") != "" && c.String("password") != "" {
		m.configFile.ClientId = c.String("clientid")
		m.configFile.ClientSecret = c.String("clientsecret")
		m.configFile.Mail = c.String("email")
		m.configFile.Pass = c.String("password")
		m.getNewRefreshtoken()
		return true
	}
	return false
}

// chkCfg : Check config file.
func (m *materials) chkCfg(c *cli.Context) error {
	var err error
	var cfg []byte
	if !m.chkParamsForTokens(c) {
		if cfg, err = ioutil.ReadFile(filepath.Join(m.para.WorkDir, cfgFile)); err == nil {
			if err = json.Unmarshal(cfg, &m.configFile); err == nil {
				if (m.para.pstart.Unix()-m.configFile.tokens.EndTime) > 0 || m.configFile.tokens.Accesstoken == "" {
					err = m.getAccesstokenByRefreshtoken()
				} else if c.String("googleapikey") != "" {
					m.configFile.GoogleApiKey = c.String("googleapikey")
					m.makecfgfile()
				}
				return err
			} else {
				return err
			}
		} else {
			if !m.chkParamsForTokens(c) {
				return errors.New("No parameters for retrieving refresh token. Please run with the parameters of client id, client secret, mail address and password for Netatmo, again.\nYou can see HELP by\n\n $ gonetatmo --help\n\nCommand for retrieving access token of Netatmo is\n\n $ gonetatmo --clientid ### --clientsecret ### --email ### --password ###\n")
			}
			return nil
		}
	}
	return nil
}

// initParams : Initialize parameters
func initParams() *materials {
	var err error
	var cfgDir string
	cfgDir, err = filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	envDir := os.Getenv(cfgpathenv)
	if envDir != "" {
		cfgDir = envDir
	}
	m := &materials{
		&para{},
		&tokens{},
		&configFile{
			tokens: &tokens{},
		},
	}
	m.para.pstart = time.Now()
	m.para.WorkDir = cfgDir
	return m
}
