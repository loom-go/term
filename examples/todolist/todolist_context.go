package main

import (
	"slices"
	"time"

	. "github.com/AnatoleLucet/loom/components"
)

var todos, TodoListContext = NewContext[*TodoStore](nil)

type Todo struct {
	id    int
	title string
	done  *Writable[bool]
}

func NewTodo(title string) Todo {
	return Todo{
		id:    time.Now().Nanosecond(),
		title: title,
		done:  NewWritable(false),
	}
}

func (t *Todo) Done() bool {
	return t.done.Get()
}

func (t *Todo) Toggle() {
	t.done.Set(!t.done.Get())
}

type TodoStore struct {
	list *Writable[[]*Todo]
}

func NewTodoStore() *TodoStore {
	return &TodoStore{NewWritable([]*Todo{})}
}

func (t *TodoStore) All() []*Todo {
	return t.list.Get()
}

func (t *TodoStore) Add(title string) int {
	todo := NewTodo(title)
	t.list.Set(append(t.list.Get(), &todo))
	return todo.id
}

func (t TodoStore) Toggle(id int) {
	for _, todo := range t.list.Get() {
		if todo.id == id {
			todo.Toggle()
			break
		}
	}
}

func (t TodoStore) Remove(id int) {
	t.list.Set(slices.DeleteFunc(t.list.Get(), func(t *Todo) bool {
		return t.id == id
	}))
}
