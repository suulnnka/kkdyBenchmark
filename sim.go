// GOOS=windows GOARCH=amd64 go build sim.go
// GOOS=darwin GOARCH=arm64 go build sim.go
	
package main

import (
	"os"
	"fmt"
	"time"
	"math"
	"runtime"
	"sync/atomic"
)

var log bool = true
var debug bool = false
var times int32 = 2148 // 2148
var step int32 = 1000000 // 1000000
var fout *os.File

type Random struct{
	inext int
	inextp int
	SeedArray [56]int32
	seed int32

	mp int
}
func (random *Random) init( Seed int32 ){
    const MBIG int32 = 2147483647;
    const MSEED int32 = 161803398;

    const Int32_min int32 = -2147483648;
    const Int32_max int32 = 2147483647;

	var mj, mk int32

	mj = MSEED - Seed;
	random.SeedArray[55] = mj;
	mk = 1;
	for i := 1; i < 55; i++ {
		ii := (21 * i) % 55;
		random.SeedArray[ii] = mk;
		mk = mj - mk;
		if (mk < 0) { mk = mk + MBIG }
		mj = random.SeedArray[ii];
	}
	for k := 1; k < 5; k++ {
		for i := 1; i < 56; i++ {
			random.SeedArray[i] = random.SeedArray[i] - random.SeedArray[1 + (i + 30) % 55]
			if random.SeedArray[i] < 0 { random.SeedArray[i] = random.SeedArray[i] + MBIG }
		}
	}
	random.inext = 0
	random.inextp = 21
	random.mp = 0
	random.seed = Seed
}
func (random *Random) NextDouble() float64{
	const MBIG int32 = 2147483647;

	random.inext = random.inext + 1
	random.inextp = random.inextp + 1
	if random.inext >= 56 {
		random.inext = 1
	}
	if random.inextp >= 56 {
		random.inextp = 1
	}
	var retVal int32 = random.SeedArray[random.inext] - random.SeedArray[random.inextp]

	if (retVal == MBIG) {retVal--}    
	if (retVal < 0){
	  retVal = retVal + MBIG
	}

	random.SeedArray[random.inext] = retVal;

	return float64(retVal) * (1.0 / float64(MBIG));
}
func (random *Random) attack (def int, hp int) int {
	const prob float64 = 0.20
	const kkdyAtk float64 = 446.0

	var normal int = int(math.Round(kkdyAtk))
	var cirt int = int(math.Round(kkdyAtk * 1.6))
	var skill int = int(math.Round(kkdyAtk * 1.4))
	var skill_crit int = int(math.Round(kkdyAtk * 1.6 * 1.4))

	damage := 0
	if random.mp == 4 {
	  if random.NextDouble() > prob {
		damage = skill - def;
	  } else {
		damage = skill_crit - def;
	  }
	  if hp > damage {
		if random.NextDouble() > prob {
			damage = damage + skill - def;
		} else {
			damage = damage + skill_crit - def;
		}
	  }

	} else {
	  if random.NextDouble() > prob {
		damage = normal - def;
	  } else {
		damage = cirt - def;
	  }
	}

	random.mp = random.mp + 1
	if random.mp > 4 {
	  random.mp = 0;
	}

	return damage;
}

const (
	def_g int = 0
	hp_g int = 1700
	
	def_s int = 100
	hp_s int = 2000

	def_t int = 85
	hp_t int = 2000

	def_d int = 250
	hp_d int = 2050

	def_l int = 50
	hp_l int = 5000

	def_w int = 100
	hp_w int = 10000
)
var (
	def []int
	max_hp []int
	
	process chan int32
	wait chan bool
)

func print(args ...interface{}){
	//fmt.Println(args...)
	if debug {
		fmt.Println(args...)
	}else{
		fmt.Fprintln(fout,args...)
	}
}

// 统计
var (
	统计_阶段_1 int32 = 0
	统计_阶段_3 int32 = 0
	统计_死亡_13 int32 = 0
	统计_阶段_4 int32 = 0
	统计_阶段_5_1 int32 = 0
	统计_阶段_5_2 int32 = 0
	统计_阶段_5_1_分支 int32 = 0
	统计_阶段_5_2_分支 int32 = 0
	统计_死亡_18 int32 = 0
	统计_速死_18 int32 = 0
	统计_速死_19 int32 = 0
	统计_死亡_20 int32 = 0
	统计_死亡_21 int32 = 0
	统计_速死_21 int32 = 0
	统计_死亡_22 int32 = 0
	统计_速死_23 int32 = 0
	统计_死亡_23 int32 = 0
	统计_死亡_w int32 = 0
	统计_阶段_w int32 = 0
	
	统计_流未死 int32 = 0
	统计_盾未死 int32 = 0
	
	TODO_T1 int32 = 0
	TODO_T2 int32 = 0
	TODO_T3 int32 = 0
	TODO_T4 int32 = 0

	DONE int32 = 0
)

// timeline
var (
	阶段_2 [][]int = [][]int{
		{7,0,0},
		{7,0,0},
		{7,0,0},
		{7,0,0},
		{8,7,0},
		{8,7,0},
		{8,7,0},
		{8,7,0},
		{8,7,9},
		{7,9,10},
		{7,9,10},
		{7,9,10},
		{7,9,10},
		{9,10,0},
		{9,10,0},
		{9,10,0},
	}
	阶段_3 [][]int = [][]int{
		{11,0,0},
		{11,0,0},
		{11,12,0},
		{11,13,12},
		{11,13,12},
		{11,13,12},
		{13,11,12},
		{11,12,0},
		{12,0,0},
		{12,0,0},
		{12,0,0},
		{12,0,0},
		{12,0,0},
		{12,0,0},
		{12,0,0},
		{12,0,0},
	}
	阶段_4 [][]int = [][]int{
		{14,0,0},
		{14,0,0},
		{14,15,0},
		{14,15,0},
		{14,15,16},
		{14,15,16},
		{14,15,16},
		{14,15,16},
		{15,17,16},
		{15,17,16},
		{17,16,0},
		{17,16,0},
		{16,0,0},
		{16,0,0},
		{16,0,0},
		{16,18,0},
		{16,18,19},
	}
	阶段_5_1_S [][]int = [][]int{
		{18,0,0},
		{18,0,0},
		{18,19,0},
		{18,19,0},
		{18,19,0},
		{19,0,0}, // {19,20,0}, 
		{19,20,21},
		{19,20,21},
		{19,20,21},
		{20,21,0},
		{21,0,0},
		{21,0,0},
		{21,0,0},
	}
	阶段_5_1_M [][]int = [][]int{
		{22,21, 0, 0},
		{22,21,23, 0},
		{22,21,23, 0},
		{22,21,23,24},
		{22,21,23,24},
	}
	阶段_5_1_燃烧_不变 [][]int = [][]int{
		{22,21},
		{22,21,23},
		{22,21,23},
		{22,21,23,24},
		{22,21,24,23},
		{22,25,21,24,23},
		{22,25,21,24,23},
		{22,25,21,24,23},
		{22,25,21,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{21,22,25,24,23},
		{22,21,24,25,23},
		{22,24,25,26,23},
		{24,26,25,23},
		{24,26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{25,23},
		{25,27,23},
		{27,25,23,28},
		{27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,23,28},
		{25,29,23,28},
		{29,23,28},
		{29,23,28},
		{29,23,28},
		{23,28},
		{23,28},
		{28,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_1_燃烧_变 [][]int = [][]int{
		{24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{24,25,23},
		{24,25,23},
		{24,25,26,23},
		{24,26,25,23},
		{24,26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{25,27,23},
		{27,25,23},
		{27,25,23,28},
		{27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,28,23},
		{25,29,28,23},
		{29,28,23},
		{29,28,23},
		{29,28,23},
		{28,23},
		{28,23},
		{28,23},
		{23,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_1_双持_不变 [][]int = [][]int{
		{22,21},
		{22,21,23},
		{22,21,23},
		{22,21,23},
		{22,21,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{21,22,25,23,24},
		{22,21,25,24,23},
		{22,25,26,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,24,25,23},
		{24,25,23},
		{24,25,27,23},
		{24,27,25,23,28},
		{24,27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,23,28},
		{25,29,23,28},
		{29,23,28},
		{29,23,28},
		{29,23,28},
		{23,28},
		{23,28},
		{28,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_1_双持_变 [][]int = [][]int{
		{24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,26,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,24,25},
		{26,24,25},
		{24,25,27},
		{24,27,25},
		{24,27,25,28},
		{27,25,29,28},
		{27,25,29,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{25,29,28},
		{25,29,28},
		{29,28},
		{29,28},
		{29,28},
		{28},
		{28},
		{28},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_1_w_不变 [][]int = [][]int{
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,25,26,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,24,25,23},
		{24,25,23},
		{24,25,27,23},
		{24,27,25,23,28},
		{24,27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,28,23},
		{29,25,28,23},
		{29,25,23,28},
		{29,28,25,23},
		{29,28,23,25},
		{28,23,25},
		{23,28,25},
		{28,25,23,30},
		{23,25,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_1_w_变 [][]int = [][]int{
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,26,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,24,25,23},
		{24,25,23},
		{24,25,27,23},
		{24,27,25,23,28},
		{24,27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,28,23},
		{25,29,28,23},
		{29,28},
		{29,28},
		{29,28},
		{28},
		{28},
		{28,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_21速死_w_不变 [][]int = [][]int{{22},{22,23},{22,23},{22,23},{22,23,24},{22,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,26,24,23},{22,25,26,24,23},{26,25,24,23},{26,25,24,23},{26,25,24,23},{26,25,24,23},{26,24,25,23},{26,24,25,23},{24,25,27,23},{24,27,25,23},{24,27,25,23,28},{27,25,29,23,28},{27,25,29,23,28},{29,25,23,28},{29,25,23,28},{29,25,23,28},{29,25,28,23},{29,25,28,23},{29,25,28,23},{29,28,25,23},{29,28,23,25},{28,23,25},{28,23,25},{28,25,23,30},{23,25,30},{30},{30},{30},{30},{30},{30},{30},{30},{30},{30},{30}}
	阶段_5_21速死_w_变 [][]int = [][]int{{22},{22,23},{22,23},{22,23},{22,23,24},{22,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,24,23},{22,25,26,24,23},{22,25,26,24,23},{26,25,24,23},{26,25,24,23},{26,25,24,23},{26,25,24,23},{26,24,25,23},{26,24,25,23},{24,25,27,23},{24,27,25,23},{24,27,25,23,28},{27,25,29,23,28},{27,25,29,23,28},{29,25,23,28},{29,25,23,28},{29,25,23,28},{25,29,28},{25,29,28},{29,28},{29,28},{29,28},{28},{28},{28,30},{30},{30},{30},{30},{30},{30},{30},{30},{30},{30},{30},{30}}
	阶段_5_21速死_燃烧_不变 [][]int = [][]int{
		{22},
		{22,23},
		{22,23},
		{22,23},
		{22,23,24},
		{22,24,23},
		{22,25,24,23},
		{22,25,24,23},
		{22,25,24,23},
		{22,25,24,23},
		{22,25,24,23},
		{22,25,24,23},
		{22,24,25,23},
		{22,24,25,26,23},
		{22,24,25,26,23},
		{24,26,25,23},
		{24,26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{25,27,23},
		{27,25,23},
		{27,25,23,28},
		{27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,23,28},
		{25,29,23,28},
		{29,28,23},
		{29,23,28},
		{29,23,28},
		{23,28},
		{23,28},
		{23,28},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_21速死_燃烧_变 [][]int = [][]int{
		{23},
		{23},
		{23},
		{24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{24,25,23},
		{24,25,26,23},
		{24,26,25,23},
		{24,26,25,23},
		{24,26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{25,23},
		{25,27,23},
		{27,25,23,28},
		{27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,23,28},
		{25,29,23,28},
		{29,23,28},
		{29,23,28},
		{29,23,28},
		{23,28},
		{23,28},
		{28,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_21速死_双持_不变 [][]int = [][]int{
		{22},
		{22,23},
		{22,23},
		{22,23},
		{22,23,24},
		{22,23,24},
		{22,25,23,24},
		{22,25,23,24},
		{22,25,23,24},
		{22,25,23,24},
		{22,25,23,24},
		{22,25,23,24},
		{22,25,24,23},
		{22,25,24,23},
		{22,25,26,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,24,25,23},
		{26,24,25,23},
		{24,25,27,23},
		{24,27,25,23},
		{24,27,25,23,28},
		{27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,23,28},
		{25,29,23,28},
		{29,23,28},
		{29,23,28},
		{29,23,28},
		{23,28},
		{23,28},
		{23,28},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_21速死_双持_变 [][]int = [][]int{
		{24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,26,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,24,25},
		{24,25},
		{24,25,27},
		{24,27,25,28},
		{24,27,25,29,28},
		{27,25,29,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{25,29,28},
		{25,29,28},
		{29,28},
		{29,28},
		{29,28},
		{28},
		{28},
		{28,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_2_S [][]int = [][]int{ {18,19},{18,19},{19},{19,20},{19,20,21},{19,20,21},{19,20,21},{21},{21},{21},{21} }
	阶段_5_2_M [][]int = [][]int{
		{21,22,23},
		{22,21,23},
		{22,21,23},
		{22,21,23,24},
		{22,21,24,23},
	}
	阶段_5_2_燃烧_不变 [][]int = [][]int{
		{21,22,23},
		{22,21,23},
		{22,21,23},
		{22,21,23,24},
		{22,21,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,21,24,25,23},
		{22,24,25,26,23},
		{24,26,25,23},
		{24,26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{25,23},
		{27,25,23},
		{27,25,23,28},
		{27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,28,23},
		{25,29,28,23},
		{29,23,28},
		{29,28,23},
		{29,28,23},
		{28,23},
		{23,28},
		{28,23},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_2_燃烧_变 [][]int = [][]int{
		{24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{25,24,23},
		{24,25,23},
		{24,25,26,23},
		{24,25,26,23},
		{24,26,25,23},
		{24,26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{26,25,23},
		{25,27,23},
		{27,25,23},
		{27,25,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,28,23},
		{25,29,28,23},
		{29,23,28},
		{29,28,23},
		{29,28,23},
		{28,23},
		{23,28},
		{28,23,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_2_双持_不变 [][]int = [][]int{
		{21,22,23},
		{22,21,23},
		{22,21,23},
		{22,21,23,24},
		{22,21,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,23,24},
		{22,21,25,24,23},
		{22,21,25,24,23},
		{22,25,26,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,25,24,23},
		{26,24,25,23},
		{24,25,23},
		{24,27,25,23},
		{24,27,25,23,28},
		{27,25,29,23,28},
		{27,25,29,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{29,25,23,28},
		{25,29,23,28},
		{25,29,23,28},
		{29,23,28},
		{29,23,28},
		{29,23,28},
		{23,28},
		{23,28},
		{23,28},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_2_w_不变 [][]int = [][]int{
		{22,21,25,24},
		{22,21,25,24},
		{22,21,25,24},
		{22,21,25,24},
		{22,21,25,24},
		{22,21,25,24},
		{22,21,25,24},
		{22,21,25,24},
		{22,25,26,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,24,25},
		{24,25},
		{24,27,25},
		{24,27,25,28},
		{27,25,29,28},
		{27,25,29,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{29,28,25},
		{29,28,25},
		{28,25},
		{28,25},
		{28,25},
		{23,25,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_2_w_变 [][]int = [][]int{
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,26,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,25,24},
		{26,24,25},
		{24,25},
		{24,27,25},
		{24,27,25,28},
		{27,25,29,28},
		{27,25,29,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{25,29,28},
		{29,28},
		{29,28},
		{29,28},
		{28},
		{28},
		{28},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
	阶段_5_23速死 [][]int = [][]int{
		{24},
		{24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{25,24},
		{24,25},
		{24,25,26},
		{24,26,25},
		{24,26,25},
		{26,25},
		{26,25},
		{26,25},
		{26,25},
		{25},
		{25,27},
		{27,25,28},
		{27,25,29,28},
		{27,25,29,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{29,25,28},
		{25,29,28},
		{25,29,28},
		{29,28},
		{29,28},
		{29,28},
		{28},
		{28},
		{28,30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
		{30},
	}
)

func main(){
	sim()
}

func sim() {

	file, fileErr := os.Create("output.txt")
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}
	fout = file

	def = []int{0,def_g,def_g,def_g,def_g,def_g,def_g,def_s,def_g,def_g,def_g,def_s,def_d,def_g,def_t,def_g,def_s,def_g,def_g,def_s,def_g,def_s,def_s,def_t,def_s,def_w,def_s,def_g,def_d,def_l,def_t}
	max_hp = []int{0,hp_g,hp_g,hp_g,hp_g,hp_g,hp_g,hp_s,hp_g,hp_g,hp_g,hp_s,hp_d,hp_g,hp_t,hp_g,hp_s,hp_g,hp_g,hp_s,hp_g,hp_s,hp_s,hp_t,hp_s,hp_w,hp_s,hp_g,hp_d,hp_l,hp_t}

	startTime := time.Now().Unix()

	if debug == false {
		coreNumber := runtime.NumCPU()
		print("cpu core number is",coreNumber)
		process = make(chan int32 ,coreNumber)
		wait = make(chan bool)

		for i := int32(0) ; i < times ; i ++ {
			var min int32 = i * step
			var max int32 = min + step - 1
			var test_min int64 = int64(i) * int64(step)
			var test_max int64 = int64(i) * int64(step) + int64(step)
			if test_min > 2147483647 {
				min = 2147483647
			}
			if test_max > 2147483647 {
				max = 2147483647
			}
			process <- i
			go run_range(min,max)
		}

		for j := 0 ; j < int(times) ; j ++ {
			<-wait
		}

	}else{
		// run(44990697) // 被打死
		run(2124733276) // 吃一个暴击

		// run(1079385285) // +2 攻击
		// run(1097943659) // +1 攻击
		// run(1094706438) // +2% 暴率
		// run(1813350023) // +2% 暴率
		// run(1454282831) // +2% 暴率
		// run(1748821225) // +1% 暴率 被打死

	}

	if log {
		print("统计_阶段_1",统计_阶段_1)
		print("统计_阶段_3",统计_阶段_3)
		print("统计_死亡_13",统计_死亡_13)
		print("统计_阶段_4",统计_阶段_4)
		print("统计_阶段_5_1",统计_阶段_5_1)
		print("统计_阶段_5_2",统计_阶段_5_2)
		print("统计_阶段_5_1_分支",统计_阶段_5_1_分支)
		print("统计_阶段_5_2_分支",统计_阶段_5_2_分支)
		print("统计_死亡_18",统计_死亡_18)
		print("统计_速死_18",统计_速死_18)
		print("统计_速死_19",统计_速死_19)
		print("统计_死亡_20",统计_死亡_20)
		print("统计_死亡_21",统计_死亡_21)
		print("统计_速死_21",统计_速死_21)
		print("统计_死亡_22",统计_死亡_22)
		print("统计_速死_23",统计_速死_23)
		print("统计_死亡_23",统计_死亡_23)
		print("统计_死亡_w",统计_死亡_w)
		print("统计_阶段_w",统计_阶段_w)
		print("统计_流未死",统计_流未死)
		print("统计_盾未死",统计_盾未死)
		print("TODO T1",TODO_T1)
		print("TODO T2",TODO_T2)
		print("TODO T3",TODO_T3)
		print("TODO T4",TODO_T4)
		print("DONE",DONE)

		endTime := time.Now().Unix()

		print("共耗时:",endTime - startTime,"秒")
	}

}

func run_range (min,max int32){
	if min == max {
		<-process
		wait <- true
		return
	}

	for t := min ; ; t ++ {
		run(t)

		if t == max {
			break
		}
	}

	if log {
		//if min % 10000000 == 0 {
			print("complete",min,"->",max)
		//}
	}

	<-process
	wait <- true
}
func run(seed int32){
	var rand Random
	rand.init(seed)

	var hp [31]int
	copy(hp[:],max_hp[:])

	// 阶段 1
	// 1 - 6 狗 狗 狗 狗 狗 狗
	// 三帧索敌的不影响次序
	atomic.AddInt32(&统计_阶段_1, 1)
	for enemy := 1 ; enemy <= 6 ; enemy ++ {
		for hp[enemy] > 0 {
			hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
		}
	}

	// 阶段 2
	// 7 - 10 双 狗 狗 狗
	// 三帧索敌的不影响次序
	for i := 0 ; i < len(阶段_2) ; i ++ {
		for j := 0 ; j < 3 ; j ++ {
			enemy := 阶段_2[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	// 阶段 3
	// 11 - 13 双 狗 盾
	// 三帧索敌的不影响次序
	// 7次攻击以内打死 8狗 且 14次攻击未打死 10 狗
	// 15次攻击打死 10 狗
	// 15次攻击未打死 10 狗 , 不会出现
	atomic.AddInt32(&统计_阶段_3, 1)
	for i := 0 ; i < len(阶段_3) ; i ++ {
		for j := 0 ; j < 3 ; j ++ {
			enemy := 阶段_3[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	if hp[13] > 0 {
		atomic.AddInt32(&统计_死亡_13, 1)
		return
	}

	// 阶段 4
	// 14 - 17 投 狗 双 狗
	atomic.AddInt32(&统计_阶段_4, 1)
	阶段_4_攻击次数 := 0

	for i := 0 ; i < len(阶段_4) ; i ++ {
		select_enemy := -1
		for j := 0 ; j < 3 ; j ++ {
			enemy := 阶段_4[i][j]
			if hp[enemy] > 0 {
				// if enemy == 14 {
					// print(enemy,i,hp[enemy],rand.seed)
				// }
				阶段_4_攻击次数 = i + 1
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 && 阶段_4_攻击次数 >= 10 {
			// 16及之前怪速死仅影响三帧索敌,对次序无影响
			// >= 10 用于合并情况
			break
		}
	}

	// 阶段 5
	// S 18 狗 19 双 20 狗 21 双
	// M 22 双 23 投 24 双
	// L 25 W 26 双 27 狗 28 盾 29 流 30 投
    if 阶段_4_攻击次数 < 15 {
		// 如果四阶段14下及以内解决，状况1

		atomic.AddInt32(&统计_阶段_5_1, 1)

		// S 18 狗 19 双 20 狗 21 双
		for i := 0 ; i < len(阶段_5_1_S) ; i ++ {
			select_enemy := -1
			for j := 0 ; j < 3 ; j ++ {
				enemy := 阶段_5_1_S[i][j]
				if hp[enemy] > 0 {
					select_enemy = enemy
					hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
					break
				}
			}
			if select_enemy == -1 {
				// 快速死亡判断
				if i == 1 {
					atomic.AddInt32(&统计_速死_18, 1)
					run_dead_18(rand,hp)
					return
				}
				if i == 5 {
					// 根据三帧索敌,有一种情况有可能会出现19速死,其他两种可能性不会出现这种情况
					atomic.AddInt32(&统计_速死_19, 1)
					// 特殊判断三帧索敌中的一种情况
					// run_dead_19(rand,hp)
					atomic.AddInt32(&TODO_T3, 1)
					// 其他两种情况攻击20,可不影响
					hp[20] = hp[20] - rand.attack(def[20],hp[20])
				}
				if i > 9 {
					run_dead_21(rand,hp)
					return
				}
				
			}
		}

		if hp[20] > 0 {
			atomic.AddInt32(&统计_死亡_20, 1)
			return
		}

		run_part_5_1(rand,hp)
	}else{
		// 如果四阶段17下解决，状况2

		atomic.AddInt32(&统计_阶段_5_2, 1)

		if hp[19] != max_hp[19] {
			// 如果第16次攻击打死 18 狗
			// 根据三帧索敌,会空闲3帧,原三帧索敌的三种情况变四种
			// 只有最早寻敌的一帧可能会变晚

			// 直接跟5_1_18速死轴合并
			atomic.AddInt32(&统计_速死_18, 1)
			//run_dead_18(rand,hp)
			atomic.AddInt32(&TODO_T3, 1)
			//return
		}

		for i := 0 ; i < len(阶段_5_2_S) ; i ++ {
			select_enemy := -1
			for j := 0 ; j < len(阶段_5_2_S[i]) ; j ++ {
				enemy := 阶段_5_2_S[i][j]
				if hp[enemy] > 0 {
					select_enemy = enemy
					hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
					break
				}
			}
			if select_enemy == -1 {
				run_dead_21(rand,hp)
				return
			}
		}
		if hp[18] > 0 {
			atomic.AddInt32(&统计_死亡_18, 1)
			return
		}
		if hp[20] > 0 {
			atomic.AddInt32(&统计_死亡_20, 1)
			return
		}

		run_part_5_2(rand,hp)
	}
}
func run_dead_18(rand Random, hp [31]int){
	// 进入条件
	// 18速死
	// 可能跟5-2可混轴

	if hp[19] == max_hp[19] {
		// 来自5-1
		hp[19] = hp[19] - rand.attack(def[19],hp[19])
	}

	for i := 0 ; i < len(阶段_5_2_S) ; i ++ {
		select_enemy := -1
		for j := 0 ; j < len(阶段_5_2_S[i]) ; j ++ {
			enemy := 阶段_5_2_S[i][j]
			if hp[enemy] > 0 {
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 {
			run_dead_21(rand,hp)
			return
		}
	}
	if hp[18] > 0 {
		atomic.AddInt32(&统计_死亡_18, 1)
		return
	}
	if hp[20] > 0 {
		atomic.AddInt32(&统计_死亡_20, 1)
		return
	}

	run_part_5_2(rand,hp)
}
func run_dead_19(rand Random, hp [31]int){
	// 进入条件
	// 19速死

	// 只有三帧索敌的这一种情况
	// 可能跟5-1可混轴

	atomic.AddInt32(&TODO_T3, 1)
}
func run_dead_21(rand Random, hp [31]int){
	// 进入条件
	// 21速死

	if hp[20] > 0 {
		atomic.AddInt32(&统计_死亡_20, 1)
		return
	}

	atomic.AddInt32(&统计_速死_21, 1)

	run_dead_21_w_不变(rand,hp)
	run_dead_21_双持_不变(rand,hp)
	run_dead_21_燃烧_不变(rand,hp)
	run_dead_21_w_变(rand,hp)
	run_dead_21_双持_变(rand,hp)
	run_dead_21_燃烧_变(rand,hp)
}

// 44990697
func run_dead_21_w_不变(rand Random, hp [31]int){
	for i := 0 ; i < 6 ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_w_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_w_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	if hp[22] > 0 {
		atomic.AddInt32(&统计_死亡_22, 1)
		return
	}
	if hp[23] > 1000 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000
	hp[25] = hp[25] - 500

	for i := 6 ; i < len(阶段_5_21速死_w_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_w_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_w_不变[i][j]
			if hp[enemy] > 0 {
				damage := rand.attack(def[enemy],hp[enemy])
				//print(i,enemy,hp[enemy],damage,hp[enemy] - damage)
				hp[enemy] = hp[enemy] - damage
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_1(rand Random, hp [31]int){
	atomic.AddInt32(&统计_阶段_5_1_分支, 1)

	run_part_5_1_双持_不变(rand,hp)
	run_part_5_1_燃烧_不变(rand,hp)
	run_part_5_1_w_不变(rand,hp)

	run_part_5_1_双持_变(rand,hp)
	run_part_5_1_燃烧_变(rand,hp)
	run_part_5_1_w_变(rand,hp)
}

func run_part_5_1_双持_不变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双

	for i := 0 ; i < 4 ; i ++ {
		for j := 0 ; j < len(阶段_5_1_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_1_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	hp[21] = hp[21] - 1000
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000

	for i := 4 ; i < 5 ; i ++ {
		select_enemy := -1
		for j := 0 ; j < len(阶段_5_1_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_1_双持_不变[i][j]
			if hp[enemy] > 0 {
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 {
			// 没有这种情况
			atomic.AddInt32(&TODO_T4, 1)
		}
	}

	if hp[23] > 0 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 5 ; i < len(阶段_5_1_双持_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_1_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_1_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}
func run_part_5_1_燃烧_不变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双

	for j := 0 ; j < len(阶段_5_1_燃烧_不变[0]) ; j ++ {
		enemy := 阶段_5_1_燃烧_不变[0][j]
		if hp[enemy] > 0 {
			hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
			break
		}
	}

	hp[21] = hp[21] - 1000
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000

	for i := 1 ; i < 5 ; i ++ {
		select_enemy := -1
		for j := 0 ; j < len(阶段_5_1_燃烧_不变[i]) ; j ++ {
			enemy := 阶段_5_1_燃烧_不变[i][j]
			if hp[enemy] > 0 {
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 {
			// 没有这种情况
			atomic.AddInt32(&TODO_T4, 1)
		}
	}
	
	if hp[23] > 0 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 5 ; i < len(阶段_5_1_燃烧_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_1_燃烧_不变[i]) ; j ++ {
			if i == 32 {
				// 只吃一个暴击的情况
				// rand.NextDouble()
			}
			enemy := 阶段_5_1_燃烧_不变[i][j]
			if hp[enemy] > 0 {
				damage := rand.attack(def[enemy],hp[enemy])
				//print(i,enemy,hp[enemy],damage,hp[enemy] - damage)
				hp[enemy] = hp[enemy] - damage
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_1_w_不变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双
	for i := 0 ; i < len(阶段_5_1_M) ; i ++ {
		for j := 0 ; j < 4 ; j ++ {
			enemy := 阶段_5_1_M[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}
	if hp[21] > 0 {
		atomic.AddInt32(&统计_死亡_21, 1)
		return
	}
	if hp[22] > 0 {
		atomic.AddInt32(&统计_死亡_22, 1)
		return
	}
	if hp[23] > 1000 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}
	
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000
	hp[25] = hp[25] - 500

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_1_w_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_1_w_不变[i]) ; j ++ {
			enemy := 阶段_5_1_w_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_1_双持_变(rand Random, hp [31]int){
	// 0 - 3
	for i := 0 ; i <= 3 ; i ++ {
		for j := 0 ; j < 3 ; j ++ {
			enemy := 阶段_5_1_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	if hp[22] > 0 {
		atomic.AddInt32(&统计_死亡_22, 1)
		return
	}
	if hp[23] > 1000 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_1_双持_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_1_双持_变[i]) ; j ++ {
			enemy := 阶段_5_1_双持_变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
	
}
func run_part_5_1_燃烧_变(rand Random, hp [31]int){
	hp[22] = hp[22] - rand.attack(def[22],hp[22])
	if hp[22] < 1000 {
		run_part_5_1_燃烧_变_22(rand,hp,2)
	}
	hp[22] = hp[22] - rand.attack(def[22],hp[22])
	if hp[22] <= 0 {
		return
	}
	if hp[22] < 1000 {
		run_part_5_1_燃烧_变_22(rand,hp,1)
	}
	hp[22] = hp[22] - rand.attack(def[22],hp[22])
	if hp[22] <= 0 {
		return
	}
	if hp[22] < 1000 {
		run_part_5_1_燃烧_变_22(rand,hp,0)
	}
}

func run_part_5_1_燃烧_变_22(rand Random, hp [31]int, times int){
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000
	for i := 1 ; i <= times ; i ++ {
		hp[23] = hp[23] - rand.attack(def[23],hp[23])
	}
	if hp[23] > 0 {
		hp[23] = hp[23] - rand.attack(def[23],hp[23])
	}else{
		hp[24] = hp[24] - rand.attack(def[24],hp[24])
	}
	if hp[23] > 0{
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_1_燃烧_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_1_燃烧_变[i]) ; j ++ {
			enemy := 阶段_5_1_燃烧_变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_1_w_变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双
	for i := 0 ; i < len(阶段_5_1_M) ; i ++ {
		for j := 0 ; j < 4 ; j ++ {
			enemy := 阶段_5_1_M[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}
	if hp[21] > 0 {
		atomic.AddInt32(&统计_死亡_21, 1)
		return
	}
	if hp[22] > 0 {
		atomic.AddInt32(&统计_死亡_22, 1)
		return
	}
	if hp[23] > 1000 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}
	
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000
	hp[25] = hp[25] - 500

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_1_w_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_1_w_变[i]) ; j ++ {
			enemy := 阶段_5_1_w_变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_2(rand Random, hp [31]int){
	atomic.AddInt32(&统计_阶段_5_2_分支, 1)
	run_part_5_2_双持_不变(rand,hp)
	run_part_5_2_燃烧_不变(rand,hp)
	run_part_5_2_w_不变(rand,hp)

	run_part_5_2_双持_变(rand,hp)
	run_part_5_2_燃烧_变(rand,hp)
	run_part_5_2_w_变(rand,hp)
}

func run_part_5_2_w_不变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双
	for i := 0 ; i < len(阶段_5_2_M) ; i ++ {
		for j := 0 ; j < len(阶段_5_2_M[i]) ; j ++ {
			enemy := 阶段_5_2_M[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}
	if hp[21] > 0 {
		atomic.AddInt32(&统计_死亡_21, 1)
		return
	}
	if hp[22] > 0 {
		atomic.AddInt32(&统计_死亡_22, 1)
		return
	}
	if hp[23] > 1000 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}
	
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000
	hp[25] = hp[25] - 500

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_2_w_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_2_w_不变[i]) ; j ++ {
			enemy := 阶段_5_2_w_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_2_双持_不变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双

	for i := 0 ; i < 3 ; i ++ {
		for j := 0 ; j < len(阶段_5_2_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_2_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	hp[21] = hp[21] - 1000
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000

	for i := 3 ; i < 5 ; i ++ {
		select_enemy := -1
		for j := 0 ; j < len(阶段_5_2_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_2_双持_不变[i][j]
			if hp[enemy] > 0 {
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 {
			// 没有这种情况
			atomic.AddInt32(&TODO_T4, 1)
		}
	}
	
	if hp[23] > 0 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 5 ; i < len(阶段_5_2_双持_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_2_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_2_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}
func run_part_5_2_燃烧_不变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双

	hp[21] = hp[21] - 1000
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000

	for i := 0 ; i < 5 ; i ++ {
		select_enemy := -1
		for j := 0 ; j < len(阶段_5_2_燃烧_不变[i]) ; j ++ {
			enemy := 阶段_5_2_燃烧_不变[i][j]
			if hp[enemy] > 0 {
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 {
			// 没有这种情况
			atomic.AddInt32(&TODO_T4, 1)
		}
	}
	
	if hp[23] > 0 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 5 ; i < len(阶段_5_2_燃烧_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_2_燃烧_不变[i]) ; j ++ {
			enemy := 阶段_5_2_燃烧_不变[i][j]
			if hp[enemy] > 0 {
				damage := rand.attack(def[enemy],hp[enemy])
				//print(enemy,damage,hp[enemy],hp[enemy]-damage)
				hp[enemy] = hp[enemy] - damage
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_2_双持_变(rand Random, hp [31]int){
	// 查无此轴
	return
}
func run_part_5_2_燃烧_变(rand Random, hp [31]int){
	for i := 0 ; i <= 2 ; i ++ {
		for j := 0 ; j < len(阶段_5_2_燃烧_不变[i]) ; j ++ {
			enemy := 阶段_5_2_燃烧_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if hp[22] <= 0 {
			return
		}
		if hp[22] <= 1000 {
			run_part_5_2_燃烧_变_22(rand,hp,2 - i)
		}
	}
}

func run_part_5_2_燃烧_变_22(rand Random, hp [31]int, times int){

	hp[21] = hp[21] - 1000
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000

	for i := 0 ; i < times ; i ++ {
		hp[23] = hp[23] - rand.attack(def[23],hp[23])
	}
	if hp[23] > 0 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_2_燃烧_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_2_燃烧_变[i]) ; j ++ {
			enemy := 阶段_5_2_燃烧_变[i][j]
			if hp[enemy] > 0 {
				damage := rand.attack(def[enemy],hp[enemy])
				//print(enemy,damage,hp[enemy],hp[enemy]-damage)
				hp[enemy] = hp[enemy] - damage
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_part_5_2_w_变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双
	for i := 0 ; i < len(阶段_5_2_M) ; i ++ {
		for j := 0 ; j < len(阶段_5_2_M[i]) ; j ++ {
			enemy := 阶段_5_2_M[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}
	if hp[21] > 0 {
		atomic.AddInt32(&统计_死亡_21, 1)
		return
	}
	if hp[22] > 0 {
		atomic.AddInt32(&统计_死亡_22, 1)
		return
	}
	if hp[23] > 1000 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}
	
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000
	hp[25] = hp[25] - 500

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_2_w_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_2_w_变[i]) ; j ++ {
			enemy := 阶段_5_2_w_变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_dead_21_双持_不变(rand Random, hp [31]int){
	// M 22 双 23 投 24 双
	for i := 0 ; i < 4 ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000

	for i := 4 ; i < 6 ; i ++ {
		select_enemy := -1
		for j := 0 ; j < len(阶段_5_21速死_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_双持_不变[i][j]
			if hp[enemy] > 0 {
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 {
			// 没有这种情况
			atomic.AddInt32(&TODO_T4, 1)
		}
	}

	if hp[23] > 0 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 6 ; i < len(阶段_5_21速死_双持_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}
func run_dead_21_燃烧_不变(rand Random, hp [31]int){
	hp[22] = hp[22] - rand.attack(def[22],hp[22])

	// M 22 双 23 投 24 双

	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000

	for i := 1 ; i < 6 ; i ++ {
		select_enemy := -1
		for j := 0 ; j < len(阶段_5_21速死_燃烧_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_燃烧_不变[i][j]
			if hp[enemy] > 0 {
				select_enemy = enemy
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
		if select_enemy == -1 {
			// 22 23 早死
			run_dead_23(rand,hp)
			return
		}
	}
	
	if hp[23] > 0 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 6 ; i < len(阶段_5_21速死_燃烧_不变) ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_燃烧_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_燃烧_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_dead_21_双持_变(rand Random, hp [31]int){

	for i := 0 ; i < 4 ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_双持_不变[i]) ; j ++ {
			enemy := 阶段_5_21速死_双持_不变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	if hp[23] <= 1000 {
		run_dead_21_双持_变_22(rand,hp,1)
	}

	for j := 0 ; j < len(阶段_5_21速死_双持_不变[4]) ; j ++ {
		enemy := 阶段_5_21速死_双持_不变[4][j]
		if hp[enemy] > 0 {
			hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
			break
		}
	}

	if hp[23] <= 1000 {
		run_dead_21_双持_变_22(rand,hp,0)
	}

}

func run_dead_21_双持_变_22(rand Random, hp [31]int,times int){
	hp[21] = hp[21] - 1000
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000

	atomic.AddInt32(&统计_阶段_w, 1)
	
	for i := 1 - times ; i < len(阶段_5_21速死_双持_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_双持_变[i]) ; j ++ {
			enemy := 阶段_5_21速死_双持_变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_dead_21_燃烧_变(rand Random, hp [31]int){
	hp[22] = hp[22] - rand.attack(def[22],hp[22])
	if hp[22] < 1000 {
		run_dead_21_燃烧_变_22(rand,hp,2)
	}
	hp[22] = hp[22] - rand.attack(def[22],hp[22])
	if hp[22] <= 0 {
		return
	}
	if hp[22] < 1000 {
		run_dead_21_燃烧_变_22(rand,hp,1)
	}
	hp[22] = hp[22] - rand.attack(def[22],hp[22])
	if hp[22] <= 0 {
		return
	}
	if hp[22] < 1000 {
		run_dead_21_燃烧_变_22(rand,hp,0)
	}
}
func run_dead_21_燃烧_变_22(rand Random, hp [31]int, times int){
	hp[21] = hp[21] - 1000
	hp[22] = hp[22] - 1000
	hp[23] = hp[23] - 1000

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 2 - times ; i < len(阶段_5_21速死_燃烧_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_燃烧_变[i]) ; j ++ {
			enemy := 阶段_5_21速死_燃烧_变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 21 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_dead_21_w_变(rand Random, hp [31]int){
	for i := 0 ; i < 6 ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_w_变[i]) ; j ++ {
			enemy := 阶段_5_21速死_w_变[i][j]
			if hp[enemy] > 0 {
				hp[enemy] = hp[enemy] - rand.attack(def[enemy],hp[enemy])
				break
			}
		}
	}

	if hp[22] > 0 {
		atomic.AddInt32(&统计_死亡_22, 1)
		return
	}
	if hp[23] > 1000 {
		atomic.AddInt32(&统计_死亡_23, 1)
		return
	}

	hp[23] = hp[23] - 1000
	hp[24] = hp[24] - 1000
	hp[25] = hp[25] - 500

	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 6 ; i < len(阶段_5_21速死_w_变) ; i ++ {
		for j := 0 ; j < len(阶段_5_21速死_w_变[i]) ; j ++ {
			enemy := 阶段_5_21速死_w_变[i][j]
			if hp[enemy] > 0 {
				damage := rand.attack(def[enemy],hp[enemy])
				//print(i,enemy,hp[enemy],damage,hp[enemy] - damage)
				hp[enemy] = hp[enemy] - damage
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func run_dead_23(rand Random, hp [31]int){

	atomic.AddInt32(&统计_速死_23, 1)
	atomic.AddInt32(&统计_阶段_w, 1)

	for i := 0 ; i < len(阶段_5_23速死) ; i ++ {
		for j := 0 ; j < len(阶段_5_23速死[i]) ; j ++ {
			enemy := 阶段_5_23速死[i][j]
			if hp[enemy] > 0 {
				damage := rand.attack(def[enemy],hp[enemy])
				hp[enemy] = hp[enemy] - damage
				break
			}
		}
	}

	stat(&rand,&hp)

	for i := 24 ; i <= 30 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&DONE, 1)
	print("done",runFuncName(),rand.seed)
}

func stat(rand *Random,hp *[31]int){
	//return
	if hp[25] > 0 {
		return
	}

	for i := 1 ; i <= 27 ; i ++ {
		if hp[i] > 0 {
			return
		}
	}

	atomic.AddInt32(&统计_死亡_w, 1)

	if hp[29] > 0 && hp[28] <= 0 {
		atomic.AddInt32(&统计_流未死, 1)
		//print("流未死. seed:",rand.seed,runFuncName3(),"w hp:",hp[25],"流 hp:",hp[29],"盾 hp:",hp[28])
	}
	if hp[28] > 0 && hp[29] <= 0 {
		atomic.AddInt32(&统计_盾未死, 1)
		//print("盾未死. seed:",rand.seed,runFuncName3(),"w hp:",hp[25],"流 hp:",hp[29],"盾 hp:",hp[28])
	}
}

func runFuncName() string{
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func runFuncName3() string{
	pc := make([]uintptr,1)
	runtime.Callers(3,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}