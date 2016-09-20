# wb-bolid

Relay control
------------

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
/devices/c2000-sp1/127/1/status/relay on
/devices/c2000-sp1/127/2/status/relay off
```
Bad reply if connection problem
```
/devices/c2000-sp1/127/2/status/relay none
```

Activate the relay for a time specified in the settings
------------

Example Device 127 Relay-2:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/127/2/on" -m "3"
```
Success reply:
```
/devices/c2000-sp1/127/2/status/relay while
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

Request ADC input Voltage
-------------------------

(U1 = 1, U2 = 2,  maximum from two inputs=0)

for example request DeviceID=5,  Voltage Input=1:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/5/1/adc" -m ""
```
Reply:
```
/devices/c2000-sp1/5/1/status/adc 24.4
```

for example request DeviceID=5,  the maximum voltage at both inputs:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/5/0/adc" -m ""
```
Reply:
```
/devices/c2000-sp1/5/0/status/adc 24.4
```


Change DeviceID
-------------------------
(Address 1...127)

For example change address from 127 to 2
```
root@wirenboard:~#  mosquitto_pub -t "/devices/c2000-sp1/127/2/changeaddress" -m ""
```

Set default relay mode
-------------------------
(on, off, blink, pcn)
```
root@wirenboard:~#  mosquitto_pub -t "/devices/c2000-sp1/7/1/setrelaydefaultmode" -m "on"
root@wirenboard:~#  mosquitto_pub -t "/devices/c2000-sp1/7/1/setrelaydefaultmode" -m "off"
```

Set relay time
-------------------------
(time in seconds 1...60)
```
root@wirenboard:~#  mosquitto_pub -t "/devices/c2000-sp1/7/1/setrelaytime" -m "1"
root@wirenboard:~#  mosquitto_pub -t "/devices/c2000-sp1/7/2/setrelaytime" -m "60"
```


Subscribe
---------

Subscribe to all devices:

```
root@wirenboard:~# mosquitto_sub -t '/devices/c2000-sp1/#' -v
```
