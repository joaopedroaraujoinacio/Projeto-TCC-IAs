<template>
  <div class="chat-container">
    <h1>AI Chat Interface</h1>
    
    <!-- Chat Messages Display -->
    <div class="messages-container" ref="messagesContainer">
      <div 
        v-for="(msg, index) in messages" 
        :key="index" 
        :class="['message', msg.type]"
      >
        <div class="message-content">
          <strong>{{ msg.type === 'user' ? 'You:' : 'AI:' }}</strong>
          <div class="message-text" v-html="formatMessage(msg.content)"></div>
          <small class="timestamp">{{ msg.timestamp }}</small>
        </div>
      </div>
      
      <!-- Loading indicator -->
      <div v-if="isLoading" class="message ai loading">
        <div class="message-content">
          <strong>AI:</strong>
          <div class="typing-indicator">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      </div>
    </div>

    <!-- Input Form -->
    <div class="input-container">
      <form @submit.prevent="sendMessage" class="message-form">
        <input 
          v-model="currentMessage" 
          type="text" 
          placeholder="Type your message..." 
          :disabled="isLoading"
          class="message-input"
          ref="messageInput"
        />
        <button 
          type="submit" 
          :disabled="isLoading || !currentMessage.trim()"
          class="send-button"
        >
          {{ isLoading ? 'Sending...' : 'Send' }}
        </button>
      </form>
    </div>

    <!-- Connection Status -->
    <div class="status-bar">
      <span :class="['status-indicator', connectionStatus]"></span>
      {{ connectionStatus === 'connected' ? 'Connected to AI' : 'Disconnected' }}
      <span v-if="lastModel" class="model-info">| Model: {{ lastModel }}</span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChatInterface',
  data() {
    return {
      messages: [],
      currentMessage: '',
      isLoading: false,
      connectionStatus: 'disconnected',
      lastModel: null,
      apiUrl: 'https://localhost/api/chat'
    }
  },
  mounted() {
    // Test connection on component mount
    this.testConnection();
    // Focus on input
    this.$refs.messageInput?.focus();
  },
  methods: {
    async sendMessage() {
      if (!this.currentMessage.trim() || this.isLoading) return;

      const userMessage = this.currentMessage.trim();
      this.currentMessage = '';

      // Add user message to chat
      this.addMessage(userMessage, 'user');

      this.isLoading = true;

      try {
        const response = await fetch(this.apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            message: userMessage
          })
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (data.success && data.data) {
          // Add AI response to chat
          this.addMessage(data.data.response, 'ai');
          this.lastModel = data.data.model;
          this.connectionStatus = 'connected';
        } else {
          throw new Error('Invalid response format');
        }

      } catch (error) {
        console.error('Error sending message:', error);
        this.addMessage(`Error: ${error.message}`, 'error');
        this.connectionStatus = 'disconnected';
      } finally {
        this.isLoading = false;
        this.scrollToBottom();
        // Focus back on input
        this.$nextTick(() => {
          this.$refs.messageInput?.focus();
        });
      }
    },

    addMessage(content, type) {
      this.messages.push({
        content,
        type,
        timestamp: new Date().toLocaleTimeString()
      });
      this.$nextTick(() => {
        this.scrollToBottom();
      });
    },

    formatMessage(message) {
      // Convert \n to <br> tags for proper line breaks
      return message.replace(/\n/g, '<br>');
    },

    scrollToBottom() {
      const container = this.$refs.messagesContainer;
      if (container) {
        container.scrollTop = container.scrollHeight;
      }
    },

    async testConnection() {
      try {
        const response = await fetch(this.apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            message: 'Connection test'
          })
        });

        if (response.ok) {
          this.connectionStatus = 'connected';
        }
      } catch (error) {
        console.log('Connection test failed:', error);
        this.connectionStatus = 'disconnected';
      }
    },

    clearChat() {
      this.messages = [];
    }
  }
}
</script>

<style scoped>
.chat-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

h1 {
  text-align: center;
  color: #333;
  margin-bottom: 20px;
}

.messages-container {
  flex: 1;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  overflow-y: auto;
  background-color: #f9f9f9;
  margin-bottom: 20px;
}

.message {
  margin-bottom: 15px;
}

.message.user .message-content {
  background-color: #007bff;
  color: white;
  margin-left: 20%;
  border-radius: 18px 18px 5px 18px;
}

.message.ai .message-content {
  background-color: #e9ecef;
  color: #333;
  margin-right: 20%;
  border-radius: 18px 18px 18px 5px;
}

.message.error .message-content {
  background-color: #dc3545;
  color: white;
  margin-right: 20%;
  border-radius: 18px 18px 18px 5px;
}

.message-content {
  padding: 12px 16px;
  word-wrap: break-word;
}

.message-text {
  margin: 5px 0;
  line-height: 1.4;
}

.timestamp {
  opacity: 0.7;
  font-size: 0.8em;
}

.input-container {
  margin-bottom: 10px;
}

.message-form {
  display: flex;
  gap: 10px;
}

.message-input {
  flex: 1;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 16px;
}

.message-input:focus {
  outline: none;
  border-color: #007bff;
}

.message-input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.send-button {
  padding: 12px 24px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
}

.send-button:hover:not(:disabled) {
  background-color: #0056b3;
}

.send-button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.status-bar {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #666;
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 8px;
}

.status-indicator.connected {
  background-color: #28a745;
}

.status-indicator.disconnected {
  background-color: #dc3545;
}

.model-info {
  margin-left: 10px;
  font-weight: 500;
}

/* Loading animation */
.loading .typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #666;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

/* Responsive design */
@media (max-width: 600px) {
  .chat-container {
    padding: 10px;
  }
  
  .message.user .message-content {
    margin-left: 10%;
  }
  
  .message.ai .message-content,
  .message.error .message-content {
    margin-right: 10%;
  }
}
</style>
