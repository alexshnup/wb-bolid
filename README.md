# wb-bolid

Device 127 Relay-1 ON:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/127/1/on" -m "1"
```

Device 127 Relay-2 OFF:
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/127/2/on" -m "0"
```


Request status device 127, Relay-3
```
root@wirenboard:~# mosquitto_pub -t "/devices/c2000-sp1/127/3/state" -m ""

```

reply:
```
/devices/c2000-sp1/127/1/status on
/devices/c2000-sp1/127/2/status off
/devices/c2000-sp1/127/3/status off
/devices/c2000-sp1/127/4/status off
```

Subscribe to all devices:

```
root@wirenboard:~# mosquitto_sub -t '/devices/c2000-sp1/#' -v
```
