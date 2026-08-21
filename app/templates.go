package main

import "html/template"

// ── Templates ─────────────────────────────────────────────────────────────────

var loginTmpl = template.Must(template.New("login").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>LogVault · Sign In</title>
  <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>⬡</text></svg>">
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&display=swap" rel="stylesheet">
  <style>
    :root{--bg:#0c0c0c;--surface:#161616;--border:#242424;--accent:#10b981;--text:#d4d4d8;--muted:#52525b;--danger:#ef4444}
    *{margin:0;padding:0;box-sizing:border-box}
    body{background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace;min-height:100vh;display:flex;align-items:center;justify-content:center}
    .card{width:100%;max-width:380px;margin:1rem;background:var(--surface);border:1px solid var(--border);border-radius:8px;padding:2.2rem}
    .logo-row{display:flex;align-items:center;gap:.7rem;margin-bottom:.3rem}
    .logo-icon{width:34px;height:34px;background:var(--accent);border-radius:6px;display:flex;align-items:center;justify-content:center;font-size:.95rem;color:#000;font-weight:700}
    h1{font-weight:700;font-size:1.4rem;color:var(--text)}
    .tagline{font-size:.7rem;color:var(--muted);margin-bottom:1.8rem}
    label{display:block;font-size:.68rem;color:var(--muted);letter-spacing:.08em;text-transform:uppercase;margin-bottom:.35rem}
    input{width:100%;padding:.65rem .85rem;background:var(--bg);border:1px solid var(--border);border-radius:5px;color:var(--text);font-family:'JetBrains Mono',monospace;font-size:.82rem;outline:none;transition:border-color .15s;margin-bottom:1rem}
    input:focus{border-color:var(--accent)}
    button{width:100%;padding:.75rem;background:var(--accent);border:none;border-radius:5px;color:#000;font-family:'JetBrains Mono',monospace;font-size:.82rem;font-weight:700;cursor:pointer;transition:opacity .15s;margin-top:.2rem}
    button:hover{opacity:.85}
    .error{background:rgba(239,68,68,.08);border:1px solid rgba(239,68,68,.25);color:var(--danger);border-radius:5px;padding:.6rem .85rem;font-size:.75rem;margin-bottom:1rem}
    .hint{font-size:.65rem;color:var(--muted);text-align:center;margin-top:1.2rem}
  </style>
</head>
<body>
  <div class="card">
    <div class="logo-row"><div class="logo-icon">⬡</div><h1>LogVault</h1></div>
    <p class="tagline">sign in to continue</p>
    {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
    <form method="POST" action="{{.LoginAction}}">
      <label>Username</label>
      <input type="text" name="username" autocomplete="username" autofocus placeholder="username">
      <label>Password</label>
      <input type="password" name="password" autocomplete="current-password" placeholder="password">
      <button type="submit">Sign In</button>
    </form>
    <p class="hint">session expires after 8 hours</p>
  </div>
</body>
</html>`))

var funcMap = template.FuncMap{"notDir": func(b bool) bool { return !b }}

var browserTmpl = template.Must(template.New("browser").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>LogVault{{if .SubPath}} · /{{.SubPath}}{{end}}</title>
  <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>⬡</text></svg>">
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&display=swap" rel="stylesheet">
  <style>
    :root{--bg:#0c0c0c;--surface:#161616;--border:#242424;--accent:#10b981;--text:#d4d4d8;--muted:#52525b;--folder:#d97706;--danger:#ef4444}
    *{margin:0;padding:0;box-sizing:border-box}
    body{background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace;min-height:100vh}
    .container{max-width:960px;margin:0 auto;padding:2rem 1.5rem}
    .topbar{display:flex;align-items:center;justify-content:space-between;margin-bottom:1.75rem}
    .logo-row{display:flex;align-items:center;gap:.7rem}
    .logo-icon{width:32px;height:32px;background:var(--accent);border-radius:5px;display:flex;align-items:center;justify-content:center;font-size:.88rem;color:#000;font-weight:700}
    h1{font-weight:700;font-size:1.3rem;color:var(--text)}
    .logout-btn{padding:.35rem .8rem;background:transparent;border:1px solid var(--border);border-radius:4px;color:var(--muted);font-family:'JetBrains Mono',monospace;font-size:.68rem;cursor:pointer;transition:all .15s}
    .logout-btn:hover{border-color:var(--danger);color:var(--danger)}
    .breadcrumb{display:flex;align-items:center;gap:.35rem;flex-wrap:wrap;padding:.55rem .9rem;background:var(--surface);border:1px solid var(--border);border-radius:5px;font-size:.72rem;margin-bottom:1rem}
    .breadcrumb a{color:var(--accent);text-decoration:none}
    .breadcrumb a:hover{text-decoration:underline}
    .breadcrumb .sep{color:var(--muted)}
    .breadcrumb .current{color:var(--text)}
    .status-bar{display:flex;align-items:center;gap:.6rem;margin-bottom:1.25rem;padding:.5rem .9rem;background:var(--surface);border:1px solid var(--border);border-radius:5px;font-size:.7rem;color:var(--muted)}
    .dot{width:5px;height:5px;border-radius:50%;background:var(--accent)}
    .status-bar span{color:var(--text)}
    .tag{display:inline-block;padding:.1rem .4rem;background:rgba(16,185,129,.08);border:1px solid rgba(16,185,129,.2);border-radius:3px;color:var(--accent);font-size:.63rem}
    .log-table{background:var(--surface);border:1px solid var(--border);border-radius:7px;overflow:hidden}
    .table-header{display:grid;grid-template-columns:1fr 90px 150px auto;padding:.55rem 1rem;border-bottom:1px solid var(--border);font-size:.6rem;text-transform:uppercase;letter-spacing:.1em;color:var(--muted)}
    .log-row{display:grid;grid-template-columns:1fr 90px 150px auto;padding:.75rem 1rem;border-bottom:1px solid var(--border);align-items:center;transition:background .12s;position:relative}
    .log-row:last-child{border-bottom:none}
    .log-row:hover{background:rgba(255,255,255,.02)}
    .log-row::before{content:'';position:absolute;left:0;top:0;bottom:0;width:2px;background:var(--accent);opacity:0;transition:opacity .12s}
    .log-row:hover::before{opacity:.6}
    .log-row.is-dir::before{background:var(--folder)}
    .entry-name{display:flex;align-items:center;gap:.55rem;font-size:.8rem;overflow:hidden}
    .entry-icon{flex-shrink:0;font-size:.8rem;width:20px;text-align:center;color:var(--muted)}
    .entry-link{white-space:nowrap;overflow:hidden;text-overflow:ellipsis;text-decoration:none;color:inherit;transition:color .12s}
    .entry-link:hover{color:var(--accent)}
    .dir-link{color:var(--folder)}
    .dir-link:hover{color:#fbbf24}
    .up-link{color:var(--muted)}
    .file-size{font-size:.72rem;color:var(--muted)}
    .file-modified{font-size:.7rem;color:var(--muted)}
    .actions{display:flex;align-items:center;gap:.4rem}
    .download-btn,.tail-btn,.view-btn{display:inline-flex;align-items:center;gap:.3rem;padding:.28rem .65rem;background:transparent;border-radius:4px;font-family:'JetBrains Mono',monospace;font-size:.67rem;text-decoration:none;transition:all .12s;white-space:nowrap}
    .view-btn{border:1px solid var(--border);color:var(--muted)}
    .view-btn:hover{border-color:var(--text);color:var(--text)}
    .download-btn{border:1px solid var(--border);color:var(--muted)}
    .download-btn:hover{border-color:var(--accent);color:var(--accent)}
    .tail-btn{border:1px solid rgba(16,185,129,.3);color:var(--accent)}
    .tail-btn:hover{background:rgba(16,185,129,.15);border-color:var(--accent)}
    .th-sort{cursor:pointer;user-select:none;display:flex;align-items:center;gap:.3rem;transition:color .12s}
    .th-sort:hover{color:var(--text)}
    .sort-icon{font-size:.6rem;opacity:.45;transition:opacity .12s}
    .th-sort.active .sort-icon{opacity:1;color:var(--accent)}
    .empty-state{text-align:center;padding:4rem 2rem;color:var(--muted)}
    .empty-state h3{font-size:1rem;color:var(--text);margin-bottom:.4rem}
    .empty-state p{font-size:.75rem}
    .footer{margin-top:1.5rem;text-align:center;font-size:.65rem;color:var(--muted)}
    @media(max-width:640px){.table-header,.log-row{grid-template-columns:1fr 70px 130px}.table-header>*:last-child,.log-row>*:last-child{display:none}}
  </style>
</head>
<body>
  <div class="container">
    <div class="topbar">
      <div class="logo-row"><div class="logo-icon">⬡</div><h1>LogVault</h1></div>
      {{if .AuthEnabled}}
      <form method="POST" action="{{.LogoutAction}}" style="margin:0">
        <button class="logout-btn" type="submit">Sign out</button>
      </form>
      {{end}}
    </div>

    <div class="breadcrumb">
      <a href="{{.BrowseRoot}}">root</a>
      {{range .Crumbs}}
        <span class="sep">/</span>
        {{if .IsLast}}<span class="current">{{.Name}}</span>
        {{else}}<a href="{{.URL}}">{{.Name}}</a>{{end}}
      {{end}}
    </div>

    <div class="status-bar">
      <div class="dot"></div>
      <span>{{.Count}} item{{if ne .Count 1}}s{{end}}</span>
      &nbsp;·&nbsp;<span class="tag">/app/logs{{if .SubPath}}/{{.SubPath}}{{end}}</span>
    </div>

    {{if .Entries}}
    <div class="log-table" id="tbl">
      <div class="table-header">
        <div class="th-sort active" id="th-name" onclick="sortBy('name')">Name<span class="sort-icon" id="si-name">▲</span></div>
        <div class="th-sort" id="th-size" onclick="sortBy('size')">Size<span class="sort-icon" id="si-size">⇅</span></div>
        <div class="th-sort" id="th-modified" onclick="sortBy('modified')">Modified<span class="sort-icon" id="si-modified">⇅</span></div>
        <div>Actions</div>
      </div>
      {{if .SubPath}}
      <div class="log-row" data-parent="1">
        <div class="entry-name">
          <div class="entry-icon">↑</div>
          <a class="entry-link up-link" href="{{.ParentURL}}">..</a>
        </div>
        <div></div><div></div><div></div>
      </div>
      {{end}}
      {{range .Entries}}
      <div class="log-row{{if .IsDir}} is-dir{{end}}" data-name="{{.Name}}" data-size="{{.SizeBytes}}" data-modified="{{.Modified}}" data-isdir="{{if .IsDir}}1{{else}}0{{end}}">
        <div class="entry-name">
          <div class="entry-icon">{{if .IsDir}}▶{{else}}≡{{end}}</div>
          {{if .IsDir}}
            <a class="entry-link dir-link" href="{{.BrowseURL}}">{{.Name}}</a>
          {{else}}
            <span class="entry-link" style="cursor:default">{{.Name}}</span>
          {{end}}
        </div>
        <div class="file-size">{{if notDir .IsDir}}{{.Size}}{{else}}—{{end}}</div>
        <div class="file-modified">{{.Modified}}</div>
        <div class="actions">
          {{if notDir .IsDir}}
          <a class="view-btn" href="{{.ViewURL}}">≡ View</a>
          <a class="tail-btn" href="{{.TailURL}}">⊞ Tail</a>
          <a class="download-btn" href="{{.DownloadURL}}">↓ Download</a>
          {{end}}
        </div>
      </div>
      {{end}}
    </div>
    {{else}}
    <div class="empty-state">
      <h3>Empty directory</h3>
      <p>No files or folders here.</p>
    </div>
    {{end}}

    <div class="footer"><a href="https://github.com/lokendrabhat/logvault" target="_blank" rel="noopener noreferrer"><strong>LogVault</strong></a> · Minimal Log Server · by <a href="https://lokendrabhat.com.np" target="_blank" rel="noopener noreferrer"><strong>Lokendra Bhat</strong></a></div>
  </div>
  <script>
    var sortCol = 'name', sortDir = 1;
    function sortBy(col) {
      if (sortCol === col) sortDir *= -1; else { sortCol = col; sortDir = 1; }
      ['name','size','modified'].forEach(function(c) {
        var th = document.getElementById('th-' + c);
        var si = document.getElementById('si-' + c);
        if (!th || !si) return;
        var active = c === sortCol;
        th.classList.toggle('active', active);
        si.textContent = active ? (sortDir === 1 ? '▲' : '▼') : '⇅';
      });
      var tbl = document.getElementById('tbl');
      if (!tbl) return;
      var parent = tbl.querySelector('[data-parent]');
      var rows = Array.from(tbl.querySelectorAll('.log-row:not([data-parent])'));
      rows.sort(function(a, b) {
        var aDir = a.dataset.isdir === '1', bDir = b.dataset.isdir === '1';
        if (aDir !== bDir) return aDir ? -1 : 1;
        if (sortCol === 'size') return sortDir * (parseInt(a.dataset.size||0) - parseInt(b.dataset.size||0));
        if (sortCol === 'modified') return sortDir * a.dataset.modified.localeCompare(b.dataset.modified);
        return sortDir * a.dataset.name.toLowerCase().localeCompare(b.dataset.name.toLowerCase());
      });
      if (parent) tbl.insertBefore(parent, tbl.children[1]);
      rows.forEach(function(r) { tbl.appendChild(r); });
    }
  </script>
</body>
</html>`))

var tailTmpl = template.Must(template.New("tail").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>tail · {{.FileName}}</title>
  <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>⬡</text></svg>">
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&display=swap" rel="stylesheet">
  <style>
    :root{--bg:#0c0c0c;--surface:#141414;--border:#222;--accent:#10b981;--text:#d4d4d8;--muted:#52525b;--danger:#ef4444}
    *{margin:0;padding:0;box-sizing:border-box}
    html,body{height:100%;background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace;display:flex;flex-direction:column}
    .topbar{flex-shrink:0;display:flex;align-items:center;justify-content:space-between;padding:.6rem 1.2rem;background:var(--surface);border-bottom:1px solid var(--border)}
    .topbar-left{display:flex;align-items:center;gap:.8rem}
    .logo-icon{width:28px;height:28px;background:var(--accent);border-radius:4px;display:flex;align-items:center;justify-content:center;font-size:.8rem;color:#000;font-weight:700;flex-shrink:0}
    .file-label{font-size:.65rem;color:var(--muted)}
    .file-name{font-size:.8rem;color:var(--text);font-weight:600}
    .topbar-right{display:flex;align-items:center;gap:.5rem}
    .status-pill{display:flex;align-items:center;gap:.4rem;padding:.25rem .65rem;border-radius:3px;font-size:.68rem;border:1px solid var(--border)}
    .status-pill.connected{border-color:rgba(16,185,129,.3);color:var(--accent)}
    .status-pill.disconnected{border-color:rgba(239,68,68,.3);color:var(--danger)}
    .status-dot{width:5px;height:5px;border-radius:50%;background:currentColor}
    .status-pill.connected .status-dot{animation:blink 1.4s infinite}
    @keyframes blink{0%,100%{opacity:1}50%{opacity:.25}}
    .btn{padding:.28rem .7rem;background:transparent;border:1px solid var(--border);border-radius:4px;color:var(--muted);font-family:'JetBrains Mono',monospace;font-size:.68rem;cursor:pointer;text-decoration:none;display:inline-flex;align-items:center;gap:.3rem;transition:all .12s}
    .btn:hover{border-color:var(--text);color:var(--text)}
    .btn.accent{border-color:rgba(16,185,129,.35);color:var(--accent)}
    .btn.accent:hover{background:rgba(16,185,129,.15);border-color:var(--accent)}
    .log-viewport{flex:1;overflow-y:auto;padding:.75rem 1.2rem}
    .log-viewport::-webkit-scrollbar{width:4px}
    .log-viewport::-webkit-scrollbar-track{background:transparent}
    .log-viewport::-webkit-scrollbar-thumb{background:var(--border);border-radius:2px}
    #log-output{font-size:.76rem;line-height:1.7;white-space:pre-wrap;word-break:break-all}
    .log-line{color:var(--text)}
    .log-line.new{animation:fadeIn .15s ease}
    @keyframes fadeIn{from{opacity:0}to{opacity:1}}
    .log-line.has-match mark{background:#78350f;color:#fde68a;border-radius:2px;padding:0 1px}
    .log-line.no-match{opacity:.3}
    .empty-msg{color:var(--muted);font-size:.75rem;margin-top:1rem}
    .cursor{display:inline-block;width:7px;height:13px;background:var(--accent);vertical-align:middle;margin-left:2px;animation:blink 1s infinite}
    .statusbar{flex-shrink:0;display:flex;align-items:center;justify-content:space-between;padding:.35rem 1.2rem;background:var(--surface);border-top:1px solid var(--border);font-size:.65rem;color:var(--muted)}
    .search-input{padding:.25rem .6rem;background:var(--bg);border:1px solid var(--border);border-radius:4px;color:var(--text);font-family:'JetBrains Mono',monospace;font-size:.68rem;outline:none;width:160px;transition:border-color .15s}
    .search-input:focus{border-color:var(--accent)}
    .match-badge{font-size:.65rem;color:var(--muted);white-space:nowrap;min-width:4rem}
  </style>
</head>
<body>
  <div class="topbar">
    <div class="topbar-left">
      <div class="logo-icon">⬡</div>
      <div>
        <div class="file-label">tailing</div>
        <div class="file-name">{{.FilePath}}</div>
      </div>
    </div>
    <div class="topbar-right">
      <input class="search-input" id="search" type="text" placeholder="Filter lines...">
      <span class="match-badge" id="match-badge"></span>
      <div id="status-pill" class="status-pill connected">
        <div class="status-dot"></div>
        <span id="status-text">live</span>
      </div>
      <button class="btn accent" onclick="clearAndRestart()">↺ Clear</button>
      <a class="btn" href="{{.BrowseURL}}">← Back</a>
      {{if .AuthEnabled}}
      <form method="POST" action="{{.LogoutAction}}" style="margin:0">
        <button class="btn" type="submit">Sign out</button>
      </form>
      {{end}}
    </div>
  </div>

  <div class="log-viewport" id="viewport">
    <div id="log-output"></div>
    <span id="cursor" class="cursor"></span>
    <div id="empty-msg" class="empty-msg">Waiting for new log lines... (click ↺ Clear to restart stream)</div>
  </div>

  <div class="statusbar">
    <span id="line-count">0 lines</span>
    <span>{{.FilePath}} · refresh page to clear &amp; restart</span>
  </div>

  <script>
    const streamURL = {{.StreamURL}};
    const MAX_LINES = 2000;
    const output = document.getElementById('log-output');
    const emptyMsg = document.getElementById('empty-msg');
    const statusPill = document.getElementById('status-pill');
    const statusText = document.getElementById('status-text');
    const lineCountEl = document.getElementById('line-count');
    const matchBadge = document.getElementById('match-badge');
    const viewport = document.getElementById('viewport');
    const searchEl = document.getElementById('search');
    let lineCount = 0;
    let allLines = [];   // capped at MAX_LINES, each {text, el}
    let matchCount = 0;  // incremental — never scan allLines on append
    let searchQuery = '';
    let es = null;
    let autoScroll = true;
    let lineBuffer = []; // incoming lines waiting for next animation frame
    let flushPending = false;

    viewport.addEventListener('scroll', () => {
      const atBottom = viewport.scrollHeight - viewport.scrollTop - viewport.clientHeight < 40;
      autoScroll = atBottom;
    });

    function setStatus(connected) {
      statusPill.className = 'status-pill ' + (connected ? 'connected' : 'disconnected');
      statusText.textContent = connected ? 'live' : 'disconnected';
    }

    function escHtml(s) {
      return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
    }

    function highlight(text, q) {
      const ql = q.toLowerCase(), tl = text.toLowerCase();
      const parts = []; let pos = 0, idx;
      while ((idx = tl.indexOf(ql, pos)) !== -1) {
        parts.push(escHtml(text.slice(pos, idx)));
        parts.push('<mark>' + escHtml(text.slice(idx, idx + q.length)) + '</mark>');
        pos = idx + q.length;
      }
      parts.push(escHtml(text.slice(pos)));
      return parts.join('');
    }

    function styleEl(el, text) {
      const q = searchQuery;
      if (!q) {
        el.className = 'log-line new';
        el.textContent = text;
      } else if (text.toLowerCase().includes(q.toLowerCase())) {
        el.className = 'log-line new has-match';
        el.innerHTML = highlight(text, q);
      } else {
        el.className = 'log-line new no-match';
        el.textContent = text;
      }
    }

    function updateCounter() {
      const atCap = allLines.length >= MAX_LINES;
      lineCountEl.textContent = lineCount + ' lines' + (atCap ? ' (last ' + MAX_LINES + ')' : '');
      if (!searchQuery) {
        matchBadge.textContent = '';
      } else {
        matchBadge.textContent = matchCount ? matchCount + ' match' + (matchCount !== 1 ? 'es' : '') : 'no matches';
      }
    }

    // Flush all buffered lines in one animation frame:
    // - builds a DocumentFragment (single reflow)
    // - evicts oldest lines when over MAX_LINES
    // - scrolls once at the end
    function flushLines() {
      flushPending = false;
      if (!lineBuffer.length) return;

      const frag = document.createDocumentFragment();
      const toRemove = [];

      lineBuffer.forEach(text => {
        if (lineCount === 0) emptyMsg.style.display = 'none';
        const el = document.createElement('div');
        styleEl(el, text);
        if (searchQuery && text.toLowerCase().includes(searchQuery.toLowerCase())) matchCount++;
        frag.appendChild(el);
        allLines.push({text, el});
        lineCount++;
        if (allLines.length > MAX_LINES) {
          const old = allLines.shift();
          if (searchQuery && old.text.toLowerCase().includes(searchQuery.toLowerCase())) matchCount--;
          toRemove.push(old.el);
        }
      });
      lineBuffer = [];

      toRemove.forEach(el => output.removeChild(el));
      output.appendChild(frag);
      updateCounter();
      if (autoScroll) viewport.scrollTop = viewport.scrollHeight;
    }

    function appendLine(text) {
      lineBuffer.push(text);
      if (!flushPending) {
        flushPending = true;
        requestAnimationFrame(flushLines);
      }
    }

    function startStream() {
      if (es) es.close();
      es = new EventSource(streamURL);
      es.onopen = () => setStatus(true);
      es.onmessage = (e) => appendLine(e.data);
      es.onerror = () => {
        setStatus(false);
        setTimeout(() => { if (es.readyState === EventSource.CLOSED) startStream(); }, 2000);
      };
    }

    function clearAndRestart() {
      if (es) es.close();
      lineBuffer = []; flushPending = false;
      output.innerHTML = '';
      allLines = []; lineCount = 0; matchCount = 0;
      lineCountEl.textContent = '0 lines';
      matchBadge.textContent = '';
      emptyMsg.style.display = '';
      autoScroll = true;
      startStream();
    }

    let searchTimer;
    searchEl.addEventListener('input', e => {
      clearTimeout(searchTimer);
      searchTimer = setTimeout(() => {
        searchQuery = e.target.value.trim();
        matchCount = 0;
        allLines.forEach(({text, el}) => {
          styleEl(el, text);
          if (searchQuery && text.toLowerCase().includes(searchQuery.toLowerCase())) matchCount++;
        });
        updateCounter();
        if (autoScroll) {
          const last = output.querySelector('.log-line.has-match:last-of-type');
          if (last) last.scrollIntoView({block:'end'});
        }
      }, 150);
    });

    startStream();
  </script>
</body>
</html>`))

var viewerTmpl = template.Must(template.New("viewer").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>view · {{.FileName}}</title>
  <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>⬡</text></svg>">
  <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&display=swap" rel="stylesheet">
  <style>
    :root{--bg:#0c0c0c;--surface:#141414;--border:#222;--accent:#10b981;--text:#d4d4d8;--muted:#52525b;--danger:#ef4444}
    *{margin:0;padding:0;box-sizing:border-box}
    html,body{height:100%;background:var(--bg);color:var(--text);font-family:'JetBrains Mono',monospace;display:flex;flex-direction:column}
    .topbar{flex-shrink:0;display:flex;align-items:center;gap:.8rem;padding:.6rem 1.2rem;background:var(--surface);border-bottom:1px solid var(--border);flex-wrap:wrap}
    .topbar-left{display:flex;align-items:center;gap:.8rem;flex:1;min-width:0}
    .logo-icon{width:28px;height:28px;background:var(--accent);border-radius:4px;display:flex;align-items:center;justify-content:center;font-size:.8rem;color:#000;font-weight:700;flex-shrink:0}
    .file-label{font-size:.65rem;color:var(--muted)}
    .file-name{font-size:.8rem;color:var(--text);font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
    .search-row{display:flex;align-items:center;gap:.4rem;flex-shrink:0}
    .search-input{padding:.28rem .65rem;background:var(--bg);border:1px solid var(--border);border-radius:4px;color:var(--text);font-family:'JetBrains Mono',monospace;font-size:.72rem;outline:none;width:200px;transition:border-color .15s}
    .search-input:focus{border-color:var(--accent)}
    .match-info{font-size:.65rem;color:var(--muted);white-space:nowrap;min-width:5rem}
    .nav-btn{padding:.25rem .5rem;background:transparent;border:1px solid var(--border);border-radius:3px;color:var(--muted);font-family:'JetBrains Mono',monospace;font-size:.72rem;cursor:pointer;transition:all .12s}
    .nav-btn:hover{border-color:var(--text);color:var(--text)}
    .filter-btn{padding:.25rem .55rem;background:transparent;border:1px solid var(--border);border-radius:3px;color:var(--muted);font-family:'JetBrains Mono',monospace;font-size:.68rem;cursor:pointer;transition:all .12s}
    .filter-btn:hover{border-color:var(--accent);color:var(--accent)}
    .filter-btn.active{border-color:var(--accent);color:var(--accent);background:rgba(16,185,129,.08)}
    .topbar-right{display:flex;align-items:center;gap:.4rem;flex-shrink:0}
    .btn{padding:.28rem .7rem;background:transparent;border:1px solid var(--border);border-radius:4px;color:var(--muted);font-family:'JetBrains Mono',monospace;font-size:.68rem;cursor:pointer;text-decoration:none;display:inline-flex;align-items:center;gap:.3rem;transition:all .12s}
    .btn:hover{border-color:var(--text);color:var(--text)}
    .warn-bar{flex-shrink:0;padding:.4rem 1.2rem;background:rgba(217,119,6,.12);border-bottom:1px solid rgba(217,119,6,.25);font-size:.7rem;color:#fbbf24}
    .warn-bar a{color:#fbbf24;text-underline-offset:2px}
    .viewer{flex:1;overflow-y:auto}
    .viewer::-webkit-scrollbar{width:4px}
    .viewer::-webkit-scrollbar-track{background:transparent}
    .viewer::-webkit-scrollbar-thumb{background:var(--border);border-radius:2px}
    #content{padding:.5rem 0}
    .line{display:flex;min-height:1.4rem}
    .line:hover{background:rgba(255,255,255,.02)}
    .line.has-match{background:rgba(120,53,15,.08)}
    .line.no-match{opacity:.25}
    .ln{flex-shrink:0;width:4.5rem;padding:.0 .75rem 0 1rem;text-align:right;color:var(--muted);font-size:.72rem;line-height:1.7;user-select:none}
    .lt{flex:1;padding-right:1.2rem;font-size:.76rem;line-height:1.7;white-space:pre-wrap;word-break:break-all;color:var(--text)}
    mark{background:#78350f;color:#fde68a;border-radius:2px;padding:0 1px}
    mark.active{background:#d97706;color:#000}
    .statusbar{flex-shrink:0;display:flex;align-items:center;justify-content:space-between;padding:.3rem 1.2rem;background:var(--surface);border-top:1px solid var(--border);font-size:.65rem;color:var(--muted)}
  </style>
</head>
<body>
  <div class="topbar">
    <div class="topbar-left">
      <div class="logo-icon">⬡</div>
      <div style="min-width:0">
        <div class="file-label">viewing</div>
        <div class="file-name">{{.FilePath}}</div>
      </div>
    </div>
    <div class="search-row">
      <input class="search-input" id="search" type="text" placeholder="Search keyword... (Ctrl+F)">
      <span class="match-info" id="match-info"></span>
      <button class="nav-btn" onclick="navigate(-1)" title="Previous (Shift+Enter)">↑</button>
      <button class="nav-btn" onclick="navigate(1)" title="Next (Enter)">↓</button>
      <button class="filter-btn" id="filter-btn" onclick="toggleFilter()">Filter</button>
    </div>
    <div class="topbar-right">
      <a class="btn" href="{{.BrowseURL}}">← Back</a>
      {{if .AuthEnabled}}
      <form method="POST" action="{{.LogoutAction}}" style="margin:0">
        <button class="btn" type="submit">Sign out</button>
      </form>
      {{end}}
    </div>
  </div>

  {{if .Truncated}}
  <div class="warn-bar">⚠ Large file — showing last 5 MB. <a href="{{.DownloadURL}}">Download</a> for full content.</div>
  {{end}}

  <div class="viewer" id="viewer">
    <div id="content"></div>
  </div>

  <div class="statusbar">
    <span id="line-info">{{.LineCount}} lines</span>
    <span>{{.FilePath}}</span>
  </div>

  <script>
    const rawContent = {{.ContentJSON}};
    const lines = rawContent.split('\n');
    if (lines.length && lines[lines.length-1] === '') lines.pop();

    const contentEl = document.getElementById('content');
    const searchEl  = document.getElementById('search');
    const matchInfoEl = document.getElementById('match-info');
    const lineInfoEl  = document.getElementById('line-info');
    const filterBtnEl = document.getElementById('filter-btn');
    const viewer = document.getElementById('viewer');

    let query = '', filterMode = false, matchIndices = [], currentIdx = -1;

    function escHtml(s) {
      return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
    }

    function highlightLine(text, q) {
      const ql = q.toLowerCase(), tl = text.toLowerCase();
      const parts = []; let pos = 0, idx;
      while ((idx = tl.indexOf(ql, pos)) !== -1) {
        parts.push(escHtml(text.slice(pos, idx)));
        parts.push('<mark>' + escHtml(text.slice(idx, idx + q.length)) + '</mark>');
        pos = idx + q.length;
      }
      parts.push(escHtml(text.slice(pos)));
      return parts.join('') || '&nbsp;';
    }

    function render() {
      matchIndices = [];
      const parts = [];
      const q = query;

      for (let i = 0; i < lines.length; i++) {
        const has = q && lines[i].toLowerCase().includes(q.toLowerCase());
        if (filterMode && q && !has) continue;
        if (has) matchIndices.push(i);
        const cls = !q ? 'line' : has ? 'line has-match' : 'line no-match';
        const body = has ? highlightLine(lines[i], q) : (escHtml(lines[i]) || '&nbsp;');
        parts.push('<div class="' + cls + '" data-i="' + i + '"><span class="ln">' + (i+1) + '</span><span class="lt">' + body + '</span></div>');
      }

      contentEl.innerHTML = parts.join('');

      const n = matchIndices.length;
      if (!q) {
        matchInfoEl.textContent = '';
        lineInfoEl.textContent = lines.length + ' lines';
      } else {
        matchInfoEl.textContent = n ? n + ' match' + (n !== 1 ? 'es' : '') : 'no matches';
        lineInfoEl.textContent = lines.length + ' lines' + (filterMode ? ' · ' + n + ' shown' : '');
      }

      currentIdx = n > 0 ? 0 : -1;
      activateCurrent();
    }

    function activateCurrent() {
      document.querySelectorAll('mark.active').forEach(m => m.classList.remove('active'));
      if (currentIdx < 0 || !matchIndices.length) return;
      const el = contentEl.querySelector('[data-i="' + matchIndices[currentIdx] + '"]');
      if (!el) return;
      const m = el.querySelector('mark');
      if (m) m.classList.add('active');
      el.scrollIntoView({block:'center', behavior:'smooth'});
      matchInfoEl.textContent = (currentIdx+1) + ' / ' + matchIndices.length + ' match' + (matchIndices.length !== 1 ? 'es' : '');
    }

    function navigate(dir) {
      if (!matchIndices.length) return;
      currentIdx = (currentIdx + dir + matchIndices.length) % matchIndices.length;
      activateCurrent();
    }

    function toggleFilter() {
      filterMode = !filterMode;
      filterBtnEl.classList.toggle('active', filterMode);
      render();
    }

    let timer;
    searchEl.addEventListener('input', e => {
      clearTimeout(timer);
      timer = setTimeout(() => { query = e.target.value.trim(); render(); }, 150);
    });

    searchEl.addEventListener('keydown', e => {
      if (e.key === 'Enter') { e.preventDefault(); navigate(e.shiftKey ? -1 : 1); }
      if (e.key === 'Escape') { searchEl.value = ''; query = ''; render(); }
    });

    document.addEventListener('keydown', e => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
        e.preventDefault(); searchEl.focus(); searchEl.select();
      }
    });

    render();
  </script>
</body>
</html>`))
