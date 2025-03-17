package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/ppriyankuu/seiban/pkg/files"
)

// represents the board name and a slice of list
type Data struct {
	boardName string
	lists     []List
	fileName  string
}

// represents the title of list and a list of items inside it.
type List struct {
	listTitle string
	listItems []ListItem
}

// represents the name of item and it's description
type ListItem struct {
	ItemName        string
	ItemDescription string
}

func (d *Data) SetFileName(fileName string) {
	d.fileName = fileName
}

func (d *Data) GetContentFromFile() []string {
	if !files.CheckFile(d.fileName) {
		files.CreateFile(d.fileName)
	}
	return files.OpenFile(d.fileName)
}

// parses the contents of the file to custom type Data
func (d *Data) ParseData(fileContent []string) error {
	for lineNumber, line := range fileContent {
		line = strings.TrimSpace(line)
		// skipping empty lines
		if len(line) < 1 {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			boardNameStartingIndex := strings.Index(line, " ") + 1
			boardName := line[boardNameStartingIndex:]
			d.boardName = boardName
		} else if strings.HasPrefix(line, "## ") {
			listNameStartIndex := strings.Index(line, " ") + 1
			listTitle := line[listNameStartIndex:]
			d.lists = append(d.lists, List{
				listTitle: listTitle,
			})
		} else if strings.HasPrefix(line, "- ") {
			listCount := d.GetListCount()
			if listCount < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList := d.lists[listCount-1]
			itemNameStartIndex := strings.Index(line, " ") + 1
			itemName := line[itemNameStartIndex:]
			currentList.listItems = append(currentList.listItems, ListItem{
				ItemName: itemName,
			})
			d.lists[listCount-1] = currentList
		} else if strings.HasPrefix(line, "> ") {
			listCount := d.GetListCount()
			if listCount < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList := d.lists[listCount-1]
			itemDescStartIndex := strings.Index(line, " ") + 1
			itemDesc := line[itemDescStartIndex:]
			listItemLen := len(currentList.listItems)
			if listItemLen < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList.listItems[listItemLen-1].ItemDescription = itemDesc
			d.lists[listCount-1] = currentList
		} else {
			return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
		}
	}
	return nil
}

// returns the name of board
func (d *Data) GetBoardName() string {
	return d.boardName
}

// returns the list based on index
func (d *Data) GetList(listIdx int) (*List, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return nil, err
	}
	return &d.lists[listIdx], nil
}

// returns a list of all the list names
func (d *Data) GetListNames() []string {
	var listNames []string
	for _, list := range d.lists {
		listNames = append(listNames, list.listTitle)
	}
	return listNames
}

// gives the task it's title and description given the list and task index.
func (d *Data) GetTask(listIdx, taskIdx int) (*ListItem, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return nil, err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return nil, err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return nil, err
	}
	return &list.listItems[taskIdx], nil
}

// returns a list of all the tasks of a particular list
func (d *Data) GetTasks(listIdx int) ([]string, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return nil, err
	}
	var tasks []string
	for _, item := range list.listItems {
		tasks = append(tasks, item.ItemName)
	}
	return tasks, nil
}

// adds a new task to a list provided the list index, task (title, description), and task index.
func (d *Data) AddNewTask(listIdx int, taskTitle, taskDesc string, taskIdx int) error {
	if err := checkBounds(listIdx, d.GetListCount()); err != nil {
		return err
	}
	newTask := ListItem{
		ItemName:        taskTitle,
		ItemDescription: taskDesc,
	}
	err := d.insertTask(listIdx, newTask, taskIdx)
	if err != nil {
		return err
	}
	d.Save()
	return nil
}

// edits a task
func (d *Data) EditTask(listIdx, taskIdx int, taskTitle, taskDesc string) error {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return err
	}
	task.ItemName = taskTitle
	task.ItemDescription = taskDesc
	return nil
}

func (d *Data) MoveTask(taskIdx, sourceListIdx, destListIdx int) error {
	sourceTask, err := d.GetTask(sourceListIdx, taskIdx)
	if err != nil {
		return err
	}
	newListTaskCount, err := d.GetTaskCount(destListIdx)
	if err != nil {
		return err
	}
	taskTitle := sourceTask.ItemName
	taskDesc := sourceTask.ItemDescription
	err = d.AddNewTask(destListIdx, taskTitle, taskDesc, newListTaskCount)
	if err != nil {
		return nil
	}
	_, err = d.RemoveTask(sourceListIdx, taskIdx)
	return err
}

// removes a task given the index of list and the task.
func (d *Data) RemoveTask(listIdx, taskIdx int) (ListItem, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return ListItem{}, err
	}
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return ListItem{}, err
	}
	taskData := *task
	list.listItems = append(list.listItems[:taskIdx], list.listItems[taskIdx+1:]...)
	d.Save()
	return taskData, nil
}

func (d *Data) insertTask(listIdx int, task ListItem, taskIdx int) error {
	list, err := d.GetList(listIdx)
	if err != nil {
		return err
	}
	if len(list.listItems) < 1 {
		list.listItems = append(list.listItems, task)
		return nil
	}

	list.listItems = append(list.listItems, ListItem{})
	copy(list.listItems[(taskIdx+1):], list.listItems[taskIdx:])
	list.listItems[taskIdx] = task
	return nil
}

func (d *Data) Save() {
	if d.fileName == "" {
		return
	}
	var fileContent []string
	fileContent = append(fileContent, "# "+d.boardName+"\n")
	for _, list := range d.lists {
		fileContent = append(fileContent, "## "+list.listTitle)
		for _, listItem := range list.listItems {
			fileContent = append(fileContent, "\t- "+listItem.ItemName)
			if len(listItem.ItemDescription) > 0 {
				fileContent = append(fileContent, "\t\t> "+listItem.ItemDescription)
			}
		}
		fileContent = append(fileContent, "\n")
	}
	err := files.WriteFile(fileContent, d.fileName)
	if err != nil {
		log.Fatal(err)
	}
}

// swaps a task of one list with another list
func (d *Data) SwapListItems(listIdx, firstTaskIdx, secondTaskIdx int) error {
	firstTask, err := d.GetTask(listIdx, firstTaskIdx)
	if err != nil {
		return err
	}
	secondTask, err := d.GetTask(listIdx, secondTaskIdx)
	if err != nil {
		return err
	}
	*firstTask, *secondTask = *secondTask, *firstTask
	d.Save()
	return nil
}

// returns the count of lists
func (d *Data) GetListCount() int {
	return len(d.lists)
}

// gives the task count of a particular task
func (d *Data) GetTaskCount(listIdx int) (int, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return 0, err
	}
	return len(list.listItems), nil
}

func checkBounds(idx, boundary int) error {
	if idx < 0 || idx >= boundary {
		return fmt.Errorf("Index out of bounds: got %v, length %v", idx, boundary)
	}
	return nil
}
