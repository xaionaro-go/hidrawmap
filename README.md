This project just allows to map `/dev/hidrawX` events to keycode (which you can then bind to something using `xbindkeys`).

```
xaionaro@void:~/go/src/github.com/xaionaro-go/hidrawmap$ go run ./cmd/elgatopedal/ --log-level trace
{"level":"debug","ts":1704860304.7139666,"caller":"hidraw/hid_raw.go:59","msg":"waiting for a HID event"}
{"level":"debug","ts":1704860305.4448748,"caller":"hidraw/hid_raw.go:68","msg":"received a HID event 0100030000010000"}
{"level":"debug","ts":1704860305.4449453,"caller":"elgatopedal/elgato_pedal.go:61","msg":"sending KeyPress for keycode 198"}
{"level":"debug","ts":1704860305.4450274,"caller":"hidraw/hid_raw.go:59","msg":"waiting for a HID event"}
{"level":"debug","ts":1704860306.3538744,"caller":"hidraw/hid_raw.go:68","msg":"received a HID event 0100030000000000"}
{"level":"debug","ts":1704860306.3540092,"caller":"elgatopedal/elgato_pedal.go:70","msg":"sending KeyPress for keycode 197"}
{"level":"debug","ts":1704860306.354193,"caller":"hidraw/hid_raw.go:59","msg":"waiting for a HID event"}
```

Example of a `.xbindkeysrc`:
```
xaionaro@void:~$ cat .xbindkeysrc
"/home/xaionaro/bin/pedal-center-down.sh"
	c:206
"/home/xaionaro/bin/pedal-center-up.sh"
	c:205
```
