package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Reactant struct {
	ChemName  string
	QuanInput int
}

type Formula struct {
	Name       string
	QuanOutput int
	Reactants  []Reactant
}

var (
	Formulas map[string]*Formula
	Stock    map[string]int
	OreUsed  int
)

func main() {
	Formulas = make(map[string]*Formula)
	Stock = make(map[string]int)
	buildFromFile("formulas.txt")
	Stock["ORE"] = 1000000000000
	total := 0
	for Stock["ORE"] > 0 {
		GetChemical(1, "FUEL", 1)
		total += 1
	}
	fmt.Println(OreUsed)
	fmt.Println(Stock["ORE"])
	fmt.Println(total - 1)
}

func GetChemical(depth int, ResultName string, Quan int) {
	//indent := strings.Repeat("\t", depth-1)
	//fmt.Printf("%sGet %d of %s\n", indent, Quan, ResultName)
	if stock, found := Stock[ResultName]; found {
		if stock >= Quan {
			stock -= Quan
			Stock[ResultName] = stock
			return
		} else {
			Quan -= stock
			Stock[ResultName] = 0
		}
	}
	//fmt.Printf("%sProduce %d of %s\n", indent, Quan, ResultName)
	if ResultName == "ORE" {
		OreUsed += Quan
		return
	}
	//indent = strings.Repeat("\t", depth)
	f, found := Formulas[ResultName]
	if !found {
		panic(fmt.Sprintf("Can't find formula for chemical %s", ResultName))
	}
	n := Quan / f.QuanOutput
	if Quan%f.QuanOutput != 0 {
		n++
	}
	for _, r := range f.Reactants {
		GetChemical(depth+1, r.ChemName, n*r.QuanInput)
	}
	Stock[ResultName] = n*f.QuanOutput - Quan
}

func buildFromFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var (
			resultQuan int
			resultName string
		)
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, "=>")
		inputs := strings.Split(parts[0], ", ")
		fmt.Sscanf(parts[1], "%d %s", &resultQuan, &resultName)
		AddFormula(resultQuan, resultName, inputs)
	}
}

func buildTest1() {
	AddFormula2(10, "A", 10, "ORE")
	AddFormula2(1, "B", 1, "ORE")
	AddFormula2(1, "C", 7, "A", 1, "B")
	AddFormula2(1, "D", 7, "A", 1, "C")
	AddFormula2(1, "E", 7, "A", 1, "D")
	AddFormula2(1, "FUEL", 7, "A", 1, "E")
}

func buildTest2() {
	AddFormula2(2, "A", 9, "ORE")
	AddFormula2(3, "B", 8, "ORE")
	AddFormula2(5, "C", 7, "ORE")
	AddFormula2(1, "AB", 3, "A", 4, "B")
	AddFormula2(1, "BC", 5, "B", 7, "C")
	AddFormula2(1, "CA", 4, "C", 1, "A")
	AddFormula2(1, "FUEL", 2, "AB", 3, "BC", 4, "CA")
}

// Inputs are ##, NNNN, ##, NNNN
func AddFormula2(ResultQuan int, ResultName string, inputs ...interface{}) {
	a := make([]string, 0)
	for j := 0; j < len(inputs); j += 2 {
		a = append(a, fmt.Sprintf("%d %s", inputs[j].(int), inputs[j+1].(string)))
	}
	AddFormula(ResultQuan, ResultName, a)
}

func AddFormula(ResultQuan int, ResultName string, inputs []string) {
	f := new(Formula)
	f.Name = ResultName
	f.QuanOutput = ResultQuan
	numR := len(inputs)
	f.Reactants = make([]Reactant, numR)
	for i := 0; i < numR; i++ {
		var (
			q int
			n string
		)
		fmt.Sscanf(inputs[i], "%d %s", &q, &n)
		f.Reactants[i] = Reactant{
			ChemName:  n,
			QuanInput: q,
		}
	}
	Formulas[ResultName] = f
}
