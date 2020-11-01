# zzcache
Yet another cache module

## Installation

```
go get github.com/saromanov/zzcache
```

## Usage

```go
import (
    "github.com/saromanov/zzcache"
)

func main(){
    c, _ := New(10, "")
	if err := c.Set([]byte("key"), []byte("value"), 2*time.Second); err != nil {
        panic(err)
    }
    c.Get([]byte("key")) // value
}

```
