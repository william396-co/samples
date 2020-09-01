package models

import (
	"../util"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var DefaultTaskList *TaskManager

type Task struct {
	ID    int64  // Unique identifier
	Title string // Description
	Done  bool   // Is this task done?
}

func (this *Task) String() string {
	return fmt.Sprintf("[ID: %d, Title: %s ,Done:%t]",this.ID,this.Title, this.Done)
}

func (this *Task) TableName() string  {
	return TableName("Task")
}

func (this* Task) Save()  {
	o := orm.NewOrm()
	res, err := o.Raw(fmt.Sprintf("INSERT %s values(%d,'%s',%d)",
		this.TableName(), this.ID, this.Title, util.Btoi(this.Done))).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		logs.Debug("mysql row affected nums: ", num)
		logs.Info("Save Task: %s",this.String())
	} else {
		logs.Error(err)
	}
}

func (this *Task) Update() {
	o := orm.NewOrm()

	r :=  o.Raw(fmt.Sprintf("UPDATE %s SET title = '%s',done = %d WHERE id = %d",
		this.TableName(), this.Title, util.Btoi(this.Done),this.ID))
	if r != nil {
		logs.Debug(r)
		res, err := r.Exec()
		if err == nil {
			num, _ := res.RowsAffected()
			logs.Debug("mysql row affected nums: ", num)
			logs.Info("Update Task: %s",this.String())
		} else {
			logs.Error(err)
		}
	}
}

// NewTask creates a new task given a title, that can't be empty.
func NewTask(title string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("empty title")
	}
	return &Task{0, title, false}, nil
}

// TaskManager manages a list of tasks in memory.
type TaskManager struct {
	tasks  []*Task
	lastID int64
}

// NewTaskManager returns an empty TaskManager.
func NewTaskManager() *TaskManager {
	return &TaskManager{}
}

type Item struct{
	Id int64
	Title string
	Done int
}

func (this *TaskManager) loaddata() {
	o := orm.NewOrm()
	var lists []Item
	num, err := o.Raw(fmt.Sprintf("SELECT id,title,done FROM %s ORDER BY id",
		TableName("Task"))).QueryRows(&lists)
	if err == nil && num > 0 {
		for _, item := range lists {
			if this.lastID < item.Id{
				this.lastID = item.Id
			}
			this.tasks = append(this.tasks, &Task{item.Id, item.Title, util.Itob(item.Done)})
		}
		logs.Info("loaddata todo count: %d, max taskID: %d",num,this.lastID)
	}
}

// Save saves the given Task in the TaskManager.
func (m *TaskManager) Save(task *Task) error {
	if task.ID == 0 {
		m.lastID++
		task.ID = m.lastID
		m.tasks = append(m.tasks, cloneTask(task))
		task.Save()
		return nil
	}

	for i, t := range m.tasks {
		if t.ID == task.ID {
		//m.tasks[i] = cloneTask(task)
			assignTask(m.tasks[i],task)
			m.tasks[i].Update()
			return nil
		}
	}
	return fmt.Errorf("unknown task")
}

func assignTask(dst *Task,src *Task)  {
	dst.Title = src.Title
	dst.Done = src.Done
}

// cloneTask creates and returns a deep copy of the given Task.
func cloneTask(t *Task) *Task {
	c := new(Task)
	c.ID = t.ID
	c.Title = t.Title
	c.Done = t.Done
	return c
}

// All returns the list of all the Tasks in the TaskManager.
func (m *TaskManager) All() []*Task {
	return m.tasks
}

// Find returns the Task with the given id in the TaskManager and a boolean
// indicating if the id was found.
func (m *TaskManager) Find(ID int64) (*Task, bool) {
	for _, t := range m.tasks {
		if t.ID == ID {
			return t, true
		}
	}
	return nil, false
}

func init() {
	DefaultTaskList = NewTaskManager()
	DefaultTaskList.loaddata()
}
