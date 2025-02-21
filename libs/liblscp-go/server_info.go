package liblscp

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/hcl/strconv"
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
	var err error
	for _, v := range resultSet {
		ln, f := strings.CutPrefix(v, "DESCRIPTION: ")
		if f {
			s.Desc, err = strconv.Unquote(ln)
			if err != nil {
				return fmt.Errorf("can't parse index: '%s' %w", ln, err)
			}
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
