package main

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/mskrha/oauth2-token-refresher"
)

type Status struct {
	Total  uint32
	Unseen uint32
}

/*
	Definition of the IMAP server
*/
type Server struct {
	/*
		Internal name used as reference from account specification
	*/
	Name string `json:"name"`

	/*
		Host name where to connect
	*/
	Host string `json:"host"`

	/*
		TCP port of where the IMAP server is listening
	*/
	Port uint `json:"port"`

	/*
		Whether use or not use TLS
	*/
	TLS bool `json:"tls"`
}

/*
	Definition of the IMAP account
*/
type Account struct {
	/*
		Name shown in the result line
	*/
	Name string `json:"name"`

	/*
		Which of defined servers to use
	*/
	Server string `json:"server"`

	/*
		Authentication type
		- password
		- oauth2
	*/
	Auth string `json:"auth"`

	/*
		Account credentials
	*/
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Checker struct {
	/*
		Instance name
	*/
	name string

	/*
		Connection string (host:port)
	*/
	addr string
	/*
		Use or use not secure connection
	*/
	tls bool

	/*
		Authentication type
		- password
		- oauth2
	*/
	auth string

	/*
		Account login credentials
	*/
	username string
	password string
	token    string

	/*
		IMAP client connection
	*/
	conn *client.Client
}

/*
	Create a new Checker instance
*/
func NewChecker(a Account) (*Checker, error) {
	var c Checker

	/*
		Check if all requested specifications have been set
	*/
	if len(a.Name) == 0 {
		return nil, fmt.Errorf("No account name specified")
	}
	if len(a.Server) == 0 {
		return nil, fmt.Errorf("No server specified for account %s", a.Name)
	}
	if len(a.Username) == 0 {
		return nil, fmt.Errorf("No username specified for account %s", a.Name)
	}

	switch a.Auth {
	case "password":
		if len(a.Password) == 0 {
			return nil, fmt.Errorf("No password specified for account %s", a.Name)
		}
		c.password = a.Password
	case "oauth2":
		if len(a.Token) == 0 {
			return nil, fmt.Errorf("No token specified for account %s", a.Name)
		}
		c.token = a.Token
	default:
		return nil, fmt.Errorf("Authentication type %s not supported", a.Auth)
	}
	c.auth = a.Auth

	/*
		Find server specification by the given name
	*/
	s, ok := config.servers[a.Server]
	if !ok {
		return nil, fmt.Errorf("Server with name %s not defined", a.Server)
	}

	c.tls = s.TLS
	p := s.Port
	/*
		Use default IMAP(S) port if not defined by the user
	*/
	if p == 0 {
		if c.tls {
			p = 993
		} else {
			p = 143
		}
	}
	c.addr = fmt.Sprintf("%s:%d", s.Host, p)

	c.name = a.Name
	c.username = a.Username

	return &c, nil
}

/*
	Connect and login to the IMAP server
*/
func (c *Checker) Connect() error {
	if err := c.connect(); err != nil {
		return err
	}

	if err := c.login(); err != nil {
		return err
	}

	return nil
}

/*
	Try to connect to the IMAP server
*/
func (c *Checker) connect() (err error) {
	if c.tls {
		c.conn, err = client.DialTLS(c.addr, nil)
	} else {
		c.conn, err = client.Dial(c.addr)
	}
	return
}

/*
	Try to login to the IMAP server
*/
func (c *Checker) login() error {
	switch c.auth {
	case "password":
		return c.conn.Login(c.username, c.password)
	case "oauth2":
		o2r, err := refresher.New("microsoft", c.token, "", time.Now())
		if err != nil {
			return err
		}
		t, err := o2r.GetToken()
		if err != nil {
			return err
		}
		return c.conn.Authenticate(sasl.NewXoauth2Client(c.username, t))
	default:
		return fmt.Errorf("BUG bad authentication type")
	}
}

/*
	Disconnect from the IMAP server
*/
func (c *Checker) Disconnect() error {
	return c.conn.Logout()
}

/*
	Gather requested statistics from INBOX folder on the IMAP server
*/
func (c *Checker) GetStatus() (ret Status, err error) {
	f := []imap.StatusItem{imap.StatusMessages, imap.StatusUnseen}
	s, err := c.conn.Status("INBOX", f)
	if err != nil {
		return
	}
	ret.Total = s.Messages
	ret.Unseen = s.Unseen
	return
}
