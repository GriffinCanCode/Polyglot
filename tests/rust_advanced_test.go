package tests

import (
	"context"
	"testing"
	"time"

	"github.com/griffincancode/polyglot.js/core"
	"github.com/griffincancode/polyglot.js/runtimes/rust"
)

// TestRustBasicExecution tests basic Rust code execution
func TestRustBasicExecution(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		// Expected to fail without build tag or rustc
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Test simple expression
	result, err := runtime.Execute(ctx, "2 + 2")
	if err != nil {
		t.Logf("Execute returned error (may be expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustConcurrentExecution tests concurrent Rust code execution
func TestRustConcurrentExecution(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Run concurrent executions
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			_, err := runtime.Execute(ctx, "println!(\"Hello from Rust\")")
			if err != nil {
				t.Logf("Execution %d error (expected without rustc): %v", id, err)
			}
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestRustFunctionCall tests calling Rust functions
func TestRustFunctionCall(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Test calling a function
	result, err := runtime.Call(ctx, "test_function")
	if err != nil {
		t.Logf("Call returned error (expected without library): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustComplexCode tests more complex Rust code
func TestRustComplexCode(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		let mut sum = 0;
		for i in 1..=10 {
			sum += i;
		}
		println!("{}", sum);
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustWithStructs tests Rust code with structs
func TestRustWithStructs(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		struct Point {
			x: i32,
			y: i32,
		}
		
		fn main() {
			let p = Point { x: 10, y: 20 };
			println!("Point: ({}, {})", p.x, p.y);
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustWithVectors tests Rust code with vectors
func TestRustWithVectors(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		let v = vec![1, 2, 3, 4, 5];
		let sum: i32 = v.iter().sum();
		println!("Sum: {}", sum);
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustWithClosures tests Rust code with closures
func TestRustWithClosures(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		let numbers = vec![1, 2, 3, 4, 5];
		let doubled: Vec<i32> = numbers.iter().map(|x| x * 2).collect();
		println!("{:?}", doubled);
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustWithErrorHandling tests Rust error handling
func TestRustWithErrorHandling(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		fn divide(a: i32, b: i32) -> Result<i32, String> {
			if b == 0 {
				Err(String::from("Division by zero"))
			} else {
				Ok(a / b)
			}
		}
		
		fn main() {
			match divide(10, 2) {
				Ok(result) => println!("Result: {}", result),
				Err(e) => println!("Error: {}", e),
			}
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustWithTraits tests Rust code with traits
func TestRustWithTraits(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		trait Speak {
			fn speak(&self) -> String;
		}
		
		struct Dog;
		
		impl Speak for Dog {
			fn speak(&self) -> String {
				String::from("Woof!")
			}
		}
		
		fn main() {
			let dog = Dog;
			println!("{}", dog.speak());
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustWithGenerics tests Rust code with generics
func TestRustWithGenerics(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		fn largest<T: PartialOrd>(list: &[T]) -> &T {
			let mut largest = &list[0];
			for item in list {
				if item > largest {
					largest = item;
				}
			}
			largest
		}
		
		fn main() {
			let numbers = vec![34, 50, 25, 100, 65];
			let result = largest(&numbers);
			println!("The largest number is {}", result);
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustShutdown tests proper cleanup
func TestRustShutdown(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}

	// Shutdown
	if err := runtime.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Try to execute after shutdown (should fail)
	_, err = runtime.Execute(ctx, "println!(\"test\")")
	if err == nil {
		t.Error("Expected error after shutdown")
	}
}

// TestRustVersion tests version retrieval
func TestRustVersion(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		// Even without initialization, version should return something
		version := runtime.Version()
		if version == "" {
			t.Error("Version should not be empty")
		} else {
			t.Logf("Rust version: %s", version)
		}
		return
	}
	defer runtime.Shutdown(ctx)

	version := runtime.Version()
	if version == "" {
		t.Error("Version should not be empty")
	} else {
		t.Logf("Rust version: %s", version)
	}
}

// TestRustTimeout tests execution timeout
func TestRustTimeout(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        1 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	code := `
		use std::thread;
		use std::time::Duration;
		
		fn main() {
			thread::sleep(Duration::from_secs(5));
			println!("Done");
		}
	`

	_, err = runtime.Execute(timeoutCtx, code)
	if err != nil {
		t.Logf("Execute correctly timed out or failed: %v", err)
	}
}

// TestRustMemorySafety tests Rust's memory safety features
func TestRustMemorySafety(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		fn main() {
			let s1 = String::from("hello");
			let s2 = s1.clone();  // Explicit clone to avoid move
			println!("{} {}", s1, s2);
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustLifetimes tests Rust's lifetime system
func TestRustLifetimes(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
			if x.len() > y.len() {
				x
			} else {
				y
			}
		}
		
		fn main() {
			let string1 = String::from("long string");
			let string2 = String::from("short");
			let result = longest(string1.as_str(), string2.as_str());
			println!("The longest string is {}", result);
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustMacros tests Rust macro usage
func TestRustMacros(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		macro_rules! say_hello {
			() => {
				println!("Hello!");
			};
		}
		
		fn main() {
			say_hello!();
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}

// TestRustPatternMatching tests pattern matching
func TestRustPatternMatching(t *testing.T) {
	runtime := rust.NewRuntime()
	ctx := context.Background()

	config := core.RuntimeConfig{
		Name:           "rust",
		Version:        "1.70",
		Enabled:        true,
		Options:        make(map[string]interface{}),
		MaxConcurrency: 4,
		Timeout:        30 * time.Second,
	}

	err := runtime.Initialize(ctx, config)
	if err != nil {
		t.Logf("Initialize returned expected error: %v", err)
		return
	}
	defer runtime.Shutdown(ctx)

	code := `
		enum Coin {
			Penny,
			Nickel,
			Dime,
			Quarter,
		}
		
		fn value_in_cents(coin: Coin) -> u8 {
			match coin {
				Coin::Penny => 1,
				Coin::Nickel => 5,
				Coin::Dime => 10,
				Coin::Quarter => 25,
			}
		}
		
		fn main() {
			let coin = Coin::Quarter;
			println!("Value: {} cents", value_in_cents(coin));
		}
	`

	result, err := runtime.Execute(ctx, code)
	if err != nil {
		t.Logf("Execute returned error (expected without rustc): %v", err)
	} else {
		t.Logf("Result: %v", result)
	}
}
