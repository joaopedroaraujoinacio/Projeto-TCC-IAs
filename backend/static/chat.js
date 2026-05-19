const form         = document.getElementById('chatForm');
const input        = document.getElementById('messageInput');
const messagesDiv  = document.getElementById('messages');
const sendBtn      = document.getElementById('sendBtn');
const ragChatTab   = document.getElementById('ragChatTab');
const webSearchTab = document.getElementById('webSearchTab');
const cloudTab     = document.getElementById('cloudTab');
const cloudBar     = document.getElementById('cloudBar');
const gptBtn       = document.getElementById('gptBtn');
const geminiBtn    = document.getElementById('geminiBtn');

let chatMode     = 'normal';
let cloudProvider = 'openai'; // tracks which cloud provider is selected
let conversationHistory = [];

input.addEventListener('input', function() {
    this.style.height = 'auto';
    const newHeight = Math.min(this.scrollHeight, 150);
    this.style.height = newHeight + 'px';
    this.style.overflowY = newHeight >= 150 ? 'auto' : 'hidden';
});

input.addEventListener('keydown', function(e) {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendBtn.click(); }
});

ragChatTab.addEventListener('click', () => {
    if (chatMode === 'rag') { chatMode = 'normal'; ragChatTab.classList.remove('active'); }
    else {
        chatMode = 'rag';
        ragChatTab.classList.add('active');
        webSearchTab.classList.remove('active');
        cloudTab.classList.remove('active');
        cloudBar.style.display = 'none';
    }
});

webSearchTab.addEventListener('click', () => {
    if (chatMode === 'websearch') { chatMode = 'normal'; webSearchTab.classList.remove('active'); }
    else {
        chatMode = 'websearch';
        webSearchTab.classList.add('active');
        ragChatTab.classList.remove('active');
        cloudTab.classList.remove('active');
        cloudBar.style.display = 'none';
    }
});

cloudTab.addEventListener('click', () => {
    if (chatMode === 'cloud') {
        chatMode = 'normal';
        cloudTab.classList.remove('active');
        cloudBar.style.display = 'none';
    } else {
        chatMode = 'cloud';
        cloudTab.classList.add('active');
        ragChatTab.classList.remove('active');
        webSearchTab.classList.remove('active');
        cloudBar.style.display = 'flex';
    }
});

gptBtn.addEventListener('click', () => {
    cloudProvider = 'openai';
    gptBtn.classList.add('active');
    geminiBtn.classList.remove('active');
});

geminiBtn.addEventListener('click', () => {
    cloudProvider = 'gemini';
    geminiBtn.classList.add('active');
    gptBtn.classList.remove('active');
});

document.getElementById('clearChatBtn').addEventListener('click', () => {
    if (confirm('Deseja realmente limpar toda a conversa?')) {
        conversationHistory = [];
        messagesDiv.innerHTML = '';
    }
});

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const message = input.value.trim();
    if (!message) return;

    const userMsg = document.createElement('div');
    userMsg.className = 'message user';
    userMsg.setAttribute('data-avatar', 'U');
    userMsg.innerHTML = `<div>${escapeHtml(message)}</div>`;
    messagesDiv.appendChild(userMsg);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;

    input.value = '';
    input.style.height = 'auto';
    input.style.overflowY = 'hidden';
    sendBtn.disabled = true;
    sendBtn.textContent = 'Enviando...';

    try {
        await handleStreamingChat(message);
    } catch (err) {
        const errorMsg = document.createElement('div');
        errorMsg.className = 'message assistant';
        errorMsg.innerHTML = `<div>Erro ao processar mensagem: ${err.message}</div>`;
        messagesDiv.appendChild(errorMsg);
        conversationHistory.push({ role: 'user', content: message });
        conversationHistory.push({ role: 'assistant', content: 'Error: ' + err.message });
    } finally {
        sendBtn.disabled = false;
        sendBtn.textContent = 'Enviar';
    }
});

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function buildEndpoint(message) {
    // resolve cloud provider endpoint
    const cloudURL = cloudProvider === 'gemini' ? '/api/chat/gemini' : '/api/chat/openai';
    const endpoints = {
        rag:      { url: '/api/chat/rag',       body: { message, history: conversationHistory.slice(0, -1) } },
        websearch:{ url: '/api/chat/web-search', body: { query: message } },
        cloud:    { url: cloudURL,               body: { message, history: conversationHistory.slice(0, -1) } },
        normal:   { url: '/api/chat',            body: { message, history: conversationHistory.slice(0, -1) } },
    };
    return endpoints[chatMode] || endpoints.normal;
}

async function handleStreamingChat(message) {
    conversationHistory.push({ role: 'user', content: message });

    const aiMsg = document.createElement('div');
    aiMsg.className = 'message streaming';
    aiMsg.innerHTML = `<div>
        <span class="streaming-content"></span>
        <div class="sources-section" style="display:none;"></div>
        <div class="token-stats streaming-stats"></div>
    </div>`;
    messagesDiv.appendChild(aiMsg);

    const contentSpan    = aiMsg.querySelector('.streaming-content');
    const sourcesSection = aiMsg.querySelector('.sources-section');
    const statsSpan      = aiMsg.querySelector('.streaming-stats');

    let fullResponse = '', tokenCount = 0, startTime = Date.now(), firstTokenTime = null;
    const { url, body } = buildEndpoint(message);

    try {
        const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);

        const reader = res.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            buffer += decoder.decode(value, { stream: true });
            const messages = buffer.split('\n\n');
            buffer = messages.pop() || '';

            for (const msg of messages) {
                if (!msg.trim()) continue;
                let eventType = 'message', data = '';
                for (const line of msg.split('\n')) {
                    if (line.startsWith('event:')) eventType = line.slice(6).trim();
                    else if (line.startsWith('data:')) { const d = line.slice(5); data += d.startsWith(' ') ? d.slice(1) : d; }
                }
                handleSSEEvent(eventType, data, {
                    aiMsg, contentSpan, sourcesSection, statsSpan, startTime,
                    firstTokenTime: () => firstTokenTime,
                    setFirstTokenTime: (t) => { firstTokenTime = t; },
                    tokenCount: () => tokenCount,
                    incrementToken: () => { tokenCount++; },
                    fullResponse: () => fullResponse,
                    appendResponse: (chunk) => { fullResponse += chunk; },
                });
            }
        }

        if (buffer.trim()) {
            let data = '';
            for (const line of buffer.split('\n')) {
                if (line.startsWith('data:')) { const d = line.slice(5); data += d.startsWith(' ') ? d.slice(1) : d; }
            }
            if (data) { fullResponse += data; contentSpan.textContent = fullResponse; }
        }

        if (fullResponse && !conversationHistory.some(m => m.role === 'assistant' && m.content === fullResponse)) {
            conversationHistory.push({ role: 'assistant', content: fullResponse });
            if (aiMsg.className === 'message streaming') {
                finalizeMessage(aiMsg, statsSpan, tokenCount, firstTokenTime, startTime, 'Stream completed');
            }
        }
    } catch (error) {
        contentSpan.innerHTML = `<span style="color:#ff4f4f;">Erro no streaming: ${error.message}</span>`;
        throw error;
    }
}

function handleSSEEvent(eventType, data, ctx) {
    const { aiMsg, contentSpan, sourcesSection, statsSpan, startTime } = ctx;

    if (eventType === 'message' && data) {
        if (!ctx.firstTokenTime()) ctx.setFirstTokenTime(Date.now());
        ctx.incrementToken();
        ctx.appendResponse(data);
        contentSpan.textContent = ctx.fullResponse();
        const elapsed = (Date.now() - ctx.firstTokenTime()) / 1000;
        const tokensPerSec = elapsed > 0 ? (ctx.tokenCount() / elapsed).toFixed(2) : '0.00';
        const timeToFirst = ((ctx.firstTokenTime() - startTime) / 1000).toFixed(2);
        statsSpan.innerHTML = `<span>${tokensPerSec} tok/sec</span><span>${ctx.tokenCount()} tokens</span><span>${timeToFirst}s to first token</span>`;
        aiMsg.parentElement && (aiMsg.parentElement.scrollTop = aiMsg.parentElement.scrollHeight);

    } else if (eventType === 'sources' && data) {
        try {
            const sources = JSON.parse(data);
            if (sources && sources.length > 0) {
                let html = '<div class="sources-title">📚 Fontes:</div>';
                sources.forEach((s, i) => {
                    html += `<div class="source-item"><strong>[${i+1}]</strong> <a href="${escapeHtml(s.url)}" target="_blank">${escapeHtml(s.title || s.url)}</a></div>`;
                });
                sourcesSection.innerHTML = html;
                sourcesSection.style.display = 'block';
            }
        } catch (e) { console.error('Error parsing sources:', e); }

    } else if (eventType === 'done') {
        conversationHistory.push({ role: 'assistant', content: ctx.fullResponse() });
        finalizeMessage(aiMsg, statsSpan, ctx.tokenCount(), ctx.firstTokenTime(), startTime, 'EOS Token Found');

    } else if (eventType === 'error') {
        throw new Error('Streaming error received from server');
    }
}

function finalizeMessage(aiMsg, statsSpan, tokenCount, firstTokenTime, startTime, stopReason) {
    const avgSpeed    = firstTokenTime ? (tokenCount / ((Date.now() - firstTokenTime) / 1000)).toFixed(2) : '0.00';
    const timeToFirst = firstTokenTime ? ((firstTokenTime - startTime) / 1000).toFixed(2) : '0.00';
    aiMsg.className = 'message assistant';
    statsSpan.innerHTML = `<span>${avgSpeed} tok/sec</span><span>${tokenCount} tokens</span><span>${timeToFirst}s to first token</span><span>Stop: ${stopReason}</span>`;
}

document.getElementById('openUpload').onclick = () => { document.getElementById('uploadBox').style.display = 'block'; };
document.getElementById('cancelUpload').onclick = () => { document.getElementById('uploadBox').style.display = 'none'; document.getElementById('uploadForm').reset(); };

document.getElementById('uploadForm').addEventListener('submit', async e => {
    e.preventDefault();
    const file = document.getElementById('fileInput').files[0];
    if (!file) return;
    const formData = new FormData();
    formData.append('file', file);
    try {
        await fetch('/api/rag', { method: 'POST', body: formData });
        const msg = document.createElement('div');
        msg.className = 'message assistant';
        msg.innerHTML = `<div>O arquivo "${escapeHtml(file.name)}" foi enviado!</div>`;
        messagesDiv.appendChild(msg);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
        document.getElementById('uploadBox').style.display = 'none';
        e.target.reset();
    } catch (error) { alert('Erro ao enviar arquivo: ' + error.message); }
});

document.getElementById('helpBtn').onclick   = () => { document.getElementById('helpBox').style.display = 'block'; };
document.getElementById('closeHelp').onclick = () => { document.getElementById('helpBox').style.display = 'none'; };

input.focus();
