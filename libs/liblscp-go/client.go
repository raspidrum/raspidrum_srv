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

func (c *Client) Disconnect() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed disconnect from LinuxSampler: %w", err)
	}
	return nil
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

func (c *Client) retrieveIndex(lscpCmd string) (int, error) {
	rs, err := c.retrieveInfo(lscpCmd, false)
	if err != nil {
		return 0, err
	}
	return rs.Index, nil
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

func (c *Client) getIntegerList(lscpCmd string) ([]int, error) {
	rs, err := c.retrieveInfo(lscpCmd, false)
	if err != nil {
		return nil, fmt.Errorf("failed execute: %s : %w", lscpCmd, err)
	}
	return ParseIntList(rs.Message)
}

// LSCP common commands

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

// LSCP AUDIO commands

// Creates a new audio output device for the desired audio output system.
// adrv The desired audio output system
// paramList An optional list of driver specific parameters. <code>Parameter</code>
// Return the numerical ID of the newly created device.
func (c *Client) CreateAudioOutputDevice(adrv string, params ...Parameter[any]) (int, error) {
	cmd := "CREATE AUDIO_OUTPUT_DEVICE"
	plist := make([]string, len(params))
	for i, v := range params {
		plist[i] = fmt.Sprintf("%s=%s", v.Name, v.GetStringValue())
	}
	return c.retrieveIndex(fmt.Sprintf("%s %s", cmd, strings.Join(plist, " ")))
}

// Gets a list of all created audio output devices.
func (c *Client) GetAudioOutputDevices() ([]AudioOutputDevice, error) {
	ids, err := c.GetAudioOutputDeviceIDs()
	if err != nil {
		return nil, err
	}

	audDevs := make([]AudioOutputDevice, len(ids))
	for i, v := range ids {
		audDevs[i], err = c.GetAudioOutputDeviceInfo(v)
		if err != nil {
			return nil, err
		}
	}
	return audDevs, nil
}

// Gets a list of numerical IDs of all created audio output devices.
func (c *Client) GetAudioOutputDeviceIDs() ([]int, error) {
	return c.getIntegerList("LIST AUDIO_OUTPUT_DEVICES")
}

// Gets the current settings of a specific, already created audio output device.
// devId Specifies the numerical ID of the audio output device.
func (c *Client) GetAudioOutputDeviceInfo(devId int) (AudioOutputDevice, error) {
	cmd := "GET AUDIO_OUTPUT_DEVICE INFO"
	rs, err := c.retrieveInfo(cmd, true)
	if err != nil {
		return AudioOutputDevice{}, fmt.Errorf("failed lscp command: %s : %w", cmd, err)
	}
	aod, err := ParseAudioOutputDevice(devId, rs.MultiLineResult)
	if err != nil {
		return AudioOutputDevice{}, err
	}
	return aod, nil
}

// Alters a specific setting of an audio output channel.
// chn The audio channel number.
// prm A <code>Parameter</code> instance containing the name of the parameter
// and the new value for this parameter.
func (c *Client) SetAudioOutputChannelParameter(devId int, chn int, prm Parameter[any]) error {
	cmd := fmt.Sprintf("SET AUDIO_OUTPUT_CHANNEL_PARAMETER %d %d %s=%s", devId, chn, prm.Name, prm.GetStringValue())
	_, err := c.retrieveIndex(cmd)
	if err != nil {
		return err
	}
	return nil
}

// LSCP MIDI commands
// Gets all MIDI input drivers currently available for the LinuxSampler instance.
func (c *Client) GetMidiInputDriverNames() ([]string, error) {
	cmd := "LIST AVAILABLE_MIDI_INPUT_DRIVERS"
	rs, err := c.retrieveInfo(cmd, false)
	if err != nil {
		return nil, err
	}
	return strings.Split(rs.Message, ","), nil
}

// Creates a new MIDI input device.
// miDriver The desired MIDI input system.
// paramList An optional list of driver specific parameters. <code>Parameter</code>
// instances can be easily created using {@link ParameterFactory} factory.
// Return the numerical ID of the newly created device.
func (c *Client) CreateMidiInputDevice(miDriver string, params ...Parameter[any]) (int, error) {
	cmd := "CREATE MIDI_INPUT_DEVICE"
	plist := make([]string, len(params))
	for i, v := range params {
		plist[i] = fmt.Sprintf("%s=%s", v.Name, v.GetStringValue())
	}
	return c.retrieveIndex(fmt.Sprintf("%s %s", cmd, strings.Join(plist, " ")))
}

// Destroys already created MIDI input device.
// devId The numerical ID of the MIDI input device to be destroyed.
func (c *Client) DestroyMidiInputDevice(devId int) error {
	_, err := c.retrieveIndex(fmt.Sprintf("DESTROY MIDI_INPUT_DEVICE %d", devId))
	return err
}

// Gets a list of numerical IDs of all created MIDI input devices.
// An <code>Integer</code> array providing the numerical IDs of all created MIDI input devices.
func (c *Client) GetMidiInputDeviceIDs() ([]int, error) {
	return c.getIntegerList("LIST MIDI_INPUT_DEVICES")
}

// Gets detailed information about a specific MIDI input port.
// devId The numerical ID of the MIDI input device.
// midiPort The MIDI input port number.
// Return an <code>MidiPort</code> instance containing information about the specified MIDI input port.
func (c *Client) GetMidiInputPortInfo(devId int, midiPort int) (MidiPort, error) {
	cmd := "GET MIDI_INPUT_PORT INFO"
	rs, err := c.retrieveInfo(cmd, true)
	if err != nil {
		return MidiPort{}, fmt.Errorf("failed lscp command: %s : %w", cmd, err)
	}
	mp, err := ParseMidiPort(rs.MultiLineResult)
	if err != nil {
		return MidiPort{}, err
	}
	return mp, nil
}
