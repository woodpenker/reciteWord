package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/manifoldco/promptui"
)

const (
	chooseCount = 4
)

var (
	red     = Color("\033[1;31m%s\033[0m")
	green   = Color("\033[1;32m%s\033[0m")
	yellow  = Color("\033[1;33m%s\033[0m")
	purple  = Color("\033[1;34m%s\033[0m")
	magenta = Color("\033[1;35m%s\033[0m")
	teal    = Color("\033[1;36m%s\033[0m")
	white   = Color("\033[1;37m%s\033[0m")

	num        int
	dataPath   string
	dataType   string
	lineMaxLen int
)

func main() {
	flag.Parse()
	wds := loadData(dataPath, dataType)
	CallClear()
	fmt.Println(promptui.Styler(promptui.BGWhite)(green("开始背单词: ")))
	qts := wds.randQuestion(num)
	cur := 0
	l := len(qts)
	wrongCount := l
	for wrongCount > 0 {
		if !qts[cur].answered {
			prompt := promptui.Select{
				Label:        teal(qts[cur].content),
				Items:        qts[cur].chooses,
				Size:         4,
				HideSelected: true,
				HideHelp:     true,
			}

			choose, _, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			CallClear()
			if choose != qts[cur].ansIdx {
				fmt.Printf(red("Wrong!!! "))
				ansIdx := shuffle(qts[cur].chooses, qts[cur].ansIdx)
				qts[cur].ansIdx = ansIdx
			} else {
				fmt.Printf(green("Correct! "))
				qts[cur].answered = true
				wrongCount--
			}
			fmt.Printf("%v/%v\t%v : %v \n", purple(l-wrongCount), purple(l), teal(qts[cur].content), yellow(qts[cur].chooses[qts[cur].ansIdx]))
		}
		cur++
		if cur == l {
			cur = 0
		}
	}
}

func init() {
	rand.Seed(time.Now().Unix())

	flag.IntVar(&num, "n", 20, "the number of words to recite")
	flag.IntVar(&lineMaxLen, "l", 180, "the max byte to show one line")
	flag.StringVar(&dataPath, "f", "我的生词本.json", "the data file path")
	flag.StringVar(&dataType, "t", "json", "the data file type")

	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Color print color
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var clear map[string]func() //create a map for storing clear funcs

// CallClear clear terminal screen
func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

type word struct {
	Word    string `json:"word"`
	Chinese string `json:"chinese"`
}

type words struct {
	wds []word
	len int
}

type question struct {
	content  string
	chooses  []string
	ansIdx   int
	answered bool
}

func loadData(file string, tp string) *words {
	var wd words
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	switch tp {
	case "json":
		err = json.Unmarshal(data, &wd.wds)
		if err != nil {
			log.Fatalln(err)
		}
		wd.len = len(wd.wds)
	default:
		log.Fatalln("Unknown data type")
	}
	return &wd
}

func shuffle(s []string, ansIdx int) int {
	for i := len(s) - 1; i >= 0; i-- {
		rnd := rand.Intn(i + 1)
		if rnd == ansIdx {
			ansIdx = i
		} else if i == ansIdx {
			ansIdx = rnd
		}
		s[rnd], s[i] = s[i], s[rnd]
	}
	return ansIdx
}

func toShort(s string) string {
	if len(s) > lineMaxLen {
		return s[:lineMaxLen] + "..."
	}
	return s
}

func (w *words) randChoose(idx int) ([]string, int) {
	res := make([]string, 0, chooseCount)
	res = append(res, toShort(w.wds[idx].Chinese))
	rd := rand.Perm(w.len)[0 : chooseCount-1]
	for _, v := range rd {
		res = append(res, toShort(w.wds[v].Chinese))
	}
	rd = nil
	ansIdx := shuffle(res, 0)
	return res, ansIdx
}

func (w *words) randQuestion(cnt int) []question {
	res := make([]question, 0, cnt)
	rd := rand.Perm(w.len)[0:cnt]
	var chooses []string
	var ansIdx int
	for _, v := range rd {
		chooses, ansIdx = w.randChoose(v)
		res = append(res, question{
			content: w.wds[v].Word,
			chooses: chooses,
			ansIdx:  ansIdx,
		})
	}
	return res
}
