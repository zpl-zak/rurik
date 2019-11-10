package core

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/system"
)

const (
	QsInProgress = iota
	QsFinished
	QsFailed
)

type Quest struct {
	ID               int64
	name             string
	runsInBackground bool // used by events, they don't count as an actual Quest
	state            int
	timers           map[string]QuestTimer
	stages           map[int]QuestStage
	tasks            []QuestTask
	activeQuestTask  *QuestTask
	QuestDef
}

const (
	kindNumber = iota
	kindVector
)

type QuestVarData interface {
	Str() string
}

type QuestVar struct {
	kind  int
	value QuestVarData
}

type QuestTask struct {
	variables map[string]QuestVar
	QuestTaskDef
}

type QuestTimer struct {
	time     float32
	duration float32
}

type QuestStage struct {
	step  string
	state int
}

func (qs *Quest) Printf(qt *QuestTask, format string, args ...interface{}) {
	log.Printf("Quest '%s':'%s'(%d): %s", qs.name, qt.Name, qt.ProgramCounter, fmt.Sprintf(format, args...))
}

func (qs *Quest) GetResource(id string) (*QuestResource, bool) {
	val, err := strconv.Atoi(id)

	if err != nil {
		return nil, false
	}

	res, ok := qs.Resources[val]

	if !ok {
		return nil, false
	}

	return &res, true
}

func (qs *Quest) GetNumberOrVariable(arg string) (float64, bool) {
	val, err := strconv.ParseFloat(arg, 64)

	if err != nil {
		exprStr := qs.ResolveVariables(arg)
		expr, err := govaluate.NewEvaluableExpression(exprStr)

		if err != nil {
			return 0, false
		}

		res, err := expr.Evaluate(nil)

		if err != nil {
			return 0, false
		}

		return res.(float64), true
	}

	return val, true
}

func (qs *Quest) GetRelevantVariables() (a map[string]QuestVar) {
	if qs.activeQuestTask == &qs.tasks[0] {
		return qs.tasks[0].variables
	}

	a = map[string]QuestVar{}

	for k, v := range qs.tasks[0].variables {
		a[k] = v
	}

	for k, v := range qs.activeQuestTask.variables {
		a[k] = v
	}

	return a
}

func (qs *Quest) ProcessText(content string) string {
	for k, v := range qs.GetRelevantVariables() {
		content = strings.ReplaceAll(content, fmt.Sprintf("%%%s%%", k), v.value.Str())
	}

	return content
}

func (qs *Quest) ResolveVariables(expr string) string {
	for k, v := range qs.GetRelevantVariables() {
		expr = strings.ReplaceAll(expr, k, v.value.Str())
	}

	return expr
}

func (qs *Quest) GetTaskOverride(name string) *QuestTask {
	aq := qs.activeQuestTask

	_, ok := qs.tasks[0].variables[name]

	if ok {
		aq = &qs.tasks[0]
	}

	return aq
}

func (qs *Quest) SetVariable(name string, val float64) {
	qs.GetTaskOverride(name).variables[name] = QuestVar{
		kind:  kindNumber,
		value: &QuestVarNumber{Value: val},
	}
}

func (qs *Quest) SetVector(name string, val rl.Vector2) {
	qs.GetTaskOverride(name).variables[name] = QuestVar{
		kind:  kindVector,
		value: &QuestVarVector{Value: val},
	}
}

func (qs *Quest) GetVariable(name string) (float64, bool) {
	vars := qs.GetRelevantVariables()

	val, ok := vars[name]

	if !ok {
		return 0, false
	}

	return val.value.(*QuestVarNumber).Value, true
}

func (qs *Quest) GetVector(name string) (rl.Vector2, bool) {
	vars := qs.GetRelevantVariables()

	val, ok := vars[name]

	if !ok {
		return rl.Vector2{}, false
	}

	return val.value.(*QuestVarVector).Value, true
}

func (qs *Quest) ProcessTimers() {
	for k, v := range qs.timers {
		if v.time >= 0 {
			v.time -= system.FrameTime

			if v.time < 0 {
				v.time = 0
			}

			qs.timers[k] = v
		}

		qs.SetVariable(k, float64(RoundFloatToInt32(v.time)))
	}
}

func (qs *Quest) ProcessTask(q *QuestManager, qt *QuestTask) bool {
	if qt.ProgramCounter >= len(qt.Commands) {
		qt.IsDone = true
		return false
	}

	qs.activeQuestTask = qt

	qs.ProcessVariables()

	cmd := qt.Commands[qt.ProgramCounter]
	ok, err := q.DispatchCommand(qs, qt, cmd.Name, cmd.Args)

	if err {
		qt.IsDone = true
		return false
	}

	if !ok {
		return false
	}

	qt.ProgramCounter++
	return true
}

func (qs *Quest) ProcessTasks(q *QuestManager) {
	for i := range qs.tasks {
		v := &qs.tasks[i]

		if v.IsDone || v.IsEvent {
			continue
		}

		for qs.ProcessTask(q, v) {
			// task is being processed
		}

		state := 0

		if v.IsDone {
			state = 1
		}

		qs.SetVariable(v.Name, float64(state))
	}
}

func (qs *Quest) CallEvent(q *QuestManager, name string, args []float64) {
	for i := range qs.tasks {
		v := &qs.tasks[i]

		if v.Name != name {
			continue
		}

		v.IsDone = false
		v.EventArgs = args[:]

		for qs.ProcessTask(q, v) {
			// task is being processed
		}
	}
}

func (qs *Quest) ProcessVariables() {
	qt := qs.activeQuestTask
	qs.activeQuestTask = &qs.tasks[0]
	qs.SetVariable("$random", float64(rand.Int()))
	qs.SetVariable("$frandom", rand.Float64())
	qs.SetVariable("$step", float64(stepCounter))
	qs.SetVariable("$time", float64(rl.GetTime()))

	// player
	qs.SetVector("$pc.position", LocalPlayer.Position)

	// user-defined vars
	if ProcessCustomVariables != nil {
		ProcessCustomVariables(qs)
	}

	qs.activeQuestTask = qt
}
