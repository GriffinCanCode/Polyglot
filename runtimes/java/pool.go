//go:build runtime_java
// +build runtime_java

package java

/*
#include <jni.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// EnvPool manages JNI environment instances
type EnvPool struct {
	jvm      *C.JavaVM
	envs     chan *C.JNIEnv
	size     int
	mu       sync.Mutex
	shutdown bool
}

// NewEnvPool creates an environment pool
func NewEnvPool(size int) *EnvPool {
	return &EnvPool{
		envs: make(chan *C.JNIEnv, size),
		size: size,
	}
}

// Initialize prepares the pool
func (p *EnvPool) Initialize(jvm *C.JavaVM, size int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.jvm = jvm
	p.size = size

	// Create initial environments
	for i := 0; i < size; i++ {
		env, err := p.attachThread()
		if err != nil {
			return fmt.Errorf("failed to create env %d: %w", i, err)
		}
		p.envs <- env
	}

	return nil
}

// Acquire gets an environment from the pool
func (p *EnvPool) Acquire() *C.JNIEnv {
	return <-p.envs
}

// Release returns an environment to the pool
func (p *EnvPool) Release(env *C.JNIEnv) {
	if !p.shutdown {
		p.envs <- env
	}
}

// Close shuts down the pool
func (p *EnvPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.shutdown = true
	close(p.envs)

	// Detach all threads
	for env := range p.envs {
		if env != nil {
			p.detachThread()
		}
	}
}

// attachThread attaches the current thread to JVM
func (p *EnvPool) attachThread() (*C.JNIEnv, error) {
	var env *C.JNIEnv
	res := C.jvm_AttachCurrentThread(p.jvm, (*unsafe.Pointer)(unsafe.Pointer(&env)), nil)
	if res != C.JNI_OK {
		return nil, fmt.Errorf("failed to attach thread: %d", res)
	}
	return env, nil
}

// detachThread detaches the current thread from JVM
func (p *EnvPool) detachThread() {
	C.jvm_DetachCurrentThread(p.jvm)
}

// jvm_AttachCurrentThread wraps thread attachment
//
//export jvm_AttachCurrentThread
func jvm_AttachCurrentThread(jvm *C.JavaVM, penv *unsafe.Pointer, args unsafe.Pointer) C.jint {
	return (*jvm).AttachCurrentThread(jvm, penv, args)
}

// jvm_DetachCurrentThread wraps thread detachment
//
//export jvm_DetachCurrentThread
func jvm_DetachCurrentThread(jvm *C.JavaVM) C.jint {
	return (*jvm).DetachCurrentThread(jvm)
}
