# self-delete

> **DISCLAIMER.** All information contained in this repository is for educational and research purposes only. The owner is not responsible for any illegal use of included code snippets.

The `selfdelete` package allows an executable to perform self deletion while running. This concept was initially discovered by [@jonasLyk](https://twitter.com/jonasLyk/status/1350401461985955840).

## Usage
Import the package and call the method `SelfDelete()`.

```go
import selfdelete "github.com/thesh1n/self-delete"

package main

func main(){
    println("Hello world!")

    selfdelete.SelfDelete()
}
```

## References

- https://github.com/LloydLabs/delete-self-poc
- https://twitter.com/jonasLyk/status/1350401461985955840