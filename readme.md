# MADCOLOR

Set each glyph in a text string with a randomly selected color with the
total brightness of the color less than 160*3 (dark enough
to show up against white). Currently, this is an arbitrary
limit. Writes a &lt;div&gt; with the colorized text to STDOUT.

`madcolor` can also take input from a file, and write it to 
a file. The `--paste` option will write the output to the
clipboard, if supported.

`madcolor` can function as a pipe (read from `STDIN`, write to `STDOUT`) with `--pipe`

Unless `--invent` is specified, the random colors are selected
from a preexisting list of HTML colors drawn from web sources,
starting with web-safe and some other colors such as PANTONE
color of the year and other sources.

Create a logfile `madcolor.log` in the working directory; all error
and verbose information is written to stderr and the logfile. `--quiet`
suppresses this output to stderr (it does not suppress logfile output).

Relative Luminance is used to calculate and determine contrast. There is
a minimum color distance (as grays have distracting/confusing contrast levels)
as well. Flags to set these exist, but have not had any serious testing; the
resulting color combinations (or colors against the supplied background) are
tested for a minimum level of contrast. The defaults appear to work well.

## INSTALLING
Prerequisites:
* `git` installed and configured with GitHub credentials
* working GO compiler
* `GOPATH` environment variable correctly set

To install and compile:
 * `git clone github.com/onyx-tao/madcolor.git`
 * `cd madcolor`
 * `go install madcolor`


## USAGE
madcolor --text "randomly color a string"

## OUTPUT
This is example output from one run. Since colors are created/assigned randomly, each run
will (and should) differ.


<blockquote>&lt;div&gt;&lt;span style="color:#6b8e23;"&gt;r&lt;/span&gt;&lt;span style="color:#696969;"&gt;a&lt;/span&gt;&lt;span style="color:#595ca1;"&gt;n&lt;/span&gt;&lt;span style="color:#8a2be2;"&gt;d&lt;/span&gt;&lt;span style="color:#6667ab;"&gt;o&lt;/span&gt;&lt;span style="color:#cd5c5c;"&gt;m&lt;/span&gt;&lt;span style="color:#41b6ab;"&gt;l&lt;/span&gt;&lt;span style="color:#5f4b8b;"&gt;y&lt;/span&gt;&lt;span style="color:#7b68ee;"&gt; &lt;/span&gt;&lt;span style="color:#939597;"&gt;c&lt;/span&gt;&lt;span style="color:#ff6f61;"&gt;o&lt;/span&gt;&lt;span style="color:#41b6ab;"&gt;l&lt;/span&gt;&lt;span style="color:#daa520;"&gt;o&lt;/span&gt;&lt;span style="color:#483d8b;"&gt;r&lt;/span&gt;&lt;span style="color:#228b22;"&gt; &lt;/span&gt;&lt;span style="color:#008b8b;"&gt;a&lt;/span&gt;&lt;span style="color:#00ced1;"&gt; &lt;/span&gt;&lt;span style="color:#800000;"&gt;s&lt;/span&gt;&lt;span style="color:#20b2aa;"&gt;t&lt;/span&gt;&lt;span style="color:#0000ff;"&gt;r&lt;/span&gt;&lt;span style="color:#c94476;"&gt;i&lt;/span&gt;&lt;span style="color:#c94476;"&gt;n&lt;/span&gt;&lt;span style="color:#b565a7;"&gt;g&lt;/span&gt;&lt;/div&gt;</blockquote>

## TODO:
* ~~Remove min/max brightness levels, replace with contrast control~~
  * Done, see `--contrast`
* ~~random background color~~
  * Done, see `--anti`
* ~~Complementary background color?~~
* ~~Select darkness / brightness levels?~~
  * Done, see `--max` and `--min`
* Add more colors?
  * Added PANTONE color-of-year colors 2000&ndash;2024
* ~~Generate random HTML colors?~~
  * Done see `--invent` flag
* Return properly capitalized color names?
* ~~Input file option?~~
  * Done see `--input`
* ~~Output file option?~~
  * Done see `--output`
* Add usage() directions
* Create an external color list option?
  * Check for madcolor.csv?
* force color of whitespace (default white)?
  * TODO as `--whitespace <string>` where string matches a hex color identifier or name`
    * Hex color identifier: `#[a-fA-F0-9]{6}` (don't bother with three-hex-digit colors)
    * color name: "aliceblue", case ignored
* ~~copy to clipboard~~
  * Done `--nopaste` will disable
* ~~Suppress output to stdout if writing to a file or clipboard~~
  * Done
* Create release YAML for packages on GitHub
  * Not sure how to do this ... must research

## FLAGS

#### -a, --anti
Adds  background color for color (r, g, b) of (255-r, 255-g, 255-b),
but if this has insufficient contrast, invent a color with sufficient
contrast.

#### -b, --background-color
Assume the background color (for contrast calculation). Takes a string
which may be either a six-digit hex value (such as "#AA3388") or the
name of a web color. All web-safe colors are accepted, as well as some
other pantone and other color names. If a color name is unrecognized,
the program terminates. A string matching the regular expression:
`#?([\da-fA-F]{6}|?[\da-fA-F]{3})`
to be a hex value. A three-digit hex string is ALWAYS expanded
to a six digit string by doubling the hex digit. `#1D8` is
equivalent to `#11DD88`.

#### -c, --contrast
This defines the minimum contrast between foreground and
background as an integer from 0 (no contrast) 
to 1000 (max contrast)

#### -d, --debug
Enable debug logic.

#### -h, --help
Help message and usage. Flags are explained, other notes might be
present.

#### -i, --input
Input file to read 

#### -I, --invent
Invent colors, with minimum contrast

#### -o, --output
Write output to a file instead of stdout

#### -nopaste
Suppress output to the clipboard in addition to stdout or input file.
By default, output is **always** copied to the clipboard.

#### -p, --pipe
Function in pipe mode, from STDIN to STDOUT. `--input`, `--output`, 
are disabled. All output to STDOUT is disabled. Output is not placed
in the clipboard.

#### -q, --quiet
By default, debug / verbose output goes to both stderr and the logfile;
this flag suppresses output to logfile.

#### --stdout
Always send output to stdout, even when writing to a file.

#### -t, --text 
Supply a string to decorate. Otherwise, the default string is decorated and returned.


