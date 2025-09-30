package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	pythonRuntime "github.com/griffincancode/polyglot.js/runtimes/python"
	"github.com/griffincancode/polyglot.js/webview"
)

// AppState holds application state
type AppState struct {
	pythonRuntime *pythonRuntime.Runtime
	counter       int
	tasks         []Task
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	DueDate     string    `json:"dueDate"`
}

var (
	appState  *AppState
	startTime = time.Now()
)

func main() {
	// Lock to OS thread - required for webview
	runtime.LockOSThread()

	// Initialize application state
	appState = &AppState{
		counter: 0,
		tasks: []Task{
			{
				ID:          1,
				Title:       "Learn Polyglot",
				Description: "Explore Python + Go + JS integration",
				Priority:    "high",
				Status:      "in-progress",
				CreatedAt:   time.Now(),
				DueDate:     time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
			},
			{
				ID:          2,
				Title:       "Build Demo App",
				Description: "Create a demo showcasing all features",
				Priority:    "medium",
				Status:      "pending",
				CreatedAt:   time.Now(),
				DueDate:     time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
			},
		},
	}

	// Initialize Python runtime
	if err := initializePython(); err != nil {
		log.Printf("Warning: Python runtime initialization failed: %v", err)
		log.Println("Continuing without Python support...")
	} else {
		defer shutdownPython()
		log.Println("✓ Python runtime initialized successfully")
	}

	// Setup bridge with demo functions
	bridge := setupBridge()

	// Create webview configuration
	config := core.WebviewConfig{
		Title:     "Polyglot Python + JS + Webview Demo",
		Width:     1400,
		Height:    900,
		Resizable: true,
		Debug:     true,
		URL:       generateDemoHTML(),
	}

	// Create and initialize webview
	wv := webview.New(config, bridge)
	if wv == nil {
		log.Fatal("Failed to create webview")
	}

	if err := wv.Initialize(); err != nil {
		log.Fatalf("Failed to initialize webview: %v", err)
	}
	defer wv.Terminate()

	// Log startup
	log.Println("╔════════════════════════════════════════════════════════════════╗")
	log.Println("║  Polyglot Python + JS + Webview Demo                          ║")
	log.Println("╚════════════════════════════════════════════════════════════════╝")
	log.Println("Features:")
	log.Println("  • Python runtime for backend calculations")
	log.Println("  • JavaScript for interactive UI")
	log.Println("  • Go for system operations and state management")
	log.Println("  • Real-time data processing")
	log.Println("")
	log.Println("Controls:")
	log.Println("  • Press F12 for developer tools")
	log.Println("  • Close window to exit")
	log.Println("")

	// Run the webview (blocks until window is closed)
	if err := wv.Run(); err != nil {
		log.Fatalf("Failed to run webview: %v", err)
	}

	log.Println("Application closed")
}

// initializePython initializes the Python runtime
func initializePython() error {
	appState.pythonRuntime = pythonRuntime.NewRuntime()

	config := core.RuntimeConfig{
		Name:           "python",
		Enabled:        true,
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	ctx := context.Background()
	return appState.pythonRuntime.Initialize(ctx, config)
}

// shutdownPython cleans up Python runtime
func shutdownPython() {
	if appState.pythonRuntime != nil {
		ctx := context.Background()
		appState.pythonRuntime.Shutdown(ctx)
	}
}

// setupBridge creates and configures the bridge with all demo functions
func setupBridge() core.Bridge {
	bridge := core.NewBridge()

	// Python calculation demos
	bridge.Register("pythonCalculate", pythonCalculate)
	bridge.Register("pythonFibonacci", pythonFibonacci)
	bridge.Register("pythonStatistics", pythonStatistics)
	bridge.Register("pythonTextAnalysis", pythonTextAnalysis)
	bridge.Register("pythonDataTransform", pythonDataTransform)
	bridge.Register("pythonMathOperations", pythonMathOperations)
	bridge.Register("pythonListProcessing", pythonListProcessing)

	// Task management functions
	bridge.Register("getTasks", getTasks)
	bridge.Register("addTask", addTask)
	bridge.Register("updateTask", updateTask)
	bridge.Register("deleteTask", deleteTask)
	bridge.Register("filterTasks", filterTasks)

	// System info functions
	bridge.Register("getSystemInfo", getSystemInfo)
	bridge.Register("getPythonVersion", getPythonVersion)
	bridge.Register("getUptime", getUptime)

	// Counter demo
	bridge.Register("increment", increment)
	bridge.Register("getCounter", getCounter)

	return bridge
}

// Python calculation functions

func pythonCalculate(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return nil, fmt.Errorf("Python runtime not available")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("requires 1 argument (expression)")
	}

	expr, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("expression must be a string")
	}

	result, err := appState.pythonRuntime.Execute(ctx, expr)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}

	return result, nil
}

func pythonFibonacci(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return nil, fmt.Errorf("Python runtime not available")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("requires 1 argument (n)")
	}

	var n int
	switch v := args[0].(type) {
	case float64:
		n = int(v)
	case int:
		n = v
	default:
		return nil, fmt.Errorf("n must be a number")
	}

	code := fmt.Sprintf(`
def fibonacci(n):
    if n <= 1:
        return n
    a, b = 0, 1
    for _ in range(n - 1):
        a, b = b, a + b
    return b

fibonacci(%d)
`, n)

	result, err := appState.pythonRuntime.Execute(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("fibonacci calculation failed: %w", err)
	}

	return result, nil
}

func pythonStatistics(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return nil, fmt.Errorf("Python runtime not available")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("requires 1 argument (numbers array)")
	}

	numbers, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("argument must be an array")
	}

	// Convert numbers for Python
	numStr := "["
	for i, n := range numbers {
		if i > 0 {
			numStr += ", "
		}
		numStr += fmt.Sprintf("%v", n)
	}
	numStr += "]"

	code := fmt.Sprintf(`
import statistics

data = %s

{
    'mean': statistics.mean(data),
    'median': statistics.median(data),
    'stdev': statistics.stdev(data) if len(data) > 1 else 0,
    'min': min(data),
    'max': max(data),
    'sum': sum(data),
    'count': len(data)
}
`, numStr)

	result, err := appState.pythonRuntime.Execute(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("statistics calculation failed: %w", err)
	}

	return result, nil
}

func pythonTextAnalysis(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return nil, fmt.Errorf("Python runtime not available")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("requires 1 argument (text)")
	}

	text, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("text must be a string")
	}

	// Escape quotes in text
	text = fmt.Sprintf("%q", text)

	code := fmt.Sprintf(`
text = %s

words = text.split()
chars = len(text)
words_count = len(words)
sentences = text.count('.') + text.count('!') + text.count('?')
avg_word_length = sum(len(word) for word in words) / words_count if words_count > 0 else 0

{
    'characters': chars,
    'words': words_count,
    'sentences': max(sentences, 1),
    'avgWordLength': round(avg_word_length, 2),
    'longestWord': max(words, key=len) if words else '',
    'uniqueWords': len(set(word.lower() for word in words))
}
`, text)

	result, err := appState.pythonRuntime.Execute(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("text analysis failed: %w", err)
	}

	return result, nil
}

func pythonDataTransform(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return nil, fmt.Errorf("Python runtime not available")
	}

	if len(args) != 2 {
		return nil, fmt.Errorf("requires 2 arguments (data, operation)")
	}

	data, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("data must be an array")
	}

	operation, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("operation must be a string")
	}

	// Convert data for Python
	dataStr := "["
	for i, item := range data {
		if i > 0 {
			dataStr += ", "
		}
		dataStr += fmt.Sprintf("%v", item)
	}
	dataStr += "]"

	var code string
	switch operation {
	case "square":
		code = fmt.Sprintf(`[x * x for x in %s]`, dataStr)
	case "double":
		code = fmt.Sprintf(`[x * 2 for x in %s]`, dataStr)
	case "reverse":
		code = fmt.Sprintf(`list(reversed(%s))`, dataStr)
	case "sort":
		code = fmt.Sprintf(`sorted(%s)`, dataStr)
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	result, err := appState.pythonRuntime.Execute(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("data transform failed: %w", err)
	}

	return result, nil
}

func pythonMathOperations(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return nil, fmt.Errorf("Python runtime not available")
	}

	if len(args) != 3 {
		return nil, fmt.Errorf("requires 3 arguments (a, b, operation)")
	}

	var a, b float64
	switch v := args[0].(type) {
	case float64:
		a = v
	case int:
		a = float64(v)
	default:
		return nil, fmt.Errorf("a must be a number")
	}

	switch v := args[1].(type) {
	case float64:
		b = v
	case int:
		b = float64(v)
	default:
		return nil, fmt.Errorf("b must be a number")
	}

	operation, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf("operation must be a string")
	}

	code := fmt.Sprintf(`
import math

a, b = %f, %f
`, a, b)

	switch operation {
	case "add":
		code += "a + b"
	case "subtract":
		code += "a - b"
	case "multiply":
		code += "a * b"
	case "divide":
		code += "a / b if b != 0 else None"
	case "power":
		code += "a ** b"
	case "sqrt":
		code += "math.sqrt(a)"
	case "sin":
		code += "math.sin(a)"
	case "cos":
		code += "math.cos(a)"
	case "log":
		code += "math.log(a) if a > 0 else None"
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}

	result, err := appState.pythonRuntime.Execute(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("math operation failed: %w", err)
	}

	return result, nil
}

func pythonListProcessing(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return nil, fmt.Errorf("Python runtime not available")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("requires 1 argument (size)")
	}

	var size int
	switch v := args[0].(type) {
	case float64:
		size = int(v)
	case int:
		size = v
	default:
		return nil, fmt.Errorf("size must be a number")
	}

	code := fmt.Sprintf(`
# Generate list comprehensions and demonstrate Python's data processing
size = %d

result = {
    'range': list(range(size)),
    'squares': [x*x for x in range(size)],
    'evens': [x for x in range(size) if x %% 2 == 0],
    'odds': [x for x in range(size) if x %% 2 != 0],
    'sum_squares': sum(x*x for x in range(size)),
    'fibonacci': [
        (lambda n: 
            n if n <= 1 else 
            sum((
                (lambda f, i: f(f, i))(
                    lambda f, i: i if i <= 1 else f(f, i-1) + f(f, i-2),
                    n-j
                )
                for j in range(2)
            ))
        )(i) 
        for i in range(min(size, 15))
    ][:size]
}

result
`, size)

	result, err := appState.pythonRuntime.Execute(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("list processing failed: %w", err)
	}

	return result, nil
}

// Task management functions

func getTasks(ctx context.Context, args ...interface{}) (interface{}, error) {
	return appState.tasks, nil
}

func addTask(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("requires at least 2 arguments (title, description)")
	}

	title, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("title must be a string")
	}

	description, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("description must be a string")
	}

	priority := "medium"
	if len(args) > 2 {
		if p, ok := args[2].(string); ok {
			priority = p
		}
	}

	dueDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	if len(args) > 3 {
		if d, ok := args[3].(string); ok {
			dueDate = d
		}
	}

	task := Task{
		ID:          len(appState.tasks) + 1,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      "pending",
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
	}

	appState.tasks = append(appState.tasks, task)
	return task, nil
}

func updateTask(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("requires at least 2 arguments (id, field, value)")
	}

	var id int
	switch v := args[0].(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	default:
		return nil, fmt.Errorf("id must be a number")
	}

	for i := range appState.tasks {
		if appState.tasks[i].ID == id {
			if len(args) >= 2 {
				field, ok := args[1].(string)
				if ok && len(args) >= 3 {
					value := args[2]
					switch field {
					case "status":
						if v, ok := value.(string); ok {
							appState.tasks[i].Status = v
						}
					case "priority":
						if v, ok := value.(string); ok {
							appState.tasks[i].Priority = v
						}
					case "title":
						if v, ok := value.(string); ok {
							appState.tasks[i].Title = v
						}
					case "description":
						if v, ok := value.(string); ok {
							appState.tasks[i].Description = v
						}
					}
				}
			}
			return appState.tasks[i], nil
		}
	}

	return nil, fmt.Errorf("task not found")
}

func deleteTask(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("requires 1 argument (id)")
	}

	var id int
	switch v := args[0].(type) {
	case float64:
		id = int(v)
	case int:
		id = v
	default:
		return nil, fmt.Errorf("id must be a number")
	}

	for i := range appState.tasks {
		if appState.tasks[i].ID == id {
			appState.tasks = append(appState.tasks[:i], appState.tasks[i+1:]...)
			return true, nil
		}
	}

	return false, fmt.Errorf("task not found")
}

func filterTasks(ctx context.Context, args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("requires 2 arguments (field, value)")
	}

	field, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("field must be a string")
	}

	value, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("value must be a string")
	}

	var filtered []Task
	for _, task := range appState.tasks {
		match := false
		switch field {
		case "status":
			match = task.Status == value
		case "priority":
			match = task.Priority == value
		}
		if match {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

// System info functions

func getSystemInfo(ctx context.Context, args ...interface{}) (interface{}, error) {
	pythonVersion := "not available"
	if appState.pythonRuntime != nil {
		pythonVersion = appState.pythonRuntime.Version()
	}

	return map[string]interface{}{
		"platform":      runtime.GOOS,
		"arch":          runtime.GOARCH,
		"goVersion":     runtime.Version(),
		"pythonVersion": pythonVersion,
		"numCPU":        runtime.NumCPU(),
		"numGoroutine":  runtime.NumGoroutine(),
		"uptime":        time.Since(startTime).Seconds(),
		"timestamp":     time.Now().Format(time.RFC3339),
	}, nil
}

func getPythonVersion(ctx context.Context, args ...interface{}) (interface{}, error) {
	if appState.pythonRuntime == nil {
		return "Python runtime not available", nil
	}

	return appState.pythonRuntime.Version(), nil
}

func getUptime(ctx context.Context, args ...interface{}) (interface{}, error) {
	return time.Since(startTime).Seconds(), nil
}

// Counter functions

func increment(ctx context.Context, args ...interface{}) (interface{}, error) {
	appState.counter++
	return appState.counter, nil
}

func getCounter(ctx context.Context, args ...interface{}) (interface{}, error) {
	return appState.counter, nil
}

// generateDemoHTML creates the embedded HTML page
func generateDemoHTML() string {
	return "data:text/html," + `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Polyglot Python + JS + Webview Demo</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            padding: 20px;
            min-height: 100vh;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
        }

        h1 {
            color: white;
            text-align: center;
            margin-bottom: 10px;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
        }

        .subtitle {
            color: rgba(255,255,255,0.9);
            text-align: center;
            margin-bottom: 30px;
            font-size: 1.1em;
        }

        .demo-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }

        .demo-section {
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }

        .demo-section.full-width {
            grid-column: 1 / -1;
        }

        h2 {
            color: #667eea;
            margin-bottom: 15px;
            font-size: 1.3em;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        h2::before {
            content: '🐍';
            font-size: 1.2em;
        }

        .description {
            color: #666;
            margin-bottom: 15px;
            font-size: 0.9em;
        }

        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 5px;
            cursor: pointer;
            font-size: 14px;
            margin: 5px 5px 5px 0;
            transition: transform 0.2s, box-shadow 0.2s;
        }

        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.2);
        }

        button:active {
            transform: translateY(0);
        }

        button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            transform: none;
        }

        input[type="text"],
        input[type="number"],
        textarea,
        select {
            padding: 10px;
            border: 2px solid #ddd;
            border-radius: 5px;
            font-size: 14px;
            margin: 5px 5px 5px 0;
            font-family: inherit;
        }

        input[type="text"]:focus,
        input[type="number"]:focus,
        textarea:focus,
        select:focus {
            outline: none;
            border-color: #667eea;
        }

        textarea {
            width: 100%;
            min-height: 80px;
            resize: vertical;
        }

        .result-box {
            background: #f8f9fa;
            border: 2px solid #e9ecef;
            border-radius: 5px;
            padding: 15px;
            margin-top: 10px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            max-height: 300px;
            overflow-y: auto;
        }

        .result-box.success {
            background: #d4edda;
            border-color: #c3e6cb;
            color: #155724;
        }

        .result-box.error {
            background: #f8d7da;
            border-color: #f5c6cb;
            color: #721c24;
        }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 10px;
            margin: 15px 0;
        }

        .stat-item {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 15px;
            border-radius: 5px;
            text-align: center;
        }

        .stat-value {
            font-size: 2em;
            font-weight: bold;
            margin: 5px 0;
        }

        .stat-label {
            font-size: 0.9em;
            opacity: 0.9;
        }

        .task-list {
            list-style: none;
        }

        .task-item {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            border-radius: 5px;
            padding: 15px;
            margin: 10px 0;
            transition: all 0.3s;
        }

        .task-item:hover {
            background: #e9ecef;
            transform: translateX(5px);
        }

        .task-item.high {
            border-left-color: #dc3545;
        }

        .task-item.medium {
            border-left-color: #ffc107;
        }

        .task-item.low {
            border-left-color: #28a745;
        }

        .task-item.completed {
            opacity: 0.6;
            text-decoration: line-through;
        }

        .task-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }

        .task-title {
            font-weight: bold;
            font-size: 1.1em;
        }

        .task-meta {
            font-size: 0.8em;
            color: #666;
            margin: 5px 0;
        }

        .badge {
            display: inline-block;
            padding: 3px 8px;
            border-radius: 3px;
            font-size: 0.8em;
            font-weight: bold;
            margin: 0 5px;
        }

        .badge.status {
            background: #17a2b8;
            color: white;
        }

        .badge.priority {
            background: #6c757d;
            color: white;
        }

        .badge.priority.high {
            background: #dc3545;
        }

        .badge.priority.medium {
            background: #ffc107;
            color: #333;
        }

        .badge.priority.low {
            background: #28a745;
        }

        .loading {
            display: inline-block;
            width: 16px;
            height: 16px;
            border: 2px solid rgba(255,255,255,.3);
            border-radius: 50%;
            border-top-color: white;
            animation: spin 1s ease-in-out infinite;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        code {
            background: #f8f9fa;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            color: #e83e8c;
            font-size: 0.9em;
        }

        .tabs {
            display: flex;
            gap: 5px;
            margin-bottom: 15px;
            border-bottom: 2px solid #e9ecef;
        }

        .tab {
            padding: 10px 20px;
            background: transparent;
            color: #667eea;
            border: none;
            border-bottom: 2px solid transparent;
            cursor: pointer;
            transition: all 0.2s;
            margin: 0 0 -2px 0;
        }

        .tab:hover {
            background: rgba(102, 126, 234, 0.1);
            transform: none;
        }

        .tab.active {
            border-bottom-color: #667eea;
            font-weight: bold;
        }

        .tab-content {
            display: none;
        }

        .tab-content.active {
            display: block;
        }

        .input-group {
            display: flex;
            gap: 10px;
            align-items: center;
            margin: 10px 0;
        }

        .input-group label {
            min-width: 100px;
            font-weight: 500;
        }

        .input-group input,
        .input-group select {
            flex: 1;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Polyglot Python + JS + Webview Demo</h1>
        <p class="subtitle">Demonstrating seamless integration between Python, JavaScript, and Go</p>

        <div class="demo-grid">
            <!-- Python Calculator -->
            <div class="demo-section">
                <h2>Python Calculator</h2>
                <p class="description">Execute Python expressions directly from JavaScript</p>
                <input type="text" id="calcInput" placeholder="e.g., 2**10, math.sqrt(144)" style="width: 100%;">
                <button onclick="calculate()">Calculate</button>
                <button onclick="setExample('calcInput', 'import math; math.pi * 2')">π × 2</button>
                <button onclick="setExample('calcInput', 'sum(range(1, 101))')">Σ(1-100)</button>
                <div id="calcResult" class="result-box" style="display:none;"></div>
            </div>

            <!-- Fibonacci -->
            <div class="demo-section">
                <h2>Fibonacci Generator</h2>
                <p class="description">Generate Fibonacci numbers using Python</p>
                <input type="number" id="fibInput" value="10" min="1" max="50" style="width: 150px;">
                <button onclick="fibonacci()">Generate</button>
                <div id="fibResult" class="result-box" style="display:none;"></div>
            </div>

            <!-- Statistics -->
            <div class="demo-section">
                <h2>Statistical Analysis</h2>
                <p class="description">Analyze number arrays with Python's statistics module</p>
                <input type="text" id="statsInput" placeholder="e.g., 1,2,3,4,5" style="width: 100%;" value="12,15,18,20,22,25,28,30">
                <button onclick="analyzeStats()">Analyze</button>
                <button onclick="generateRandomStats()">Random Data</button>
                <div id="statsResult" class="result-box" style="display:none;"></div>
            </div>

            <!-- Text Analysis -->
            <div class="demo-section">
                <h2>Text Analysis</h2>
                <p class="description">Analyze text using Python string processing</p>
                <textarea id="textInput" placeholder="Enter text to analyze...">The quick brown fox jumps over the lazy dog. This sentence contains every letter of the alphabet.</textarea>
                <button onclick="analyzeText()">Analyze</button>
                <div id="textResult" class="result-box" style="display:none;"></div>
            </div>

            <!-- Data Transformation -->
            <div class="demo-section">
                <h2>Data Transformation</h2>
                <p class="description">Transform arrays using Python list comprehensions</p>
                <input type="text" id="transformInput" placeholder="e.g., 1,2,3,4,5" style="width: 200px;" value="1,2,3,4,5,6,7,8,9,10">
                <select id="transformOp">
                    <option value="square">Square</option>
                    <option value="double">Double</option>
                    <option value="reverse">Reverse</option>
                    <option value="sort">Sort</option>
                </select>
                <button onclick="transformData()">Transform</button>
                <div id="transformResult" class="result-box" style="display:none;"></div>
            </div>

            <!-- List Processing -->
            <div class="demo-section">
                <h2>List Processing</h2>
                <p class="description">Demonstrate Python's powerful list comprehensions</p>
                <input type="number" id="listSize" value="10" min="1" max="20" style="width: 150px;">
                <button onclick="processList()">Process</button>
                <div id="listResult" class="result-box" style="display:none;"></div>
            </div>

            <!-- Task Manager -->
            <div class="demo-section full-width">
                <h2>Task Manager (Go + JS)</h2>
                <p class="description">Full CRUD operations with Go backend state management</p>
                
                <div class="tabs">
                    <button class="tab active" onclick="switchTab('tasks')">Tasks</button>
                    <button class="tab" onclick="switchTab('addTask')">Add Task</button>
                    <button class="tab" onclick="switchTab('filter')">Filter</button>
                </div>

                <div id="tab-tasks" class="tab-content active">
                    <button onclick="loadTasks()">Refresh Tasks</button>
                    <ul id="taskList" class="task-list"></ul>
                </div>

                <div id="tab-addTask" class="tab-content">
                    <div class="input-group">
                        <label>Title:</label>
                        <input type="text" id="newTaskTitle" placeholder="Task title">
                    </div>
                    <div class="input-group">
                        <label>Description:</label>
                        <input type="text" id="newTaskDesc" placeholder="Task description">
                    </div>
                    <div class="input-group">
                        <label>Priority:</label>
                        <select id="newTaskPriority">
                            <option value="low">Low</option>
                            <option value="medium" selected>Medium</option>
                            <option value="high">High</option>
                        </select>
                    </div>
                    <button onclick="addNewTask()">Add Task</button>
                </div>

                <div id="tab-filter" class="tab-content">
                    <div class="input-group">
                        <label>Filter by:</label>
                        <select id="filterField">
                            <option value="status">Status</option>
                            <option value="priority">Priority</option>
                        </select>
                        <input type="text" id="filterValue" placeholder="e.g., pending, high">
                        <button onclick="filterTasksBy()">Filter</button>
                        <button onclick="loadTasks()">Show All</button>
                    </div>
                    <ul id="filteredTaskList" class="task-list"></ul>
                </div>
            </div>

            <!-- System Info -->
            <div class="demo-section full-width">
                <h2>System Information</h2>
                <p class="description">Runtime information from Go backend</p>
                <button onclick="loadSystemInfo()">Load Info</button>
                <div id="systemInfo" class="stats-grid" style="display:none;"></div>
            </div>
        </div>
    </div>

    <script>
        let tasks = [];

        // Initialize
        window.addEventListener('DOMContentLoaded', async () => {
            console.log('Polyglot Python + JS + Webview Demo initialized');
            await loadTasks();
            showMessage('Application ready! Try the Python features above.', 'success', 'calcResult');
        });

        // Utility functions
        function showResult(elementId, content, isError = false) {
            const element = document.getElementById(elementId);
            element.style.display = 'block';
            element.className = 'result-box ' + (isError ? 'error' : 'success');
            element.textContent = typeof content === 'object' ? JSON.stringify(content, null, 2) : content;
        }

        function showMessage(text, type, elementId) {
            showResult(elementId, text, type === 'error');
        }

        function setExample(inputId, value) {
            document.getElementById(inputId).value = value;
        }

        // Python functions
        async function calculate() {
            const expr = document.getElementById('calcInput').value;
            if (!expr.trim()) {
                showMessage('Please enter an expression', 'error', 'calcResult');
                return;
            }

            try {
                const result = await window.polyglot.call('pythonCalculate', expr);
                showResult('calcResult', 'Result: ' + result);
            } catch (error) {
                showMessage('Error: ' + error, 'error', 'calcResult');
            }
        }

        async function fibonacci() {
            const n = parseInt(document.getElementById('fibInput').value);
            if (isNaN(n) || n < 1) {
                showMessage('Please enter a valid number', 'error', 'fibResult');
                return;
            }

            try {
                const result = await window.polyglot.call('pythonFibonacci', n);
                showResult('fibResult', ` + "`Fibonacci(${n}) = ${result}`" + `);
            } catch (error) {
                showMessage('Error: ' + error, 'error', 'fibResult');
            }
        }

        async function analyzeStats() {
            const input = document.getElementById('statsInput').value;
            const numbers = input.split(',').map(n => parseFloat(n.trim())).filter(n => !isNaN(n));
            
            if (numbers.length === 0) {
                showMessage('Please enter valid numbers', 'error', 'statsResult');
                return;
            }

            try {
                const stats = await window.polyglot.call('pythonStatistics', numbers);
                showResult('statsResult', stats);
            } catch (error) {
                showMessage('Error: ' + error, 'error', 'statsResult');
            }
        }

        function generateRandomStats() {
            const numbers = Array.from({length: 20}, () => Math.floor(Math.random() * 100) + 1);
            document.getElementById('statsInput').value = numbers.join(',');
            analyzeStats();
        }

        async function analyzeText() {
            const text = document.getElementById('textInput').value;
            if (!text.trim()) {
                showMessage('Please enter text to analyze', 'error', 'textResult');
                return;
            }

            try {
                const analysis = await window.polyglot.call('pythonTextAnalysis', text);
                showResult('textResult', analysis);
            } catch (error) {
                showMessage('Error: ' + error, 'error', 'textResult');
            }
        }

        async function transformData() {
            const input = document.getElementById('transformInput').value;
            const numbers = input.split(',').map(n => parseFloat(n.trim())).filter(n => !isNaN(n));
            const operation = document.getElementById('transformOp').value;
            
            if (numbers.length === 0) {
                showMessage('Please enter valid numbers', 'error', 'transformResult');
                return;
            }

            try {
                const result = await window.polyglot.call('pythonDataTransform', numbers, operation);
                showResult('transformResult', ` + "`Operation: ${operation}\\nResult: ${JSON.stringify(result)}`" + `);
            } catch (error) {
                showMessage('Error: ' + error, 'error', 'transformResult');
            }
        }

        async function processList() {
            const size = parseInt(document.getElementById('listSize').value);
            if (isNaN(size) || size < 1) {
                showMessage('Please enter a valid size', 'error', 'listResult');
                return;
            }

            try {
                const result = await window.polyglot.call('pythonListProcessing', size);
                showResult('listResult', result);
            } catch (error) {
                showMessage('Error: ' + error, 'error', 'listResult');
            }
        }

        // Task management
        async function loadTasks() {
            try {
                tasks = await window.polyglot.call('getTasks');
                renderTasks(tasks, 'taskList');
            } catch (error) {
                console.error('Load tasks failed:', error);
            }
        }

        async function addNewTask() {
            const title = document.getElementById('newTaskTitle').value;
            const desc = document.getElementById('newTaskDesc').value;
            const priority = document.getElementById('newTaskPriority').value;

            if (!title.trim()) {
                alert('Please enter a task title');
                return;
            }

            try {
                await window.polyglot.call('addTask', title, desc, priority);
                document.getElementById('newTaskTitle').value = '';
                document.getElementById('newTaskDesc').value = '';
                switchTab('tasks');
                await loadTasks();
            } catch (error) {
                alert('Failed to add task: ' + error);
            }
        }

        async function updateTaskStatus(id, newStatus) {
            try {
                await window.polyglot.call('updateTask', id, 'status', newStatus);
                await loadTasks();
            } catch (error) {
                console.error('Update task failed:', error);
            }
        }

        async function deleteTask(id) {
            if (!confirm('Are you sure you want to delete this task?')) {
                return;
            }

            try {
                await window.polyglot.call('deleteTask', id);
                await loadTasks();
            } catch (error) {
                console.error('Delete task failed:', error);
            }
        }

        async function filterTasksBy() {
            const field = document.getElementById('filterField').value;
            const value = document.getElementById('filterValue').value;

            if (!value.trim()) {
                alert('Please enter a filter value');
                return;
            }

            try {
                const filtered = await window.polyglot.call('filterTasks', field, value);
                renderTasks(filtered, 'filteredTaskList');
            } catch (error) {
                alert('Filter failed: ' + error);
            }
        }

        function renderTasks(taskList, elementId) {
            const container = document.getElementById(elementId);
            
            if (taskList.length === 0) {
                container.innerHTML = '<li style="padding: 20px; text-align: center; color: #999;">No tasks found</li>';
                return;
            }

            container.innerHTML = taskList.map(task => ` + "`" + `
                <li class="task-item ${task.priority}">
                    <div class="task-header">
                        <div class="task-title">${task.title}</div>
                        <div>
                            <span class="badge priority ${task.priority}">${task.priority}</span>
                            <span class="badge status">${task.status}</span>
                        </div>
                    </div>
                    <div>${task.description}</div>
                    <div class="task-meta">
                        📅 Due: ${task.dueDate} | Created: ${new Date(task.createdAt).toLocaleDateString()}
                    </div>
                    <div style="margin-top: 10px;">
                        ${task.status !== 'completed' ? ` + "`" + `
                            <button onclick="updateTaskStatus(${task.id}, 'completed')">✓ Complete</button>
                        ` + "`" + ` : ` + "`" + `
                            <button onclick="updateTaskStatus(${task.id}, 'pending')">↺ Reopen</button>
                        ` + "`" + `}
                        <button onclick="deleteTask(${task.id})">🗑 Delete</button>
                    </div>
                </li>
            ` + "`" + `).join('');
        }

        function switchTab(tabName) {
            // Update tab buttons
            document.querySelectorAll('.tab').forEach(tab => {
                tab.classList.remove('active');
            });
            event.target.classList.add('active');

            // Update tab content
            document.querySelectorAll('.tab-content').forEach(content => {
                content.classList.remove('active');
            });
            document.getElementById(` + "`tab-${tabName}`" + `).classList.add('active');
        }

        // System info
        async function loadSystemInfo() {
            try {
                const info = await window.polyglot.call('getSystemInfo');
                const container = document.getElementById('systemInfo');
                container.style.display = 'grid';
                container.innerHTML = ` + "`" + `
                    <div class="stat-item">
                        <div class="stat-label">Platform</div>
                        <div class="stat-value" style="font-size: 1.5em;">${info.platform}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Architecture</div>
                        <div class="stat-value" style="font-size: 1.5em;">${info.arch}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Go Version</div>
                        <div class="stat-value" style="font-size: 1em;">${info.goVersion}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Python Version</div>
                        <div class="stat-value" style="font-size: 0.9em;">${info.pythonVersion.split(' ')[0]}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">CPU Cores</div>
                        <div class="stat-value">${info.numCPU}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Goroutines</div>
                        <div class="stat-value">${info.numGoroutine}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-label">Uptime</div>
                        <div class="stat-value">${Math.floor(info.uptime)}s</div>
                    </div>
                ` + "`" + `;
            } catch (error) {
                alert('Failed to load system info: ' + error);
            }
        }

        // Keyboard shortcuts
        document.getElementById('calcInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') calculate();
        });

        document.getElementById('fibInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') fibonacci();
        });

        document.getElementById('statsInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') analyzeStats();
        });
    </script>
</body>
</html>`
}
