package windows

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/fatih/color"

	"0xFalcon/pkg"
)

var (
	setConsoleTitleWProc uintptr
	KernelLib            syscall.Handle
)

func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

// MessageBoxPlain of Win32 API.
func MsgBox(title, caption string) int {
	const (
		NULL  = 0
		MB_OK = 0
	)
	return MessageBox(NULL, caption, title, MB_OK)
}

func setMaxStdio() (int, error) {

	var err error
	var crt syscall.Handle
	crt, err = syscall.LoadLibrary("msvcrt.dll")
	if err != nil {
		crt, err = syscall.LoadLibrary("crtdll.dll")
		if err != nil {
			crt, err = syscall.LoadLibrary("crt.dll")
			if err != nil {
				return 0, err
			}
		}
	}
	defer syscall.FreeLibrary(crt)
	_getmaxstdioProc, err := syscall.GetProcAddress(crt, "_getmaxstdio")
	_setmaxstdioProc, err := syscall.GetProcAddress(crt, "_setmaxstdio")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(_getmaxstdioProc, 0, 0, 0, 0)
	maxstdio := int(ret)
	//fmt.Println(maxstdio)

	if maxstdio != 4294967295 {
		ret, _, _ := syscall.Syscall(_setmaxstdioProc, 18, 4294967295, 0, 0)
		if int(ret) == 4294967295 {
			return int(ret), nil
		}
	}
	return maxstdio, nil
}

func MaxStdio() {
	ret, err := setMaxStdio()
	if err != nil {
		panic(err)
	}
	if ret == 4294967295 {
		pkg.PPrint(pkg.WHITE, " + ", pkg.GREEN, "Successfully Change File Descriptors Limit To "+color.CyanString("(4294967295)")+"", true)

	} else {
		fmt.Println(ret)
		pkg.PPrint(pkg.RED, " ! ", pkg.RED, "Failed Fetched Files ", true)
	}

}

func FreeKernelLib() {
	syscall.FreeLibrary(KernelLib)
}

func LoadKernelAndProc() error {
	KernelLib, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return err
	}
	setConsoleTitleWProc, err = syscall.GetProcAddress(KernelLib, "SetConsoleTitleW")
	if err != nil {
		return err
	}
	return nil
}

func SetWindowSize(x string, y string) {
	cmd := exec.Command("cmd.exe", "/c", fmt.Sprintf("mode con: cols=%s lines=%s", x, y))
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func SetConsoleTitle(title string) (int, error) {
	r, _, err := syscall.Syscall(setConsoleTitleWProc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}
