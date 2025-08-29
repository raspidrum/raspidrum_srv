# Service Configuration

jackd must work in realtime mode and memory limitations need to be removed.

Since you cannot specify limits when launching a transient unit, you need to create a systemd drop-in file:

```bash
sudo mkdir -p /etc/systemd/system/jack.service.d/
sudo tee /etc/systemd/system/jack.service.d/10-limits.conf << EOF
[Service]
LimitRTPRIO=95
LimitMEMLOCK=infinity
EOF
```
 
and reload systemd:

```bash
sudo systemctl daemon-reload
```

**Configuration Verification**

Check service status:

```bash
$ systemctl status jack.service 
● jack.service - JACK Audio Server (hw:1,0)
     Loaded: loaded (/run/systemd/transient/jack.service; transient)
  Transient: yes
    Drop-In: /etc/systemd/system/jack.service.d
             └─10-limits.conf
     Active: active (running) since Fri 2025-08-29 19:03:28 BST; 6min ago
   Main PID: 1015 (jackd)
      Tasks: 4 (limit: 1582)
        CPU: 1.321s
     CGroup: /user.slice/jack.service
             └─1015 /usr/bin/jackd -t 2000 -R -P 95 -d alsa -d hw:1,0 -r 48000 -p 512 -n 2

Aug 29 19:03:28 raspidrum-aabf jackd[1015]: This is free software, and you are welcome to redistribute it
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: under certain conditions; see the file COPYING for details
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: JACK server starting in realtime mode with priority 95
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: self-connect-mode is "Don't restrict self connect requests"
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: creating alsa driver ... hw:1,0|hw:1,0|512|2|48000|0|0|nomon|swmeter|-|32bit
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: configuring for 48000Hz, period = 512 frames (10.7 ms), buffer = 2 periods
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: ALSA: final selected sample format for capture: 16bit little-endian
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: ALSA: use 2 periods for capture
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: ALSA: final selected sample format for playback: 16bit little-endian
Aug 29 19:03:28 raspidrum-aabf jackd[1015]: ALSA: use 2 periods for playback
```

1. There should be a line with the applied drop-in: `Drop-In: /etc/systemd/system/jack.service.d`
2. There should be no errors in the logs like `Cannot use real-time scheduling (RR/95) (1: Operation not permitted)` and `Cannot lock down 107350048 byte memory area (Cannot allocate memory)`
  

Check limits:

```bash
cat /proc/$(pgrep jackd)/limits
Limit                     Soft Limit           Hard Limit           Units     
Max cpu time              unlimited            unlimited            seconds   
Max file size             unlimited            unlimited            bytes     
Max data size             unlimited            unlimited            bytes     
Max stack size            8388608              unlimited            bytes     
Max core file size        0                    unlimited            bytes     
Max resident set          unlimited            unlimited            bytes     
Max processes             5275                 5275                 processes 
Max open files            1024                 524288               files     
Max locked memory         unlimited            unlimited            bytes     
Max address space         unlimited            unlimited            bytes     
Max file locks            unlimited            unlimited            locks     
Max pending signals       5275                 5275                 signals   
Max msgqueue size         819200               819200               bytes     
Max nice priority         0                    0                    
Max realtime priority     95                   95                   
Max realtime timeout      unlimited            unlimited            us    
```

1. The `Max realtime priority` line should show 95
2. All `* memory` lines should show `unlimited`
