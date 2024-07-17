# Madcolor

Set each glyph in a text string with a randomly selected color with the
total brightness of the color less than 160*3 (dark enough
to show up against white). Currently, this is an arbitrary
limit. Writes a &lt;div&gt; with the colorized text to STDOUT.

Unless `--invent` is specified, the random colors are selected
from a preexisting list of HTML colors.

Create a logfile `madcolor.log` in the working directory.

## USAGE
madcolor --text "randomly color a string"

## OUTPUT
This is example output from one run. Since colors are created/assigned randomly, each run
will (and should) differ.

<blockquote>&lt;div&gt;&lt;span style="color:#6b8e23;"&gt;r&lt;/span&gt;&lt;span style="color:#696969;"&gt;a&lt;/span&gt;&lt;span style="color:#595ca1;"&gt;n&lt;/span&gt;&lt;span style="color:#8a2be2;"&gt;d&lt;/span&gt;&lt;span style="color:#6667ab;"&gt;o&lt;/span&gt;&lt;span style="color:#cd5c5c;"&gt;m&lt;/span&gt;&lt;span style="color:#41b6ab;"&gt;l&lt;/span&gt;&lt;span style="color:#5f4b8b;"&gt;y&lt;/span&gt;&lt;span style="color:#7b68ee;"&gt; &lt;/span&gt;&lt;span style="color:#939597;"&gt;c&lt;/span&gt;&lt;span style="color:#ff6f61;"&gt;o&lt;/span&gt;&lt;span style="color:#41b6ab;"&gt;l&lt;/span&gt;&lt;span style="color:#daa520;"&gt;o&lt;/span&gt;&lt;span style="color:#483d8b;"&gt;r&lt;/span&gt;&lt;span style="color:#228b22;"&gt; &lt;/span&gt;&lt;span style="color:#008b8b;"&gt;a&lt;/span&gt;&lt;span style="color:#00ced1;"&gt; &lt;/span&gt;&lt;span style="color:#800000;"&gt;s&lt;/span&gt;&lt;span style="color:#20b2aa;"&gt;t&lt;/span&gt;&lt;span style="color:#0000ff;"&gt;r&lt;/span&gt;&lt;span style="color:#c94476;"&gt;i&lt;/span&gt;&lt;span style="color:#c94476;"&gt;n&lt;/span&gt;&lt;span style="color:#b565a7;"&gt;g&lt;/span&gt;&lt;/div&gt;</blockquote>

## TODO:
* random background color?
  * Complementary background color?
* Select darkness / brightness levels?
* Add more colors?
* ~~Generate random HTML colors?~~
  * Done see `--invent` flag
* Return properly capitalized color names?
* Output file option?
* Add usage() directions
* Create an external colorlist option?
  * Check for madcolor.csv?
* force color of whitespace (default white)?

## FLAGS

#### -a, --anti
Adds  background color for color (r, g, b) of (255-r, 255-g, 255-b). Might not work well with grays; probably still needs some overall tweaking to guarantee a reasonable foreground/background contrast.

#### -d, --debug
Enable debug logic.

#### -h, --help
Help message and usage. Flags are explained, other notes might be
present.

#### -i, --invent
Invent colors, between a minimum/maximum total brightness

#### -q, --quiet
Send debugging / verbose text to the logfile 

#### -t, --text 
Supply a string to decorate. Otherwise, the default string is decorated and returned.


