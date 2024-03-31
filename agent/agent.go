package agent

import (
	"arithmetic_operations/orchestrator/models"
	"arithmetic_operations/orchestrator/topostfix"
	"errors"
	"log"
	"log/slog"
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

func NewAgent(n int) *Agent {
	return &Agent{
		Id:          n,
		Task:        nil,
		IsCompleted: true,
	}
}

func (c *Calculator) CheckerForNewTasks(expressionUpdater func(expression *models.Expression) error) {

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				for c.NumberOfTasks > 0 {
					agent := c.getFreeAgent()
					if agent != nil {
						go c.solveExpression(agent, c.Tasks[0], expressionUpdater)
						c.Tasks = c.Tasks[1:]
						c.NumberOfTasks--
						break
					} else {
						slog.Info("No free agents available. Waiting for a free agent...")
						break
					}

				}
				slog.Debug("Checking for new task")
			}
		}
	}()
}

func (c *Calculator) newTaskChannel() chan *Task {
	taskChannel := make(chan *Task)
	go func() {
		for {
			if len(c.Tasks) > 0 {
				task := c.Tasks[0]
				c.Tasks = c.Tasks[1:]
				c.NumberOfTasks--
				taskChannel <- task
			}
		}
	}()
	return taskChannel
}

func (c *Calculator) getFreeAgent() *Agent {
	for _, agent := range c.Agents {
		if agent.IsCompleted {
			return agent
		}
	}
	return nil
}

func (c *Calculator) addTask(task *Task) {
	c.Tasks = append(c.Tasks, task)
	c.NumberOfTasks++
}

func (c *Calculator) AddAgent() {
	agent := NewAgent(c.NumberOfAgents)
	c.NumberOfAgents++
	c.Agents = append(c.Agents, agent)
}

func (c *Calculator) RemoveAgent() error {
	if len(c.Agents) > 2 {
		c.Agents = c.Agents[:c.NumberOfAgents-1]
		c.NumberOfAgents--
		return nil
	} else {
		return errors.New("you have only one agent")
	}

}

func (c *Calculator) CreateTask(expression *models.Expression, operations []*models.Operation) {
	task := &Task{Expression: expression, Operation: operations}
	c.Tasks = append(c.Tasks, task)
	c.NumberOfTasks++
}

func (c *Calculator) solveExpression(agent *Agent, task *Task, expressionUpdater func(expression *models.Expression) error) {
	var solvedSubexpressions []models.SubExpression
	divisionByZero := false

	agent.IsCompleted = false
	agent.Task = task.Expression

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
			agent.IsCompleted = true
			agent.Task = nil
			log.Printf("task %d completed", task.Expression.Id)
			if err != nil {
				log.Println("error with updating database", err)
			}
			break
		}
	}
}

func InitializeAgents(num int) (*Calculator, error) {
	if num <= 0 {
		return nil, errors.New("num of agents should be bigger than 0")
	}

	var agents []*Agent
	for i := 0; i < num; i++ {
		agents = append(agents, NewAgent(i))
	}

	return &Calculator{
		Agents:         agents,
		NumberOfAgents: num,
	}, nil
}
