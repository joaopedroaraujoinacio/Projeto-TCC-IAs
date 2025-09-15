<template>
  <div class="app">
    <h1>🤖 RAG API Test</h1>
    
    <div class="test-section">
      <h2>Hello Route Test</h2>
      <button @click="testHello" :disabled="loading">
        {{ loading ? 'Loading...' : 'Test /api/hello' }}
      </button>
      
      <div v-if="error" class="error">
        <h3>❌ Error:</h3>
        <pre>{{ error }}</pre>
      </div>
      
      <div v-if="response" class="response">
        <h3>✅ Response from /api/hello:</h3>
        <pre>{{ JSON.stringify(response, null, 2) }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const loading = ref(false)
const response = ref(null)
const error = ref(null)

const testHello = async () => {
  loading.value = true
  error.value = null
  response.value = null
  
  try {
    const result = await axios.get('/api/hello')
    response.value = result.data
    console.log('Hello route response:', result.data)
  } catch (err) {
    error.value = {
      message: err.message,
      status: err.response?.status,
      statusText: err.response?.statusText,
      data: err.response?.data
    }
    console.error('Hello route error:', err)
  } finally {
    loading.value = false
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f5f5f5;
  color: #333;
}

.app {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

h1 {
  text-align: center;
  margin-bottom: 2rem;
  color: #2c3e50;
}

.test-section {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

h2 {
  margin-bottom: 1rem;
  color: #34495e;
}

button {
  background: #3498db;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  transition: background 0.2s;
}

button:hover:not(:disabled) {
  background: #2980b9;
}

button:disabled {
  background: #bdc3c7;
  cursor: not-allowed;
}

.response {
  margin-top: 1.5rem;
  padding: 1rem;
  background: #d5f4e6;
  border-radius: 6px;
  border-left: 4px solid #27ae60;
}

.error {
  margin-top: 1.5rem;
  padding: 1rem;
  background: #fadbd8;
  border-radius: 6px;
  border-left: 4px solid #e74c3c;
}

pre {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 4px;
  overflow-x: auto;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 14px;
  margin-top: 0.5rem;
}

h3 {
  margin-bottom: 0.5rem;
}
</style>
