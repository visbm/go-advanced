package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap Heap
}

func NewScheduler() Scheduler {
	return Scheduler{
		heap: Heap{tasks: make([]Task, 0)},
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.heap.AddTask(task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.heap.ChangeTaskPriority(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	return s.heap.GetTask()
}

type Heap struct {
	tasks []Task
}

func (s *Heap) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
	s.shiftUp(len(s.tasks) - 1)
}

func (s *Heap) GetTask() Task {
	t := s.tasks[0]
	s.tasks = s.tasks[1:]
	s.shiftDown(0)
	return t
}

func (s *Heap) ChangeTaskPriority(taskID int, newPriority int) {
	index := s.findTask(taskID)
	if index == -1 {
		return
	}
	oldPriority := s.tasks[index].Priority
	s.tasks[index].Priority = newPriority

	if newPriority > oldPriority {
		s.shiftUp(index)
	} else if newPriority < oldPriority {
		s.shiftDown(index)
	}

}

func (s *Heap) shiftUp(index int) {
	if index >= len(s.tasks) || index < 0 {
		return
	}

	for index > 0 {
		parentIndex := s.parent(index)
		if parentIndex == -1 {
			break
		}

		if s.tasks[index].Priority > s.tasks[parentIndex].Priority {
			s.tasks[index], s.tasks[parentIndex] = s.tasks[parentIndex], s.tasks[index]
			index = parentIndex
		} else {
			break
		}
	}
}

func (s *Heap) shiftDown(index int) {
	if index >= len(s.tasks) || index < 0 {
		return
	}
	for index < len(s.tasks) {
		largest := index
		leftIndex := s.left(index)
		rightIndex := s.right(index)

		if leftIndex != -1 && s.tasks[leftIndex].Priority > s.tasks[largest].Priority {
			largest = leftIndex
		}

		if rightIndex != -1 && s.tasks[rightIndex].Priority > s.tasks[largest].Priority {
			largest = rightIndex
		}
		if largest == index {
			break
		}

		s.tasks[index], s.tasks[largest] = s.tasks[largest], s.tasks[index]
		index = largest

	}

}

func (s *Heap) findTask(taskID int) int {
	for i := range s.tasks {
		if s.tasks[i].Identifier == taskID {
			return i
		}
	}

	return -1
}

func (s *Heap) left(index int) int {
	if index*2+1 >= len(s.tasks) {
		return -1
	}
	return index*2 + 1
}

func (s *Heap) right(index int) int {
	if index*2+2 >= len(s.tasks) {
		return -1
	}
	return index*2 + 2
}

func (s *Heap) parent(id int) int {
	if (id-1)/2 < 0 {
		return -1
	}
	return (id - 1) / 2
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
