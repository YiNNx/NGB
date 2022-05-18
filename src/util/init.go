package util

var hub = newHub()

func init() {
	initLog()
	initEs()
	
	go hub.run()
}
