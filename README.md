# Building microservices with Go

Course made by [Nick Jackson](https://www.youtube.com/c/NicJackson).
You can find YouTube
playlist [here](https://www.youtube.com/watch?v=VzBGi_n65iU&list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_).
The code made by the author is [here](https://github.com/nicholasjackson/building-microservices-youtube/tree/main).

# Some useful resources found on web

## Pointers in Go

***What***: A pointer is simply a variable containing the address of another variable.

[Link](https://medium.com/@meeusdylan/when-to-use-pointers-in-go-44c15fe04eac)
[Link](https://medium.com/@annapeterson89/whats-the-point-of-golang-pointers-everything-you-need-to-know-ac5e40581d4d)

```
// & = generating pointer
// * = access to the values held at a location in memory
func TestOne() {
	println("* TEST ONE *")
	helloWorld1 := GenerateHelloWorld()
	pointer := &helloWorld1

	println(pointer) // We're expecting something like 0x...

	value := *pointer
	println(value) // We're expecting Hello world!
}

func TestTwo() {
	println("* TEST TWO *")

	helloWorld2 := GeneratePointerHelloWorld()
	println(helloWorld2)  // we're expecting an address memory 0x..
	println(*helloWorld2) // we're expecting the value
}

func GenerateHelloWorld() string {
	return "hello world!"
}

func GeneratePointerHelloWorld() *string {
	value := string("hello world!")
	return &value
}
```

- They're not like c/c++.
- They're not always faster than passing by value, use it for example: copying large structs.
- Needed for mutability (see example "code#1")
- True absence, when we're using pointer our "zero-default" is `nil`

```
// code#1

type person struct {
    name string
}

func main() {
    p := person{"Richard"}
    rename(p)
    fmt.Println(p)
}
func rename(p person) {
    p.name = "test"
}
```

*NB.* pointer are not safety for concurrency. Avoid on go-routine.



