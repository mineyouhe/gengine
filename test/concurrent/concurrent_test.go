package concurrent

import (
	"fmt"
	"gengine/base"
	"gengine/builder"
	"gengine/context"
	"gengine/engine"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"testing"
	"time"
)


func readAll() string {
	f, err := os.Open("/Users/renyunyi/go/src/gengine/test/rule.gengine")
	if err != nil {
		logrus.Errorf("read file err: %+v", err)
	}

	b,e := ioutil.ReadAll(f)
	if e != nil {
		logrus.Errorf("read file err: %+v", e)
	}
	return string(b)

}



func Test_concurrent(t *testing.T)  {
	user := &User{
		Name: "Calo",
		Age:  0,
		Male: true,
	}

	s := readAll()
	exe_concurrent(user , s)
}

type User struct {
	Name string
	Age  int64
	Male bool
}

func (u *User)GetNum(i int64) int64 {
	return i
}

func (u *User)Print(s string){
	fmt.Println(s)
}

func (u *User)Say(){
	fmt.Println("hello world")
}


func Hello()  {
	fmt.Println("hello")
}

func PrintReal(real float64){
	fmt.Println(real)
}


func exe_concurrent(user *User, s string){

	dataContext := context.NewDataContext()
	//inject struct
	dataContext.Add("User", user)
	//rename and inject
	dataContext.Add("Sout",fmt.Println)
	//inject
	dataContext.Add("Hello",Hello)
	dataContext.Add("PrintReal", PrintReal)

	//init rule engine
	knowledgeContext := base.NewKnowledgeContext()
	ruleBuilder := builder.NewRuleBuilder(knowledgeContext, dataContext)

	//read rule
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(s)
	end1 := time.Now().UnixNano()

	logrus.Infof("rules num:%d, load rules cost time:%d ns", len(knowledgeContext.RuleEntities), end1-start1 )

	if err != nil{
		logrus.Errorf("err:%s ", err)
	}else{
		eng := engine.NewGengine()

		for i :=0 ; i< 10 ; i ++  {
			start := time.Now().UnixNano()
			// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
			eng.ExecuteConcurrent(ruleBuilder)
			end := time.Now().UnixNano()
			logrus.Infof("execute rule cost %d ns",end-start)
		}


		if err != nil{
			logrus.Errorf("execute rule error: %v", err)
		}

		logrus.Infof("user.Age=%d,Name=%s,Male=%t", user.Age, user.Name, user.Male)
	}
}
