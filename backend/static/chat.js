const form         = document.getElementById('chatForm');
const input        = document.getElementById('messageInput');
const messagesDiv  = document.getElementById('messages');
const sendBtn      = document.getElementById('sendBtn');
const normalTab    = document.getElementById('normalTab');
const ragChatTab   = document.getElementById('ragChatTab');
const webSearchTab = document.getElementById('webSearchTab');
const cloudTab     = document.getElementById('cloudTab');
const cloudBar     = document.getElementById('cloudBar');
const gptBtn       = document.getElementById('gptBtn');
const geminiBtn    = document.getElementById('geminiBtn');

let chatMode      = 'normal';
let cloudProvider = 'openai';
let conversationHistory = [];

const modeLabels  = { normal: 'Chat Local', rag: 'RAG', websearch: 'Web Search', cloud: 'Cloud AI' };
const modelLabels = { normal: 'llama3.2:3b', rag: 'llama3.2:3b', websearch: 'llama3.2:3b', openai: 'gpt-4.1-mini', gemini: 'gemini-2.5-flash' };
const protectedModes = ['rag', 'cloud'];

// ── auth helpers ────────────────────────────────────────────────
function getToken()   { return localStorage.getItem('token'); }
function clearToken() { localStorage.removeItem('token'); }
function isLoggedIn() { return !!getToken(); }

function authHeaders() {
    const h = { 'Content-Type': 'application/json' };
    const t = getToken();
    if (t) h['Authorization'] = 'Bearer ' + t;
    return h;
}

function handleUnauthorized(res) {
    if (res.status === 401) { clearToken(); refreshHeaderAuth(); return true; }
    return false;
}

function refreshHeaderAuth() {
    const loggedIn = isLoggedIn();
    document.getElementById('logoutBtn').style.display  = loggedIn ? '' : 'none';
    document.getElementById('loginBtn').style.display   = loggedIn ? 'none' : '';
    document.getElementById('signupBtn').style.display  = loggedIn ? 'none' : '';
    const userEl = document.getElementById('headerUser');
    userEl.textContent = '';
}

// ── status bar ──────────────────────────────────────────────────
function updateStatusBar() {
    const model = chatMode === 'cloud' ? modelLabels[cloudProvider] : (modelLabels[chatMode] || 'llama3.2:3b');
    document.getElementById('modeIndicator').textContent = 'Modo: ' + (modeLabels[chatMode] || 'Chat Local');
    document.getElementById('modelIndicator').textContent = model;
}

// ── welcome screen ──────────────────────────────────────────────
function showWelcome() {
    messagesDiv.innerHTML = `
        <div id="welcomeScreen">
            <div id="welcomeLogo">
                <img src="/static/images/farol.png" alt="FAROL" style="width:64px;opacity:.85;">
            </div>
            <p id="welcomeTitle">Como posso ajudar?</p>
            <div id="suggestions">
                <button class="suggestion-btn" data-icon="🤖">O que é um modelo de linguagem?</button>
                <button class="suggestion-btn" data-icon="📚">Como funciona o RAG?</button>
                <button class="suggestion-btn" data-icon="🔍">Como funciona a busca vetorial?</button>
                <button class="suggestion-btn" data-icon="⚡">Diferença entre modelos locais e cloud?</button>
            </div>
        </div>`;
    bindSuggestions();
}

function hideWelcome() {
    const ws = document.getElementById('welcomeScreen');
    if (ws) ws.remove();
}

// ── auth gate ───────────────────────────────────────────────────
function requireAuthPrompt(mode) {
    hideWelcome();
    const existing = document.getElementById('authPrompt');
    if (existing) return;
    const el = document.createElement('div');
    el.id = 'authPrompt';
    el.innerHTML = `
        <div id="authPromptBox">
            <p>O modo <strong>${modeLabels[mode]}</strong> requer uma conta.</p>
            <div style="display:flex;gap:10px;justify-content:center;margin-top:16px;">
                <button onclick="window.location.href='/login'"  class="prompt-login-btn">Entrar</button>
                <button onclick="window.location.href='/signup'" class="prompt-signup-btn">Criar conta</button>
            </div>
        </div>`;
    messagesDiv.appendChild(el);
}

function clearAuthPrompt() {
    const el = document.getElementById('authPrompt');
    if (el) el.remove();
}

// ── tab logic ───────────────────────────────────────────────────
function setTab(mode) {
    chatMode = mode;
    [normalTab, ragChatTab, webSearchTab, cloudTab].forEach(b => b && b.classList.remove('active'));
    cloudBar.style.display = 'none';

    if (mode === 'normal')    normalTab.classList.add('active');
    if (mode === 'rag')       ragChatTab.classList.add('active');
    if (mode === 'websearch') webSearchTab.classList.add('active');
    if (mode === 'cloud')     { cloudTab.classList.add('active'); cloudBar.style.display = 'flex'; }

    updateStatusBar();

    if (protectedModes.includes(mode) && !isLoggedIn()) {
        requireAuthPrompt(mode);
    } else {
        clearAuthPrompt();
        if (messagesDiv.children.length === 0) showWelcome();
    }
}

normalTab.addEventListener('click',    () => chatMode === 'normal'    ? setTab('normal')    : setTab('normal'));
ragChatTab.addEventListener('click',   () => chatMode === 'rag'       ? setTab('normal')    : setTab('rag'));
webSearchTab.addEventListener('click', () => chatMode === 'websearch' ? setTab('normal')    : setTab('websearch'));
cloudTab.addEventListener('click',     () => chatMode === 'cloud'     ? setTab('normal')    : setTab('cloud'));

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

// ── header ──────────────────────────────────────────────────────
document.getElementById('logoutBtn').addEventListener('click', () => {
    clearToken();
    refreshHeaderAuth();
    setTab('normal');
});

document.getElementById('clearChatBtn').addEventListener('click', () => {
    if (confirm('Limpar toda a conversa?')) {
        conversationHistory = [];
        clearAuthPrompt();
        showWelcome();
    }
});

// ── form submit ─────────────────────────────────────────────────
form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const message = input.value.trim();
    if (!message) return;

    if (protectedModes.includes(chatMode) && !isLoggedIn()) {
        requireAuthPrompt(chatMode);
        return;
    }

    hideWelcome();
    clearAuthPrompt();

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
        errorMsg.innerHTML = `<div>Erro: ${err.message}</div>`;
        messagesDiv.appendChild(errorMsg);
    } finally {
        sendBtn.disabled = false;
        sendBtn.textContent = 'Enviar';
    }
});

input.addEventListener('input', function() {
    this.style.height = 'auto';
    const h = Math.min(this.scrollHeight, 150);
    this.style.height = h + 'px';
    this.style.overflowY = h >= 150 ? 'auto' : 'hidden';
});
input.addEventListener('keydown', function(e) {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendBtn.click(); }
});

function escapeHtml(text) {
    const d = document.createElement('div');
    d.textContent = text;
    return d.innerHTML;
}

function buildEndpoint(message) {
    const cloudURL = cloudProvider === 'gemini' ? '/api/chat/gemini' : '/api/chat/openai';
    const endpoints = {
        rag:       { url: '/api/rag/rag_chat',       body: { message, history: conversationHistory.slice(0, -1) } },
        websearch: { url: '/api/chat/web-search', body: { query: message } },
        cloud:     { url: cloudURL,              body: { message, history: conversationHistory.slice(0, -1) } },
        normal:    { url: '/api/chat',            body: { message, history: conversationHistory.slice(0, -1) } },
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
        appendResponse:    (c) => { fullResponse += c; },
        getRealStats:      () => realStats,
        setRealStats:      (s) => { realStats = s; },
    };

    try {
        const res = await fetch(url, {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify(body),
        });

        if (handleUnauthorized(res)) {
            aiMsg.remove();
            requireAuthPrompt(chatMode);
            return;
        }
        if (!res.ok) throw new Error(`HTTP ${res.status}`);

        const reader  = res.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            buffer += decoder.decode(value, { stream: true });
            const parts = buffer.split('\n\n');
            buffer = parts.pop() || '';
            for (const part of parts) {
                if (!part.trim()) continue;
                let eventType = 'message', data = '';
                for (const line of part.split('\n')) {
                    if (line.startsWith('event:'))      eventType = line.slice(6).trim();
                    else if (line.startsWith('data:')) { const d = line.slice(5); data += d.startsWith(' ') ? d.slice(1) : d; }
                }
                handleSSEEvent(eventType, data, ctx);
            }
        }

        if (fullResponse && !conversationHistory.some(m => m.role === 'assistant' && m.content === fullResponse)) {
            conversationHistory.push({ role: 'assistant', content: fullResponse });
            if (aiMsg.className === 'message streaming')
                finalizeMessage(aiMsg, statsSpan, wordCount, firstTokenTime, startTime, 'Stream completed', realStats);
        }
    } catch (error) {
        contentSpan.innerHTML = `<span style="color:var(--error)">Erro: ${error.message}</span>`;
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
        const elapsed      = (Date.now() - ctx.firstTokenTime()) / 1000;
        const wordsPerSec  = elapsed > 0 ? (ctx.wordCount() / elapsed).toFixed(1) : '0.0';
        const timeToFirst  = ((ctx.firstTokenTime() - startTime) / 1000).toFixed(2);
        statsSpan.className = 'token-stats streaming-stats';
        statsSpan.innerHTML = `<span>~${wordsPerSec} palavras/s</span><span>~${ctx.wordCount()} palavras</span><span>${timeToFirst}s até 1ª resposta</span><span class="stat-label-approx">estimativa</span>`;
        aiMsg.parentElement && (aiMsg.parentElement.scrollTop = aiMsg.parentElement.scrollHeight);

    } else if (eventType === 'token_stats' && data) {
        try {
            const stats = JSON.parse(data);
            ctx.setRealStats(stats);
            const elapsed     = ctx.firstTokenTime() ? (Date.now() - ctx.firstTokenTime()) / 1000 : 0;
            const tokPerSec   = elapsed > 0 ? (stats.token_count / elapsed).toFixed(1) : '—';
            const timeToFirst = ctx.firstTokenTime() ? ((ctx.firstTokenTime() - startTime) / 1000).toFixed(2) : '—';
            statsSpan.className = 'token-stats streaming-stats real-stats';
            statsSpan.innerHTML = `<span class="stat-confirmed">✓</span><span>${tokPerSec} tok/s</span><span>${stats.token_count} tokens</span><span>${stats.prompt_tokens} prompt</span><span>${timeToFirst}s</span>`;
        } catch {}

    } else if (eventType === 'sources' && data) {
        try {
            const sources = JSON.parse(data);
            if (sources?.length) {
                let html = '<div class="sources-title">Fontes:</div>';
                sources.forEach((s, i) => { html += `<div class="source-item"><strong>[${i+1}]</strong> <a href="${escapeHtml(s.url)}" target="_blank">${escapeHtml(s.title || s.url)}</a></div>`; });
                sourcesSection.innerHTML = html;
                sourcesSection.style.display = 'block';
            }
        } catch {}

    } else if (eventType === 'done') {
        conversationHistory.push({ role: 'assistant', content: ctx.fullResponse() });
        finalizeMessage(ctx.aiMsg, statsSpan, ctx.wordCount(), ctx.firstTokenTime(), startTime, 'EOS Token Found', ctx.getRealStats());

    } else if (eventType === 'error') {
        console.error('[SSE error]', data);
        throw new Error(data || 'Streaming error');
    }
}

function finalizeMessage(aiMsg, statsSpan, wordCount, firstTokenTime, startTime, stopReason, realStats) {
    aiMsg.className = 'message assistant';
    const timeToFirst = firstTokenTime ? ((firstTokenTime - startTime) / 1000).toFixed(2) : '—';
    if (realStats) {
        const elapsed   = firstTokenTime ? (Date.now() - firstTokenTime) / 1000 : 0;
        const tokPerSec = elapsed > 0 ? (realStats.token_count / elapsed).toFixed(1) : '—';
        statsSpan.className = 'token-stats real-stats';
        statsSpan.innerHTML = `<span class="stat-confirmed">✓ dados reais</span><span>${tokPerSec} tok/s</span><span>${realStats.token_count} tokens</span><span>${realStats.prompt_tokens} prompt</span><span>${timeToFirst}s</span><span>Parada: ${stopReason}</span>`;
    } else {
        const avgSpeed = firstTokenTime ? (wordCount / ((Date.now() - firstTokenTime) / 1000)).toFixed(1) : '0.0';
        statsSpan.className = 'token-stats';
        statsSpan.innerHTML = `<span>~${avgSpeed} palavras/s</span><span>~${wordCount} palavras</span><span>${timeToFirst}s</span><span class="stat-label-approx">estimativa</span><span>Parada: ${stopReason}</span>`;
    }
}

// ── upload ──────────────────────────────────────────────────────
document.getElementById('openUpload').onclick = () => {
    if (!isLoggedIn()) { requireAuthPrompt('rag'); return; }
    document.getElementById('uploadBox').style.display = 'block';
};

document.getElementById('cancelUpload').onclick = () => {
    document.getElementById('uploadBox').style.display = 'none';
    document.getElementById('uploadForm').reset();
    document.getElementById('selectedFileName').textContent = '';
    document.getElementById('fileDropZone').classList.remove('has-file');
};

document.getElementById('fileInput').addEventListener('change', function() {
    const label = document.getElementById('selectedFileName');
    const zone  = document.getElementById('fileDropZone');
    if (this.files[0]) {
        label.textContent = '📄 ' + this.files[0].name;
        zone.classList.add('has-file');
    } else {
        label.textContent = '';
        zone.classList.remove('has-file');
    }
});

const dropZone = document.getElementById('fileDropZone');
dropZone.addEventListener('dragover',  e => { e.preventDefault(); dropZone.classList.add('drag-over'); });
dropZone.addEventListener('dragleave', () => dropZone.classList.remove('drag-over'));
dropZone.addEventListener('drop', e => {
    e.preventDefault();
    dropZone.classList.remove('drag-over');
    const file = e.dataTransfer.files[0];
    if (file) {
        document.getElementById('fileInput').files = e.dataTransfer.files;
        document.getElementById('selectedFileName').textContent = '📄 ' + file.name;
        dropZone.classList.add('has-file');
    }
});

document.getElementById('uploadForm').addEventListener('submit', async e => {
    e.preventDefault();
    const file      = document.getElementById('fileInput').files[0];
    const submitBtn = document.getElementById('uploadSubmitBtn');
    if (!file) { alert('Selecione um arquivo.'); return; }

    const allowed = ['.txt', '.md', '.csv', '.pdf'];
    const ext = '.' + file.name.split('.').pop().toLowerCase();
    if (!allowed.includes(ext)) {
        alert('Formato não suportado. Use: PDF, TXT, MD ou CSV.');
        return;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = 'Indexando...';

    const formData = new FormData();
    formData.append('file', file);

    const headers = {};
    const t = getToken();
    if (t) headers['Authorization'] = 'Bearer ' + t;

    try {
        const res = await fetch('/api/rag/add_rag_data', {
            method: 'POST',
            headers,       
            body: formData,
        });
        if (handleUnauthorized(res)) return;
        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.error || 'Upload falhou');
        }

        hideWelcome();
        const msg = document.createElement('div');
        msg.className = 'message assistant';
        msg.innerHTML = `<div>✓ "<strong>${escapeHtml(file.name)}</strong>" indexado no RAG com sucesso.</div>`;
        messagesDiv.appendChild(msg);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;

        document.getElementById('uploadBox').style.display = 'none';
        document.getElementById('uploadForm').reset();
        document.getElementById('selectedFileName').textContent = '';
        document.getElementById('fileDropZone').classList.remove('has-file');
    } catch (err) {
        alert('Erro: ' + err.message);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Indexar';
    }
});

// ── help ────────────────────────────────────────────────────────
document.getElementById('helpBtn').onclick   = () => { document.getElementById('helpBox').style.display = 'block'; };
document.getElementById('closeHelp').onclick = () => { document.getElementById('helpBox').style.display = 'none'; };

function bindSuggestions() {
    document.querySelectorAll('.suggestion-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            input.value = btn.textContent.trim();
            input.dispatchEvent(new Event('input'));
            sendBtn.click();
        });
    });
}

// ── init ────────────────────────────────────────────────────────
refreshHeaderAuth();
setTab('normal');
showWelcome();
updateStatusBar();
window.addEventListener('beforeunload', (e) => {
  if (conversationHistory. length === 0) return;
  e.preventDefault();
  e.preventDefault() = '';
});
window.addEventListener('unload', () => {
  conversationHistory = [];
});
input.focus();
