# eep
simple expression evaluator written in golang

## examples
```
import (
	"fmt"
	"github.com/zhiruchen/eep"
)

func main() {
	fmt.Println(eep.Eval("1")) // print 1
	fmt.Println(eep.Eval("(100*6)")) // print 600
	fmt.Println(eep.Eval("1==1")) // print true
	fmt.Println(eep.Eval(`"abc"=="def"`)) // print false
}
```
