package agent

import (
	"arithmetic_operations/orchestrator/models"
	"arithmetic_operations/orchestrator/topostfix"
	"errors"
	"log"
	"log/slog"
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
	Id          int
	Task        *models.Expression
	IsCompleted bool
}

type Task struct {
	Expression *models.Expression
	Operation  []*models.Operation
}

func NewAgent(iid int) *Agent {
	return &Agent{
		Id:          iid,
		Task:        nil,
		IsCompleted: true,
	}
}

func (c *Calculator) CheckerForNewTasks(expressionUpdater func(expression *models.Expression) error) {
	var wg sync.WaitGroup

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if c.NumberOfTasks > 0 && len(IsNotBusy(c)) > 0 {
					for i := 0; i < c.NumberOfTasks; i++ {
						wg.Add(1)
						go func(i int) {
							defer wg.Done()
							SolveExpression(c, c.Tasks[i], expressionUpdater)

							c.Tasks = c.Tasks[1:]
							c.NumberOfTasks--
						}(i)
					}
				}
				slog.Debug("Checking for new task")
			}
		}
	}()
	wg.Wait()
}

func InitializeAgents(num int) (*Calculator, error) {
	var tasks []*Task
	if num <= 0 {
		return &Calculator{}, errors.New("num of agents should be bigger than 0")
	}
	calc := make([]*Agent, num)
	for i := 0; i < num; i++ {
		calc[i] = NewAgent(i)
	}

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
func CreateTask(calc *Calculator, expression *models.Expression, operations []*models.Operation) {
	task := &Task{Expression: expression, Operation: operations}
	calc.Tasks = append(calc.Tasks, task)
	calc.NumberOfTasks++
}

func SolveExpression(calc *Calculator, task *Task, expressionUpdater func(expression *models.Expression) error) {
	var solvedSubexpressions []models.SubExpression
	divisionByZero := false

	b := IsNotBusy(calc)
	if len(b) == 0 {
		return
	}
	b[0].Task = task.Expression
	b[0].IsCompleted = false

	postfixExpression := topostfix.ToPostfix(task.Expression.Expression)

	solvedSubexpressionChan := make(chan models.SubExpression)
	errorChan := make(chan error)
	for {
		subexpressions, expressionToInsert := topostfix.GetSubExpressions(postfixExpression)

		for _, subexpr := range subexpressions {
			go func(subexpr models.SubExpression) {
				solvedSubExpr, err := topostfix.CountSubExpressions(subexpr, task.Operation)
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
				divisionByZero = true
			}
		}

		postfixExpression = topostfix.InsertSubExpressions(solvedSubexpressions, expressionToInsert)
		solvedSubexpressions = nil
		if len(expressionToInsert) == 1 {
			if divisionByZero {
				task.Expression.Answer = ""
				task.Expression.Status = models.Invalid
			} else {
				task.Expression.Answer = postfixExpression
				task.Expression.Status = models.Completed
			}
			timeCompleted := time.Now()
			task.Expression.CompletedAt = &timeCompleted
			err := expressionUpdater(task.Expression)
			b[0].IsCompleted = true
			b[0].Task = nil
			log.Printf("task %d completed", task.Expression.Id)
			if err != nil {
				log.Println("error with updating database", err)
			}
			break
		}
	}
}
