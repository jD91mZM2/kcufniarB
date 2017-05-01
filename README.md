# kcufniarB
*Brainfuck... reversed?*

kcufniarB is a simple command line Brainfuck interpreter.

To run a script, type
```
kcufniarB run test.bf
```
but replace `test.bf` with your Brainfuck file.

## Features

 - Running Brainfuck normally
 - Debugging output (extra information, such as character code and index)
 - Simplify program to understand exactly what's happening
 - Terminal raw mode (skips having to press enter when using the `,` command)
 - Generate Brainfuck code from int (`kcufniarB genval 10` for example returns `+++[->+++<]>+`)
 - Generate Brainfuck code from string (`kcufniarB genstr Hi` for example returns `++++++++[->+++++++++<]>.>+++++++[->+++++++++++++++<]>.`)

## Screenshot
![Imgur](http://i.imgur.com/yufbYp9.png)
![Imgur](http://i.imgur.com/MljwcTt.png)
