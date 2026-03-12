package main

import (
	"math/rand"
	"slices"

	. "github.com/loom-go/loom/components"
)

var todos, TodoListContext = NewContext[*TodoStore](nil)

type Todo struct {
	id    uint32
	title string
	done  *Writable[bool]
}

func NewTodo(title string) Todo {
	return Todo{
		id:    rand.Uint32(),
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

func (t *TodoStore) Add(title string) uint32 {
	todo := NewTodo(title)
	t.list.Set(append(t.list.Get(), &todo))
	return todo.id
}

func (t *TodoStore) Toggle(id uint32) {
	for _, todo := range t.list.Get() {
		if todo.id == id {
			todo.Toggle()
			break
		}
	}
}

func (t *TodoStore) Remove(id uint32) {
	t.list.Set(slices.DeleteFunc(t.list.Get(), func(t *Todo) bool {
		return t.id == id
	}))
}
