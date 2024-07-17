package htmlcolors

import (
	"fmt"
	"math/rand"
	"os"
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

func HtmlColorsInitialize() {
	setupLock.Lock()
	defer setupLock.Unlock()
	if setup {
		return
	}
	setup = true
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

func getDarkness(hex string) (val int) {
	return hexByteToInt(hex[1:3]) +
		hexByteToInt(hex[3:5]) +
		hexByteToInt(hex[5:7])
}

func hexByteToInt(hex string) (val int) {
	i, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		msg := fmt.Sprintf("huh? could not convert %s into an int because %s", hex, err.Error())
		fmt.Fprintln(os.Stderr, msg)
		panic(msg)
	}
	return int(i)
}

func RandomColor(mb ...int) (name string, hex string) {
	var max = 0xFF + 0xFF + 0xFF
	ixStart := rand.Intn(len(ColorNames))
	if len(mb) > 0 {
		max = mb[0]
	}
	if max <= 0 {
		return "black", "#000000"
	}
	h := &htmlColorArray[ixStart]
	ix := ixStart
	for h.bright >= max {
		ix = (ix + 1) % len(ColorNames)
		if ix == ixStart {
			return "black", "#000000"
		}
		h = &htmlColorArray[ix]
	}
	return h.name, h.hex
}

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
	"pantone turquoise": "#41b6ab",
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

// colornames
const AliceBlue = "#F0F8FF"
const AntiqueWhite = "#FAEBD7"
const Aqua = "#00FFFF"
const Aquamarine = "#7FFFD4"
const Azure = "#F0FFFF"
const Beige = "#F5F5DC"
const Bisque = "#FFE4C4"
const Black = "#000000"
const BlanchedAlmond = "#FFEBCD"
const Blue = "#0000FF"
const BlueViolet = "#8A2BE2"
const Brown = "#A52A2A"
const BurlyWood = "#DEB887"
const CadetBlue = "#5F9EA0"
const Chartreuse = "#7FFF00"
const Chocolate = "#D2691E"
const Coral = "#FF7F50"
const CornflowerBlue = "#6495ED"
const Cornsilk = "#FFF8DC"
const Crimson = "#DC143C"
const Cyan = "#00FFFF"
const DarkBlue = "#00008B"
const DarkCyan = "#008B8B"
const DarkGoldenrod = "#B8860B"
const DarkGray = "#A9A9A9"
const DarkGreen = "#006400"
const DarkGrey = "#A9A9A9"
const DarkKhaki = "#BDB76B"
const DarkMagenta = "#8B008B"
const DarkOliveGreen = "#556B2F"
const DarkOrange = "#FF8C00"
const DarkRed = "#8B0000"
const DarkSalmon = "#E9967A"
const DarkSeaGreen = "#8FBC8F"
const DarkSlateBlue = "#483D8B"
const DarkSlateGray = "#2F4F4F"
const DarkSlateGrey = "#2F4F4F"
const DarkTurquoise = "#00CED1"
const DarkViolet = "#9400D3"
const DeepPink = "#FF1493"
const DeepSkyBlue = "#00BFFF"
const DimGray = "#696969"
const DimGrey = "#696969"
const DodgerBlue = "#1E90FF"
const FireBrick = "#B22222"
const FloralWhite = "#FFFAF0"
const ForestGreen = "#228B22"
const Fuchsia = "#FF00FF"
const Gainsboro = "#DCDCDC"
const GhostWhite = "#F8F8FF"
const Gold = "#FFD700"
const Goldenrod = "#DAA520"
const Gray = "#808080"
const Grey = "#808080"
const Green = "#008000"
const GreenYellow = "#ADFF2F"
const HoneyDew = "#F0FFF0"
const HotPink = "#FF69B4"
const IndianRed = "#CD5C5C"
const Indigo = "#4B0082"
const Ivory = "#FFFFF0"
const Khaki = "#F0E68C"
const Lavender = "#E6E6FA"
const LavenderBlush = "#FFF0F5"
const LawnGreen = "#7CFC00"
const LemonChiffon = "#FFFACD"
const LightBlue = "#ADD8E6"
const LightCoral = "#F08080"
const LightCyan = "#E0FFFF"
const LightGoldenrodYellow = "#FAFAD2"
const LightGray = "#D3D3D3"
const LightGreen = "#90EE90"
const LightGrey = "#D3D3D3"
const LightPink = "#FFB6C1"
const LightSalmon = "#FFA07A"
const LightSeaGreen = "#20B2AA"
const LightSkyBlue = "#87CEFA"
const LightSlateGray = "#778899"
const LightSlateGrey = "#778899"
const LightSteelBlue = "#B0C4DE"
const LightYellow = "#FFFFE0"
const Lime = "#00FF00"
const LimeGreen = "#32CD32"
const Linen = "#FAF0E6"
const Magenta = "#FF00FF"
const Maroon = "#800000"
const MediumAquamarine = "#66CDAA"
const MediumBlue = "#0000CD"
const MediumOrchid = "#BA55F3"
const MediumPurple = "#9370DB"
const MediumSeaGreen = "#3CB371"
const MediumSlateBlue = "#7B68EE"
const MediumSpringGreen = "#00FA9A"
const MediumTurquoise = "#48D1CC"
const MediumVioletRed = "#C71585"
const MidnightBlue = "#191970"
const MintCream = "#F5FFFA"
const MistyRose = "#FFE4E1"
const Moccasin = "#FFE4B5"
const NavajoWhite = "#FFDEAD"
const Navy = "#000080"
const OldLace = "#FDF5E6"
const Olive = "#808000"
const OliveDrab = "#6B8E23"
const Orange = "#FFA500"
const OrangeRed = "#FF4500"
const Orchid = "#DA70D6"
const PaleGoldenrod = "#EEE8AA"
const PaleGreen = "#98FB98"
const PaleTurquoise = "#AFEEEE"
const PaleVioletRed = "#DB7093"
const PapayaWhip = "#FFEFD5"
const PeachPuff = "#FFDAB9"
const Peru = "#CD853F"
const Pink = "#FFC0CB"
const Plum = "#DDA0DD"
const PowderBlue = "#B0E0E6"
const Purple = "#800080"
const RebeccaPurple = "#663399"
const Red = "#FF0000"
const RosyBrown = "#BC8F8F"
const RoyalBlue = "#4169E1"
const SaddleBrown = "#8B4513"
const Salmon = "#FA8072"
const SandyBrown = "#F4A460"
const SeaGreen = "#2E8B57"
const SeaShell = "#FFF5EE"
const Sienna = "#A0522D"
const Silver = "#C0C0C0"
const SkyBlue = "#87CEEB"
const SlateBlue = "#6A5ACD"
const SlateGray = "#708090"
const SlateGrey = "#708090"
const Snow = "#FFFAFA"
const SpringGreen = "#00FF7F"
const SteelBlue = "#4682B4"
const Tan = "#D2B48C"
const Teal = "#008080"
const Thistle = "#D8BFD8"
const Tomato = "#FF6347"
const Turquoise = "#40E0D0"
const Violet = "#EE82EE"
const Wheat = "#F5DEB3"
const White = "#FFFFFF"
const WhiteSmoke = "#F5F5F5"
const Yellow = "#FFFF00"
const YellowGreen = "#9ACD32"
