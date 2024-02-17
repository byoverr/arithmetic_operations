package agent

import (
	"arithmetic_operations/orchestrator/models"
	"arithmetic_operations/orchestrator/topostfix"
	"errors"
	"log"
	"sync"
	"time"
)

type Calculator struct {
	Agents         []*Agent
	Tasks          []*Task
	NumberOfAgents int
	NumberOfTasks  int
}
type Agent struct {
	id          int
	mu          sync.Mutex
	task        *models.Expression
	IsCompleted bool
}

type Task struct {
	queue     *models.Expression
	operation []*models.Operation
}

func NewAgent(iid int) *Agent {
	return &Agent{
		id:          iid,
		mu:          sync.Mutex{},
		task:        nil,
		IsCompleted: true,
	}
}

func InitializeAgents(num int) (*Calculator, error) {
	if num <= 0 {
		return &Calculator{}, errors.New("num of agents should be bigger than 0")
	}
	calc := make([]*Agent, num)
	for i := 0; i < num; i++ {
		calc[i] = NewAgent(i)
	}
	tasks := make([]*Task, 1)

	return &Calculator{
		Agents:         calc,
		NumberOfAgents: num,
		Tasks:          tasks,
		NumberOfTasks:  0,
	}, nil
}

func AddAgent(calculator Calculator) error {
	agent := NewAgent(calculator.NumberOfAgents + 1)
	calculator.Agents = append(calculator.Agents, agent)
	return nil
}

func RemoveAgent(calculator Calculator) error {
	calculator.Agents = calculator.Agents[:calculator.NumberOfAgents-1]
	calculator.NumberOfAgents--
	return nil
}

func IsNotBusy(calculator *Calculator) []*Agent {
	var freeAgents []*Agent
	for _, i := range calculator.Agents {
		if i.IsCompleted {
			freeAgents = append(freeAgents, i)
		}
	}
	return freeAgents
}
func CreateTask(calc *Calculator, expression *models.Expression, operations []*models.Operation, expressionUpdater func(expression *models.Expression) error) {
	task := &Task{queue: expression, operation: operations}
	calc.Tasks = append(calc.Tasks, task)
	calc.NumberOfTasks++
	SolveExpression(calc, expression, operations, expressionUpdater)
}

func SolveExpression(calc *Calculator, expression *models.Expression, operations []*models.Operation, expressionUpdater func(expression *models.Expression) error) {
	var solvedSubexpressions []models.SubExpression
	b := IsNotBusy(calc)
	b[0].task = expression
	b[0].IsCompleted = false
	postfixExpression := topostfix.ToPostfix(expression.Expression)
	solvedSubexpressionChan := make(chan models.SubExpression)
	errorChan := make(chan error)
	for {
		subexpressions, expressionToInsert := topostfix.GetSubExpressions(postfixExpression)
		for _, subexpr := range subexpressions {
			go func(subexpr models.SubExpression) {
				solvedSubExpr, err := topostfix.CountSubExpressions(subexpr, operations)
				if err != nil {
					errorChan <- err
					return
				}
				solvedSubexpressionChan <- solvedSubExpr
			}(subexpr)
		}

		for range subexpressions {
			select {
			case solvedSubexpr := <-solvedSubexpressionChan:
				solvedSubexpressions = append(solvedSubexpressions, solvedSubexpr)
			case err := <-errorChan:
				log.Println("Error occurred while counting subexpressions:", err)
			}
		}
		postfixExpression = topostfix.InsertSubExpressions(solvedSubexpressions, expressionToInsert)
		if len(expressionToInsert) == 1 {
			expression.Answer = postfixExpression
			expression.Status = models.Completed
			timeCompleted := time.Now()
			expression.CompletedAt = &timeCompleted
			err := expressionUpdater(expression)
			if err != nil {
				log.Println("error with updating database", err)
			}
			// TODO: update database and sleep dont work
			break
		}
	}
}
