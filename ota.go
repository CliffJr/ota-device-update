package main

import (
	"fmt"
	"time"
)

var crctab8 = [256]byte{ /* reversed, 8-bit, poly=0x07 */
	0x00, 0x91, 0xE3, 0x72, 0x07, 0x96, 0xE4, 0x75, 0x0E, 0x9F, 0xED, 0x7C, 0x09, 0x98, 0xEA, 0x7B,
	0x1C, 0x8D, 0xFF, 0x6E, 0x1B, 0x8A, 0xF8, 0x69, 0x12, 0x83, 0xF1, 0x60, 0x15, 0x84, 0xF6, 0x67,
	0x38, 0xA9, 0xDB, 0x4A, 0x3F, 0xAE, 0xDC, 0x4D, 0x36, 0xA7, 0xD5, 0x44, 0x31, 0xA0, 0xD2, 0x43,
	0x24, 0xB5, 0xC7, 0x56, 0x23, 0xB2, 0xC0, 0x51, 0x2A, 0xBB, 0xC9, 0x58, 0x2D, 0xBC, 0xCE, 0x5F,
	0x70, 0xE1, 0x93, 0x02, 0x77, 0xE6, 0x94, 0x05, 0x7E, 0xEF, 0x9D, 0x0C, 0x79, 0xE8, 0x9A, 0x0B,
	0x6C, 0xFD, 0x8F, 0x1E, 0x6B, 0xFA, 0x88, 0x19, 0x62, 0xF3, 0x81, 0x10, 0x65, 0xF4, 0x86, 0x17,
	0x48, 0xD9, 0xAB, 0x3A, 0x4F, 0xDE, 0xAC, 0x3D, 0x46, 0xD7, 0xA5, 0x34, 0x41, 0xD0, 0xA2, 0x33,
	0x54, 0xC5, 0xB7, 0x26, 0x53, 0xC2, 0xB0, 0x21, 0x5A, 0xCB, 0xB9, 0x28, 0x5D, 0xCC, 0xBE, 0x2F,
	0xE0, 0x71, 0x03, 0x92, 0xE7, 0x76, 0x04, 0x95, 0xEE, 0x7F, 0x0D, 0x9C, 0xE9, 0x78, 0x0A, 0x9B,
	0xFC, 0x6D, 0x1F, 0x8E, 0xFB, 0x6A, 0x18, 0x89, 0xF2, 0x63, 0x11, 0x80, 0xF5, 0x64, 0x16, 0x87,
	0xD8, 0x49, 0x3B, 0xAA, 0xDF, 0x4E, 0x3C, 0xAD, 0xD6, 0x47, 0x35, 0xA4, 0xD1, 0x40, 0x32, 0xA3,
	0xC4, 0x55, 0x27, 0xB6, 0xC3, 0x52, 0x20, 0xB1, 0xCA, 0x5B, 0x29, 0xB8, 0xCD, 0x5C, 0x2E, 0xBF,
	0x90, 0x01, 0x73, 0xE2, 0x97, 0x06, 0x74, 0xE5, 0x9E, 0x0F, 0x7D, 0xEC, 0x99, 0x08, 0x7A, 0xEB,
	0x8C, 0x1D, 0x6F, 0xFE, 0x8B, 0x1A, 0x68, 0xF9, 0x82, 0x13, 0x61, 0xF0, 0x85, 0x14, 0x66, 0xF7,
	0xA8, 0x39, 0x4B, 0xDA, 0xAF, 0x3E, 0x4C, 0xDD, 0xA6, 0x37, 0x45, 0xD4, 0xA1, 0x30, 0x42, 0xD3,
	0xB4, 0x25, 0x57, 0xC6, 0xB3, 0x22, 0x50, 0xC1, 0xBA, 0x2B, 0x59, 0xC8, 0xBD, 0x2C, 0x5E, 0xCF,
}

func CalcCrc(buffer []byte) byte {
	crc := byte(0xff)
	for _, val := range buffer {
		crc = crctab8[crc^val]
	}
	// ones complement
	return 0xff - crc
}

func InitCrc(crc *byte) {
	*crc = 0xff
}

func EndCrc(crc *byte) {
	*crc = 0xff - *crc
}

func CalcCrcByByte(crc *byte, curdata byte) {
	*crc = crctab8[*crc^curdata]
}

func main() {
	// Sample usage for the file
	data := []byte{0x01, 0x23, 0x45, 0x67, 0x89}
	calculatedCrc := CalcCrc(data)
	println("Calculated CRC:", calculatedCrc)

	// Define state machine constants
	const (
		StateStart = iota
		StateCheckOTA6
		StateCheckOTA5
		StateCheckOTA5Response
		StateCheckOTA1
		StateCheckOTA1Response
		StateCheckOTA2
		StateCheckOTA2Response
		StateCheckOTA2ResponseResumable
		StateResendOTA2
		StateCheckOTA3
		StateCheckOTA3Response
		StateCheckOTA4
		StateSubprocess2
		StateOTAUpgradeFailed
		StateOTAUpgradeCanceled
		StateDeviceResumable
		StateOTAUpgradeSuccess
	)

	// Set initial state,a
	state := StateStart

	// Simulate the process with a loop
	for {
		switch state {
		case StateStart:
			fmt.Println("Send OTA,6 command to Device")
			time.Sleep(5 * time.Second)
			state = StateCheckOTA6

		case StateCheckOTA6:
			fmt.Println("Query current firmware of the device")
			time.Sleep(5 * time.Second)
			if getFirmwareVersion() == "20230322" {
				fmt.Println("Firmware device update success")
				state = StateOTAUpgradeSuccess
			} else if getFirmwareVersion() != "20230322" {
				state = StateCheckOTA5
			}

		case StateCheckOTA5:
			fmt.Println("Send OTA,5 command to Device")
			time.Sleep(5 * time.Second)
			state = StateCheckOTA5Response

		case StateCheckOTA5Response:
			fmt.Println("Get OTA,5 command device response")
			// Simulate OTA,5 command response here
			// For simplicity, let's assume the response is always (DEVICEID,1,001,OTA,5,1,X)
			response := "(DEVICEID,1,001,OTA,5,1,X)"
			if response == "(DEVICEID,1,001,OTA,5,1,X)" {
				fmt.Println("Send OTA,2 command to Device")
				state = StateCheckOTA2
			} else if response == "(DEVICEID,1,001,OTA,5,0,1)" {
				fmt.Println("Send OTA,1 command to Device")
				state = StateCheckOTA1
			} else {
				fmt.Println("Resend OTA,5 command to Device")
				state = StateCheckOTA5
			}

		case StateCheckOTA1:
			fmt.Println("Send OTA,1 command to Device")
			time.Sleep(5 * time.Second)
			state = StateCheckOTA1Response

		case StateCheckOTA1Response:
			fmt.Println("OTA,1 command response")
			// Simulate OTA,1 command response here
			// For simplicity, let's assume the response is always (DEVICEID,1,001,OTA,1,0,20210619)
			response := "(DEVICEID,1,001,OTA,1,1,20210619)"
			if response == "(DEVICEID,1,001,OTA,1,1,20210619)" {
				fmt.Println("Send OTA,2 command to Device")
				time.Sleep(5 * time.Second)
				state = StateCheckOTA2
			} else if response == "(DEVICEID,1,001,OTA,1,0,20210619)" {
				fmt.Println("Firmware device update failed")
				time.Sleep(5 * time.Second)
				state = StateOTAUpgradeFailed
			} else {
				fmt.Println("Send OTA,1 command to Device")
				time.Sleep(5 * time.Second)
				state = StateCheckOTA1
			}

		case StateSubprocess2:
			fmt.Println("Send OTA,3 command to Device")
			time.Sleep(5 * time.Second)
			state = StateCheckOTA3

		case StateCheckOTA3:
			time.Sleep(5 * time.Second)
			state = StateCheckOTA3Response

		case StateCheckOTA3Response:
			fmt.Println("Send OTA,3 command to device response")
			response := "(DEVICEID,1,001,OTA,3,120230322)"

			if response == "(DEVICEID,1,001,OTA,3,120230322)" {
				fmt.Println("OTA upgrading Canceled")
				state = StateOTAUpgradeCanceled
			} else {
				fmt.Println("Resend OTA,3 command to Device")
				state = StateCheckOTA3
			}

		case StateCheckOTA2:
			time.Sleep(5 * time.Second)
			state = StateCheckOTA2Response

		case StateCheckOTA2Response:
			// Simulate OTA,2 command response here
			// For simplicity, let's assume the response is always (7000313309,1,001,OTA,2,0,310,311)
			response := "(7000313309,1,001,OTA,2,1,310,311)"
			if response == "(7000313309,1,001,OTA,2,1,1,2)" {
				fmt.Println("Responce OTA,2 command to device: (7000313309,1,001,OTA,2,1,1,2)")
				fmt.Println("Send OTA,2,2 the second OTA package")
				time.Sleep(5 * time.Second)
				state = StateResendOTA2
			} else if response == "(7000313309,1,001,OTA,2,0,1,1)" {
				fmt.Println("Responce OTA,2 command to device: (7000313309,1,001,OTA,2,0,1,1)")
				fmt.Println("Resend OTA,2,1 the first OTA package")
				time.Sleep(5 * time.Second)
				state = StateResendOTA2
			} else if response == "(7000313309,1,001,OTA,2,0,310,311)" {
				fmt.Println("Responce OTA,2 command to device: (7000313309,1,001,OTA,2,0,310,311)")
				fmt.Println("Resend OTA,2,310 OTA package")
				time.Sleep(5 * time.Second)
				state = StateResendOTA2
			} else if response == "(7000313309,1,001,OTA,2,1,310,311)" {
				fmt.Println("Responce OTA,2 command to device: (7000313309,1,001,OTA,2,1,310,311)")
				time.Sleep(5 * time.Second)
				state = StateCheckOTA2ResponseResumable
			}

		case StateCheckOTA2ResponseResumable:
			// checking 1 saving success && 311 > 310 OTA package to send
			if StateCheckOTA2ResponseResumable == StateCheckOTA2ResponseResumable {
				// Perform the necessary actions for Subprocess 1 here
				// Simulate a condition where Subprocess 1 is successful and there are more packages to send
				fmt.Println("Saving success && 311 > 310 total OTA package to send")
				time.Sleep(2 * time.Second)
				state = StateCheckOTA4
				fmt.Println("Wait device reboot...")
				time.Sleep(20 * time.Second)
				state = StateCheckOTA6
			} else {
				state = StateResendOTA2
			}

		case StateResendOTA2:
			fmt.Println("Resend OTA,2 packet")
			// Resend the OTA,2 packet here
			time.Sleep(5 * time.Second)
			state = StateCheckOTA2

		}

	}
}

func getFirmwareVersion() string {
	// Simulate querying the current firmware version of the device
	// For simplicity, return a fixed value here.
	return "020230322"
}
