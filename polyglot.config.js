// Polyglot Framework Configuration
// This file determines which language runtimes to include in the build

export default {
  languages: {
    python: { 
      enabled: true, 
      version: '3.11',
      packages: ['numpy', 'pandas', 'scikit-learn']
    },
    rust: { 
      enabled: true, 
      features: ['tokio', 'serde'],
      crates: ['tokio', 'serde', 'serde_json']
    },
    java: { 
      enabled: false, // Not needed for initial implementation
      version: '17',
      graalvm: true
    },
    php: { 
      enabled: false, // Not needed for initial implementation
      version: '8.2'
    },
    go: { 
      enabled: true, 
      version: '1.21',
      modules: ['github.com/gorilla/websocket']
    },
    zig: { 
      enabled: false, // Not needed for initial implementation
      version: '0.11'
    },
    cpp: { 
      enabled: true, 
      std: 'c++17',
      libraries: ['opencv', 'ffmpeg']
    }
  },
  
  build: {
    target: 'native', // 'native' | 'wasm' | 'both'
    optimization: 'release',
    compression: true,
    codeSigning: false
  },
  
  webview: {
    enabled: true,
    library: 'wry',
    size: { width: 1200, height: 800 }
  },
  
  memory: {
    sharedBuffers: true,
    referenceCounting: true,
    isolation: true
  },
  
  security: {
    sandboxing: true,
    codeSigning: false,
    permissions: {
      network: true,
      filesystem: true,
      native: true
    }
  }
}
