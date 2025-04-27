# Wafer

Wafer is a simple, stack-based language inspired by Forth. It follows a concatenative style and uses Reverse Polish Notation (RPN) to evaluate expressions.

## Building
Pretty simple:
```bash
make
```

## Running
An input file is required:
```bash
wafer yourfile.w
```

## Project layout
* `src/` - Go source files
* `builtins.tsv` - defines basic builtins
* `generate_builtins.py` - generates Go code to implement builtins.tsv
* `std.w` - standard library

## Syntax

### Numbers
Numbers are 64-bit floats.
```py
42
3.14
-7
```
Numbers are pushed onto the stack.

---

### Strings
Strings are written with double quotes:
```py
"Hello, world!"
"123"
```
Strings are pushed onto the stack.

---

### Comments
Anything after `#` on a line is a comment:
```
# This is a comment
42 # push 42 onto the stack
```
There is currently no support for multiline comments.

---

### Defining words
To define a new word that can execute a block of operations:
```py
: name definition ;
```
Example:
```py
: double
	2 *
;
```
Calling `double` will double the top value on the stack.

---

### Control flow

Blocks are written with curly braces `{}`:
```py
{
	1 2 +
	print
}
```
When a block is encountered, the top of the stack is popped and, if it's non-zero, the block is executed.
This repeats again every time the block finishes.
```py
5 dup {									# since the initial index will be consumed, we need to dup first
	"This will execute " print
	dup print
	" more times" println
	-- dup								# same thing, we dup so the loop still has something to consume
} drop									# we can drop the loop index
```
Conditionals can be thought of as loops that execute either once or never:
```py
myAge 18 >= {
	"You could legally drink in most countries!" println
0 } # for conditionals, make sure 0 is on top when the block ends
```
---

### I/O
Any value on the stack can be printed:
```py
"Hello, world!" print
42 print
```
Files can be loaded as strings:
```py
"message.txt" loadfile print
```
Strings and files can both be ran as subroutines:
```py
"2 3 *" runstring
"other_file.w" runfile
```

## FAQs

### Why "Wafer"?

A long time ago, I created a terminal-based calculator with nearly identical syntax also called Wafer. Every line in the CLI started with `::` for some reason I can't recall-- but I just thought it looked like a waffle so Wafer it was.

I’ve used Wafer (the calculator) in some form or another for nearly a decade, and eventually realized it was nearly Turing-complete. Just for fun, I started tweaking it, eventually porting it from C# to Go, and adding strings, file I/O, and more.

---

### Why not just use GForth, ColorForth, Joy, etc?

Because making something myself is more fun :)

---

### What statements are built-in?
Check `builtins.tsv` and `src/stdlib.go`. It's changed so often it's hard to keep track but it's relatively self-documenting

---

### Why isn’t there a built-in function to do X?

The Go-defined functions are intentionally minimal.
Even `println` is just defined in `src/stdlib.go` as `print "\n" print`.
Higher-level behavior is built from small pieces.

---

### What is Reverse Polish Notation? What is the stack?

Reverse Polish Notation (RPN) is a way of writing expressions where operators come after their operands. Instead of using parentheses to control order of operations, everything is evaluated left-to-right using a stack. For example, `3 4 +` means "add 3 and 4". Or, to be more accurate at the cost of verbosity:
1. Push 3 onto the stack
1. Push 4 onto the stack
1. Pop the top two numbers off the stack, add them, and push their result back on top.

The stack is exactly what it sounds like: a stack of values. Every operation in Wafer "pops" some number of values off the stack and "pushes" some number back on. There’s one global stack that persists across the program, and all operations interact with it.

---

### What data types are supported?

Currently:
- Numbers (64-bit floats, also acting as booleans where appropriate)
- Strings
- Words

Words (variables) internally just reference blocks of code or values.
Eventually, blocks (and arrays?) may become first-class values too.

---

### Why did you do loops/conditionals like that?
I don't know why I made a Turing tarpit out of the only control structure in the language but I did. If you're familiar with Brainfuck, they're basically the same concept: a jump-if-zero in disguise.