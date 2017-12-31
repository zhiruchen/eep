# eep
simple expression evaluator written in golang

## examples
```golang
import (
	"fmt"
	"github.com/zhiruchen/eep"
)

func main() {
	fmt.Println(eep.Eval("1")) // print 1
	fmt.Println(eep.Eval("(100*6)")) // print 600
	fmt.Println(eep.Eval("1==1")) // print true
	fmt.Println(eep.Eval(`"abc"=="def"`)) // print false

	add := func(args ...interface{}) interface{} {
		n1, n2 := args[0].(int), args[1].(int)
		return n1 + n2
	}

	exponent := func(args ...interface{}) interface{} {
		n1, n2 := args[0].(float64), args[1].(float64)
		var n float64 = 1
		var i int64 = 1
		for i = 1; i <= int64(n2); i++ {
			n *= n1
		}
		return n
	}

	len := func(args ...interface{}) interface{} {
		return len(args[0].([]int))
	}

	v, _ := EvalWithEnv("add(1000, 24)", map[string]interface{}{"add": add})
	fmt.Printf("type: %T, value: %f\n", v, v) // float64, 1024

	v, _ = EvalWithEnv("exponent(2, 10)", map[string]interface{}{"exponent": exponent})
	fmt.Printf("type: %T, value: %f\n", v, v) // float64, 1024

	v, _ = EvalWithEnv("len(x)", map[string]interface{}{"len": len, "x":[]int{1,2,3,4,5,6}})
	fmt.Printf("type: %T, value: %f\n", ) // float64, 6
}
```
