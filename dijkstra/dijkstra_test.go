package dijkstra

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestNoPath(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4, -1)
}

func TestLoop(t *testing.T) {
	testSolution(t, BestPath{}, newErrLoop(2, 1), "testdata/J.txt", 0, 4, -1)
}

func TestCorrect(t *testing.T) {
	testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, -1)
	testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, -1)
}

func TestCorrectSolutionsAll(t *testing.T) {
	graph := NewGraph()
	//Add the 3 verticies
	graph.AddVertex(0)
	graph.AddVertex(1)
	graph.AddVertex(2)
	graph.AddVertex(3)

	//Add the arcs
	graph.AddArc(0, 1, 1)
	graph.AddArc(0, 2, 1)
	graph.AddArc(1, 3, 0)
	graph.AddArc(2, 3, 0)
	testGraphSolutionAll(t, BestPaths{BestPath{1, []int{0, 2, 3}}, BestPath{1, []int{0, 1, 3}}}, nil, *graph, 0, 3)
}

func TestCorrectAllLists(t *testing.T) {
	for i := 0; i <= 3; i++ {
		testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, i)
		testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, i)
	}
}

func TestCorrectAutoLargeList(t *testing.T) {
	g := NewGraph()
	for i := 0; i < 2000; i++ {
		v := g.AddNewVertex()
		v.AddArc(i+1, 1)
	}
	g.AddNewVertex()
	_, err := g.Shortest(0, 2000)
	testErrors(t, nil, err, "manual test")
}

var benchNames = []string{"github.com/RyanCarrier-ALL", "github.com/RyanCarrier", "github.com/ProfessorQ", "github.com/albertorestifo"}
var listNames = []string{"PQShort", "LLShort"}

func BenchmarkSetup(b *testing.B) {
	nodeIterations := 6
	nodes := 1
	for j := 0; j < nodeIterations; j++ {
		nodes *= 4
		b.Run("setup/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
			filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
			if _, err := os.Stat(filename); err != nil {
				g := Generate(nodes)
				err := g.ExportToFile(filename)
				if err != nil {
					log.Fatal(err)
				}
			}
			g, _ := Import(filename)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				g.setup(0, -1)
			}
		})
	}
}

func BenchmarkLists(b *testing.B) {
	nodeIterations := 6

	for i, n := range listNames {
		nodes := 1
		for j := 0; j < nodeIterations; j++ {
			nodes *= 4
			b.Run("Short/"+n+"/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
				benchmarkList(b, nodes, i)
			})
		}
	}
}

func benchmarkList(b *testing.B, nodes, list int) {
	filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
	if _, err := os.Stat(filename); err != nil {
		g := Generate(nodes)
		err := g.ExportToFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.setup(src, list)
		graph.postSetupEvaluate(src, dest)
	}
}

func BenchmarkAll(b *testing.B) {
	nodeIterations := 6
	for i, n := range benchNames {
		nodes := 1
		for j := 0; j < nodeIterations; j++ {
			nodes *= 4
			b.Run(n+"/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
				benchmarkAlt(b, nodes, i)
			})

		}
	}
	//Cleanup
	nodes := 1
	for j := 0; j < nodeIterations; j++ {
		nodes *= 4
		os.Remove("testdata/bench/" + strconv.Itoa(nodes) + ".txt")
	}
}

/*
//Mattomatics does not work.
func BenchmarkMattomaticNodes4(b *testing.B)    { benchmarkAlt(b, 4, 3) }
*/
func benchmarkAlt(b *testing.B, nodes, i int) {
	filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
	if _, err := os.Stat(filename); err != nil {
		g := Generate(nodes)
		err := g.ExportToFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	switch i {
	case 0:
		benchmarkRCall(b, filename)
	case 1:
		benchmarkRC(b, filename)
	default:
		b.Error("You're retarded")
	}
}

func benchmarkRC(b *testing.B, filename string) {
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Shortest(src, dest)
	}
}
func benchmarkRCall(b *testing.B, filename string) {
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.ShortestAll(src, dest)
	}
}

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int, list int) {
	var err error
	var graph Graph
	graph, err = Import(filename)
	if err != nil {
		t.Fatal(err, filename)
	}
	var got BestPath
	var gotAll BestPaths
	if list >= 0 {
		graph.setup(from, list)
		got, err = graph.postSetupEvaluate(from, to)
	} else {
		got, err = graph.Shortest(from, to)
	}
	testErrors(t, wanterr, err, filename)
	testResults(t, got, best, filename)
	if list >= 0 {
		graph.setup(from, list)
		gotAll, err = graph.postSetupEvaluateAll(from, to)
	} else {
		gotAll, err = graph.ShortestAll(from, to)
	}
	testErrors(t, wanterr, err, filename)
	if len(gotAll) == 0 {
		gotAll = BestPaths{BestPath{}}
	}
	testResults(t, gotAll[0], best, filename)
}

func testGraphSolutionAll(t *testing.T, best BestPaths, wanterr error, graph Graph, from, to int) {
	gotAll, err := graph.ShortestAll(from, to)

	testErrors(t, wanterr, err, "From graph")
	if len(gotAll) == 0 {
		gotAll = BestPaths{BestPath{}}
	}
	testResultsGraphAll(t, gotAll, best)
}

func testResultsGraphAll(t *testing.T, got, best BestPaths) {
	distmethod := "Shortest"

	if len(got) != len(best) {
		t.Error(distmethod, " amount of solutions incorrect\ngot: ", len(got), "\nwant: ", len(best))
		return
	}
	for i := range got {
		if got[i].Distance != best[i].Distance {
			t.Error(distmethod, " distance incorrect\ngot: ", got[i].Distance, "\nwant: ", best[i].Distance)
		}
	}
	for i := range got {
		found := false
		j := -1
		for j = range best {
			if reflect.DeepEqual(got[i].Path, best[j].Path) {
				//delete found result
				best = append(best[:j], best[j+1:]...)
				found = true
				break
			}
		}
		if found == false {
			t.Error(distmethod, " could not find path in solution\ngot:", got[i].Path)
		}
	}
}

func testResults(t *testing.T, got, best BestPath, filename string) {
	distmethod := "Shortest"

	if got.Distance != best.Distance {
		t.Error(distmethod, " distance incorrect\n", filename, "\ngot: ", got.Distance, "\nwant: ", best.Distance)
	}
	if !reflect.DeepEqual(got.Path, best.Path) {
		t.Error(distmethod, " path incorrect\n\n", filename, "got: ", got.Path, "\nwant: ", best.Path)
	}
}

func getKSolShort() BestPath {
	return BestPath{
		2,
		[]int{
			0, 3, 4,
		},
	}
}
