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

let chatMode      = 'normal';
let cloudProvider = 'openai';
let conversationHistory = [];

const modeLabels = {
    normal:    'Chat Local',
    rag:       'RAG',
    websearch: 'Web Search',
    cloud:     'Cloud AI',
};
const modelLabels = {
    normal:    'llama3.2:3b',
    rag:       'llama3.2:3b',
    websearch: 'llama3.2:3b',
    openai:    'gpt-4.1-mini',
    gemini:    'gemini-2.5-flash',
};

function getToken()    { return localStorage.getItem('token'); }
function clearToken()  { localStorage.removeItem('token'); }
function authHeaders() {
    const h = { 'Content-Type': 'application/json' };
    const t = getToken();
    if (t) h['Authorization'] = 'Bearer ' + t;
    return h;
}

function handleUnauthorized(res) {
    if (res.status === 401) { clearToken(); window.location.href = '/login'; return true; }
    return false;
}

function updateStatusBar() {
    const model = chatMode === 'cloud' ? modelLabels[cloudProvider] : (modelLabels[chatMode] || 'llama3.2:3b');
    document.getElementById('modeIndicator').textContent  = 'Modo: ' + (modeLabels[chatMode] || 'Chat Local');
    document.getElementById('modelIndicator').textContent = model;
}

function hideWelcome() {
    const ws = document.getElementById('welcomeScreen');
    if (ws) ws.style.display = 'none';
}

function showWelcome() {
    messagesDiv.innerHTML = `
        <div id="welcomeScreen">
            <p id="welcomeTitle">Farol, como posso ajudar?</p>
            <div id="suggestions">
                <button class="suggestion-btn">O que é RAG?</button>
                <button class="suggestion-btn">O que é um modelo de linguagem?</button>
                <button class="suggestion-btn">Qual a diferença entre modelos locais e cloud?</button>
                <button class="suggestion-btn">Como funciona a busca vetorial?</button>
            </div>
        </div>
    `;
    bindSuggestions();
}

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
    updateStatusBar();
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
    updateStatusBar();
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
    updateStatusBar();
});

gptBtn.addEventListener('click', () => {
    cloudProvider = 'openai';
    gptBtn.classList.add('active');
    geminiBtn.classList.remove('active');
    updateStatusBar();
});

geminiBtn.addEventListener('click', () => {
    cloudProvider = 'gemini';
    geminiBtn.classList.add('active');
    gptBtn.classList.remove('active');
    updateStatusBar();
});

document.getElementById('clearChatBtn').addEventListener('click', () => {
    if (confirm('Deseja realmente limpar toda a conversa?')) {
        conversationHistory = [];
        showWelcome();
    }
});

document.getElementById('logoutBtn').addEventListener('click', () => {
    clearToken();
    window.location.href = '/login';
});

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const message = input.value.trim();
    if (!message) return;

    hideWelcome();

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
    const cloudURL = cloudProvider === 'gemini' ? '/api/chat/gemini' : '/api/chat/openai';
    const endpoints = {
        rag:       { url: '/api/chat/rag',        body: { message, history: conversationHistory.slice(0, -1) } },
        websearch: { url: '/api/chat/web-search',  body: { query: message } },
        cloud:     { url: cloudURL,                body: { message, history: conversationHistory.slice(0, -1) } },
        normal:    { url: '/api/chat',             body: { message, history: conversationHistory.slice(0, -1) } },
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

    let fullResponse = '', wordCount = 0, startTime = Date.now(), firstTokenTime = null, realStats = null;
    const { url, body } = buildEndpoint(message);

    const ctx = {
        aiMsg, contentSpan, sourcesSection, statsSpan, startTime,
        firstTokenTime:    () => firstTokenTime,
        setFirstTokenTime: (t) => { firstTokenTime = t; },
        wordCount:         () => wordCount,
        incrementWord:     () => { wordCount++; },
        fullResponse:      () => fullResponse,
        appendResponse:    (chunk) => { fullResponse += chunk; },
        getRealStats:      () => realStats,
        setRealStats:      (s) => { realStats = s; },
    };

    try {
        const res = await fetch(url, {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify(body)
        });

        if (handleUnauthorized(res)) return;
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
                handleSSEEvent(eventType, data, ctx);
            }
        }

        if (fullResponse && !conversationHistory.some(m => m.role === 'assistant' && m.content === fullResponse)) {
            conversationHistory.push({ role: 'assistant', content: fullResponse });
            if (aiMsg.className === 'message streaming') {
                finalizeMessage(aiMsg, statsSpan, ctx.wordCount(), firstTokenTime, startTime, 'Stream completed', ctx.getRealStats());
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
        ctx.incrementWord();
        ctx.appendResponse(data);
        contentSpan.textContent = ctx.fullResponse();
        const elapsed = (Date.now() - ctx.firstTokenTime()) / 1000;
        const wordsPerSec = elapsed > 0 ? (ctx.wordCount() / elapsed).toFixed(1) : '0.0';
        const timeToFirst = ((ctx.firstTokenTime() - startTime) / 1000).toFixed(2);
        statsSpan.className = 'token-stats streaming-stats';
        statsSpan.innerHTML = `
            <span>~${wordsPerSec} palavras/s</span>
            <span>~${ctx.wordCount()} palavras</span>
            <span>${timeToFirst}s até 1ª resposta</span>
            <span class="stat-label-approx">estimativa</span>
        `;
        aiMsg.parentElement && (aiMsg.parentElement.scrollTop = aiMsg.parentElement.scrollHeight);

    } else if (eventType === 'token_stats' && data) {
        try {
            const stats = JSON.parse(data);
            ctx.setRealStats(stats);
            const elapsed = ctx.firstTokenTime() ? (Date.now() - ctx.firstTokenTime()) / 1000 : 0;
            const tokPerSec = elapsed > 0 ? (stats.token_count / elapsed).toFixed(1) : '—';
            const timeToFirst = ctx.firstTokenTime() ? ((ctx.firstTokenTime() - startTime) / 1000).toFixed(2) : '—';
            statsSpan.className = 'token-stats streaming-stats real-stats';
            statsSpan.innerHTML = `
                <span class="stat-confirmed">✓</span>
                <span>${tokPerSec} tok/s</span>
                <span>${stats.token_count} tokens gerados</span>
                <span>${stats.prompt_tokens} tokens prompt</span>
                <span>${timeToFirst}s até 1ª resposta</span>
            `;
        } catch (e) { console.error('Erro ao parsear token_stats:', e); }

    } else if (eventType === 'sources' && data) {
        try {
            const sources = JSON.parse(data);
            if (sources && sources.length > 0) {
                let html = '<div class="sources-title">Fontes:</div>';
                sources.forEach((s, i) => {
                    html += `<div class="source-item"><strong>[${i+1}]</strong> <a href="${escapeHtml(s.url)}" target="_blank">${escapeHtml(s.title || s.url)}</a></div>`;
                });
                sourcesSection.innerHTML = html;
                sourcesSection.style.display = 'block';
            }
        } catch (e) {}

    } else if (eventType === 'done') {
        conversationHistory.push({ role: 'assistant', content: ctx.fullResponse() });
        finalizeMessage(ctx.aiMsg, statsSpan, ctx.wordCount(), ctx.firstTokenTime(), startTime, 'EOS Token Found', ctx.getRealStats());

    } else if (eventType === 'error') {
        console.error('[SSE error]', data);
        throw new Error(data || 'Streaming error received from server');
    }
}

function finalizeMessage(aiMsg, statsSpan, wordCount, firstTokenTime, startTime, stopReason, realStats) {
    aiMsg.className = 'message assistant';
    const timeToFirst = firstTokenTime ? ((firstTokenTime - startTime) / 1000).toFixed(2) : '—';
    if (realStats) {
        const elapsed = firstTokenTime ? (Date.now() - firstTokenTime) / 1000 : 0;
        const tokPerSec = elapsed > 0 ? (realStats.token_count / elapsed).toFixed(1) : '—';
        statsSpan.className = 'token-stats real-stats';
        statsSpan.innerHTML = `
            <span class="stat-confirmed">✓ dados reais</span>
            <span>${tokPerSec} tok/s</span>
            <span>${realStats.token_count} tokens gerados</span>
            <span>${realStats.prompt_tokens} tokens prompt</span>
            <span>${timeToFirst}s até 1ª resposta</span>
            <span>Parada: ${stopReason}</span>
        `;
    } else {
        const avgSpeed = firstTokenTime ? (wordCount / ((Date.now() - firstTokenTime) / 1000)).toFixed(1) : '0.0';
        statsSpan.className = 'token-stats';
        statsSpan.innerHTML = `
            <span>~${avgSpeed} palavras/s</span>
            <span>~${wordCount} palavras</span>
            <span>${timeToFirst}s até 1ª resposta</span>
            <span class="stat-label-approx">estimativa</span>
            <span>Parada: ${stopReason}</span>
        `;
    }
}

document.getElementById('openUpload').onclick  = () => { document.getElementById('uploadBox').style.display = 'block'; };
document.getElementById('cancelUpload').onclick = () => {
    document.getElementById('uploadBox').style.display = 'none';
    document.getElementById('uploadForm').reset();
};

document.getElementById('uploadForm').addEventListener('submit', async e => {
    e.preventDefault();
    const file = document.getElementById('fileInput').files[0];
    if (!file) return;

    // read file as text and send as JSON with auth token
    const text = await file.text();
    try {
        const res = await fetch('/api/rag/add_rag_data', {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify({ content: text, content_name: file.name }),
        });
        if (handleUnauthorized(res)) return;
        if (!res.ok) throw new Error('Upload falhou: ' + res.status);

        hideWelcome();
        const msg = document.createElement('div');
        msg.className = 'message assistant';
        msg.innerHTML = `<div>O arquivo "<strong>${escapeHtml(file.name)}</strong>" foi enviado e indexado com sucesso.</div>`;
        messagesDiv.appendChild(msg);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
        document.getElementById('uploadBox').style.display = 'none';
        e.target.reset();
    } catch (error) {
        alert('Erro ao enviar arquivo: ' + error.message);
    }
});

document.getElementById('helpBtn').onclick   = () => { document.getElementById('helpBox').style.display = 'block'; };
document.getElementById('closeHelp').onclick = () => { document.getElementById('helpBox').style.display = 'none'; };

function bindSuggestions() {
    document.querySelectorAll('.suggestion-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            input.value = btn.textContent;
            input.dispatchEvent(new Event('input'));
            sendBtn.click();
        });
    });
}

bindSuggestions();
updateStatusBar();
input.focus();
