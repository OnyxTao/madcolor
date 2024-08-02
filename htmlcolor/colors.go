package htmlcolors

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"

	"madcolor/misc"
)

// var modeDebug = false

const regExpHexB = "[\\da-fA-F]{2}"     // match 2-digit hex byte (only)
const regExpHex6 = "#?([\\da-fA-F]{6})" // match 6-digit hex value (submatch the hex)
const regExpHex3 = "#?([\\da-fA-F])([\\da-fA-F])([\\da-fA-F])"

// match 3-digit hex value, submatch each 4-byte digit

var rxHexB *regexp.Regexp
var rxHex6 *regexp.Regexp
var rxHex3 *regexp.Regexp

type htmlColor struct {
	name string
	hex  string
}

var htmlColorArray []htmlColor

// init is a function that is automatically called before the main function at the start of program execution.
// It initializes the regular expression patterns rxHexB, rxHex6, and rxHex3 with their corresponding regular expression strings.
// It also initializes the htmlColorArray with the ColorNames map values, randomly filling the array.
// This function does not return any value.
func init() {
	rxHexB = regexp.MustCompile(regExpHexB)
	rxHex6 = regexp.MustCompile(regExpHex6)
	rxHex3 = regexp.MustCompile(regExpHex3)
	htmlColorArray = make([]htmlColor, len(ColorNames), len(ColorNames))
	ix := 0
	// this fills the array randomly. Doesn't matter for our purposes.
	for key, val := range ColorNames {
		htmlColorArray[ix].name = key
		htmlColorArray[ix].hex = val
		ix++
	}
}

// StringToColor takes a string and converts it to a hexadecimal color value.
// It first checks if the setup has been done by calling the Initialize function.
// If the string matches a 6-digit hexadecimal pattern, it extracts the digits
// and returns the corresponding color value in the format "#RRGGBB". If the string
// matches a 3-digit hexadecimal pattern, it duplicates each digit and returns
// the color value in the format "#RRGGBB". If the string matches a color name in
// the ColorNames map, it retrieves the hexadecimal value from the map and returns
// it. If none of the above conditions are met, it returns the default color value
// "#888888" and false to indicate that the conversion was unsuccessful.
//
// The function relies on the rxHex6 and rxHex3 regular expression patterns for
// validating the hexadecimal strings. The function also converts the input string
// to lowercase before processing. The function uses a strings.Builder to efficiently
// build the resulting color value by appending characters.
//
// This function returns the hexadecimal color value as a string and a boolean flag
// indicating if the conversion was successful or not.
func StringToColor(s string) (hex string, ok bool) {
	var sb strings.Builder

	if rxHex6.MatchString(s) {
		zx := rxHex6.FindSubmatch([]byte(s))
		sb.WriteRune('#')
		sb.WriteString(string(zx[1]))
		return sb.String(), true
	}

	if rxHex3.MatchString(s) {
		zx := rxHex3.FindSubmatch([]byte(s))
		sb.WriteRune('#')
		sb.WriteByte(zx[1][0])
		sb.WriteByte(zx[1][0])
		sb.WriteByte(zx[2][0])
		sb.WriteByte(zx[2][0])
		sb.WriteByte(zx[3][0])
		sb.WriteByte(zx[3][0])
		return sb.String(), true
	}

	s = strings.ToLower(s)

	hex, ok = ColorNames[s]
	if ok {
		return hex, true
	}

	return "#888888", false
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

	return (0.2126 * RGB[0]) + (0.7152 * RGB[1]) + (0.0722 * RGB[2])
}

func InventColor(backColor string, minContrast int, minDistance int) (fg, bg string) {
	var contrast = float64(minContrast) / 100.0
	var distance = float64(float64(3*0xFF)*float64(minDistance)) / 100.0
	var cnt, dst float64
	var ix int
	var ok bool

	if misc.IsStringSet(&backColor) {
		bg, ok = StringToColor(backColor)
		if !ok {
			backColor = "white"
			bg = "#FFFFFF"
		}
	} else {
		bg = RandColor()
	}

	cnt, dst = 0.0, 0.0
	for ix < 500 && (cnt < contrast || dst < distance) {
		fg = RandColor()
		dst, cnt = ColorDistance(fg, bg)
		ix++
	}
	if ix >= 500 {
		fg = "#FFFFFF" // white
		bg = "#000000" // black
	}
	return fg, bg
}

func RandColor() (color string) {
	_, r, g, b := randColorBytes()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func randColorBytes() (sum, r, g, b int) {
	bits := make([]byte, 3)
	_, _ = rand.Read(bits)
	return int(bits[0] + bits[1] + bits[2]), int(bits[0]), int(bits[1]), int(bits[2])
}

func hexByteToInt(hex string) (val int) {

	// paranoia check. can remove in 2025.
	if !rxHexB.MatchString(hex) {
		panic("hex byte conversion failed for " + hex)
	}

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
	if !rxHex6.MatchString(hex) {
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
// Note that the maximum distance is 255 + 255 + 255 == 765

func ColorDistance(a string, b string) (dist float64, contrast float64) {
	aRed, aGreen, aBlue := getRGB(a)
	bRed, bGreen, bBlue := getRGB(b)

	r1 := relativeLuminance(uint8(aRed), uint8(aGreen), uint8(aBlue))
	r2 := relativeLuminance(uint8(bRed), uint8(bGreen), uint8(bBlue))
	if r2 < r1 {
		r1, r2 = r2, r1
	}

	contrast = 1.0 - ((r1 + 0.05) / (r2 + 0.05))

	dist = math.Sqrt(
		math.Pow(float64(aRed-bRed), 2.0) +
			math.Pow(float64(aGreen-bGreen), 2.0) +
			math.Pow(float64(aBlue-bBlue), 2.0))

	return dist, contrast

}

func RandomColor(bg string, contrast int, distance int) (name string, hex string) {
	var ok bool

	minContrast := float64(contrast) / 100
	minDistance := float64(3*0xFF) * (float64(distance) / 100)

	if misc.IsStringSet(&bg) {
		bg, ok = StringToColor(bg)
		if !ok {
			bg = "#FFFFFF" // white
		}
	}

	var arrayLen = big.NewInt(int64(len(ColorNames)))
	ixBig, _ := rand.Int(rand.Reader, arrayLen)
	ixStart := int(ixBig.Int64())
	ix := ixStart

	fg := htmlColorArray[ixStart].hex
	dst, cst := ColorDistance(fg, bg)

	for cst < minContrast && dst < minDistance {
		ix++
		if ix >= len(htmlColorArray) {
			ix = 0
		}
		if ixStart == ix {
			fg, _ = InventColor(bg, contrast, distance)
			return "", fg
		}
		fg = htmlColorArray[ix].hex
		dst, cst = ColorDistance(fg, bg)
	}
	return htmlColorArray[ix].name, htmlColorArray[ix].hex
}

func RandNamedColor() (ix int, name, hex string) {
	ixBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(ColorNames))))
	if nil != err {
		msg := fmt.Sprintf(
			"huh? Failed to generate a big.Int from %d (len of ColorNames array) because %s",
			len(ColorNames), err.Error())
		panic(msg)
	}
	ix = int(ixBig.Int64())
	return ix, htmlColorArray[ix].name, htmlColorArray[ix].hex
}
