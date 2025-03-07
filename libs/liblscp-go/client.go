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

func NewClient(host, port string, timeout string) Client {
	return Client{
		host:       host,
		port:       port,
		conTimeout: timeout,
	}
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
	slog.Info("connected to LinuxSampler", slog.String("ver", si.Version))

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
		return ResultSet{}, fmt.Errorf("failed lscp command: %s : %w", lscpCmd, err)
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
	rd := bufio.NewReader(c.conn)

	ln, err := c.getLine(rd)
	if err != nil {
		return rs, err
	}

	if f := strings.HasPrefix(ln, "ERR"); f {
		if err := parseError(ln, &rs); err != nil {
			return rs, err
		}
		// it's error got from LinuxSampler
		return rs, &LscpError{rs.Code, rs.Message}
	}
	if f := strings.HasPrefix(ln, "WRN"); f {
		if err := parseWarning(ln, &rs); err != nil {
			return rs, err
		}
		// it's warning got from LinuxSampler
		slog.Warn("LinuxSampler", slog.Int("code", rs.Code), slog.String("msg", rs.Message))
		return rs, nil
	}
	if f := strings.HasPrefix(ln, "OK"); f {
		if err := parseOk(ln, &rs); err != nil {
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
		ln, err = c.getLine(rd)
		if err != nil {
			return rs, err
		}
	}
	rs.Type = ResultType.Ok
	return rs, nil
}

func (c *Client) getLine(r *bufio.Reader) (string, error) {
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		if !strings.HasPrefix(s, "NOTIFY:") {
			return strings.TrimSuffix(s, "\r\n"), nil
		}
	}
}

func (c *Client) getIntegerList(lscpCmd string) ([]int, error) {
	rs, err := c.retrieveInfo(lscpCmd, false)
	if err != nil {
		return nil, fmt.Errorf("failed execute: %s : %w", lscpCmd, err)
	}
	return parseIntList(rs.Message)
}

// LSCP common commands

// Gets information about the LinuxSampler instance.
func (c *Client) GetServerInfo() (ServerInfo, error) {
	rs, err := c.retrieveInfo("GET SERVER INFO", true)
	if err != nil {
		return ServerInfo{}, err
	}
	return NewServerInfo(rs.MultiLineResult)
}

// AUDIO commands

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
		return AudioOutputDevice{}, err
	}
	return ParseAudioOutputDevice(devId, rs.MultiLineResult)
}

// Alters a specific setting of an audio output channel.
// chn The audio channel number.
// prm A <code>Parameter</code> instance containing the name of the parameter
// and the new value for this parameter.
func (c *Client) SetAudioOutputChannelParameter(devId int, chn int, prm Parameter[any]) error {
	cmd := fmt.Sprintf("SET AUDIO_OUTPUT_CHANNEL_PARAMETER %d %d %s=%s", devId, chn, prm.Name, prm.GetStringValue())
	_, err := c.retrieveIndex(cmd)
	return err
}

// MIDI commands

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
	return ParseMidiPort(rs.MultiLineResult)
}

// Alters a specific setting of a MIDI input port.
func (c *Client) SetMidiInputPortParameter(devId int, port int, prm Parameter[any]) error {
	cmd := fmt.Sprintf("SET MIDI_INPUT_PORT_PARAMETER %d %d %s=%s", devId, port, prm.Name, prm.GetStringValue())
	_, err := c.retrieveIndex(cmd)
	return err
}

// Sends a MIDI event to this sampler channel.
// samplerChn The sampler channel number.
// type The type of MIDI message to send.
func (c *Client) SendChannelMidiData(samplerChn int, eventType string, arg1, arg2 int) error {
	cmd := fmt.Sprintf("SEND CHANNEL MIDI_DATA %s %d %d %d", eventType, samplerChn, arg1, arg2)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Common commands

// Loads and assigns an instrument to a sampler channel. Notice that this function will
// return after the instrument is fully loaded and the channel is ready to be used.
// filename The name of the instrument file on the LinuxSampler instance's host system.
// instrIdx The index of the instrument in the instrument file.
// samplerChn The number of the sampler channel the instrument should be assigned to.
func (c *Client) LoadInstrument(filename string, instrIdx int, samplerChn int) error {
	cmd := fmt.Sprintf("LOAD INSTRUMENT '%s' %d %d", filename, instrIdx, samplerChn)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Loads a sampler engine to a specific sampler channel.
// engineName The name of the engine.
// amplerChn The number of the sampler channel the deployed engine should be assigned to.
func (c *Client) LoadSamplerEngine(engineName string, samplerChn int) error {
	cmd := fmt.Sprintf("LOAD ENGINE %s %d", engineName, samplerChn)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Gets a list with numerical IDs of all created sampler channels.
// return An <code>Integer</code> array providing the numerical IDs of all created sampler channels.
func (c *Client) GetSamplerChannelIDs() ([]int, error) {
	return c.getIntegerList("LIST CHANNELS")
}

// Adds a new sampler channel. This method will increment the sampler channel count by one
// and the new sampler channel will be appended to the end of the sampler channel list.
// return The number of the newly created sampler channel.
func (c *Client) AddSamplerChannel() (int, error) {
	return c.retrieveIndex("ADD CHANNEL")
}

// Removes the specified sampler channel.
// samplerChn The numerical ID of the sampler channel to be removed.
func (c *Client) RemoveSamplerChannel(samplerChn int) error {
	cmd := fmt.Sprintf("REMOVE CHANNEL %d", samplerChn)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Gets a list of all available engines' names.
// return <code>String</code> array with all available engines' names.
func (c *Client) GetEngineNames() ([]string, error) {
	rs, err := c.retrieveInfo("LIST AVAILABLE_ENGINES", false)
	if err != nil {
		return nil, err
	}
	ls, err := parseStringList(rs.Message, ",")
	if err != nil {
		return nil, err
	}
	return ls, nil
}

// Sets the audio output device on the specified sampler channel.
// samplerChn The sampler channel number.
// devId The numerical ID of the audio output device.
func (c *Client) SetChannelAudioOutputDevice(samplerChn int, devId int) error {
	cmd := fmt.Sprintf("SET CHANNEL AUDIO_OUTPUT_DEVICE %d %d", samplerChn, devId)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Sets the audio output channel on the specified sampler channel.
// samplerChn The sampler channel number.
// audioOut The sampler channel's audio output channel which should be rerouted.
// audioIn The audio channel of the selected audio output device where <code>audioOut</code> should be routed to.
func (c *Client) SetChannelAudioOutputChannel(samplerChn int, audioOut int, audioIn int) error {
	cmd := fmt.Sprintf("SET CHANNEL AUDIO_OUTPUT_CHANNEL %d %d %d", samplerChn, audioOut, audioIn)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Sets the MIDI input device on the specified sampler channel.
// samplerChn The sampler channel number.
// devId The numerical ID of the MIDI input device.
func (c *Client) SetChannelMidiInputDevice(samplerChn int, devId int) error {
	cmd := fmt.Sprintf("SET CHANNEL MIDI_INPUT_DEVICE %d %d", samplerChn, devId)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Sets the volume of the specified sampler channel.
// samplerChn The sampler channel number.
// volume The new volume value.
func (c *Client) SetChannelVolume(samplerChn int, volume float64) error {
	cmd := fmt.Sprintf("SET CHANNEL VOLUME %d %.2f", samplerChn, volume)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Mute/unmute the specified sampler channel.
// samplerChn The sampler channel number.
// mute If <code>true</code> the specified channel is muted, else the channel is unmuted.
func (c *Client) SetChannelMute(samplerChn int, mute bool) error {
	b := 0
	if mute {
		b = 1
	}
	cmd := fmt.Sprintf("SET CHANNEL MUTE %d %d", samplerChn, b)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Solo/unsolo the specified sampler channel.
// samplerChn The sampler channel number.
// solo <code>true</code> to solo the specified channel, <code>false</code> otherwise.
func (c *Client) SetChannelSolo(samplerChn int, solo bool) error {
	b := 0
	if solo {
		b = 1
	}
	cmd := fmt.Sprintf("SET CHANNEL SOLO %d %d", samplerChn, b)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Gets the global volume of the sampler.
func (c *Client) GetVolume() (float64, error) {
	rs, err := c.retrieveInfo("GET VOLUME", false)
	if err != nil {
		return 0, err
	}
	return parseFloat(rs.Message)
}

// Sets the global volume of the sampler.
func (c *Client) SetVolume(vol float64) error {
	_, err := c.retrieveIndex(fmt.Sprintf("SET VOLUME %.2f", vol))
	return err
}

// Resets the specified sampler channel.
// samplerChn The sampler channel number.
func (c *Client) ResetChannel(samplerChn int) error {
	_, err := c.retrieveIndex(fmt.Sprintf("RESET CHANNEL %d", samplerChn))
	return err
}

// Resets the whole sampler.
func (c *Client) ResetSampler() error {
	_, err := c.retrieveInfo("RESET", false)
	return err
}

// FX

// Creates an additional effect send on the specified sampler channel.
// channel The sampler channel, on which a new effect send should be added.
// midiCtrl Defines the MIDI controller, which will be able alter the effect send level.
// name The name of the effect send entity. The name does not have to be unique.
// return The unique ID of the newly created effect send entity.
func (c *Client) CreateFxSend(channel int, midiCtrl int, name string) (int, error) {
	var cmd string
	if name == "" {
		cmd = fmt.Sprintf("CREATE FX_SEND %d %d", channel, midiCtrl)
	} else {
		cmd = fmt.Sprintf("CREATE FX_SEND %d %d '%s'", channel, midiCtrl, name)
	}
	return c.retrieveIndex(cmd)
}

// Destroys the specified effect send on the specified sampler channel.
// channel The sampler channel, from which the specified effect send should be removed.
// fxSend The ID of the effect send that should be removed.
func (c *Client) DestroyFxSend(channel int, fxSend int) error {
	_, err := c.retrieveIndex(fmt.Sprintf("DESTROY FX_SEND %d %d", channel, fxSend))
	return err
}

// Gets a list of effect sends on the specified sampler channel.
// channel The sampler channel number.
// return An <code>Integer</code> array providing the numerical IDs of all effect sends on the specified sampler channel.
func (c *Client) GetFxSendIDs(channel int) ([]int, error) {
	return c.getIntegerList(fmt.Sprintf("LIST FX_SENDS %d", channel))
}

// Gets the current settings of the specified effect send entity.
// channel The sampler channel number.
// fxSend The numerical ID of the effect send entity.
// return <code>FxSend</code> instance containing the current settings of the specified effect send entity.
func (c *Client) GetFxSendInfo(channel int, fxSend int) (FxSend, error) {
	cmd := fmt.Sprintf("GET FX_SEND INFO %d %d", channel, fxSend)
	rs, err := c.retrieveInfo(cmd, true)
	if err != nil {
		return FxSend{}, err
	}
	return ParseFxSend(fxSend, rs.MultiLineResult)
}

// Gets a list of all created effect sends on the specified sampler channel.
// channel The sampler channel number.
// return A <code>FxSend</code> array providing all created effect sends on the specified sampler channel.
func (c *Client) GetFxSends(channel int) ([]FxSend, error) {
	fids, err := c.GetFxSendIDs(channel)
	if err != nil {
		return nil, err
	}
	fxs := make([]FxSend, len(fids))
	for i, v := range fids {
		fxs[i], err = c.GetFxSendInfo(channel, v)
		if err != nil {
			return nil, err
		}
	}
	return fxs, nil
}

// Sets the destination of an effect send's audio channel in the specified sampler channel.
// channel The sampler channel number.
// fxSend The numerical ID of the effect send entity to be rerouted.
// audioSrc The numerical ID of the effect send's audio output channel, which should be rerouted.
// audioDst The audio channel of the selected audio output device where <code>audioSrc</code> should be routed to.
func (c *Client) SetFxSendAudioOutputChannel(channel int, fxSend int, audioSrc int, audioDst int) error {
	cmd := fmt.Sprintf("SET FX_SEND AUDIO_OUTPUT_CHANNEL %d %d %d %d", channel, fxSend, audioSrc, audioDst)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Assign a destination effect to an effect send.
// channel The sampler channel number.
// fxSend The numerical ID of the effect send entity.
// fxChainId The numerical ID of the destination effect chain.
// chainPos The exact effect chain position in the effect chain which hosts the actual destination effect.
func (c *Client) SetFxSendEffect(channel int, fxSend int, chainId int, chainPos int) error {
	cmd := fmt.Sprintf("SET FX_SEND EFFECT %d %d %d %d", channel, fxSend, chainId, chainPos)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Removes destination effect from an effect send.
// channel The sampler channel number.
// fxSend The numerical ID of the effect send entity.
func (c *Client) RemoveFxSendEffect(channel int, fxSend int) error {
	cmd := fmt.Sprintf("REMOVE FX_SEND EFFECT %d %d", channel, fxSend)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Sets the current send level of the specified effect send entity in the specified sampler channel.
// channel The sampler channel number.
// fxSend The numerical ID of the effect send entity.
// volume The new volume value (a value smaller than 1.0 means attenuation, whereas a value greater than 1.0 means amplification).
func (c *Client) SetFxSendLevel(channel int, fxSend int, vol float64) error {
	cmd := fmt.Sprintf("SET FX_SEND LEVEL %d %d %.2f", channel, fxSend, vol)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Retrieves the list of available internal effects.
// Note that the set of available internal effects can change at runtime.
// return An <code>Integer</code> array providing the numerical IDs of all available internal effects.
func (c *Client) GetEffectIDs() ([]int, error) {
	return c.getIntegerList("LIST AVAILABLE_EFFECTS")
}

// Gets general informations about the specified effect.
// effect The numerical ID of the effect entity.
// return <code>Effect</code> instance containing general informations about the specified effect.
func (c *Client) GetEffectInfo(effect int) (Effect, error) {
	rs, err := c.retrieveInfo(fmt.Sprintf("GET EFFECT INFO %d", effect), true)
	if err != nil {
		return Effect{}, err
	}
	return ParseEffect(effect, rs.MultiLineResult)
}

// Gets the list of internal effects available to the sampler.
// Note that the set of available internal effects can change at runtime.
// return An <code>Effect</code> array providing the current list of internal effects.
func (c *Client) GetEffects() ([]Effect, error) {
	ids, err := c.GetEffectIDs()
	if err != nil {
		return nil, err
	}
	efs := make([]Effect, len(ids))
	for i, v := range ids {
		efs[i], err = c.GetEffectInfo(v)
		if err != nil {
			return nil, err
		}
	}
	return efs, nil
}

// Creates an instance of the desired effect.
// param id The unique ID of the effect.
// return The unique ID of the newly created effect instance.
func (c *Client) CreateEffectInstanceById(id int) (int, error) {
	return c.retrieveIndex(fmt.Sprintf("CREATE EFFECT_INSTANCE %d", id))
}

// Creates an instance of the desired effect.
// The unique ID of the newly created effect instance.
func (c *Client) CreateEffectInstance(effect Effect) (int, error) {
	return c.CreateEffectInstanceById(effect.Id)
}

// Creates an instance of the desired effect.
// return The unique ID of the newly created effect instance.
func (c *Client) CreateEffectInstanceByAttrs(system, module, name string) (int, error) {
	cmd := fmt.Sprintf("CREATE EFFECT_INSTANCE %s '%s' '%s'", system, module, name)
	return c.retrieveIndex(cmd)
}

// Destroys the specified unused effect instance.
// instanceId The numerical ID of the effect instance.
func (c *Client) DestroyEffectInstance(instId int) error {
	_, err := c.retrieveIndex(fmt.Sprintf("DESTROY EFFECT_INSTANCE %d", instId))
	return err
}

// Retrieves the current list of effect instances.
// return An <code>Integer</code> array providing the numerical IDs of all available effect instances.
func (c *Client) GetEffectInscanceIDs() ([]int, error) {
	return c.getIntegerList("LIST EFFECT_INSTANCES")
}

// Gets the current informations about the specified effect instance.
// id The numerical ID of the effect instance.
// return <code>EffectInstanceInfo</code> object containing the current informations about the specified effect instance.
func (c *Client) GetEffectInstanceInfo(id int) (EffectInstance, error) {
	rs, err := c.retrieveInfo(fmt.Sprintf("GET EFFECT_INSTANCE INFO %d", id), true)
	if err != nil {
		return EffectInstance{}, err
	}
	eff, err := ParseEffectInstance(id, rs.MultiLineResult)
	if err != nil {
		return eff, err
	}
	for i, _ := range eff.Params {
		prm, err := c.GetEffectInstanceParameterInfo(id, i)
		if err != nil {
			return eff, err
		}
		eff.Params = append(eff.Params, prm)
	}
	return eff, nil
}

// Gets information about the specified effect parameter.
// id The numerical ID of the effect instance.
// parameter The parameter index.
// return <code>EffectParameter</code> object containing information about the specified effect parameter.
// Note that only the following fields are used - description, value, rangeMin, rangeMax, possibilities and default.
func (c *Client) GetEffectInstanceParameterInfo(instId int, paramId int) (Parameter[float64], error) {
	cmd := fmt.Sprintf("GET EFFECT_INSTANCE_INPUT_CONTROL INFO %d %d", instId, paramId)
	rs, err := c.retrieveInfo(cmd, true)
	if err != nil {
		return Parameter[float64]{}, err
	}
	return ParseEffectParameter(rs.MultiLineResult)
}

// Gets the current list of effect instances.
// return An <code>EffectInstanceInfo</code> array providing the current list of effect instances.
func (c *Client) GetEffectInstances() ([]EffectInstance, error) {
	ids, err := c.GetEffectInscanceIDs()
	if err != nil {
		return nil, err
	}
	efs := make([]EffectInstance, len(ids))
	for i, v := range ids {
		efs[i], err = c.GetEffectInstanceInfo(v)
		if err != nil {
			return nil, err
		}
	}
	return efs, nil
}

// Alters the current value of an effect parameter.
// instanceId The numerical ID of the effect instance.
// prmIndex The index of the parameter to alter.
// value The new value for this parameter.
func (c *Client) SetEffectInstanceParameter(instId int, paramId int, val float64) error {
	cmd := fmt.Sprintf("SET EFFECT_INSTANCE_INPUT_CONTROL VALUE %d %d %.2f", instId, paramId, val)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Retrieves the current list of send effect chains on the specified audio output device.
// audioDeviceId The numerical ID of the audio output device.
// return An <code>Integer</code> array providing the numerical IDs of all send effect chains on the specified audio output device.
func (c *Client) GetSendEffectChainIDs(audioDevId int) ([]int, error) {
	return c.getIntegerList(fmt.Sprintf("LIST SEND_EFFECT_CHAINS %d", audioDevId))
}

// Gets the current information of a send effect chain.
// audioDeviceId The numerical ID of the audio output device.
// chainId The numerical ID of the send effect chain.
// return <code>EffectChainInfo</code> object containing the current informations about the specified effect chain.
func (c *Client) GetSendEffectChainInfo(audioDevId int, chainId int) (EffectChain, error) {
	cmd := fmt.Sprintf("GET SEND_EFFECT_CHAIN INFO %d %d", audioDevId, chainId)
	rs, err := c.retrieveInfo(cmd, true)
	if err != nil {
		return EffectChain{}, nil
	}
	return ParseEffectChain(chainId, rs.MultiLineResult, c)
}

// Gets the current list of send effect chains on the specified audio output device.
// audioDeviceId The numerical ID of the audio output device.
// return An <code>EffectInstanceInfo</code> array providing the current list of effect instances.
func (c *Client) GetSendEffectChains(audioDevId int) ([]EffectChain, error) {
	ids, err := c.GetSendEffectChainIDs(audioDevId)
	if err != nil {
		return nil, err
	}
	efs := make([]EffectChain, len(ids))
	for i, v := range ids {
		efs[i], err = c.GetSendEffectChainInfo(audioDevId, v)
		if err != nil {
			return efs, nil
		}
	}
	return efs, nil
}

// Adds a send effect chain to the specified audio output device.
// audioDeviceId The numerical ID of the audio output device.
// The numerical ID of the new send effect chain.
func (c *Client) AddSendEffectChain(audioDevId int) (int, error) {
	return c.retrieveIndex(fmt.Sprintf("ADD SEND_EFFECT_CHAIN %d", audioDevId))
}

// Removes a send effect chain from an audio output device.
// audioDeviceId The numerical ID of the audio output device.
// chainId The numerical ID of the send effect chain to remove.
func (c *Client) RemoveSendEffectChain(audioDevId int, chainId int) error {
	_, err := c.retrieveIndex(fmt.Sprintf("REMOVE SEND_EFFECT_CHAIN %d %d", audioDevId, chainId))
	return err
}

// Adds an unused effect instance to the end of a send effect chain.
// audioDeviceId The numerical ID of the audio output device.
// chainId The numerical ID of the send effect chain.
// fxInstanceId The numerical ID of the effect instance to add.
func (c *Client) AppendEffectInstance(audioDevId int, chainId int, fxInstId int) error {
	cmd := fmt.Sprintf("APPEND SEND_EFFECT_CHAIN EFFECT %d %d %d", audioDevId, chainId, fxInstId)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Adds an unused effect instance at a certain position of a send effect chain.
// audioDeviceId The numerical ID of the audio output device.
// chainId The numerical ID of the send effect chain.
// pos The exact position in the effect chain where the supplied effect shall be inserted to.
// fxInstanceId The numerical ID of the effect instance to insert.
func (c *Client) InsertEffectInstance(audioDevId, chainId, pos, fxInstId int) error {
	cmd := fmt.Sprintf("INSERT SEND_EFFECT_CHAIN EFFECT %d %d %d %d", audioDevId, chainId, pos, fxInstId)
	_, err := c.retrieveIndex(cmd)
	return err
}

// Removes an effect instance from a certain position of a send effect chain.
// audioDeviceId The numerical ID of the audio output device.
// chainId The numerical ID of the send effect chain.
// pos The exact position of the effect instance to be removed from the effect chain.
func (c *Client) RemoveEffectInstanceFromChain(audioDevId, chainId, pos int) error {
	cmd := fmt.Sprintf("REMOVE SEND_EFFECT_CHAIN EFFECT %d %d %d", audioDevId, chainId, pos)
	_, err := c.retrieveIndex(cmd)
	return err
}
