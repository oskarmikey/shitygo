package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// Standard ANSI Codes

const (
	Reset = "\033[0m"
	RedBg = "\033[41m"
	WhiTx = "\033[37m"
	GreTx = "\033[32m"
	GraTx = "\033[38;5;241m" // Mörkgrå text (#626262)
)

func clear_terminal() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

type LayoutFormatBase struct {
	Selectable bool
	Content    string
	Action     func()
}

type LayoutFormat struct {
	Items []LayoutFormatBase
}

type MenuComponent interface {
	Draw(SelectionIndex int) LayoutFormat
	RunAction(SelectionIndex int)
}

func (m LayoutFormat) CountSelectable() int {
	count := 0
	for _, item := range m.Items {
		if item.Selectable {
			count++
		}
	}
	return count
}

type Calculatorlayout struct {
	Value1      float64
	HasVal1     bool
	Value2      float64
	HasVal2     bool
	Operator    string
	HasOperator bool
	Result      string
}

func (c *Calculatorlayout) Layout1Func(SelectionIndex int) LayoutFormat {

	val1Str := "Not Set"
	if c.HasVal1 {
		val1Str = fmt.Sprintf("%g", c.Value1)
	}

	opStr := "Not Set"
	if c.Operator != "" {
		opStr = c.Operator
	}

	val2Str := "Not Set"
	if c.HasVal2 {
		val2Str = fmt.Sprintf("%g", c.Value2)
	}

	resStr := c.Result
	if resStr == "" {
		resStr = "Not Ready"
	}

	var totalSelectable int
	totalSelectable = 6

	S := GetSelecterToRender(totalSelectable, SelectionIndex)
	fmt.Print(S)

	layout1 := LayoutFormat{
		Items: []LayoutFormatBase{
			{Selectable: false, Content: fmt.Sprintf("%s Use the: < ˄ v > Keys To navigate %s\n", WhiTx, Reset), Action: nil},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Value1: %s%s\n", GraTx, S[0], WhiTx, val1Str, Reset), Action: func() {
				c.setValue(1)
				MenuSelector(c.Layout1Func)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Operator: %s%s\n", GraTx, S[1], WhiTx, opStr, Reset), Action: func() {
				MenuSelector(c.Layout2OperatorFunc)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Value2: %s%s\n", GraTx, S[2], WhiTx, val2Str, Reset), Action: func() {
				c.setValue(2)
				MenuSelector(c.Layout1Func)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Run: %s%s\n", GraTx, S[3], WhiTx, resStr, Reset), Action: func() {
				if !c.HasVal1 || !c.HasVal2 || !c.HasOperator {
					fmt.Println("THats not gonna work!")
					time.Sleep(1 * time.Second)
					MenuSelector(c.Layout1Func)
				} else {
					c.calculateThatShit()
				}
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Reset: %s\n", GraTx, S[4], WhiTx, Reset), Action: func() {

				c.Value1 = 0
				c.HasVal1 = false
				c.Value2 = 0
				c.HasVal2 = false
				c.Operator = ""
				c.HasOperator = false
				c.Result = ""

				MenuSelector(c.Layout1Func)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Exit: %s\n", RedBg, S[5], GraTx, Reset), Action: func() { os.Exit(0) }},
		},
	}
	clear_terminal()
	DrawMenu(layout1)
	return layout1
}

// actualy you can actualy do !(c.HasVal1 && c.HasOperator && c.HasVal2) i did not know that before. but thanks to ai i do now.... but anyway lets look at the shait i make.
func (c *Calculatorlayout) calculateThatShit() {
	if !c.HasVal1 || !c.HasVal2 || !c.HasOperator {
		fmt.Println("HEY!! YOU cant have it at the default! change that!")

		time.Sleep(1 * time.Second)
		MenuSelector(c.Layout1Func)

	}

	switch c.Operator {
	case "+":
		c.Result = fmt.Sprintf("%g", c.Value1+c.Value2)
		MenuSelector(c.Layout1Func)

	case "-":
		c.Result = fmt.Sprintf("%g", c.Value1-c.Value2)
		MenuSelector(c.Layout1Func)

	case "*":
		c.Result = fmt.Sprintf("%g", c.Value1*c.Value2)
		MenuSelector(c.Layout1Func)

	case "/":
		if c.Value2 == 0 {
			fmt.Println("Cannot divide by zero")
			time.Sleep(1 * time.Second)
			MenuSelector(c.Layout1Func)
		} else {
			c.Result = fmt.Sprintf("%g", c.Value1/c.Value2)
			time.Sleep(1 * time.Second)
			MenuSelector(c.Layout1Func)
		}
	}

}

func (c *Calculatorlayout) Layout2OperatorFunc(SelectionIndex int) LayoutFormat {

	var totalSelectable int
	totalSelectable = 6

	S := GetSelecterToRender(totalSelectable, SelectionIndex)
	fmt.Print(S)

	layout2Operator := LayoutFormat{
		Items: []LayoutFormatBase{
			{Selectable: false, Content: fmt.Sprintf("%s Use the: < ˄ v > Keys To navigate %s\n", WhiTx, Reset), Action: nil},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Add +  %s\n", GraTx, S[0], WhiTx, Reset), Action: func() {
				c.HasOperator = true
				c.Operator = "+"
				MenuSelector(c.Layout1Func)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Subtract -  %s\n", GraTx, S[1], WhiTx, Reset), Action: func() {
				c.HasOperator = true
				c.Operator = "-"
				MenuSelector(c.Layout1Func)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Multiply * %s\n", GraTx, S[2], WhiTx, Reset), Action: func() {
				c.HasOperator = true
				c.Operator = "*"
				MenuSelector(c.Layout1Func)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Divide ÷  %s\n", GraTx, S[3], WhiTx, Reset), Action: func() {
				c.HasOperator = true
				c.Operator = "/"
				MenuSelector(c.Layout1Func)
			}},
			{Selectable: true, Content: fmt.Sprintf("%s%s%s Exit: %s\n", RedBg, S[5], GraTx, Reset), Action: func() { os.Exit(0) }},
		},
	}
	clear_terminal()
	DrawMenu(layout2Operator)
	return layout2Operator

}

// the setValue function is the only fully ai generated.
// i just really did not feel like dealing with retrying for the third time.
// so i told it to get arround the keyboard.Listen. So i just made it.. do some kind of solution..
func (c *Calculatorlayout) setValue(WhichOneQuestionmark int) {
	var inputStr string

	clear_terminal()
	fmt.Println("Yallo... Type your number and press Enter:")
	fmt.Printf("Current input: %s", inputStr)

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter:
			return true, nil // Stop listening and finish input

		case keys.Backspace, keys.Delete:
			if len(inputStr) > 0 {
				inputStr = inputStr[:len(inputStr)-1]
			}

		case keys.CtrlC:
			os.Exit(0)
			return true, nil

		default:
			// Capture typed characters (numbers, decimal points, etc.)
			s := key.String()
			if len(s) == 1 {
				inputStr += s
			}
		}

		// Redraw input prompt with current typed characters
		clear_terminal()
		fmt.Println("Yallo... Type your number and press Enter:")
		fmt.Printf("Current input: %s", inputStr)
		return false, nil
	})

	if err != nil {
		log.Println(err)
		return
	}

	// Convert the typed string into a float64
	if inputStr != "" {
		var numba float64
		_, err := fmt.Sscanf(inputStr, "%f", &numba)
		if err == nil {
			// The old "if numba == 0" block was removed here,
			// so 0 will now pass through safely!

			switch WhichOneQuestionmark {
			case 1:
				c.HasVal1 = true
				c.Value1 = numba
			case 2:
				c.HasVal2 = true
				c.Value2 = numba
			}
		}
	}
}

func DrawMenu(inputA LayoutFormat) {
	clear_terminal()
	for _, item := range inputA.Items {
		fmt.Print(item.Content)
	}
}

func GetSelecterToRender(total int, SelectionIndex int) []string {
	selectors := make([]string, total)
	for i := 0; i < total; i++ {
		if i == SelectionIndex {
			selectors[i] = "> "
		} else {
			selectors[i] = "  "
		}
	}
	return selectors
}

func /*(c *Calculatorlayout)*/ RunAction(layoutFunc func(int) LayoutFormat, selectionIndex int) {
	Layout := layoutFunc(selectionIndex)

	selectableCount := 0
	for _, item := range Layout.Items {
		if item.Selectable {
			if selectableCount == selectionIndex {
				if item.Action != nil {
					item.Action()
				}
				break
			}
			selectableCount++
		}
	}
}

type SelectionIndex struct {
	Index int
}

func MenuSelector(layoutFunc func(SelectionIndex int) LayoutFormat) {
	selction := SelectionIndex{Index: 0}

	layout := layoutFunc(selction.Index)
	Selection_length := layout.CountSelectable()

	clear_terminal()
	DrawMenu(layout)

	err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {

		switch key.Code {

		case keys.Up:
			// fmt.Println("↑ Upppil")
			if selction.Index <= 0 {
				selction.Index = Selection_length - 1
			} else {
				selction.Index--
			}

		case keys.Down:
			selction.Index = (selction.Index + 1) % Selection_length
			// fmt.Println("↓ Nerpil")
			//	if s.Index+1 < Selection_length {
			//		s.Index = s.Index + 1
			//	} else {
			//		s.Index = s.Index - (Selection_length - 1)
			//	}

		case keys.Left:
			//fmt.Println("← Vänsterpil")
			// do nothing
			// c.inputA(selction.Index)
		case keys.Right, keys.Enter:
			//fmt.Println("Enter!")
			RunAction(layoutFunc, selction.Index)
			return true, nil

		case keys.CtrlC:
			os.Exit(0)
			return true, nil // Avslutar lyssnaren

		}

		//	if layout.Items[selction.Index].Action != nil {
		//		layout.Items[selction.Index].Action()
		//	}

		clear_terminal()

		DrawMenu(layoutFunc(selction.Index))

		fmt.Printf("\rCurrent Selected Index: %d", selction.Index) //-- for testing
		return false, nil                                          // Fortsätt lyssna

	})

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("running")
	// text := "> "
	// numbah := 1
	// fmt.Printf(GraTx+text+GreTx+" Value1 : %d\n"+Reset, numbah)

	// fmt.Printf("%s%s%s Value1: %d%s\n", GraTx, text, GreTx, numbah, Reset)

	calc := Calculatorlayout{
		Value1:      0,
		HasVal1:     false,
		Value2:      0,
		HasVal2:     false,
		Operator:    "",
		HasOperator: false,
		Result:      "0",
	}

	// // Print field names alongside their values (%+v)
	// fmt.Printf("%+v\n", calc)

	// // Print individual fields
	// fmt.Printf("Calc State: %f %s %f = %s\n", calc.Value1, calc.Operator, calc.Value2, calc.Result)

	MenuSelector(calc.Layout1Func)
}
