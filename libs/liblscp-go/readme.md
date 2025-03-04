
## Required LSCP commands

+ CREATE AUDIO_OUTPUT_DEVICE JACK

+ SET AUDIO_OUTPUT_CHANNEL_PARAMETER 0 0 JACK_BINDINGS='system:playback_1'

+ CREATE MIDI_INPUT_DEVICE ALSA

+ SET MIDI_INPUT_PORT_PARAMETER 0 0 ALSA_SEQ_BINDINGS='20:0'

+ ADD CHANNEL

+ LIST CHANNELS
  use: clear prev configuration (REMOVE CHANNEL) before loading new

+ REMOVE CHANNEL <sampler-channel>

RESET CHANNEL <sampler-channel>

+ SET CHANNEL AUDIO_OUTPUT_DEVICE 0 0

+ SET CHANNEL MIDI_INPUT_DEVICE 0 0

+ LIST AVAILABLE_ENGINES

+ LOAD ENGINE sfz 0

+ LOAD INSTRUMENT '/home/drum/instruments/SamsSonor-wav/SamsSonor.sfz' 0 0

+ SET CHANNEL VOLUME 0 0.70

+ SET CHANNEL MUTE <sampler-channel> <mute>

+ SET CHANNEL SOLO <sampler-channel> <solo>

+ GET VOLUME

+ SET VOLUME <volume>

+ RESET


#### Fx

+ CREATE FX_SEND <sampler-channel> <midi-ctrl> [<name>]

+ LIST FX_SENDS <sampler-channel>

+ DESTROY FX_SEND <sampler-channel> <fx-send-id>

GET FX_SEND INFO <sampler-channel> <fx-send-id>

SET FX_SEND AUDIO_OUTPUT_CHANNEL <sampler-chan> <fx-send-id> <audio-src> <audio-dst>

SET FX_SEND EFFECT <sampler-chan> <fx-send-id> <effect-chain> <chain-pos>

REMOVE FX_SEND EFFECT <sampler-chan> <fx-send-id>

SET FX_SEND LEVEL <sampler-chan> <fx-send-id> <volume>

LIST AVAILABLE_EFFECTS

GET EFFECT INFO <effect-index>

CREATE EFFECT_INSTANCE <effect-system> <module> <effect-name>

CREATE EFFECT_INSTANCE <effect-index>

DESTROY EFFECT_INSTANCE <effect-instance>

LIST EFFECT_INSTANCES

GET EFFECT_INSTANCE INFO <effect-instance>

GET EFFECT_INSTANCE_INPUT_CONTROL INFO <effect-instance> <input-control>

SET EFFECT_INSTANCE_INPUT_CONTROL VALUE <effect-instance> <input-control> <value>

LIST SEND_EFFECT_CHAINS <audio-device>

ADD SEND_EFFECT_CHAIN <audio-device>

REMOVE SEND_EFFECT_CHAIN <audio-device> <effect-chain>

GET SEND_EFFECT_CHAIN INFO <audio-device> <effect-chain>

APPEND SEND_EFFECT_CHAIN EFFECT <audio-device> <effect-chain> <effect-instance>

INSERT SEND_EFFECT_CHAIN EFFECT <audio-device> <effect-chain> <chain-pos> <effect-instance>

REMOVE SEND_EFFECT_CHAIN EFFECT <audio-device> <effect-chain> <chain-pos>

## Optional commands

#### Common

LIST CHANNELS

GET CHANNEL INFO 0

GET SERVER INFO

#### AUDIO

GET AUDIO_OUTPUT_CHANNEL_PARAMETER INFO 0 0 JACK_BINDINGS

LIST AVAILABLE_AUDIO_OUTPUT_DRIVERS
  sample output: ALSA,JACK
  posible use: diagnostic

GET AUDIO_OUTPUT_DRIVER INFO <audio-output-driver>
  sample usage: GET AUDIO_OUTPUT_DRIVER INFO JACK
  posible use: diagnostic, get parameters list

GET AUDIO_OUTPUT_DEVICE INFO <dev-id>

GET AUDIO_OUTPUT_CHANNEL INFO <dev-id> <channel>


+ LIST AUDIO_OUTPUT_DEVICES
  posible use: diagnostic, reload config

DESTROY AUDIO_OUTPUT_DEVICE <device-id>
  posible use: reload config

+ SET CHANNEL AUDIO_OUTPUT_CHANNEL <sampler-chan> <audio-out> <audio-in>


#### MIDI

+ LIST AVAILABLE_MIDI_INPUT_DRIVERS
  posible use: configuring

GET MIDI_INPUT_DRIVER_PARAMETER INFO <midi> <param>

+ LIST MIDI_INPUT_DEVICES
  posible use: configuring

GET MIDI_INPUT_DEVICE INFO <midi-id>
  posible use: diag ACTIVE status

+ GET MIDI_INPUT_PORT INFO <midi-id> <port>

GET MIDI_INPUT_PORT_PARAMETER INFO 0 0 ALSA_SEQ_BINDINGS

DESTROY MIDI_INPUT_DEVICE <device-id>
  posible use: reconfig on USB MIDI changed

ADD CHANNEL MIDI_INPUT <sampler-channel> <midi-device-id> [<midi-input-port>]

LIST CHANNEL MIDI_INPUTS <sampler-channel>

REMOVE CHANNEL MIDI_INPUT <sampler-channel> [<midi-device-id> [<midi-input-port>]]

SEND CHANNEL MIDI_DATA <midi-msg> <sampler-chan> <arg1> <arg2>