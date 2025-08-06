//go:build linux
// +build linux

package libalsa

//cgo CXXFLAGS: -D__LINUX_ALSA__
//cgo LDFLAGS: -lasound -L/usr/local/lib
//cgo CFLAGS: -I/usr/include/alsa

/*
#cgo LDFLAGS: -lasound
#include <stdint.h>
#include <stdlib.h>
#include <alsa/asoundlib.h>
#include <alsa/pcm.h>
#include <alsa/control.h>
#include <alsa/seq.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// GetCardInfo retrieves information about a specific card
func GetCardInfo(cardNum int) (*AlsaCard, error) {
	var card *C.snd_ctl_t
	//var info *C.snd_ctl_card_info_t

	// Form device name
	deviceName := C.CString(fmt.Sprintf("hw:%d", cardNum))
	defer C.free(unsafe.Pointer(deviceName))

	// Open card
	ret := C.snd_ctl_open(&card, deviceName, 0)
	if ret < 0 {
		return nil, fmt.Errorf("cannot open control for card %d: %s", cardNum, C.GoString(C.snd_strerror(ret)))
	}
	defer C.snd_ctl_close(card)

	// Allocate memory for info structure
	//C.snd_ctl_card_info_alloca(&info)

	size := C.snd_ctl_card_info_sizeof()
	infoBytes := make([]byte, size)
	// make C-pointer
	info := (*C.snd_ctl_card_info_t)(unsafe.Pointer(&infoBytes[0]))

	// Get card info
	ret = C.snd_ctl_card_info(card, info)
	if ret < 0 {
		return nil, fmt.Errorf("cannot get card info for card %d: %s", cardNum, C.GoString(C.snd_strerror(ret)))
	}

	// Extract information
	cardInfo := &AlsaCard{
		ID:       cardNum,
		Name:     C.GoString(C.snd_ctl_card_info_get_name(info)),
		LongName: C.GoString(C.snd_ctl_card_info_get_longname(info)),
		Driver:   C.GoString(C.snd_ctl_card_info_get_driver(info)),
		Mixer:    C.GoString(C.snd_ctl_card_info_get_mixername(info)),
	}

	return cardInfo, nil
}

// GetHardwareInfo retrieves additional hardware information
func GetHardwareInfo(cardNum int) (map[string]string, error) {
	var card *C.snd_ctl_t

	deviceName := C.CString(fmt.Sprintf("hw:%d", cardNum))
	defer C.free(unsafe.Pointer(deviceName))

	ret := C.snd_ctl_open(&card, deviceName, 0)
	if ret < 0 {
		return nil, fmt.Errorf("cannot open control for card %d", cardNum)
	}
	defer C.snd_ctl_close(card)

	hwInfo := make(map[string]string)

	// Get PCM device information
	//var pcmInfo *C.snd_pcm_info_t
	//C.snd_pcm_info_alloca(&pcmInfo)
	size := C.snd_pcm_info_sizeof()
	infoBytes := make([]byte, size)
	pcmInfo := (*C.snd_pcm_info_t)(unsafe.Pointer(&infoBytes[0]))

	// Check playback devices
	playbackDevices := 0
	for device := 0; device < 32; device++ {
		C.snd_pcm_info_set_device(pcmInfo, C.uint(device))
		C.snd_pcm_info_set_subdevice(pcmInfo, 0)
		C.snd_pcm_info_set_stream(pcmInfo, C.SND_PCM_STREAM_PLAYBACK)

		ret := C.snd_ctl_pcm_info(card, pcmInfo)
		if ret >= 0 {
			playbackDevices++
		}
	}

	// Check capture devices
	captureDevices := 0
	for device := 0; device < 32; device++ {
		C.snd_pcm_info_set_device(pcmInfo, C.uint(device))
		C.snd_pcm_info_set_subdevice(pcmInfo, 0)
		C.snd_pcm_info_set_stream(pcmInfo, C.SND_PCM_STREAM_CAPTURE)

		ret := C.snd_ctl_pcm_info(card, pcmInfo)
		if ret >= 0 {
			captureDevices++
		}
	}

	hwInfo["playback_devices"] = fmt.Sprintf("%d", playbackDevices)
	hwInfo["capture_devices"] = fmt.Sprintf("%d", captureDevices)

	return hwInfo, nil
}

// GetAllCards retrieves a list of all available cards
func GetAllCards() ([]int, error) {
	var cards []int
	cardNum := C.int(-1)

	for {
		ret := C.snd_card_next(&cardNum)
		if ret < 0 {
			return nil, fmt.Errorf("cannot enumerate cards: %s", C.GoString(C.snd_strerror(ret)))
		}

		if cardNum < 0 {
			break
		}

		cards = append(cards, int(cardNum))
	}

	return cards, nil
}

// ListMidiPorts returns a list of all ALSA sequencer MIDI ports with client, port, card and names
func ListMidiPorts() ([]MidiPortInfo, error) {
	var result []MidiPortInfo
	// Open sequencer
	var seq *C.snd_seq_t
	if C.snd_seq_open(&seq, C.CString("default"), C.SND_SEQ_OPEN_DUPLEX, 0) < 0 {
		return nil, fmt.Errorf("cannot open ALSA sequencer")
	}
	defer C.snd_seq_close(seq)

	// Allocate client info
	szClient := C.snd_seq_client_info_sizeof()
	clientInfoBytes := make([]byte, szClient)
	clientInfo := (*C.snd_seq_client_info_t)(unsafe.Pointer(&clientInfoBytes[0]))
	C.snd_seq_client_info_set_client(clientInfo, -1)

	for C.snd_seq_query_next_client(seq, clientInfo) >= 0 {
		clientID := int(C.snd_seq_client_info_get_client(clientInfo))
		clientName := C.GoString(C.snd_seq_client_info_get_name(clientInfo))
		cardID := int(C.snd_seq_client_info_get_card(clientInfo))

		// Allocate port info
		szPort := C.snd_seq_port_info_sizeof()
		portInfoBytes := make([]byte, szPort)
		portInfo := (*C.snd_seq_port_info_t)(unsafe.Pointer(&portInfoBytes[0]))
		C.snd_seq_port_info_set_client(portInfo, C.int(clientID))
		C.snd_seq_port_info_set_port(portInfo, -1)

		for C.snd_seq_query_next_port(seq, portInfo) >= 0 {
			portID := int(C.snd_seq_port_info_get_port(portInfo))
			portName := C.GoString(C.snd_seq_port_info_get_name(portInfo))
			cap := C.snd_seq_port_info_get_capability(portInfo)
			isInput := (cap & C.SND_SEQ_PORT_CAP_READ) != 0
			isOutput := (cap & C.SND_SEQ_PORT_CAP_WRITE) != 0

			// Determine port type based on capabilities
			portType := C.snd_seq_port_info_get_type(portInfo)
			var pType MidiPortType
			if (portType & C.SND_SEQ_PORT_TYPE_HARDWARE) != 0 {
				pType = MidiPortTypeHardware
			} else {
				pType = MidiPortTypeSoftware
			}

			result = append(result, MidiPortInfo{
				ClientId:   clientID,
				PortId:     portID,
				CardId:     cardID,
				ClientName: clientName,
				PortName:   portName,
				PortType:   pType,
				IsInput:    isInput,
				IsOutput:   isOutput,
			})
		}
	}
	return result, nil
}
