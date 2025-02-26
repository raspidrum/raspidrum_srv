package liblscp

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"
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

	si, err := c.GetServerInfo()
	if err != nil {
		return fmt.Errorf("failed get server info: %w", err)
	}
	slog.Info("connected to LinuxSampler", slog.String("ver:", si.Version))

	defer c.conn.Close()
	return nil
}

// Gets information about the LinuxSampler instance.
func (c *Client) GetServerInfo() (*ServerInfo, error) {
	rs, err := c.retrieveInfo("GET SERVER INFO", true)
	if err != nil {
		return nil, fmt.Errorf("failed lscp command: %w", err)
	}
	si, err := NewServerInfo(rs.MultiLineResult)
	if err != nil {
		return nil, fmt.Errorf("failed lscp command: %w", err)
	}
	return &si, nil
}

func (c *Client) retrieveInfo(lscpCmd string, isMultiResult bool) (ResultSet, error) {
	_, err := fmt.Fprintf(c.conn, lscpCmd+"\r\n")
	if err != nil {
		return ResultSet{}, err
	}

	ResultSet, err := c.getResultSet(isMultiResult)
	if err != nil {
		return ResultSet, err
	}
	return ResultSet, nil
}

func (c *Client) getResultSet(isMultiResult bool) (ResultSet, error) {
	rs := ResultSet{}
	ln, err := c.getLine()
	if err != nil {
		return rs, err
	}

	if f := strings.HasPrefix(ln, "ERR"); f {
		if err := ParseError(ln, &rs); err != nil {
			return rs, err
		}
		// it's error got from LinuxSampler
		return rs, &LscpError{rs.Code, rs.Message}
	}
	if f := strings.HasPrefix(ln, "WRN"); f {
		if err := ParseWarning(ln, &rs); err != nil {
			return rs, err
		}
		// it's warning got from LinuxSampler
		slog.Warn("LinuxSampler", slog.Int("code", rs.Code), slog.String("msg", rs.Message))
		return rs, nil
	}
	if f := strings.HasPrefix(ln, "OK"); f {
		if err := ParseOk(ln, &rs); err != nil {
			return rs, err
		}
		// it's empty OK result
		return rs, nil
	}

	// It's single line result
	if !isMultiResult {
		rs.Type = ResultType.Ok
		rs.Message = ln
		return rs, nil
	}

	// it's multuline result
	for ln != "." {
		rs.AddLine(ln)
		ln, err = c.getLine()
		if err != nil {
			return rs, err
		}
	}
	rs.Type = ResultType.Ok
	return rs, nil
}

func (c *Client) getLine() (string, error) {
	for {
		s, err := bufio.NewReader(c.conn).ReadString('\r')
		if err != nil {
			return "", err
		}
		if !strings.HasPrefix(s, "NOTIFY:") {
			return s, nil
		}
	}
}
