package liblscp

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	host       string
	port       string
	conTimeout string
	conn       net.Conn
}

func (c *Client) Connect() error {
	t, err := time.ParseDuration(c.conTimeout)
	if err != nil {
		return fmt.Errorf("failed parse timeout duration: '%s' %w", c.conTimeout, err)
	}

	if c.conn != nil {
		c.conn.Close()
	}

	c.conn, err = net.DialTimeout("tcp", net.JoinHostPort(c.host, c.port), t)
	if err != nil {
		return fmt.Errorf("failed connect to: '%s:%s' %w", c.host, c.port, err)
	}

	// TODO: добавить получение версии сервера для контроля успешности соединения и логирования
	si, err := getServerInfo()

	defer c.conn.Close()
	return nil
}

// Gets information about the LinuxSampler instance.
func (c *Client) GetServerInfo() (*ServerInfo, error) {
	rs, err := retrieveInfo("GET SERVER INFO")
	if err != nil {
		return nil, fmt.Errorf("failed lscp command: %w", err)
	}
	si, err := NewServerInfo(rs.MultiLineResult)
	if err != nil {
		return nil, fmt.Errorf("failed lscp command: %w", err)
	}
	return &si, nil
}

func (c *Client) retrieveInfo(lscpCmd string) (ResultSet, error) {

	ResultSet := getMultiLineResultSet()
	return ResultSet, nil
}
