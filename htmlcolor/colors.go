package htmlcolors

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"sync"
)

const regExpHexString = "#[0-9a-fA-F]{6}"

type htmlColor struct {
	name   string
	hex    string
	bright int
}

var htmlColorArray []htmlColor
var setupLock sync.Mutex
var setup = false

var rxHex *regexp.Regexp

// Initialize initializes the package by setting up the necessary variables and data.
// It acquires a lock to ensure exclusive access while setting up.
// If the setup has already been done, the function returns early.
// If setup is in progress, the new caller waits until setup is complete to exit the function.
// The function populates the htmlColorArray with color names, hexadecimal values,
// and darkness values obtained from the ColorNames map. It calculates the darkness
// value using the getDarkness function.
//
// If the array exceeds the length of ColorNames, the function panics with an error message.
// The function also compiles the regular expression pattern for a hexadecimal string
// and assigns it to rxHex variable.
//
// This function does not return any values.

func Initialize() {
	setupLock.Lock()
	defer setupLock.Unlock()
	if setup {
		return
	}
	setup = true
	rxHex = regexp.MustCompile(regExpHexString)
	htmlColorArray = make([]htmlColor, len(ColorNames), len(ColorNames))
	ix := 0
	for key, val := range ColorNames {
		htmlColorArray[ix].name = key
		htmlColorArray[ix].hex = val
		htmlColorArray[ix].bright = getDarkness(val)
		ix++
		if ix > len(ColorNames) {
			panic("array overflow in HtmlColorsSetup")
		}
	}
}

// relativeLuminance calculates the relative luminance of a given RGB color.
// It follows the formulas provided by the W3C specifications for calculating
// the contrast ratio between two colors.
// For more information, refer to:
// https://www.w3.org/TR/WCAG20/#relativeluminancedef
//
// Parameters:
// - `rgb`: An RGB color represented by three byte values in the range [0, 255].
//
// Returns:
//   - The relative luminance of the given RGB color as a floating-point value in the range [0, 1].
//     A larger value indicates a brighter color.
//
// Panic:
//   - If the `rgb` slice does not contain exactly three elements.
//     This is to ensure correct usage and prevent calculation errors.
//
// Example usage:
//
//	rl := relativeLuminance(255, 255, 255)
//	fmt.Println(rl) // Output: 1.0
func relativeLuminance(rgb ...uint8) (rl float64) {
	var RGB [3]float64
	// https://www.omnicalculator.com/other/contrast-ratio#how-do-i-calculate-the-color-contrast-ratio-between-two-colors
	if len(rgb) != 3 {
		msg := fmt.Sprintf("missized slice passed to sRGBrl (len != 3) len == %d", len(rgb))
		panic(msg)
	}
	for ix := 0; ix < 3; ix++ {
		var a float64
		color := float64(rgb[ix]) / 255.0
		if color <= 0.04045 {
			a = color / 12.92
		} else {
			t := (color + 0.55) / 1.055
			a = math.Pow(t, 2.4)
		}
		RGB[ix] = a
	}

	return ((0.2126 * RGB[0]) + (0.7152 * RGB[1]) + (0.0722 * RGB[2]))
}

func InventColor(minl int, maxl int) (color string) {
	diff := maxl - minl
	if diff <= 20 {
		return "#000000"
	}

	// I miss do-while ...
	sum, a, b, c := rndbytes()
	for sum > maxl && sum < minl {
		sum, a, b, c = rndbytes()
	}
	msg := fmt.Sprintf("#%02x%02x%02x", a, b, c)
	// if Debug { // no real debug mode for pckg
	// _, _ = fmt.Fprintf(os.Stderr, "invented hex color: %s\n", msg)
	// }
	return msg
}

func rndbytes() (int, int, int, int) {
	a := int(rand.Int32N(0xFF + 1))
	b := int(rand.Int32N(0xFF + 1))
	c := int(rand.Int32N(0xFF + 1))
	sum := a + b + c
	return sum, a, b, c
}

// AntiColor converts a hexadecimal color representation to its anti-color.
// It expects the color to be in the format "#RRGGBB".
// If the color format is invalid, it will panic with an error message.
// The anti-color is calculated by subtracting each RGB component from 0xFF and
// then converting it back to a hexadecimal value.
// The function returns the anti-color as a string in the format "#RRGGBB".
func AntiColor(hex string) (acolor string) {
	if !rxHex.MatchString(hex) {
		panic("invalid hex format: [" + hex + "] (don't do that!)")
	}
	a := 0xff - hexByteToInt(hex[1:3])
	b := 0xff - hexByteToInt(hex[3:5])
	c := 0xff - hexByteToInt(hex[5:7])
	return fmt.Sprintf("#%02x%02x%02x", a, b, c)
}

func getDarkness(hex string) (val int) {
	return hexByteToInt(hex[1:3]) +
		hexByteToInt(hex[3:5]) +
		hexByteToInt(hex[5:7])
}

func hexByteToInt(hex string) (val int) {
	i, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		msg := fmt.Sprintf("huh? could not convert %s into an int because %s", hex, err.Error())
		_, _ = fmt.Fprintln(os.Stderr, msg)
		panic(msg)
	}
	return int(i)
}

// getRGB converts a hexadecimal color representation to its RGB values.
// It expects the color to be in the format "#RRGGBB".
// If the color format is invalid, it will panic with an error message.
// The function returns three integers representing the red, green, and blue values of the color.
func getRGB(hex string) (r, g, b int) {
	if !rxHex.MatchString(hex) {
		panic("invalid hex format: [" + hex + "] (don't do that!)")
	}
	return hexByteToInt(hex[1:3]), hexByteToInt(hex[3:5]), hexByteToInt(hex[5:7])
}

// ColorDistance calculates the Euclidean distance between two colors represented as hexadecimal strings,
// and also calculates the contrast ratio between the colors based on their relative luminance.
// The distance is calculated using the RGB values of the colors.
// The relative luminance of the colors is calculated using the sRGB color space formula.
// The function returns the distance and contrast ratio as floating-point values.
// The colors are expected to be in the format "#RRGGBB".
func ColorDistance(a string, b string) (dist float64, contrast float64) {
	aRed, aGreen, aBlue := getRGB(a)
	bRed, bGreen, bBlue := getRGB(b)

	r1 := relativeLuminance(uint8(aRed), uint8(aGreen), uint8(aBlue))
	r2 := relativeLuminance(uint8(bRed), uint8(bGreen), uint8(bBlue))
	if r2 < r1 {
		r1, r2 = r2, r1
	}

	contrast = (r1 + 0.05) / (r2 + 0.05)

	dist = math.Sqrt(
		math.Pow(float64(aRed-bRed), 2.0) +
			math.Pow(float64(aGreen-bGreen), 2.0) +
			math.Pow(float64(aBlue-bBlue), 2.0))

	return dist, contrast

}

// RandomColor returns a random color name and its hexadecimal representation.
// It takes an optional parameter, `mb`, which specifies the maximum brightness
// value for the randomly generated colors. By default, it uses the maximum
// brightness value of 0xFF + 0xFF + 0xFF.
// If the maximum brightness value is less than or equal to 0, it returns the
// color name "black" and its hexadecimal representation "#000000".
// The function selects a random starting index within the range of available
// color names and iterates over the htmlColorArray to find a color with
// brightness less than the specified maximum. If all colors in the array have
// maximum brightness, it returns "black" and "#000000".
// The function then returns the selected color name and its hexadecimal
// representation as a tuple (name, hex).
func RandomColor(mb ...int) (name string, hex string) {
	var maxColorBrightness = 0xFF + 0xFF + 0xFF
	ixStart := int(rand.Int32N(int32(len(ColorNames))))
	if len(mb) > 0 {
		maxColorBrightness = mb[0]
	}
	if maxColorBrightness <= 0 {
		return "black", "#000000"
	}
	h := &htmlColorArray[ixStart]
	ix := ixStart
	for h.bright >= maxColorBrightness {
		ix = (ix + 1) % len(ColorNames)
		if ix == ixStart {
			return "black", "#000000"
		}
		h = &htmlColorArray[ix]
	}
	return h.name, h.hex
}
