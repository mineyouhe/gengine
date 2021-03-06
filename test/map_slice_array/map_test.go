package map_slice_array

import (
	"fmt"
	"gengine/base"
	"gengine/builder"
	"gengine/context"
	"gengine/engine"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

type MS struct {
	MII *map[int]int
	MSI map[string]int
	MIS map[int]string
}

const m_1  =`
rule "map test" "m dec"
begin

//map in struct
a = 1
MS.MII[1] = 22
println("MS.MII[1]--->",MS.MII[1])
println("MS.MII[a]--->",MS.MII[a])

b = "1"
MS.MSI["1"] = 227289
println("MS.MSI[\"1\"]--->",MS.MSI["1"])
println("MS.MSI[b]---->", MS.MSI[b])

c = "2"
//
MS.MSI["2"] = 33333
println("MS.MSI[\"2\"]--->", MS.MSI["2"])
println("MS.MSI[c]--->", MS.MSI[c])

d = 1
MS.MIS[1] = "hekwld"
println("MS.MIS[1]--->", MS.MIS[1])
println("MS.MIS[d]--->", MS.MIS[d])

//single map
a = 1
MM[a] = 2222
println("MM[a]->",MM[a])
println("MM[1]->",MM[a])


//can't set value, but can get value
//MMM[1] = 11111111
println(MMM[1])

end
`

func Test_m1(t *testing.T) {
	MS := &MS{
		MII: &map[int]int{1: 1},
		MSI: map[string]int{"hello": 1},
		MIS: map[int]string{1: "helwo"},
	}

	var MM map[int]int
	MM = map[int]int{1:1000,2:1000}

	var MMM map[int]int
	MMM = map[int]int{1:1000,2:1000}

	dataContext := context.NewDataContext()
	dataContext.Add("MS", MS)
	//single map inject, must be ptr
	dataContext.Add("MM", &MM)

	dataContext.Add("MMM", MMM)
	dataContext.Add("println",fmt.Println)


	//init rule engine
	knowledgeContext := base.NewKnowledgeContext()
	ruleBuilder := builder.NewRuleBuilder(knowledgeContext, dataContext)

	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(m_1)
	end1 := time.Now().UnixNano()

	logrus.Infof("rules num:%d, load rules cost time:%d ns", len(knowledgeContext.RuleEntities), end1-start1 )


	if err != nil{
		logrus.Errorf("err:%s ", err)
	}else{
		eng := engine.NewGengine()
		start := time.Now().UnixNano()
		// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
		err := eng.Execute(ruleBuilder, true)
		end := time.Now().UnixNano()
		if err != nil{
			logrus.Errorf("execute rule error: %v", err)
		}
		logrus.Infof("execute rule cost %d ns",end-start)
	}



}