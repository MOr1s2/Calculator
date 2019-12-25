package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

func convert(str string) string {
	SimpleStr := strings.ReplaceAll(str, "sin", "s")
	SimpleStr = strings.ReplaceAll(SimpleStr, "cos", "c")
	SimpleStr = strings.ReplaceAll(SimpleStr, "tan", "t")
	SimpleStr = strings.ReplaceAll(SimpleStr, "exp", "e")
	SimpleStr = strings.ReplaceAll(SimpleStr, "ln", "l")
	return SimpleStr
}
func calc(a string) float64 {
	var currentNum string = "" //当前正在处理的数字
	stack := make([]string, 0)
	op := make([]string, 0)

	var priority map[string]int
	priority = make(map[string]int) //设置运算符优先级，用于处理运算符加入后缀表达式的顺序
	priority["+"] = 1
	priority["-"] = 1
	priority["%"] = 2
	priority["*"] = 2
	priority["/"] = 2
	priority["("] = 0
	priority["^"] = 4
	priority["s"] = 3
	priority["c"] = 3
	priority["t"] = 3
	priority["e"] = 3
	priority["l"] = 3

	for i := 0; i < len(a); i++ {
		c := string(a[i])                     //字符串原式的每一位字符
		if c >= "0" && c <= "9" || c == "." { //若属于数字，加入正在处理数字的末尾
			currentNum = currentNum + c
		} else { //处理运算符
			if currentNum != "" { //若有正在处理的数字，进栈
				stack = append(stack, currentNum)
				currentNum = ""
			}
			if c == "(" || len(op) == 0 { //若是左括号，进入运算符栈
				op = append(op, c)
			} else if c == ")" { //若是右括号，运算符出栈直到将左括号弹出
				for op[len(op)-1] != "(" {
					stack = append(stack, op[len(op)-1])
					op = op[:len(op)-1]
				}
				op = op[:len(op)-1]
			} else { //其他运算符进运算符栈时，循环比较该运算符与栈顶运算符的优先级
				for priority[op[len(op)-1]] >= priority[c] {
					stack = append(stack, op[len(op)-1])
					op = op[:len(op)-1]
					if len(op) == 0 {
						break
					}
				}
				op = append(op, c)
			}
		}
	}
	//收尾处理
	if currentNum != "" {
		stack = append(stack, currentNum)
	}
	for len(op) > 0 {
		stack = append(stack, op[len(op)-1])
		op = op[:len(op)-1]
	}
	fmt.Println(stack)
	res := make([]float64, 0) //维护结果栈，最后栈内只含一个数即结果
	for i := 0; i < len(stack); i++ {
		c := stack[i]
		if c == "+" || c == "-" || c == "%" || c == "*" || c == "/" || c == "^" {
			a1 := res[len(res)-2]
			a2 := res[len(res)-1]
			res = res[:len(res)-2]
			if c == "+" {
				res = append(res, a1+a2)
			} else if c == "-" {
				res = append(res, a1-a2)
			} else if c == "%" {
				res = append(res, float64(int64(a1)%int64(a2)))
			} else if c == "*" {
				res = append(res, a1*a2)
			} else if c == "/" {
				res = append(res, a1/a2)
			} else if c == "^" {
				res = append(res, math.Pow(a1, a2))
			}
		} else if c == "s" || c == "c" || c == "t" || c == "e" || c == "l" {
			a2 := res[len(res)-1]
			res = res[:len(res)-1]
			if c == "s" {
				res = append(res, math.Sin(a2))
			} else if c == "c" {
				res = append(res, math.Cos(a2))
			} else if c == "t" {
				res = append(res, math.Tan(a2))
			} else if c == "e" {
				res = append(res, math.Exp(a2))
			} else if c == "l" {
				res = append(res, math.Log(a2))
			}
		} else {
			num, err := strconv.ParseFloat(c, 64)
			if err != nil {
				fmt.Println("error")
			} else {
				res = append(res, num)
			}

		}
	}
	return res[0]
}

func Echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("receive failed: ", err)
			break
		}

		fmt.Println("received form client: " + reply)
		msg := strconv.FormatFloat(calc(convert(reply)), 'f', 3, 64)
		fmt.Println("sending to client: ", msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("send failed: ", err)
			break
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*将参数中缀表达式expression转为后缀表达式存放在L中，返回L
func infixToSuffix(str string)string{
    for i,_ := range str{        //对表达式中的每一个元素
        if str[i]
        {
            L.append(element)             #将这个元素追加到线性表L后
        }

        else if (element 是 一个运算符)
        {
            While (sop栈 不为空  &&  sop栈顶元素 是一个运算符  && element的优先级 <= sop栈顶运算符元素的优先级 )
            {
                L.append(sop.pop())
            }
            sop.push(element);
        }

        else if(element 是 一个分界符)
        {
            if (element is  '('  )
            {
                sop.push(element)
            }
            else if( element is  ')'  )
            {
                While (sop栈不为空 &&   sop栈顶元素 不是 '('  )    #将匹配的2个括号之间的栈元素弹出，追加到L
                {
                    L.append( sop.pop() );
                }
                if(sop栈不为空 )
                {
                    sop.pop()           #将匹配到的 '(' 弹出丢弃
                }
            }
        }

    }

    While (sop 栈 不为空)       #将sop栈中剩余的所有元素弹出，追加到L后
    {
        L.append(sop.pop())
    }

    return L
}*/
