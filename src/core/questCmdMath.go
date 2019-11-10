package core

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/raylib-go/raymath"
)

func float64to32(x float64) float32 {
	return float32(x)
}

func questInitMathCommands(q *QuestManager) {
	q.RegisterCommand("vec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) != 1 {
			return QuestCommandErrorArgCount("vec", qs, qt, len(args), 1)
		}

		vecName := args[0]
		qs.SetVector(vecName, rl.Vector2{})

		qs.Printf(qt, "vector '%s' was declared!", vecName)
		return true
	})

	q.RegisterCommand("setvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("setvec", qs, qt, len(args), 3)
		}

		vecName := args[0]

		xI, _ := qs.GetNumberOrVariable(args[1])
		x := float64to32(xI)
		yI, _ := qs.GetNumberOrVariable(args[2])
		y := float64to32(yI)

		qs.SetVector(vecName, rl.NewVector2(x, y))

		qs.Printf(qt, "vector '%s' was set to [%f, %f]!", vecName, xI, yI)
		return true
	})

	q.RegisterCommand("copyvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 2 {
			return QuestCommandErrorArgCount("copyvec", qs, qt, len(args), 2)
		}

		vecName := args[0]
		rhsVecName := args[1]

		rhs, ok := qs.GetVector(rhsVecName)

		if !ok {
			return QuestCommandErrorThing("copyvec", "vector", qs, qt, rhsVecName)
		}

		qs.SetVector(vecName, rhs)
		return true
	})

	q.RegisterCommand("getvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("getvec", qs, qt, len(args), 3)
		}

		vecName := args[0]
		xName := args[1]
		yName := args[2]

		vec, ok := qs.GetVector(vecName)

		if !ok {
			return QuestCommandErrorThing("getvec", "vector", qs, qt, vecName)
		}

		if xName != "0" {
			qs.SetVariable(xName, float64(vec.X))
		}

		if yName != "0" {
			qs.SetVariable(yName, float64(vec.Y))
		}

		return true
	})

	q.RegisterCommand("addvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("addvec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhs, rhsFound := qs.GetVector(rhsVecName)

		if !lhsFound {
			return QuestCommandErrorThing("addvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("addvec", "vector", qs, qt, rhsVecName)
		}

		qs.SetVector(destVecName, rl.NewVector2(lhs.X+rhs.X, lhs.Y+rhs.Y))
		return true
	})

	q.RegisterCommand("addivec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("addivec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhsI, rhsFound := qs.GetNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return QuestCommandErrorThing("addivec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("addivec", "number", qs, qt, args[2])
		}

		qs.SetVector(destVecName, rl.NewVector2(lhs.X+rhs, lhs.Y+rhs))
		return true
	})

	q.RegisterCommand("subvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("subvec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhs, rhsFound := qs.GetVector(rhsVecName)

		if !lhsFound {
			return QuestCommandErrorThing("subvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("subvec", "vector", qs, qt, rhsVecName)
		}

		qs.SetVector(destVecName, rl.NewVector2(lhs.X-rhs.X, lhs.Y-rhs.Y))
		return true
	})

	q.RegisterCommand("subivec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("subivec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhsI, rhsFound := qs.GetNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return QuestCommandErrorThing("subivec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("subivec", "number", qs, qt, args[2])
		}

		qs.SetVector(destVecName, rl.NewVector2(lhs.X-rhs, lhs.Y-rhs))
		return true
	})

	q.RegisterCommand("divivec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("divivec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhsI, rhsFound := qs.GetNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return QuestCommandErrorThing("divivec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("divivec", "number", qs, qt, args[2])
		}

		if rhs == 0 {
			return QuestCommandErrorDivideByZero("divivec", qs, qt)
		}

		qs.SetVector(destVecName, rl.NewVector2(lhs.X/rhs, lhs.Y/rhs))
		return true
	})

	q.RegisterCommand("mulvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("mulvec", qs, qt, len(args), 3)
		}

		destVecName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhsI, rhsFound := qs.GetNumberOrVariable(args[2])
		rhs := float64to32(rhsI)

		if !lhsFound {
			return QuestCommandErrorThing("mulvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("mulvec", "number", qs, qt, args[2])
		}

		qs.SetVector(destVecName, rl.NewVector2(lhs.X*rhs, lhs.Y*rhs))
		return true
	})

	q.RegisterCommand("dotvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("dotvec", qs, qt, len(args), 3)
		}

		destName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhs, rhsFound := qs.GetVector(rhsVecName)

		if !lhsFound {
			return QuestCommandErrorThing("dotvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("dotvec", "number", qs, qt, args[2])
		}

		res := raymath.Vector2DotProduct(lhs, rhs)
		qs.SetVariable(destName, float64(res))
		return true
	})

	q.RegisterCommand("crossvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 3 {
			return QuestCommandErrorArgCount("crossvec", qs, qt, len(args), 3)
		}

		destName := args[0]
		lhsVecName := args[1]
		rhsVecName := args[2]

		lhs, lhsFound := qs.GetVector(lhsVecName)
		rhs, rhsFound := qs.GetVector(rhsVecName)

		if !lhsFound {
			return QuestCommandErrorThing("crossvec", "vector", qs, qt, lhsVecName)
		}

		if !rhsFound {
			return QuestCommandErrorThing("crossvec", "number", qs, qt, args[2])
		}

		res := raymath.Vector2CrossProduct(lhs, rhs)
		qs.SetVariable(destName, float64(res))
		return true
	})

	q.RegisterCommand("normvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 2 {
			return QuestCommandErrorArgCount("normvec", qs, qt, len(args), 2)
		}

		destName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.GetVector(lhsVecName)

		if !lhsFound {
			return QuestCommandErrorThing("normvec", "vector", qs, qt, lhsVecName)
		}

		raymath.Vector2Normalize(&lhs)
		qs.SetVector(destName, lhs)
		return true
	})

	q.RegisterCommand("flipvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 2 {
			return QuestCommandErrorArgCount("flipvec", qs, qt, len(args), 2)
		}

		destName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.GetVector(lhsVecName)

		if !lhsFound {
			return QuestCommandErrorThing("flipvec", "vector", qs, qt, lhsVecName)
		}

		qs.SetVector(destName, rl.NewVector2(lhs.Y, -lhs.X))
		return true
	})

	q.RegisterCommand("lenvec", func(qs *Quest, qt *QuestTask, args []string) bool {
		if len(args) < 2 {
			return QuestCommandErrorArgCount("lenvec", qs, qt, len(args), 2)
		}

		destName := args[0]
		lhsVecName := args[1]

		lhs, lhsFound := qs.GetVector(lhsVecName)

		if !lhsFound {
			return QuestCommandErrorThing("lenvec", "vector", qs, qt, lhsVecName)
		}

		qs.SetVariable(destName, float64(raymath.Vector2Length(lhs)))
		return true
	})
}
