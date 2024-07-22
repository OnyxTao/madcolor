# Madcolor

Set each glyph in a text string with a randomly selected color with the
total brightness of the color less than 160*3 (dark enough
to show up against white). Currently, this is an arbitrary
limit. Writes a &lt;div&gt; with the colorized text to STDOUT.

`madcolor` can also take input from a file, and write it to 
a file. The `--paste` option will write the output to the
clipboard, if supported.

Unless `--invent` is specified, the random colors are selected
from a preexisting list of HTML colors (presently the websafe
colors &amp; PANTONE colors-of-the-year approximations)

Create a logfile `madcolor.log` in the working directory; all error
and verbose information is written to stderr and the logfile. `--quiet`
suppresses this output to stderr (it does not suppress logfile output).

This program refers to **brightness** in several parameters. This is
a rough measure of the luminance of any particular color as the average 
of the values of red, green, and blue. For a color RGB(a,b,c), the brightness
would be:
<div style="text-align: center;">
<span style="font-size: 150%; font-family: 'JetBrains Mono', monospace; color: navy; background-color: beige; padding: 10px; display: inline-block;"><sup>(a+b+c)</sup>&frasl;<sub>3</sub></span></div>

## USAGE
madcolor --text "randomly color a string"

## OUTPUT
This is example output from one run. Since colors are created/assigned randomly, each run
will (and should) differ.

<blockquote>&lt;div&gt;&lt;span style="color:#6b8e23;"&gt;r&lt;/span&gt;&lt;span style="color:#696969;"&gt;a&lt;/span&gt;&lt;span style="color:#595ca1;"&gt;n&lt;/span&gt;&lt;span style="color:#8a2be2;"&gt;d&lt;/span&gt;&lt;span style="color:#6667ab;"&gt;o&lt;/span&gt;&lt;span style="color:#cd5c5c;"&gt;m&lt;/span&gt;&lt;span style="color:#41b6ab;"&gt;l&lt;/span&gt;&lt;span style="color:#5f4b8b;"&gt;y&lt;/span&gt;&lt;span style="color:#7b68ee;"&gt; &lt;/span&gt;&lt;span style="color:#939597;"&gt;c&lt;/span&gt;&lt;span style="color:#ff6f61;"&gt;o&lt;/span&gt;&lt;span style="color:#41b6ab;"&gt;l&lt;/span&gt;&lt;span style="color:#daa520;"&gt;o&lt;/span&gt;&lt;span style="color:#483d8b;"&gt;r&lt;/span&gt;&lt;span style="color:#228b22;"&gt; &lt;/span&gt;&lt;span style="color:#008b8b;"&gt;a&lt;/span&gt;&lt;span style="color:#00ced1;"&gt; &lt;/span&gt;&lt;span style="color:#800000;"&gt;s&lt;/span&gt;&lt;span style="color:#20b2aa;"&gt;t&lt;/span&gt;&lt;span style="color:#0000ff;"&gt;r&lt;/span&gt;&lt;span style="color:#c94476;"&gt;i&lt;/span&gt;&lt;span style="color:#c94476;"&gt;n&lt;/span&gt;&lt;span style="color:#b565a7;"&gt;g&lt;/span&gt;&lt;/div&gt;</blockquote>

## TODO:
* ~~random background color~~
  * Done, see `--anti`
* ~~Complementary background color?~~
* ~~Select darkness / brightness levels?~~
  * Done, see `--max` and `--min`
* Add more colors?
* ~~Generate random HTML colors?~~
  * Done see `--invent` flag
* Return properly capitalized color names?
* ~~Input file option?~~
  * Done see `--input`
* ~~Output file option?~~
  * Done see `--output`
* Add usage() directions
* Create an external colorlist option?
  * Check for madcolor.csv?
* force color of whitespace (default white)?
  * TODO as `--whitespace <string>` where string matches a hex color identifier or name`
    * Hex color identifier: `#[a-fA-F0-9]{6}` (don't bother with three-hex-digit colors)
    * color name: "aliceblue", case ignored
* ~~copy to clipboard~~
  * Done `--paste`
* Suppress output to stdout if writing to a file or clipboard

## FLAGS

#### -a, --anti
Adds  background color for color (r, g, b) of (255-r, 255-g, 255-b). Might not work well with grays; probably still needs some overall tweaking to guarantee a reasonable foreground/background contrast.

#### -d, --debug
Enable debug logic.

#### -h, --help
Help message and usage. Flags are explained, other notes might be
present.

#### -i, --input
Input file to read 

#### -I, --invent
Invent colors, between minimum/maximum total brightness

#### --min
Set minimum brightness (default 0) for output colors.

#### --max
Set maximum brightness (default 160) for output colors.

#### -o, --output
Write output to a file instead of stdout

#### -p, --paste
Write output to the clipboard in addition to stdout or input file

#### -q, --quiet
By default, debug / verbose output goes to both stderr and the logfile;
this flag suppresses output to logfile.

#### -t, --text 
Supply a string to decorate. Otherwise, the default string is decorated and returned.


