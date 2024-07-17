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

## FLAGS

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


