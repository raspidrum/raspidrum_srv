# Детекция подключения MIDI по USB

## Вариант 1: usb -> ALSA

1. Определяем класс и подкласс устройства (если хотя бы один интерфейс имеют нужные класс и подкласс)
2. Определяем iProduct
3. Ищем ALSA устройство с совпадающим именем клиента

Проблемы:
- подключение двух устройств с одинаковым именем

## Вариант 2: usb -> udev -> ALSA

1. Определяем класс и подкласс устройства (если хотя бы один интерфейс имеют нужные класс и подкласс)
2. Определяем iProduct
3. В udev ищем устройство по соответствию ID_USB_SERIAL_SHORT = iSerial, если есть. Если нет, то по idVendor+idProduct. Если найдено несколько, то по номеру usb порта (TODO: как?)
4. Определяем номер карты (card3, например)
5. Ищем ALSA устройство с совпадающим именем клиента и номером карты

Проблемы:
- подключение двух одинаковых устройств без серийного номера можем не угадать. В таком случае нужно привязываться к usb-порту

##  Вариант 3: udev -> ALSA

1. Ищем ALSA устройство по совпадению номера card



# Определение типа устройства

**Внимание: usb устройство может одновременно быть и MIDI и Audio. Для каждого типа - отдельный интерфейс**

## По USB

- bInterfaceClass == 1 - sound (как для midi, так и для audio)

- bInterfaceSubClass == 3 - MIDI Streaming
- bInterfaceSubClass == 3 - Streaming

## В udev

### для controlCx устройства
- или найти соответствующее usb устройство и далее как в п. выше
- или ID_USB_INTERFACES
  `ID_USB_INTERFACES` является строкой и содержит несколько поледовательностей `:XXYYZZ:`, где:
    - XX — класс интерфейса (`bInterfaceClass`)
    - YY — подкласс интерфейса (`bInterfaceSubClass`)
    - ZZ — протокол интерфейса (`bInterfaceProtocol`)

  Примеры:
    - `ID_USB_INTERFACES=:010100:010300:`
    - `ID_USB_INTERFACES=:010100:010200:030000:`
  
  Правило:
  - Если ID_USB_INTERFACES содержит `:010300:`, то это MIDI
  - Если ID_USB_INTERFACES содержит `:010200:`, то это Audio

### для sound устройств

- DEVTYPE=pcm или - DEVNAME=/dev/snd/pcm* для audio (* - возможны любые символы и цифры)
- SUBSYSTEM=snd_seq или DEVNAME=/dev/midi* - MIDI (для midi DEVTYPE отсутствует)

# Определение ALSA MIDI seq 

Явно сопоставить данные из udev и usb с MIDI seq нет возможности, т.к. seq - это порядковый номер подключенного MIDI усройства в ALSA и этот номер отсутствует в udev и lsusb.

В каталоге /sys/class/sound/cardX содержится симлинк на физическое устройство в /sys/devices/..., который совпадает с DEVPATH, получаемым через udev:
```bash
$ ls "/sys/devices/platform/scb/fd500000.pcie/pci0000:00/0000:00:00.0/0000:01:00.0/usb1/1-1/1-1.3/1-1.3:1.0/sound/card3"
controlC3  device  id  number  pcmC3D0c  pcmC3D0p  power  subsystem  uevent

$ ls /sys/class/sound/card3
controlC3  device  id  number  pcmC3D0c  pcmC3D0p  power  subsystem  uevent

$ ls -l /sys/class/sound/card3
lrwxrwxrwx 1 root root 0 Jul 22 09:12 /sys/class/sound/card3 -> ../../devices/platform/scb/fd500000.pcie/pci0000:00/0000:00:00.0/0000:01:00.0/usb1/1-1/1-1.3/1-1.3:1.0/sound/card3
```

## USB

Совпадение имени клиента ALSA с именем модели USB, например:

```bash
$ aconnect -i -l
...
client 28: 'SAMSUNG_Android' [type=kernel,card=3]
    0 'SAMSUNG_Android MIDI 1'

$ lsusb -d 04e8:686c -v
...
 iProduct                2 SAMSUNG_Android
...
```

[ 3- 0]: raw midi
   3 - соответствует card3 (/proc/asound/card3/, в событиях udev .../sound/card3)
   0 - /proc/asound/card3/midi0, порт 0 в aconnect -i -l (0 'SAMSUNG_Android MIDI 1')
   в /proc/asound/card3/midi0 содержится "SAMSUNG_Android", что соответствует client.name

## UDEV

```bash
$ aconnect -i -l
...
client 28: 'SAMSUNG_Android' [type=kernel,card=3]
    0 'SAMSUNG_Android MIDI 1'

$ cat /proc/asound/devices 
...
  8: [ 3- 0]: raw midi
...
```

[ 3- 0]: raw midi
   3 - соответствует card3 (/proc/asound/card3/, в DEVPATH udev .../sound/card3)
   0 - /proc/asound/card3/midi0, порт 0 в aconnect -i -l (0 'SAMSUNG_Android MIDI 1')
   
Имя клиента "SAMSUNG_Android" соответствует ID_MODEL или ID_MODEL_ENC (скорее ID_MODEL_ENC) - только для устройства controlCx

Для устройства seq-midi-* (SUBSYSTEM=snd_seq) соответствует DEVPATH=/devices/platform/.../sound/card3/seq-midi-3-0, где 3-0 - это card3 port0


### Особенности устройств udev

 - controlC* - управляющее, всегда одно для одного физического устройства. Наиболее полный состав атрибутов, включая наименование, серийный номер и данные USB (vendor:product)
 - seq-midi-* - ALSA seq устройство. Есть только инфа о card и port. Соответствие с usb только по DEVPATH


# Сопоставление Bus:Port в lsusb и udev

Пример:

lsusb выводит такую инфу:
```
/:  Bus 02.Port 1: Dev 1, Class=root_hub, Driver=xhci_hcd/4p, 5000M
/:  Bus 01.Port 1: Dev 1, Class=root_hub, Driver=xhci_hcd/1p, 480M
    |__ Port 1: Dev 2, If 0, Class=Hub, Driver=hub/4p, 480M
        |__ Port 3: Dev 3, If 0, Class=Audio, Driver=snd-usb-audio, 480M
        |__ Port 3: Dev 3, If 1, Class=Audio, Driver=snd-usb-audio, 480M
```

- `Bus 01` соответствует корневому хабу на шине 1 (`/: Bus 01.Port 1: Dev 1, Class=root_hub`).
- `Port 1` на шине 1 подключен хаб (`Dev 2, Class=Hub`).
- `Port 3` на этом хабе подключено устройство `Dev 3`, которое имеет два интерфейса (`If 0` и `If 1`), оба с драйвером `snd-usb-audio`.

в udev:
```
ID_PATH=platform-fd500000.pcie-pci-0000:01:00.0-usb-0:1.3:1.0
```
- `platform-fd500000.pcie-pci-0000:01:00.0:` Это путь к PCI-контроллеру USB (в данном случае, xhci_hcd), к которому подключен USB-хаб или устройство.
- `usb-0:` Указывает на USB-контроллер (обычно корневой хаб, root_hub). Номер USB-шины (Bus).
Соответствует `Bus 01` в `lsusb`
- `1.3:` - путь в топологии USB, где:
  -  `1` — номер порта на корневом хабе (Bus 01, Port 1).
  -  `.3` — номер порта на следующем уровне (в данном случае, порт 3 на хабе, подключенном к порту 1 корневого хаба).

- `:1.0:` - номер конфигурации и интерфейса устройства
  - `1:` Конфигурация №1 (не отображается в `lsusb`).
  - `0:` Интерфейс №0 (поле `If 0` в `lsusb`).

**Итог**

Устройство с `ID_PATH=platform-fd500000.pcie-pci-0000:01:00.0-usb-0:1.3:1.0` в udev соответствует устройству в lsusb: `Bus 01, Port 1, Port 3, Dev 3, If 0` (интерфейс 0, драйвер `snd-usb-audio`).

Для стабильного именования создайте правило udev, например:
```
SUBSYSTEM=="sound", ATTRS{idVendor}=="VENDOR_ID", ATTRS{idProduct}=="PRODUCT_ID", ATTRS{devpath}=="1.3", SYMLINK+="my_audio_device"
```