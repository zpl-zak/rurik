package core

import "log"

func questInitBaseCommands(q *QuestManager) {
	q.RegisterCommand("variable", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("variable", qs, qt, len(args), 1)
		}

		qs.SetVariable(args[0], 0)

		qs.Printf(qt, "variable '%s' was declared", args[0])

		return true
	})

	q.RegisterCommand("setvar", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 2 {
			return QuestCommandErrorArgCount("setvar", qs, qt, len(args), 2)
		}

		val, ok := qs.GetNumberOrVariable(args[1])

		if !ok {
			return QuestCommandErrorArgType("setvar", qs, qt, args[1], "string", "integer")
		}

		qs.SetVariable(args[0], val)

		qs.Printf(qt, "variable '%s' was set to: %f", args[0], val)

		return true
	})

	q.RegisterCommand("timer", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 2 {
			return QuestCommandErrorArgCount("timer", qs, qt, len(args), 2)
		}

		duration, ok := qs.GetNumberOrVariable(args[1])

		if !ok {
			return QuestCommandErrorArgType("timer", qs, qt, args[1], "string", "integer")
		}

		qs.timers[args[0]] = QuestTimer{
			time:     -1,
			duration: float32(duration),
		}

		qs.Printf(qt, "timer '%s' was declared with duration: %f", args[0], duration)

		return true
	})

	q.RegisterCommand("stage", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("stage", qs, qt, len(args), 1)
		}

		res, ok := qs.GetResource(args[0])

		if !ok {
			return QuestCommandErrorThing("stage", "resource", qs, qt, args[0])
		}

		stageID := atoiUnsafe(args[0])

		qs.stages[stageID] = QuestStage{
			step:  res.Content,
			state: QsInProgress,
		}

		qs.Printf(qt, "stage '%d' has been added!", stageID)

		return true
	})

	q.RegisterCommand("stdone", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("stdone", qs, qt, len(args), 1)
		}

		stageID := atoiUnsafe(args[0])
		sta, ok := qs.stages[stageID]

		if !ok {
			return QuestCommandErrorThing("stdone", "resource", qs, qt, args[0])
		}

		qs.Printf(qt, "stage '%d' has succeeded!", stageID)

		sta.state = QsFinished
		qs.stages[stageID] = sta

		return true
	})

	q.RegisterCommand("stfail", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("stfail", qs, qt, len(args), 1)
		}

		stageID := atoiUnsafe(args[0])
		sta, ok := qs.stages[stageID]

		if !ok {
			return QuestCommandErrorThing("stfail", "resource", qs, qt, args[0])
		}

		qs.Printf(qt, "stage '%d' has failed!", stageID)

		sta.state = QsFailed
		qs.stages[stageID] = sta

		return true
	})

	q.RegisterCommand("repeat", func(qs *Quest, qt *QuestTask, args []string) bool {
		qt.ProgramCounter = -1

		qs.Printf(qt, "repeating task '%s'!", qt.Name)

		return true
	})

	q.RegisterCommand("fire", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("fire", qs, qt, len(args), 1)
		}

		tm, ok := qs.timers[args[0]]

		if !ok {
			return QuestCommandErrorThing("fire", "timer", qs, qt, args[0])
		}

		qs.Printf(qt, "timer '%s' was fired!", args[0])
		tm.time = tm.duration
		qs.timers[args[0]] = tm

		return true
	})

	q.RegisterCommand("stop", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("stop", qs, qt, len(args), 1)
		}

		tm, ok := qs.timers[args[0]]

		if !ok {
			return QuestCommandErrorThing("stop", "timer", qs, qt, args[0])
		}

		qs.Printf(qt, "timer '%s' was stopped!", args[0])
		tm.time = -1
		qs.timers[args[0]] = tm

		return true
	})

	q.RegisterCommand("done", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("done", qs, qt, len(args), 1)
		}

		tm, ok := qs.timers[args[0]]

		if !ok {
			return QuestCommandErrorThing("done", "timer", qs, qt, args[0])
		}

		state := tm.time == 0

		if state {
			qs.Printf(qt, "timer '%s' is done!", args[0])
		}

		return state
	})

	q.RegisterCommand("finish", func(qs *Quest, qt *QuestTask, args []string) bool {
		qs.state = QsFinished

		qs.Printf(qt, "Quest '%s' has been finished!", qs.name)

		return true
	})

	q.RegisterCommand("fail", func(qs *Quest, qt *QuestTask, args []string) bool {
		qs.state = QsFailed

		qs.Printf(qt, "Quest '%s' has been failed!", qs.name)

		return true
	})

	q.RegisterCommand("pop", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("pop", qs, qt, len(args), 1)
		}

		if len(qt.EventArgs) == 0 {
			return QuestCommandErrorEventArgsEmpty("pop", qs, qt)
		}

		val := qt.EventArgs[0]
		qt.EventArgs = qt.EventArgs[1:]

		qs.SetVariable(args[0], val)

		qs.Printf(qt, "event pop value '%f' for: '%s'", val, args[0])

		return true
	})

	q.RegisterCommand("when", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 1 {
			return QuestCommandErrorArgCount("when", qs, qt, len(args), 1)
		}

		lhs, ok := qs.GetNumberOrVariable(args[0])

		if !ok {
			return QuestCommandErrorArgType("when", qs, qt, args[0], "string", "integer")
		}

		if len(args) == 1 {
			return lhs > 0
		}

		rhs, ok2 := qs.GetNumberOrVariable(args[2])

		if !ok2 {
			return QuestCommandErrorArgType("when", qs, qt, args[2], "string", "integer")
		}

		switch args[1] {
		case KwBelow:
			return lhs < rhs
		case KwAbove:
			return lhs > rhs
		case KwEquals:
			return lhs == rhs
		case KwNotEquals:
			return lhs != rhs
		case KwAnd:
			return (lhs != 0) && (rhs != 0)
		case KwOr:
			return (lhs != 0) || (rhs != 0)
		case KwXor:
			return ((lhs != 0) || (rhs != 0)) && !((lhs != 0) && (rhs != 0))
		default:
			return QuestCommandErrorArgComp("when", qs, qt, args[2])
		}
	})

	q.RegisterCommand("invoke", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 1 {
			return QuestCommandErrorArgCount("invoke", qs, qt, len(args), 1)
		}

		FireEvent(args[0], args[1:])
		return true
	})

	questInitCommands(q)

	if QuestInitCustomCommands != nil {
		log.Printf("Custom quest commands found, adding ...")
		QuestInitCustomCommands(q)
	}
}
