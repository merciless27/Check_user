package pkg

import (
	"fmt"

	color "github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

var (
	Requests *fasthttp.Client
)

var (
	RED     = color.New(color.FgRed, color.Bold)
	BLUE    = color.New(color.FgBlue, color.Bold)
	YELLOW  = color.New(color.FgHiYellow, color.Bold)
	GREEN   = color.New(color.FgGreen, color.Bold)
	WHITE   = color.New(color.FgWhite, color.Bold)
	CYAN    = color.New(color.FgCyan, color.Bold)
	MAGENTA = color.New(color.FgMagenta, color.Bold)
)

func print(text string) string {
	s := text
	fmt.Print(s)
	return s

}

func Int(text string) int {
	strVar := text
	intValue := 0
	fmt.Sscan(strVar, &intValue)
	return intValue

}
func PPrint(COLORT *color.Color, mark string, COLORM *color.Color, text string, last bool) string {
	var Las string
	if last {
		Las = "\n"
	}
	s := text
	CYAN.Print("[")
	COLORM.Print(mark)
	CYAN.Print("] - ")
	COLORT.Print(s, Las)
	return s

}

func Input(output string, Color *color.Color) string {
	outPUT := output
	PPrint(Color, " ? ", RED, outPUT, false)
	fmt.Scanln(&outPUT)
	return outPUT

}
