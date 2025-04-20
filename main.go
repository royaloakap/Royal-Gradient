package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"royal-gradient-tool/func/version"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

type Color struct {
	R, G, B int
}

func checkVersion() {
	downloadURL, lastVersion, stableVersion, unavailableVersion, err := version.CheckVersion()
	if err != nil {
		log.Println("[\u001B[38;5;226mVERSION\u001B[38;5;230m] \u001B[38;5;196mError\u001B[38;5;230m checking version:", err)
		os.Exit(1)
	}

	log.Printf("[\u001B[38;5;226mVERSION\u001B[38;5;230m] \u001B[38;5;46mLast Version\u001B[38;5;230m: %s, \u001B[38;5;45mStable Version\u001B[38;5;230m: %s, \u001B[38;5;196mUnavailable Version\u001B[38;5;230m: %s\n", lastVersion, stableVersion, unavailableVersion)

	if version.CompareVersions(version.VersionControl, stableVersion) {
		log.Println("[\u001B[38;5;226mVERSION\u001B[38;5;230m] New version of \u001B[38;5;45mRoyal GRADIENT\u001B[38;5;230m found.")
		log.Printf("[\u001B[38;5;226mVERSION\u001B[38;5;230m] You are using \u001B[38;5;45m%s\u001B[38;5;230m. Download the \u001B[38;5;46mnew\u001B[38;5;230m version from (\u001B[38;5;226m%s\u001B[38;5;230m) or contact @\u001B[38;5;45mRoyaloakap\u001B[38;5;230m.\n", version.VersionControl, downloadURL)
		log.Println("[\u001B[38;5;226mVERSION\u001B[38;5;230m] \u001B[38;5;196mExiting\u001B[38;5;230m due to outdated version of \u001B[38;5;45mRoyal GRADIENT\u001B[38;5;230m.")
		os.Exit(1)
	} else {
		log.Println("[\u001B[38;5;226mVERSION\u001B[38;5;230m] You are using the available version of Royal GRADIENT. Version:", version.VersionControl)
	}
}

var availableFonts = []string{"basic", "shadow", "digital", "block", "big", "small", "banner", "doom", "rounded", "mini"}

// Helper function to validate font
func isFontValid(font string) bool {
	for _, f := range availableFonts {
		if f == font {
			return true
		}
	}
	return false
}

func ToRGB(h string) (c Color, err error) {
	if strings.HasPrefix(h, "#") {
		h = h[1:]
	}
	switch len(h) {
	case 6:
		_, err = fmt.Sscanf(h, "%02x%02x%02x", &c.R, &c.G, &c.B)
	default:
		err = fmt.Errorf("Invalid hex color: %s", h)
	}
	return
}

func Colorize(text string, r, g, b int) string {
	fg := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	return fg + text + "\x1b[0m"
}

func ApplyGradientToFile(filePath string, c1, c2, c3 Color, hasMiddle bool, overwrite bool) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	text := string(file)
	lines := strings.SplitAfter(text, "\n")
	var out []string
	numLines := len(lines)
	midPoint := numLines / 2

	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	variableRegex := regexp.MustCompile(`<<.*?>>`)

	for i, line := range lines {
		var coloredLine []string
		var r, g, b int

		if hasMiddle && i < midPoint {
			r = int(float64(c1.R)*(1-float64(i)/float64(midPoint))) + int(float64(c3.R)*float64(i)/float64(midPoint))
			g = int(float64(c1.G)*(1-float64(i)/float64(midPoint))) + int(float64(c3.G)*float64(i)/float64(midPoint))
			b = int(float64(c1.B)*(1-float64(i)/float64(midPoint))) + int(float64(c3.B)*float64(i)/float64(midPoint))
		} else if hasMiddle {
			r = int(float64(c3.R)*(1-float64(i-midPoint)/float64(numLines-midPoint))) + int(float64(c2.R)*float64(i-midPoint)/float64(numLines-midPoint))
			g = int(float64(c3.G)*(1-float64(i-midPoint)/float64(numLines-midPoint))) + int(float64(c2.G)*float64(i-midPoint)/float64(numLines-midPoint))
			b = int(float64(c3.B)*(1-float64(i-midPoint)/float64(numLines-midPoint))) + int(float64(c2.B)*float64(i-midPoint)/float64(numLines-midPoint))
		} else {
			r = int(float64(c1.R)*(1-float64(i)/float64(numLines))) + int(float64(c2.R)*float64(i)/float64(numLines))
			g = int(float64(c1.G)*(1-float64(i)/float64(numLines))) + int(float64(c2.G)*float64(i)/float64(numLines))
			b = int(float64(c1.B)*(1-float64(i)/float64(numLines))) + int(float64(c2.B)*float64(i)/float64(numLines))
		}

		segments := ansiRegex.Split(line, -1)
		ansiMatches := ansiRegex.FindAllString(line, -1)
		variableMatches := variableRegex.FindAllString(line, -1)

		ansiIndex, variableIndex := 0, 0
		for _, segment := range segments {
			if variableIndex < len(variableMatches) && strings.Contains(segment, variableMatches[variableIndex]) {
				coloredLine = append(coloredLine, variableMatches[variableIndex])
				variableIndex++
				continue
			}

			for _, char := range segment {
				coloredLine = append(coloredLine, Colorize(string(char), r, g, b))
			}

			if ansiIndex < len(ansiMatches) {
				coloredLine = append(coloredLine, ansiMatches[ansiIndex])
				ansiIndex++
			}
		}

		out = append(out, strings.Join(coloredLine, ""))
	}

	result := strings.Join(out, "")

	if overwrite {
		err := ioutil.WriteFile(filePath, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			return
		}
		fmt.Println("Gradient applied successfully to file:", filePath)
	} else {
		newFilePath := filePath + ".gradient"
		err := ioutil.WriteFile(newFilePath, []byte(result), 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			return
		}
		fmt.Println("Gradient applied successfully to file. New file created:", newFilePath)
	}
}

func ApplyGradientToDir(dirPath string, c1, c2, c3 Color, hasMiddle bool, overwrite bool) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())
			ApplyGradientToFile(filePath, c1, c2, c3, hasMiddle, overwrite)
		}
	}
}

func CreateCenteredAsciiText(text, font string, width, height int, c1, c2 Color, backgroundChar string) string {
	if !isFontValid(font) {
		fmt.Printf("Invalid font specified. Available fonts: %v\n", availableFonts)
		return ""
	}

	asciiArt := figure.NewFigure(text, font, true).String()
	lines := strings.Split(asciiArt, "\n")

	for i := range lines {
		if len(lines[i]) > width {
			lines[i] = lines[i][:width]
		} else {
			lines[i] = strings.Repeat(backgroundChar, (width-len(lines[i]))/2) + lines[i] + strings.Repeat(backgroundChar, (width-len(lines[i]))/2)
		}
	}

	var output strings.Builder

	for i := 0; i < (height-len(lines))/2; i++ {
		output.WriteString(strings.Repeat(backgroundChar, width) + "\n")
	}

	for _, line := range lines {
		output.WriteString(line + "\n")
	}

	for i := 0; i < (height-len(lines))/2; i++ {
		output.WriteString(strings.Repeat(backgroundChar, width) + "\n")
	}

	return output.String()
}

func PrintCredits() {
	fmt.Println("Credits:")
	fmt.Println("Developed by: [Royaloakap]")
	fmt.Println("If you enjoyed this project, you can contact me on Telegram: @Royaloakap")
	fmt.Println("Feel free to contribute or suggest improvements!")
}

func main() {
	checkVersion()
	color1 := flag.String("start", "", "Starting color in hex (e.g., #FF0000)")
	color2 := flag.String("end", "", "Ending color in hex (e.g., #0000FF)")
	color3 := flag.String("middle", "", "Middle color in hex (optional)")
	dir := flag.String("dir", "", "Directory containing files for gradient application")
	file := flag.String("file", "", "File to apply gradient to")
	width := flag.Int("width", 50, "Width of the ASCII output")
	height := flag.Int("height", 20, "Height of the ASCII output")
	text := flag.String("text", "", "Text to display in ASCII with gradient")
	font := flag.String("font", "basic", "Font to use for ASCII Art")
	backgroundChar := flag.String("char", "·", "Character for background")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing files")
	help := flag.Bool("help", false, "Display help")
	credits := flag.Bool("credits", false, "Display credits")
	flag.Parse()

	if *help {
		fmt.Println("Usage:")
		fmt.Println("  -start <color1> : Starting color in hex (e.g., #FF0000)")
		fmt.Println("  -end <color2>   : Ending color in hex (e.g., #0000FF)")
		fmt.Println("  -middle <color3>: Middle color in hex (optional)")
		fmt.Println("  -file <file>    : Text file to apply gradient to")
		fmt.Println("  -dir <directory>: Directory containing files for gradient application")
		fmt.Println("  -width <int>    : Width of ASCII output (used with -text)")
		fmt.Println("  -height <int>   : Height of ASCII output (used with -text)")
		fmt.Println("  -text <string>  : Text to display in ASCII with gradient")
		fmt.Println("  -font <string>  : Font to use for ASCII Art (available fonts: basic, shadow, digital, block, big, small, banner, doom, rounded, mini)")
		fmt.Println("  -char <string>  : Character for background")
		fmt.Println("  -overwrite      : Overwrite existing files")
		fmt.Println("  -credits        : Shows project credits")
		return
	}

	if *credits {
		PrintCredits()
		return
	}

	c1, err := ToRGB(*color1)
	if err != nil {
		fmt.Println("Error with starting color:", err)
		return
	}

	c2, err := ToRGB(*color2)
	if err != nil {
		fmt.Println("Error with ending color:", err)
		return
	}

	hasMiddle := false
	var c3 Color
	if *color3 != "" {
		c3, err = ToRGB(*color3)
		if err != nil {
			fmt.Println("Error with middle color:", err)
			return
		}
		hasMiddle = true
	}

	if *text != "" {
		asciiArt := CreateCenteredAsciiText(*text, *font, *width, *height, c1, c2, *backgroundChar)
		if asciiArt == "" {
			return
		}
		fmt.Println(asciiArt)
		fileName := *text + ".txt"

		if _, err := os.Stat(fileName); err == nil && !*overwrite {
			fmt.Printf("File '%s' already exists. Use -overwrite to replace it.\n", fileName)
			return
		}

		err := ioutil.WriteFile(fileName, []byte(asciiArt), 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			return
		}
		fmt.Printf("ASCII art file generated: %s\n", fileName)
	} else if *file != "" {
		ApplyGradientToFile(*file, c1, c2, c3, hasMiddle, *overwrite)
	} else if *dir != "" {
		ApplyGradientToDir(*dir, c1, c2, c3, hasMiddle, *overwrite)
	} else {
		fmt.Println("Please specify a file, directory, or text to display.")
	}
}
