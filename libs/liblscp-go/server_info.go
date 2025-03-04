package liblscp

import (
	"fmt"
	"strings"
)

type ServerInfo struct {
	Version         string
	ProtocolVersion string
	Desc            string
}

func NewServerInfo(resultSet []string) (ServerInfo, error) {
	si := ServerInfo{}
	err := si.parseServerInfo(resultSet)
	if err != nil {
		return si, fmt.Errorf("can't server info: %w", err)
	}
	return si, nil
}

func (s *ServerInfo) parseServerInfo(resultSet []string) error {
	for _, v := range resultSet {
		ln, f := strings.CutPrefix(v, "DESCRIPTION: ")
		if f {
			s.Desc = ln
			continue
		}
		ln, f = strings.CutPrefix(v, "VERSION: ")
		if f {
			s.Version = ln
			continue
		}
		ln, f = strings.CutPrefix(v, "PROTOCOL_VERSION: ")
		if f {
			s.ProtocolVersion = ln
			continue
		}
	}
	return nil
}
