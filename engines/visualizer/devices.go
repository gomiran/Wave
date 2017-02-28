package visualizer

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/hkparker/Wave/models"
)

func updateKnownDevices(frame models.Wireless80211Frame) {
	if _, ok := Devices[frame.Address1]; !ok {
		registerNewDevice(frame.Address1)
	} else if _, ok := Devices[frame.Address2]; !ok {
		registerNewDevice(frame.Address2)
	} else if _, ok := Devices[frame.Address3]; !ok {
		registerNewDevice(frame.Address3)
	} else if _, ok := Devices[frame.Address4]; !ok {
		registerNewDevice(frame.Address4)
	}
}

func registerNewDevice(mac string) {
	DevicesMux.Lock()
	if broadcast(mac) {
		return
	}
	device := models.Device{
		MAC: mac,
	}
	Devices[mac] = device
	device.Save()
	visualizeNewDevice(device)
	DevicesMux.Unlock()
}

func broadcast(mac string) bool {
	if mac == "ff:ff:ff:ff:ff:ff" {
		return true
	}
	return false
}

func visualizeNewDevice(device models.Device) {
	//controllers.VisualPool <-
	//log.WithFields(log.Fields{
	//      "at": "visualizeNewDevice",
	//	"mac": device.MAC,
	//}).Info("new device observed")
}