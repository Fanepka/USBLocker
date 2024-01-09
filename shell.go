package main

import (
	"log"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const WINDOWS_TITLE string = "USBLocker"
const MESSAGE_BOX_BUTTONS_OK uint32 = 0
const MESSAGE_BOX_BUTTONS_YaN uint32 = 0x04

func main1() {

	var val uint64 = CheckLockedUsb()

	if val == 3 {
		retCode, _ := SendMessage("Заблокировать подключение USB устройств", WINDOWS_TITLE, MESSAGE_BOX_BUTTONS_YaN)

		if retCode == 6 {
			LockUsb()
		}
	} else {
		retCode, _ := SendMessage("Разблокировать подключение USB устройств", WINDOWS_TITLE, MESSAGE_BOX_BUTTONS_YaN)

		if retCode == 6 {
			UnLockUsb()
		}
	}

}

func LockUsb() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\USBSTOR`, registry.SET_VALUE)

	if err != nil {
		SendMessage(err.Error(), WINDOWS_TITLE, MESSAGE_BOX_BUTTONS_OK)
	}

	err = k.SetDWordValue("Start", 4)
	if err != nil {
		log.Print(err)

	} else {
		SendMessage("Подключение USB устройств заблокировано", WINDOWS_TITLE, MESSAGE_BOX_BUTTONS_OK)
	}

	defer k.Close()
}

func UnLockUsb() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\USBSTOR`, registry.SET_VALUE)

	if err != nil {
		SendMessage(err.Error(), WINDOWS_TITLE, MESSAGE_BOX_BUTTONS_OK)
	}

	err = k.SetDWordValue("Start", 3)
	if err != nil {
		log.Print(err)

	} else {
		SendMessage("Подключение USB устройств разблокировано", WINDOWS_TITLE, MESSAGE_BOX_BUTTONS_OK)
	}

	defer k.Close()
}

func CheckLockedUsb() uint64 {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\USBSTOR`, registry.QUERY_VALUE)

	if err != nil {
		SendMessage(err.Error(), WINDOWS_TITLE, MESSAGE_BOX_BUTTONS_OK)
	}

	defer k.Close()

	var s uint64
	s, _, err = k.GetIntegerValue("Start")

	return s

}

func SendMessage(text, caption string, boxtype uint32) (uint32, error) {
	retCode, err := windows.MessageBox(
		0,
		syscall.StringToUTF16Ptr(text),
		syscall.StringToUTF16Ptr(caption),
		boxtype,
	)

	return uint32(retCode), err
}
