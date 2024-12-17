# ATTENTION

## Currently Broken
Too many changes in too little time. Fixing everything (sorry)

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
madcolor --text "randomly color a string"

## OUTPUT
This is example output from one run. Since colors are created/assigned randomly, each run
will (and should) differ.
<blockquote>
<span><span style="color: #f2451a;">W</span><span style="color: #4d7a48;">e</span><span style="color: #cc2b3a;"> </span><span style="color: #4255d8;">p</span><span style="color: #402c86;">r</span><span style="color: #2284e8;">o</span><span style="color: #970e3a;">m</span><span style="color: #a9422c;">p</span><span style="color: #1a7970;">t</span><span style="color: #69692e;">l</span><span style="color: #4c8a3f;">y</span><span style="color: #ed1b64;"> </span><span style="color: #400112;">j</span><span style="color: #3f5400;">u</span><span style="color: #1717e4;">d</span><span style="color: #ed0b85;">g</span><span style="color: #dd4b21;">e</span><span style="color: #17163f;">d</span><span style="color: #33088d;"> </span><span style="color: #653a87;">a</span><span style="color: #170bd7;">n</span><span style="color: #f71e55;">t</span><span style="color: #8b1b5c;">i</span><span style="color: #1f564d;">q</span><span style="color: #ca4a3b;">u</span><span style="color: #b2104e;">e</span><span style="color: #581ae6;"> </span><span style="color: #ee2566;">i</span><span style="color: #293655;">v</span><span style="color: #427e13;">o</span><span style="color: #421d89;">r</span><span style="color: #9a0895;">y</span><span style="color: #a30362;"> </span><span style="color: #f82b5b;">b</span><span style="color: #4839c5;">u</span><span style="color: #54776b;">c</span><span style="color: #4b342b;">k</span><span style="color: #62307a;">l</span><span style="color: #d3155d;">e</span><span style="color: #4865aa;">s</span><span style="color: #b61a01;"> </span><span style="color: #6231d0;">f</span><span style="color: #ae1625;">o</span><span style="color: #40239d;">r</span><span style="color: #c3475b;"> </span><span style="color: #3b8967;">t</span><span style="color: #0b5353;">h</span><span style="color: #513564;">e</span><span style="color: #af13b4;"> </span><span style="color: #04988a;">n</span><span style="color: #313ac4;">e</span><span style="color: #c311a6;">x</span><span style="color: #1f6d0a;">t</span><span style="color: #2f13d7;"> </span><span style="color: #781507;">p</span><span style="color: #2a17b1;">r</span><span style="color: #910883;">i</span><span style="color: #da306d;">z</span><span style="color: #15853a;">e</span><span style="color: #1158c2;"> </span><span style="color: #c03e41;">!</span><span style="color: #5b4297;">@</span><span style="color: #096d85;">#</span><span style="color: #fe2b48;">$</span><span style="color: #353864;">%</span><span style="color: #075c4e;">^</span><span style="color: #437265;">*</span><span style="color: #82460a;">(</span><span style="color: #5020d8;">)</span><span style="color: #df0e2b;">{</span><span style="color: #901e8c;">}</span><span style="color: #3a247d;">[</span><span style="color: #7913e9;">]</span><span style="color: #3f593e;">;</span><span style="color: #65202c;">:</span><span style="color: #690771;">.</span><span style="color: #021f5f;">,</span><span style="color: #3f14d7;">?</span></span>
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


