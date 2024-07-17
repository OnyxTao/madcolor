package htmlcolors

import (
	"fmt"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"sync"
)

type htmlColor struct {
	name   string
	hex    string
	bright int
}

var htmlColorArray []htmlColor

var setupLock sync.Mutex
var setup = false

const regExpHexString = "#[0-9a-fA-F]{6}"

var rxHex *regexp.Regexp

func HtmlColorsInitialize() {
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
	// if Debug
	// _, _ = fmt.Fprintf(os.Stderr, "invented hex color: %s\n", msg)
	//
	return msg
}

func rndbytes() (int, int, int, int) {
	a := int(rand.Int32N(0xFF + 1))
	b := int(rand.Int32N(0xFF + 1))
	c := int(rand.Int32N(0xFF + 1))
	sum := a + b + c
	return sum, a, b, c
}

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

// ColorNames maps names to hex value. A duplicate name
// will raise a compile-time error. Might be duplicate
// colors ... (same color, different names). Do we care?
// Is there a need for an inversion of this map?
// All names are lowercase.
var ColorNames = map[string]string{
	"aliceblue":            "#f0f8ff",
	"antiquewhite":         "#faebd7",
	"aqua":                 "#00ffff",
	"aquamarine":           "#7fffd4",
	"azure":                "#f0ffff",
	"beige":                "#f5f5dc",
	"bisque":               "#ffe4c4",
	"black":                "#000000",
	"blanchedalmond":       "#ffebcd",
	"blue":                 "#0000ff",
	"blueviolet":           "#8a2be2",
	"brown":                "#a52a2a",
	"burlywood":            "#deb887",
	"cadetblue":            "#5f9ea0",
	"chartreuse":           "#7fff00",
	"chocolate":            "#d2691e",
	"coral":                "#ff7f50",
	"cornflowerblue":       "#6495ed",
	"cornsilk":             "#fff8dc",
	"crimson":              "#dc143c",
	"cyan":                 "#00ffff",
	"darkblue":             "#00008b",
	"darkcyan":             "#008b8b",
	"darkgoldenrod":        "#b8860b",
	"darkgray":             "#a9a9a9",
	"darkgreen":            "#006400",
	"darkgrey":             "#a9a9a9",
	"darkkhaki":            "#bdb76b",
	"darkmagenta":          "#8b008b",
	"darkolivegreen":       "#556b2f",
	"darkorange":           "#ff8c00",
	"darkred":              "#8b0000",
	"darksalmon":           "#e9967a",
	"darkseagreen":         "#8fbc8f",
	"darkslateblue":        "#483d8b",
	"darkslategray":        "#2f4f4f",
	"darkslategrey":        "#2f4f4f",
	"darkturquoise":        "#00ced1",
	"darkviolet":           "#9400d3",
	"deeppink":             "#ff1493",
	"deepskyblue":          "#00bfff",
	"dimgray":              "#696969",
	"dimgrey":              "#696969",
	"dodgerblue":           "#1e90ff",
	"firebrick":            "#b22222",
	"floralwhite":          "#fffaf0",
	"forestgreen":          "#228b22",
	"fuchsia":              "#ff00ff",
	"gainsboro":            "#dcdcdc",
	"ghostwhite":           "#f8f8ff",
	"gold":                 "#ffd700",
	"goldenrod":            "#daa520",
	"gray":                 "#808080",
	"grey":                 "#808080",
	"green":                "#008000",
	"greenyellow":          "#adff2f",
	"honeydew":             "#f0fff0",
	"hotpink":              "#ff69b4",
	"indianred":            "#cd5c5c",
	"indigo":               "#4b0082",
	"ivory":                "#fffff0",
	"khaki":                "#f0e68c",
	"lavender":             "#e6e6fa",
	"lavenderblush":        "#fff0f5",
	"lawngreen":            "#7cfc00",
	"lemonchiffon":         "#fffacd",
	"lightblue":            "#add8e6",
	"lightcoral":           "#f08080",
	"lightcyan":            "#e0ffff",
	"lightgoldenrodyellow": "#fafad2",
	"lightgray":            "#d3d3d3",
	"lightgreen":           "#90ee90",
	"lightgrey":            "#d3d3d3",
	"lightpink":            "#ffb6c1",
	"lightsalmon":          "#ffa07a",
	"lightseagreen":        "#20b2aa",
	"lightskyblue":         "#87cefa",
	"lightslategray":       "#778899",
	"lightslategrey":       "#778899",
	"lightsteelblue":       "#b0c4de",
	"lightyellow":          "#ffffe0",
	"lime":                 "#00ff00",
	"limegreen":            "#32cd32",
	"linen":                "#faf0e6",
	"magenta":              "#ff00ff",
	"maroon":               "#800000",
	"mediumaquamarine":     "#66cdaa",
	"mediumblue":           "#0000cd",
	"mediumorchid":         "#ba55f3",
	"mediumpurple":         "#9370db",
	"mediumseagreen":       "#3cb371",
	"mediumslateblue":      "#7b68ee",
	"mediumspringgreen":    "#00fa9a",
	"mediumturquoise":      "#48d1cc",
	"mediumvioletred":      "#c71585",
	"midnightblue":         "#191970",
	"mintcream":            "#f5fffa",
	"mistyrose":            "#ffe4e1",
	"moccasin":             "#ffe4b5",
	"navajowhite":          "#ffdead",
	"navy":                 "#000080",
	"oldlace":              "#fdf5e6",
	"olive":                "#808000",
	"olivedrab":            "#6b8e23",
	"orange":               "#ffa500",
	"orangered":            "#ff4500",
	"orchid":               "#da70d6",
	"palegoldenrod":        "#eee8aa",
	"palegreen":            "#98fb98",
	"paleturquoise":        "#afeeee",
	"palevioletred":        "#db7093",
	"papayawhip":           "#ffefd5",
	"peachpuff":            "#ffdab9",
	"peru":                 "#cd853f",
	"pink":                 "#ffc0cb",
	"plum":                 "#dda0dd",
	"powderblue":           "#b0e0e6",
	"purple":               "#800080",
	"rebeccapurple":        "#663399",
	"red":                  "#ff0000",
	"rosybrown":            "#bc8f8f",
	"royalblue":            "#4169e1",
	"saddlebrown":          "#8b4513",
	"salmon":               "#fa8072",
	"sandybrown":           "#f4a460",
	"seagreen":             "#2e8b57",
	"seashell":             "#fff5ee",
	"sienna":               "#a0522d",
	"silver":               "#c0c0c0",
	"skyblue":              "#87ceeb",
	"slateblue":            "#6a5acd",
	"slategray":            "#708090",
	"slategrey":            "#708090",
	"snow":                 "#fffafa",
	"springgreen":          "#00ff7f",
	"steelblue":            "#4682b4",
	"tan":                  "#d2b48c",
	"teal":                 "#008080",
	"thistle":              "#d8bfd8",
	"tomato":               "#ff6347",
	// "turquoise":            "#40e0d0",
	"violet":      "#ee82ee",
	"wheat":       "#f5deb3",
	"white":       "#ffffff",
	"whitesmoke":  "#f5f5f5",
	"yellow":      "#ffff00",
	"yellowgreen": "#9acd32",
	/* Pantone Color-Of-Year */
	"cerulean":          "#9bb7d6",
	"fuschia rose":      "#c94476",
	"true red":          "#c02034",
	"aqua sky":          "#7ac5c5",
	"tigerlily":         "#e4583e",
	"blue turquoise":    "#4fb0ae",
	"sand dollar":       "#decdbf",
	"chili pepper":      "#9c1b31",
	"blue iris":         "#595ca1",
	"mimosa":            "#f0bf59",
	"pantone turquoise": "#41b6ab", // and duplicate detection proves it worth
	"honeysuckle":       "#da4f70",
	"tangerine tango":   "#f05442",
	"emerald":           "#169c78",
	"radiant orchid":    "#b565a7",
	"marsala":           "#955251",
	"rose quartz":       "#939597",
	"serenity":          "#8ca4cf",
	"greenery":          "#88b04b",
	"ultra violet":      "#5f4b8b",
	"living coral":      "#ff6f61",
	"classic blue":      "#0f4c81",
	"ultimate grey":     "#939597",
	"illuminating":      "#f5df4d",
	"very peri":         "#6667ab",
	"vivid magenta":     "#ba2649",
	"peach fuzz":        "#f87c56",
}
