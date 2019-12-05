package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"./WindowsUI"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/winlabs/gowin32"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func usb() {
	// init COM
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecNotificationQuery", "SELECT * FROM Win32_VolumeChangeEvent")
	result := resultRaw.ToIDispatch()
	defer result.Release()

	//done := make(chan bool)
	//go func() {

	for {
		// item is a SWbemObject, but really a Win32_Process
		itemRaw, _ := oleutil.CallMethod(result, "NextEvent")
		item := itemRaw.ToIDispatch()
		defer item.Release()
		asString, _ := oleutil.GetProperty(item, "DriveName")
		//Нужно узнать ID активной локальной сессии
		sessionID, err := WindowsUI.WTSGetActiveConsoleSessionId()
		if err != nil {
			fmt.Println("Error WTSGetActiveConsoleSessionId: ", err)
		}
		//декодирование UTF в ANSI Windows
		header, err := utfToAnsi("ВАС ПОСЕТИЛА КИБЕР ПОЛИЦИЯ ИЦФА")
		if err != nil {
			fmt.Println("Error Encode Utf Header: ", err)
		}
		message, err := utfToAnsi("На диске " + asString.ToString() +
			"обнаружено запрещенное аниме!!!")
		if err != nil {
			fmt.Println("Error Encode Utf Message: ", err)
		}
		//Открываем Сообщение в интерактивной сессии
		ret := WindowsUI.WTSSendMessage(sessionID, header, message, WindowsUI.MB_YESNO, 5)
		fmt.Println(ret)
		wtsServ := gowin32.OpenWTSServer("localhost")
		defer wtsServ.Close()
		//Определение текущей активной сессии
		/*wtsSessions, err := wtsServ.EnumerateSessions()
		if err != nil {
			fmt.Println("Error EnumerateSessions: ", err)
		}
		for _, sessionInfo := range wtsSessions {
			if sessionInfo.State == gowin32.WTSConnectStateActive {
				fmt.Println("Active Session ID: ", sessionInfo.SessionID)
			}
		}*/
		//Разлогиниваем клиента
		if err := wtsServ.LogoffSession(sessionID, true); err != nil {
			fmt.Println("Error logoff", err)
		}

		/*
			//Показ формы
			var answer int
			answer = robotgo.ShowAlert("VAS POSETILA CYBER POLICY ICFA", "Na diske "+
				asString.ToString()+"obnarujeno zapreschennoe anime!")

			if err := exec.Command("cmd", "/C", "logoff").Run(); err != nil {
				fmt.Println("Failed to logoff: ", err)
			}
			if answer == 0 { //Ответ ок

			} else { //Ответ отмена

			}*/

	}
	//}()

	//<-done
}

//Смена кодировки с UTF-8 в ANSI
func utfToAnsi(str string) (string, error) {
	var windows1251 *charmap.Charmap = charmap.Windows1251
	bs := []byte(str)
	readerBs := bytes.NewReader(bs)
	readerWin := transform.NewReader(readerBs, windows1251.NewEncoder())
	bWin, err := ioutil.ReadAll(readerWin)
	if err != nil {
		return "", err
	}
	return string(bWin), nil
}
