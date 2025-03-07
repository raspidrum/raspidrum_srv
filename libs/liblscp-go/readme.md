
## Required LSCP commands

+ CREATE AUDIO_OUTPUT_DEVICE JACK

+ SET AUDIO_OUTPUT_CHANNEL_PARAMETER 0 0 JACK_BINDINGS='system:playback_1'

+ CREATE MIDI_INPUT_DEVICE ALSA

+ SET MIDI_INPUT_PORT_PARAMETER 0 0 ALSA_SEQ_BINDINGS='20:0'

+ ADD CHANNEL

+ LIST CHANNELS
  use: clear prev configuration (REMOVE CHANNEL) before loading new

+ REMOVE CHANNEL <sampler-channel>

+ RESET CHANNEL <sampler-channel>

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

+ GET FX_SEND INFO <sampler-channel> <fx-send-id>

+ SET FX_SEND AUDIO_OUTPUT_CHANNEL <sampler-chan> <fx-send-id> <audio-src> <audio-dst>

+ SET FX_SEND EFFECT <sampler-chan> <fx-send-id> <effect-chain> <chain-pos>

+ REMOVE FX_SEND EFFECT <sampler-chan> <fx-send-id>

+ SET FX_SEND LEVEL <sampler-chan> <fx-send-id> <volume>

+ LIST AVAILABLE_EFFECTS

+ GET EFFECT INFO <effect-index>

+ CREATE EFFECT_INSTANCE <effect-system> <module> <effect-name>

+ CREATE EFFECT_INSTANCE <effect-index>

+ DESTROY EFFECT_INSTANCE <effect-instance>

+ LIST EFFECT_INSTANCES

+ GET EFFECT_INSTANCE INFO <effect-instance>

+ GET EFFECT_INSTANCE_INPUT_CONTROL INFO <effect-instance> <input-control>

+ SET EFFECT_INSTANCE_INPUT_CONTROL VALUE <effect-instance> <input-control> <value>

+ LIST SEND_EFFECT_CHAINS <audio-device>

+ ADD SEND_EFFECT_CHAIN <audio-device>

+ REMOVE SEND_EFFECT_CHAIN <audio-device> <effect-chain>

+ GET SEND_EFFECT_CHAIN INFO <audio-device> <effect-chain>

+ APPEND SEND_EFFECT_CHAIN EFFECT <audio-device> <effect-chain> <effect-instance>

+ INSERT SEND_EFFECT_CHAIN EFFECT <audio-device> <effect-chain> <chain-pos> <effect-instance>

+ REMOVE SEND_EFFECT_CHAIN EFFECT <audio-device> <effect-chain> <chain-pos>

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

+ SEND CHANNEL MIDI_DATA <midi-msg> <sampler-chan> <arg1> <arg2>





# Примеры ответов LinuxSampler

```
GET CHANNEL INFO 0

ENGINE_NAME: SFZ
VOLUME: 1.000
AUDIO_OUTPUT_DEVICE: 0
AUDIO_OUTPUT_CHANNELS: 2
AUDIO_OUTPUT_ROUTING: 0,1
MIDI_INPUT_DEVICE: 0
MIDI_INPUT_PORT: 0
MIDI_INPUT_CHANNEL: ALL
INSTRUMENT_FILE: /Users/art/arth/\xd1\x83\xd0\xb4\xd0\xb0\xd1\x80\xd0\xbd\xd1\x8b\xd0\xb5\x20\xd0\xb8\x20\xd0\xb3\xd0\xb8\xd1\x82\xd0\xb0\xd1\x80\xd0\xb0/drum\x20sounds/SFZ/SamsSonor-wav/SamsSonor.sfz
INSTRUMENT_NR: 0
INSTRUMENT_NAME: SamsSonor
INSTRUMENT_STATUS: 100
MUTE: false
SOLO: false
MIDI_INSTRUMENT_MAP: NONE
```



# Сценарии настройки

1. Создание выходного аудиоустройства
   Можно создать несколько выходных аудиоустройств. Но достаточно создать только одно и все выводить через одно устройство.

   Выходное аудиоустройство создается командой, к-я возвращает номер созданного устройства:

    `CREATE AUDIO_OUTPUT_DEVICE <driver_name>`

   Опционально, может потребоваться выполнить привязку каналов созданного аудио-устройства семплера к каналам аудио-устройства ОС. Например, для JACK:

   `SET AUDIO_OUTPUT_CHANNEL_PARAMETER 0 0 JACK_BINDINGS='system:playback_1'`

2. Создание входного MIDI устройства
   Можно создать несколько входных устройств. Имеет смысл, если требуется подключить несколько устройств. В простейшем случае достаточно одного

   Входное MIDI-устройство создается командой, к-я возвращает номер созданного устройства:

   `CREATE MIDI_INPUT_DEVICE <driver_name>`

   Далее необходимо привязать порт созданного MIDI-устройства к выходному порту MIDI-устройства в ОС:

   `SET MIDI_INPUT_PORT_PARAMETER <device-id> <port> <key>=<value>`

   например:

   `SET MIDI_INPUT_PORT_PARAMETER 0 0 CORE_MIDI_BINDINGS='vmpk vmpk out'`

3. Создание каналов семплера
    Необходимо добавить канал, получив в ответ его номер:

    `ADD CHANNEL`

   Далее канал необходимо привязать к выходному аудиоустройству:

   `SET CHANNEL AUDIO_OUTPUT_DEVICE <sampler-channel> <audio-device-id>`

   И привязать к входному MIDI-устройству:

   `SET CHANNEL MIDI_INPUT_DEVICE <sampler-channel> <midi-device-id>`

   Далее настраиваем инструменты. Загружаем движок:

   `LOAD ENGINE sfz <sampler-channel>`

   И загружаем семплы инструментов:

   `LOAD INSTRUMENT '<filename>' <instr-index> <sampler-channel>`

   где <instr-index> - индекс инструмента в файле (файл может содержать много наборов семплов.) Но не совсем понятно, поддерживает ли LinuxSampler несколько инструментов в sfz файле (вроде да: https://bb.linuxsampler.org/viewtopic.php?t=705). Поэтому следует использовать один инструмент в одном файле и указывать instr-index = 0