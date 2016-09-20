# wb-bolid

Device 127 Relay-1 ON:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/127/1/on" -m "1"
```

Device 127 Relay-2 OFF:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/127/2/on" -m "0"
```
Success reply:
```
/devices/c2000-sp1/127/2/status/relay on
```
Bad reply if connection problem
```
/devices/c2000-sp1/127/2/status/relay none
```

<!-- Request status device 127, Relay-3
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/127/3/state" -m ""

```

reply:
```
/devices/c2000-sp1/127/1/status/relay on
/devices/c2000-sp1/127/2/status/relay off
/devices/c2000-sp1/127/3/status/relay off
/devices/c2000-sp1/127/4/status/relay off
```

reply from bad requests
```
/devices/c2000-sp1/127/587/status/relay none
/devices/c2000-sp1/0/1/status/relay none
``` -->

Request ADC input Voltage  (U1 = 1, U2 = 2)

for example request DeviceID=5,  Voltage Input=1:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/5/1/adc" -m ""
```
Reply:
```
/devices/c2000-sp1/5/1/status/adc 24.4
```

Subscribe to all devices:

```
root@wirenboard:~# mosquitto_sub -t '/devices/c2000-sp1/#' -v
```
