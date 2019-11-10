package core

import (
	"log"
)

var (
	stepCounter = 0
	MaxQuests   = 5
)

type QuestCommandTable func(qs *Quest, qt *QuestTask, args []string) bool

type QuestManager struct {
	commands map[string]QuestCommandTable
	quests   []Quest
}

func MakeQuestManager() QuestManager {
	res := QuestManager{
		commands: map[string]QuestCommandTable{},
		quests:   []Quest{},
	}

	questInitBaseCommands(&res)

	return res
}

func (q *QuestManager) GetActiveQuests() []*Quest {
	qs := []*Quest{}

	for _, v := range q.quests {
		if v.state == QsInProgress && !v.RunsInBackground {
			qs = append(qs, &v)
		}
	}

	return qs
}

func (q *QuestManager) AddQuest(tplName string, details map[string]float64) (bool, string, int64) {
	qd := ParseQuest(tplName)

	if qd == nil {
		return false, "Quest template could not be found!", -1
	}

	if !qd.RunsInBackground && len(q.GetActiveQuests()) >= MaxQuests {
		return false, "Maximum number of quests has been reached!", -1
	}

	if details == nil {
		details = map[string]float64{}
	}

	tasks := []QuestTask{}

	for _, v := range qd.TaskDef {
		tasks = append(tasks, QuestTask{
			QuestTaskDef: v,
			variables:    map[string]QuestVar{},
		})
	}

	processedDetails := map[string]QuestVar{}

	for k, v := range details {
		processedDetails[k] = QuestVar{
			kind:  kindNumber,
			value: &QuestVarNumber{Value: v},
		}
	}

	tasks[0].variables = processedDetails

	qn := Quest{
		ID:       getNewID(),
		name:     tplName,
		QuestDef: *qd,
		state:    QsInProgress,
		timers:   map[string]QuestTimer{},
		stages:   map[int]QuestStage{},
		tasks:    tasks,
	}

	qn.activeQuestTask = &qn.tasks[0]

	for _, v := range qn.tasks {
		qn.SetVariable(v.Name, 0)
	}

	for qn.ProcessTask(q, &qn.tasks[0]) {
		// process the whole entry point
	}

	q.quests = append(q.quests, qn)

	log.Printf("Quest '%s' with title '%s' has been added!", tplName, qd.Title)

	return true, "", qn.ID
}

func (q *QuestManager) Reset() {
	q.quests = []Quest{}
}

func (q *QuestManager) RegisterCommand(name string, cb QuestCommandTable) {
	q.commands[name] = cb
}

func (q *QuestManager) DispatchCommand(qs *Quest, qt *QuestTask, name string, args []string) (bool, bool) {
	cmd, ok := q.commands[name]

	if ok {
		return cmd(qs, qt, args), false
	}

	log.Printf("Quest '%s' has unrecognized command: '%s'!\n", qs.name, name)
	return false, true
}

func (q *QuestManager) ProcessQuests() {
	for i := range q.quests {
		qs := &q.quests[i]

		if qs.state != QsInProgress {
			continue
		}

		qs.ProcessTimers()
		qs.ProcessTasks(q)
	}

	stepCounter++
}

func (q *QuestManager) CallEvent(id int64, eventName string, args []float64) {
	for i := range q.quests {
		v := &q.quests[i]

		if id != -1 && id != v.ID {
			continue
		}

		v.CallEvent(q, eventName, args)
	}
}

// todo
var (
	globalIDCounter int64
)

func getNewID() int64 {
	v := globalIDCounter
	globalIDCounter++
	return v
}
