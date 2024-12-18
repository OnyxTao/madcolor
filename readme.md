# ATTENTION

# MADCOLOR

Set each glyph in a text string with a randomly selected color with a luminance comparison against the background (default assumption is a white background). Writes a &lt;span&gt; with the colorized text to STDOUT and to the system clipboard (if supported) by default.

Randomness comes from the standard cryptographic library package,
which is several orders of magnitude more secure randomness than 
required.

`madcolor` can also take input from a file, and write it to 
a file. 

`madcolor` can function as a pipe (read from `STDIN`, write to `STDOUT`) with `--pipe`

Unless `--invent` is specified, the random colors are selected
from a preexisting list of HTML colors drawn from various sources. In practice, `--invent` often gives better contrasting results than relying on named colors.

Create a logfile `madcolor.log` in the working directory; all error
and verbose information is written to stderr and the logfile. `--quiet`
suppresses this output to stderr (it does not suppress logfile output).

Relative luminance is used to calculate and determine contrast. 
There is a minimum color distance (as grays have 
distracting/confusing contrast levels) as well. 
Flags to set these exist, but have not had any serious testing; 
the resulting color combinations (or colors against the supplied background) are tested for a minimum level of contrast. The defaults appear to work well.

## INSTALLING
Prerequisites:
* `git` installed and configured with GitHub credentials
* working GO compiler
* `GOPATH` environment variable correctly set

### To compile and install:
 * `git clone https://github.com/onyx-tao/madcolor.git`
 * `cd madcolor`
 * `go install madcolor`

### To install:
 * `go install https://github.com/onyx-tao/madcolor.git@latest`




## USAGE
<span style="font-family: monospace;">
> madcolor --invent --anti --text "Randomly color a string"</span>


## OUTPUT
This is example output from one run. Since colors are created/assigned randomly, each run
will (and should) differ.
<blockquote style="font-size: 400%;">
<span><span style="color: #7903e0; padding: 1px 0px 1px 0px; background-color: #62900d;">R</span><span style="color: #9e04ce; padding: 1px 0px 1px 0px; background-color: #13a93d;">a</span><span style="color: #52f7ce; padding: 1px 0px 1px 0px; background-color: #481732;">n</span><span style="color: 
#e1f63c; padding: 1px 0px 1px 0px; background-color: #e71acb;">d</span><span style="color: #422ee4; padding: 1px 0px 1px 0px; background-color: #66e507;">o</span><span style="color: #d5f918; padding: 1px 0px 1px 0px; background-color: #e11ca4;">m</span><span style="color: #34ff82; padding: 1px 0px 
1px 0px; background-color: #c01267;">l</span><span style="color: #f2ef44; padding: 1px 0px 1px 0px; background-color: #874ef3;">y</span><span style="color: #d7feaa; padding: 1px 0px 1px 0px; background-color: #3c3b0c;"> </span><span style="color: #02c1f6; padding: 1px 0px 1px 0px; background-color:
 #204811;">c</span><span style="color: #09097c; padding: 1px 0px 1px 0px; background-color: #dd921c;">o</span><span style="color: #2efbf3; padding: 1px 0px 1px 0px; background-color: #d410d1;">l</span><span style="color: #FFFFFF; padding: 1px 0px 1px 0px; background-color: #1688ca;">o</span><span s
tyle="color: #b9e102; padding: 1px 0px 1px 0px; background-color: #f306e5;">r</span><span style="color: #5907d0; padding: 1px 0px 1px 0px; background-color: #22ae14;"> </span><span style="color: #0ae1f2; padding: 1px 0px 1px 0px; background-color: #a9131f;">a</span><span style="color: #f7f50c; padd
ing: 1px 0px 1px 0px; background-color: #4040bb;"> </span><span style="color: #daf42a; padding: 1px 0px 1px 0px; background-color: #f302ca;">s</span><span style="color: #FFFFFF; padding: 1px 0px 1px 0px; background-color: #985b98;">t</span><span style="color: #2929f5; padding: 1px 0px 1px 0px; back
ground-color: #5ad309;">r</span><span style="color: #091004; padding: 1px 0px 1px 0px; background-color: #f8425f;">i</span><span style="color: #470bec; padding: 1px 0px 1px 0px; background-color: #f29870;">n</span><span style="color: #cf8cf2; padding: 1px 0px 1px 0px; background-color: #420b30;">g</span></span>

</blockquote>



## TODO:
* ~~Remove min/max brightness levels, replace with contrast control~~
  * Done, see `--contrast`
* ~~random background color~~
  * Done, see `--anti`
* ~~Complementary background color?~~
* ~~Convert directly from clipboard input?~~
  * Done, see `--buff`
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
Use the `--flag=false` to set a boolean flag to `false`; `--flag` sets
the flag to `true`.

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
`#?([\da-fA-F]{6}|?[\da-fA-F]{3})` is interpreted as a hex color value. 
A three-digit hex string is expanded to a six digit string by doubling
the hex digit per the W3 recommendation
[https://www.w3.org/TR/css-color-3 _section 4.2.1_](https://www.w3.org/TR/css-color-3/#numerical).

#### -buff
Colorize clipboard contents. Output is placed back in the clipboard,
and written to `STDOUT` (by default).

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
Randomly generate (invent) colors, with high minimum contrast with the background (or
invented background)

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
this flag suppresses output to stderr. Output to logfile is **not** suppressed.

#### --stdout
Always send output to stdout, even when writing to a file.

#### -t, --text 
Supply a string to decorate. Otherwise, the default string is decorated and returned.


