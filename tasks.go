package morpheus

import (
	"fmt"
	"time"
)

var (
	// TasksPath is the API endpoint for tasks
	TasksPath = "/api/tasks"
)

// Task structures for use in request and response payloads
type Task struct {
	ID        int64    `json:"id"`
	AccountId int64    `json:"accountId"`
	Name      string   `json:"name"`
	Labels    []string `json:"labels"`
	Code      string   `json:"code"`
	TaskType  struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"taskType"`
	TaskOptions struct {
		AnsibleOptions            string `json:"ansibleOptions"`
		AnsibleTags               string `json:"ansibleTags"`
		AnsiblePlaybook           string `json:"ansiblePlaybook"`
		AnsibleGitRef             string `json:"ansibleGitRef"`
		AnsibleSkipTags           string `json:"ansibleSkipTags"`
		AnsibleGitId              string `json:"ansibleGitId"`
		JsScript                  string `json:"jsScript"`
		WinrmElevated             string `json:"winrm.elevated"`
		PythonBinary              string `json:"pythonBinary"`
		PythonArgs                string `json:"pythonArgs"`
		PythonAdditionalPackages  string `json:"pythonAdditionalPackages"`
		ShellSudo                 string `json:"shell.sudo"`
		Username                  string `json:"username"`
		Host                      string `json:"host"`
		LocalScriptGitRef         string `json:"localScriptGitRef"`
		Password                  string `json:"password"`
		PasswordHash              string `json:"passwordHash"`
		WriteAttributesAttributes string `json:"writeAttributes.attributes"`
		Port                      string `json:"port"`
	} `json:"taskOptions"`
	File struct {
		ID          int64  `json:"id"`
		SourceType  string `json:"sourceType"`
		ContentRef  string `json:"contentRef"`
		ContentPath string `json:"contentPath"`
		Repository  struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"repository"`
		Content string `json:"content"`
	} `json:"file"`
	ResultType        string `json:"resultType"`
	ExecuteTarget     string `json:"executeTarget"`
	Retryable         bool   `json:"retryable"`
	RetryCount        int64  `json:"retryCount"`
	RetryDelaySeconds int64  `json:"retryDelaySeconds"`
	AllowCustomConfig bool   `json:"allowCustomConfig"`
	Credential        struct {
		Type string `json:"type"`
	} `json:"credential"`
	DateCreated time.Time `json:"dateCreated"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type ListTasksResult struct {
	Tasks *[]Task     `json:"tasks"`
	Meta  *MetaResult `json:"meta"`
}

type GetTaskResult struct {
	Task *Task `json:"task"`
}

type CreateTaskResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Task    *Task             `json:"task"`
}

type UpdateTaskResult struct {
	CreateTaskResult
}

type DeleteTaskResult struct {
	DeleteResult
}

// ListTasks lists all tasks
func (client *Client) ListTasks(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        TasksPath,
		QueryParams: req.QueryParams,
		Result:      &ListTasksResult{},
	})
}

// GetTask gets an existing task
func (client *Client) GetTask(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", TasksPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetTaskResult{},
	})
}

// CreateTask creates a new task
func (client *Client) CreateTask(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        TasksPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateTaskResult{},
	})
}

// UpdateTask updates an existing task
func (client *Client) UpdateTask(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", TasksPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateTaskResult{},
	})
}

// DeleteTask deletes an existing task
func (client *Client) DeleteTask(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", TasksPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteTaskResult{},
	})
}

// FindTaskByName gets an existing task by name
func (client *Client) FindTaskByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListTasks(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListTasksResult)
	tasksCount := len(*listResult.Tasks)
	if tasksCount != 1 {
		return resp, fmt.Errorf("found %d Tasks for %v", tasksCount, name)
	}
	firstRecord := (*listResult.Tasks)[0]
	taskID := firstRecord.ID
	return client.GetTask(taskID, &Request{})
}
