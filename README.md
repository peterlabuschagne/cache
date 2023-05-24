# cache
A concurrent in memory Go cache the uses one method for adding and reading cache entries and another method for clearing entries with support for Go Generics.

## Usage
```Go
    package main

    import (
	"fmt"
        "time"

        "github.com/peterlabuschagne/cache"
    )

    type mock struct {
        val int
    }
	
    func main() {
	c := cache.New[mock](time.Second * 1)

	m, err := c.Get(func() (mock, error) {
		return mock{val: 123}, nil
	})
	if err != nil {
		fmt.Println(err)
	}
		
	fmt.Printf("%+v\n", m)
    }
```