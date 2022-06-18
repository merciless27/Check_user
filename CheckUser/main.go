package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"

	"0xFalcon/pkg"
	"0xFalcon/pkg/windows"
)

var (
	UsersManager      int = 0
	create_validateds     = []string{"create_business_validated", "create_validated"}
	creates               = []string{"create_business", "create"}
	Attempt           uint64
	Error             uint64
	Traked            uint64
	Held              uint64
	look              sync.Mutex
	wg                sync.WaitGroup
	uuid                  = pkg.Getuuid()
	switcher          int = 1
	ProxiesAsk        int
	urlProxy          string
	black             string
	blacklist         string
)
var (
	Proxies []string
	users   []string
	Threads int
)

var System = runtime.GOOS

func main() {

	rand.Seed(time.Now().UnixNano())
	windows.LoadKernelAndProc()
	windows.FreeKernelLib()
	windows.SetConsoleTitle("0x")
	windows.SetWindowSize("80", "20")
	pkg.YELLOW.Printf("_______            __________      ______                   \n")
	pkg.YELLOW.Printf("__  __ \\___  __    ___  ____/_____ ___  /__________________ \n")
	pkg.YELLOW.Printf("_  / / /_  |/_/    __  /_   _  __ `/_  /_  ___/  __ \\_  __ \\ \n")
	pkg.YELLOW.Printf("/ /_/ /__>  <      _  __/   / /_/ /_  / / /__ / /_/ /  / / / \n")
	pkg.YELLOW.Printf("\\____/ /_/|_|      /_/      \\__,_/ /_/  \\___/ \\____//_/ /_/ \n")
	pkg.RED.Printf("\t\t./Made By 0xRayan\n\n")
	windows.MaxStdio()
	users = pkg.LoadFile("Usernames", "list.txt")
	ProxiesAsk = pkg.Int(pkg.Input("1-Url | 2- File (Skip = File) : ", pkg.WHITE))
	if ProxiesAsk == 1 {

		urlProxy = pkg.Input("Enter The Url Proxy : ", pkg.WHITE)
		GrapProxyURL(urlProxy)
		Proxies = pkg.LoadFileUpdater("PROXYURL.txt")
	} else {
		Proxies = pkg.LoadFile("Proxies", "proxies.txt")
	}

	pkg.PPrint(pkg.WHITE, " + ", pkg.GREEN, "Threads : ", false)
	fmt.Scanln(&Threads)
	time.Sleep(2 * time.Second)

	if ProxiesAsk == 1 {
		go Updater()
	}
	pkg.ClearConsole()
	go PrintCounter()
	for i := 0; i <= Threads; i++ {
		wg.Add(1)
		go ScannerS()
	}
	wg.Wait()
}
func PrintCounter() {

Retry:
	windows.SetConsoleTitle(fmt.Sprintf("Att : %d - Error : %d - Traked : %d - Held : %d - List : %d \r", Attempt, Error, Traked, Held, len(users)))
	goto Retry

}
func GrapProxyURL(URL string) {
	resp, _ := http.Get(URL)
	body, _ := ioutil.ReadAll(resp.Body)
	pkg.CreateFile("PROXYURL.txt", string(body))
}

func ScannerS() {
while:
	defer wg.Done()
	if UsersManager >= len(users)-1 || UsersManager == len(users)-1 {
		UsersManager = 0
	} else if UsersManager < len(users) {
		UsersManager++
	}
	username := ""
	if len(users) > 0 {
		username = users[UsersManager]
	} else {
		switcher = 0
		look.Lock()
		pkg.PPrint(pkg.WHITE, " # ", pkg.GREEN, "Finsh Enter To Exit :)", true)
		windows.MsgBox("", "Finsh Enter To Exit :)")
		look.Unlock()
		os.Exit(0)

	}

	switch switcher {
	case 1:
		{
			Body := pkg.NetApiInsta("POST", "accounts", create_validateds[0], "username="+username+"", "ds_user_id="+pkg.RandomString(10)+"-"+pkg.RandomStringUpper(10)+"-"+pkg.RandomString(10), Proxies)
			if strings.Contains(Body, "username_held_by_others") {
				Held++
				go Login(username)
			} else if strings.Contains(Body, "spam") || strings.Contains(Body, "wait") {
				Error++
				switcher = 2
			} else {
				Attempt++
			}
			goto while
		}
	case 2:
		{
			Body := pkg.NetApiInsta("POST", "accounts", create_validateds[1], "username="+username+"", "ds_user_id="+pkg.RandomString(10)+"-"+pkg.RandomStringUpper(10)+"-"+pkg.RandomString(10), Proxies)
			if strings.Contains(Body, "username_held_by_others") {
				Held++
				go Login(username)
			} else if strings.Contains(Body, "spam") || strings.Contains(Body, "wait") {
				Error++
				switcher = 1
			} else {
				Attempt++
			}
			goto while

		}

	}
}

func Login(username string) {
	Data := fmt.Sprintf("guid=%s&device_id=%s&login_attempt_count=0&username=%s&password=0xRayan;)11qaadddd", uuid, uuid, username)
	Body := pkg.NetApiInsta("POST", "accounts", "login", Data, "", Proxies)
	if strings.Contains(Body, "Please wait a five minutes and try again") {
		return
	} else if strings.Contains(Body, "The username you entered doesn't appear to belong to an account. Please check your username and try again.") {
		if blacklist != username {
			blacklist = username
			pkg.PPrint(pkg.WHITE, " + ", pkg.GREEN, fmt.Sprintf("%s | "+color.HiGreenString("14Day")+"", username), true)
			Traked++
			pkg.CreateFile("14Day.txt", username+"\n")
			users = pkg.Remove(users, username)
		}
		return

	} else {
		if blacklist != username {
			blacklist = username
			pkg.PPrint(pkg.WHITE, " + ", pkg.GREEN, fmt.Sprintf("%s | "+color.HiRedString("Isn't 14Day")+"", username), true)
			users = pkg.Remove(users, username)
		}
		return
	}

}

func Updater() {
Update:
	GrapProxyURL(urlProxy)
	Proxies = pkg.LoadFileUpdater("PROXYURL.txt")
	time.Sleep(5 * time.Minute)
	goto Update
}
