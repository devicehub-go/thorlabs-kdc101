/*
Author: Leonardo Rossi Leao
Created at: September 26th, 2025
Last update: September 26th, 2025
*/

package thorlabskdc101_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	kdc101 "github.com/devicehub-go/thorlabs-kdc101"
	"github.com/devicehub-go/unicomm"
	"github.com/devicehub-go/unicomm/protocol/unicommserial"
)

func TestMoveMotor(t *testing.T) {
    controller := kdc101.New(
        kdc101.MTS25Z8,
        kdc101.Brushed,
		unicomm.UnicommOptions{
			Protocol: unicomm.Serial,
			Serial: unicommserial.SerialOptions{
				PortName:     "COM6",
				BaudRate:     115200,
				DataBits:     8,
				StopBits:     unicommserial.OneStopBit,
				Parity:       unicommserial.NoParity,
				ReadTimeout:  time.Second * 5,
				WriteTimeout: time.Second * 5,
			},
        },
    )
    
    if err := controller.Connect(); err != nil {
        log.Fatal("Failed to connect: ", err)
    }
    defer controller.Disconnect()
    
    if err := controller.Enable(1, true); err != nil {
        log.Fatal("Failed to enable motor: ", err)
    }
    
    fmt.Println("Homing stage...")
    if err := controller.StartHomeMove(1); err != nil {
        log.Fatal("Failed to start homing: ", err)
    }
    
    for {
        status, err := controller.GetDCStatusUpdate(1)
        if err != nil {
            log.Fatal("Failed to get status:", err)
        }
        statusSI := controller.DCStatusUpdateToSI(status)
        if statusSI.StatusBits.IsHomed {
            fmt.Println("Homing complete!")
            break
        }
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("Moving to 10mm...")
    if err := controller.MoveAbsolutePosition(1, 10.0); err != nil {
        log.Fatal("Failed to move:", err)
    }
}