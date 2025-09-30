package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/webview"
)

// DemoState holds application state
type DemoState struct {
	counter int
	todos   []Todo
}

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}

// setupBridge creates and configures the bridge with all demo functions
func setupBridge() core.Bridge {
	bridge := core.NewBridge()
	state := &DemoState{
		counter: 0,
		todos: []Todo{
			{ID: 1, Title: "Learn Polyglot", Completed: false, CreatedAt: time.Now()},
			{ID: 2, Title: "Build awesome app", Completed: false, CreatedAt: time.Now()},
		},
	}

	// Register increment function
	bridge.Register("increment", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		state.counter++
		return state.counter, nil
	})

	// Register getCounter function
	bridge.Register("getCounter", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		return state.counter, nil
	})

	// Register greet function
	bridge.Register("greet", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("greet requires 1 argument")
		}
		name, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("name must be a string")
		}
		return fmt.Sprintf("Hello, %s! Welcome to Polyglot.", name), nil
	})

	// Register getTodos function
	bridge.Register("getTodos", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		return state.todos, nil
	})

	// Register addTodo function
	bridge.Register("addTodo", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("addTodo requires 1 argument")
		}
		title, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("title must be a string")
		}
		todo := Todo{
			ID:        len(state.todos) + 1,
			Title:     title,
			Completed: false,
			CreatedAt: time.Now(),
		}
		state.todos = append(state.todos, todo)
		return todo, nil
	})

	// Register toggleTodo function
	bridge.Register("toggleTodo", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toggleTodo requires 1 argument")
		}
		var id int
		switch v := args[0].(type) {
		case int:
			id = v
		case float64:
			id = int(v)
		default:
			return nil, fmt.Errorf("id must be a number")
		}

		for i := range state.todos {
			if state.todos[i].ID == id {
				state.todos[i].Completed = !state.todos[i].Completed
				return state.todos[i], nil
			}
		}
		return nil, fmt.Errorf("todo not found")
	})

	// Register deleteTodo function
	bridge.Register("deleteTodo", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("deleteTodo requires 1 argument")
		}
		var id int
		switch v := args[0].(type) {
		case int:
			id = v
		case float64:
			id = int(v)
		default:
			return nil, fmt.Errorf("id must be a number")
		}

		for i := range state.todos {
			if state.todos[i].ID == id {
				state.todos = append(state.todos[:i], state.todos[i+1:]...)
				return true, nil
			}
		}
		return false, fmt.Errorf("todo not found")
	})

	// Register getRandomNumber function
	bridge.Register("getRandomNumber", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		min, max := 1, 100
		if len(args) == 2 {
			if minVal, ok := args[0].(float64); ok {
				min = int(minVal)
			}
			if maxVal, ok := args[1].(float64); ok {
				max = int(maxVal)
			}
		}
		return rand.Intn(max-min+1) + min, nil
	})

	// Register getSystemInfo function
	bridge.Register("getSystemInfo", func(ctx context.Context, args ...interface{}) (interface{}, error) {
		return map[string]interface{}{
			"platform": "polyglot",
			"version":  "1.0.0",
			"time":     time.Now().Format(time.RFC3339),
			"uptime":   time.Since(startTime).Seconds(),
		}, nil
	})

	return bridge
}

var startTime = time.Now()

func main() {
	// Lock to OS thread - required for webview on macOS/Linux
	runtime.LockOSThread()

	// Create configuration
	config := core.WebviewConfig{
		Title:     "Polyglot Webview Demo",
		Width:     1280,
		Height:    720,
		Resizable: true,
		Debug:     true, // Enable developer tools
		URL:       generateDemoHTML(),
	}

	// Setup bridge with demo functions
	bridge := setupBridge()

	// Create webview
	wv := webview.New(config, bridge)
	if wv == nil {
		log.Fatal("Failed to create webview")
	}

	// Initialize
	if err := wv.Initialize(); err != nil {
		log.Fatalf("Failed to initialize webview: %v", err)
	}

	// Clean up on exit
	defer wv.Terminate()

	// Log startup
	log.Println("Starting Polyglot Webview Demo...")
	log.Println("- Press F12 to open developer tools (if debug enabled)")
	log.Println("- Close the window to exit")

	// Run the webview (this blocks until window is closed)
	if err := wv.Run(); err != nil {
		log.Fatalf("Failed to run webview: %v", err)
	}

	log.Println("Application closed")
}

// generateDemoHTML creates an embedded HTML page
func generateDemoHTML() string {
	html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Polyglot Webview Demo</title>
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
            max-width: 1200px;
            margin: 0 auto;
        }

        h1 {
            color: white;
            text-align: center;
            margin-bottom: 30px;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
        }

        .demo-section {
            background: white;
            border-radius: 10px;
            padding: 25px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }

        h2 {
            color: #667eea;
            margin-bottom: 15px;
            font-size: 1.5em;
        }

        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            margin: 5px;
            transition: transform 0.2s;
        }

        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.2);
        }

        button:active {
            transform: translateY(0);
        }

        input[type="text"] {
            padding: 10px;
            border: 2px solid #ddd;
            border-radius: 5px;
            font-size: 16px;
            width: 300px;
            margin-right: 10px;
        }

        input[type="text"]:focus {
            outline: none;
            border-color: #667eea;
        }

        .counter {
            font-size: 3em;
            color: #667eea;
            text-align: center;
            margin: 20px 0;
            font-weight: bold;
        }

        .message {
            padding: 15px;
            border-radius: 5px;
            margin: 10px 0;
            font-size: 16px;
        }

        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .message.info {
            background: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
        }

        .todo-list {
            list-style: none;
        }

        .todo-item {
            padding: 15px;
            margin: 10px 0;
            background: #f8f9fa;
            border-radius: 5px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            transition: all 0.3s;
        }

        .todo-item:hover {
            background: #e9ecef;
        }

        .todo-item.completed {
            opacity: 0.6;
            text-decoration: line-through;
        }

        .todo-content {
            flex: 1;
            cursor: pointer;
        }

        .todo-actions button {
            padding: 8px 16px;
            font-size: 14px;
            margin-left: 5px;
        }

        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }

        .stats-box {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px;
            text-align: center;
        }

        .stats-box .value {
            font-size: 2em;
            font-weight: bold;
            margin: 10px 0;
        }

        .stats-box .label {
            font-size: 0.9em;
            opacity: 0.9;
        }

        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid rgba(255,255,255,.3);
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
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Polyglot Webview Demo</h1>

        <!-- Counter Section -->
        <div class="demo-section">
            <h2>Counter Demo</h2>
            <p>This demonstrates simple state management between Go and JavaScript.</p>
            <div class="counter" id="counter">0</div>
            <div style="text-align: center;">
                <button onclick="incrementCounter()">Increment</button>
                <button onclick="refreshCounter()">Refresh</button>
            </div>
        </div>

        <!-- Greeting Section -->
        <div class="demo-section">
            <h2>Greeting Demo</h2>
            <p>Enter your name and get a greeting from the Go backend.</p>
            <div style="margin-top: 15px;">
                <input type="text" id="nameInput" placeholder="Enter your name" value="World">
                <button onclick="getGreeting()">Greet Me</button>
            </div>
            <div id="greetingMessage"></div>
        </div>

        <!-- Todo List Section -->
        <div class="demo-section">
            <h2>Todo List Demo</h2>
            <p>Full CRUD operations with Go backend managing state.</p>
            <div style="margin-top: 15px;">
                <input type="text" id="todoInput" placeholder="New todo...">
                <button onclick="addTodo()">Add Todo</button>
                <button onclick="loadTodos()">Refresh</button>
            </div>
            <ul class="todo-list" id="todoList"></ul>
        </div>

        <!-- Random Number Section -->
        <div class="demo-section">
            <h2>Random Number Demo</h2>
            <p>Generate random numbers using Go's crypto-quality random number generator.</p>
            <div style="text-align: center; margin: 20px 0;">
                <div style="font-size: 3em; color: #764ba2; font-weight: bold;" id="randomNumber">?</div>
                <button onclick="generateRandom()">Generate (1-100)</button>
                <button onclick="generateRandomRange()">Generate (1-1000)</button>
            </div>
        </div>

        <!-- System Info Section -->
        <div class="demo-section">
            <h2>System Info Demo</h2>
            <p>Retrieve system information from the Go backend.</p>
            <button onclick="getSystemInfo()">Get System Info</button>
            <div id="systemInfo"></div>
        </div>
    </div>

    <script>
        // Global state
        let todos = [];

        // Initialize on load
        window.addEventListener('DOMContentLoaded', async () => {
            console.log('Polyglot Demo initialized');
            await refreshCounter();
            await loadTodos();
            showMessage('Application ready!', 'success', 'greetingMessage');
        });

        // Utility function to show messages
        function showMessage(text, type = 'info', elementId = 'greetingMessage') {
            const element = document.getElementById(elementId);
            element.innerHTML = ` + "`<div class=\"message ${type}\">${text}</div>`" + `;
            setTimeout(() => {
                element.innerHTML = '';
            }, 3000);
        }

        // Counter functions
        async function incrementCounter() {
            try {
                const result = await window.polyglot.call('increment');
                document.getElementById('counter').textContent = result;
            } catch (error) {
                console.error('Increment failed:', error);
                showMessage('Failed to increment: ' + error, 'error', 'greetingMessage');
            }
        }

        async function refreshCounter() {
            try {
                const result = await window.polyglot.call('getCounter');
                document.getElementById('counter').textContent = result;
            } catch (error) {
                console.error('Refresh failed:', error);
            }
        }

        // Greeting functions
        async function getGreeting() {
            const name = document.getElementById('nameInput').value;
            if (!name.trim()) {
                showMessage('Please enter a name', 'error', 'greetingMessage');
                return;
            }

            try {
                const result = await window.polyglot.call('greet', name);
                showMessage(result, 'success', 'greetingMessage');
            } catch (error) {
                console.error('Greet failed:', error);
                showMessage('Failed to greet: ' + error, 'error', 'greetingMessage');
            }
        }

        // Todo functions
        async function loadTodos() {
            try {
                todos = await window.polyglot.call('getTodos');
                renderTodos();
            } catch (error) {
                console.error('Load todos failed:', error);
            }
        }

        async function addTodo() {
            const input = document.getElementById('todoInput');
            const title = input.value.trim();
            
            if (!title) {
                showMessage('Please enter a todo', 'error', 'greetingMessage');
                return;
            }

            try {
                const newTodo = await window.polyglot.call('addTodo', title);
                todos.push(newTodo);
                renderTodos();
                input.value = '';
                showMessage('Todo added!', 'success', 'greetingMessage');
            } catch (error) {
                console.error('Add todo failed:', error);
                showMessage('Failed to add todo: ' + error, 'error', 'greetingMessage');
            }
        }

        async function toggleTodo(id) {
            try {
                const updated = await window.polyglot.call('toggleTodo', id);
                const index = todos.findIndex(t => t.id === id);
                if (index !== -1) {
                    todos[index] = updated;
                    renderTodos();
                }
            } catch (error) {
                console.error('Toggle todo failed:', error);
            }
        }

        async function deleteTodo(id) {
            try {
                await window.polyglot.call('deleteTodo', id);
                todos = todos.filter(t => t.id !== id);
                renderTodos();
                showMessage('Todo deleted!', 'info', 'greetingMessage');
            } catch (error) {
                console.error('Delete todo failed:', error);
                showMessage('Failed to delete todo: ' + error, 'error', 'greetingMessage');
            }
        }

        function renderTodos() {
            const list = document.getElementById('todoList');
            if (todos.length === 0) {
                list.innerHTML = '<li style="padding: 20px; text-align: center; color: #999;">No todos yet. Add one above!</li>';
                return;
            }

            list.innerHTML = todos.map(todo => ` + "`" + `
                <li class="todo-item ${todo.completed ? 'completed' : ''}">
                    <div class="todo-content" onclick="toggleTodo(${todo.id})">
                        ${todo.completed ? '✓' : '○'} ${todo.title}
                    </div>
                    <div class="todo-actions">
                        <button onclick="event.stopPropagation(); deleteTodo(${todo.id})">Delete</button>
                    </div>
                </li>
            ` + "`" + `).join('');
        }

        // Random number functions
        async function generateRandom() {
            try {
                const result = await window.polyglot.call('getRandomNumber');
                document.getElementById('randomNumber').textContent = result;
            } catch (error) {
                console.error('Random failed:', error);
            }
        }

        async function generateRandomRange() {
            try {
                const result = await window.polyglot.call('getRandomNumber', 1, 1000);
                document.getElementById('randomNumber').textContent = result;
            } catch (error) {
                console.error('Random failed:', error);
            }
        }

        // System info function
        async function getSystemInfo() {
            try {
                const info = await window.polyglot.call('getSystemInfo');
                const infoHtml = ` + "`" + `
                    <div class="message info" style="margin-top: 15px;">
                        <strong>System Information:</strong><br>
                        <code>Platform:</code> ${info.platform}<br>
                        <code>Version:</code> ${info.version}<br>
                        <code>Current Time:</code> ${info.time}<br>
                        <code>Uptime:</code> ${info.uptime.toFixed(2)} seconds
                    </div>
                ` + "`" + `;
                document.getElementById('systemInfo').innerHTML = infoHtml;
            } catch (error) {
                console.error('System info failed:', error);
                showMessage('Failed to get system info: ' + error, 'error', 'greetingMessage');
            }
        }

        // Enable Enter key for inputs
        document.getElementById('nameInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') getGreeting();
        });

        document.getElementById('todoInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') addTodo();
        });
    </script>
</body>
</html>
`
	return "data:text/html," + html
}
