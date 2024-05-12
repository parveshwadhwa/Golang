/* Each time a defer statement executes, the function value and paramteres to the call are evaluated as usual and saved anew but the actual function is
not invoked. Instead , deffered functions are invoked immediately before the sorrouding funtion returns, in the reverse order they wer deffered. That is,
if the sorrounding function returns through an explicit statement, deffered functions are executed after any result parameters are set by that return
statement but before the function returns to its caller. if a defferd function value evaluates to nil, execution panics when the function is invoked,
not when the "defer" statement is executed

IMPORTANT :- "defer" works on stack(i.e LIFO) that line will be printed first , which comes at last
*/

package main

import "fmt"

func main() {
	defer fmt.Println("Hello world") // when we write it just cut out the line from there and just place it at last of the function before last curly brace
	defer fmt.Println("One")
	defer fmt.Println("Two")
	fmt.Println("Hello welcome to defer")
	myDefer() // here defer statements also exist and then will store in the stack alongs with 3 above lines
}

func myDefer() {
	for i := 0; i < 5; i++ {
		defer fmt.Println(i)
	}
}

/*
Output :
Hello welcome to defer
4
3
2
1
0
Two
One
Hello world
*/
