package core

import (
	"fmt"
	"log"
)

// QuestCommandErrorBase forms basis for our error handling
func QuestCommandErrorBase(cmd string, qs *Quest, qt *QuestTask) string {
	return fmt.Sprintf("Command '%s' failed at Quest '%s':'%s'(%d): ", cmd, qs.name, qt.Name, qt.ProgramCounter)
}

// QuestCommandErrorArgCount happens on arg count mismatch
func QuestCommandErrorArgCount(cmd string, qs *Quest, qt *QuestTask, has, need int) bool {
	log.Printf("%s needs '%d' arguments, got: '%d'", QuestCommandErrorBase(cmd, qs, qt), need, has)
	return false
}

// QuestCommandErrorDivideByZero happens if we divide by zero
func QuestCommandErrorDivideByZero(cmd string, qs *Quest, qt *QuestTask) bool {
	log.Printf("%s division by zero", QuestCommandErrorBase(cmd, qs, qt))
	return false
}

// QuestCommandErrorThing happens when X is expected
func QuestCommandErrorThing(cmd, thing string, qs *Quest, qt *QuestTask, resName string) bool {
	log.Printf("%s %s '%s' could not be found", QuestCommandErrorBase(cmd, qs, qt), thing, resName)
	return false
}

// QuestCommandErrorArgType happens when X is of wrong data type
func QuestCommandErrorArgType(cmd string, qs *Quest, qt *QuestTask, argName, has, need string) bool {
	log.Printf("%s argument '%s' has to be '%s', got: '%s'", QuestCommandErrorBase(cmd, qs, qt), argName, need, has)
	return false
}

// QuestCommandErrorArgComp happens when the comparator is invalid
func QuestCommandErrorArgComp(cmd string, qs *Quest, qt *QuestTask, argName string) bool {
	log.Printf("%s argument has to be either 'above,below,equals,!equals', got: '%s'", QuestCommandErrorBase(cmd, qs, qt), argName)
	return false
}

// QuestCommandErrorEventArgsEmpty happens if event argument stack is being popped
// while already being empty
func QuestCommandErrorEventArgsEmpty(cmd string, qs *Quest, qt *QuestTask) bool {
	log.Printf("%s event's arg stack is already empty!", QuestCommandErrorBase(cmd, qs, qt))
	return false
}
