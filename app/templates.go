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
    .table-header{display:grid;grid-template-columns:1fr 90px 150px 200px;padding:.55rem 1rem;border-bottom:1px solid var(--border);font-size:.6rem;text-transform:uppercase;letter-spacing:.1em;color:var(--muted)}
    .log-row{display:grid;grid-template-columns:1fr 90px 150px 200px;padding:.75rem 1rem;border-bottom:1px solid var(--border);align-items:center;transition:background .12s;position:relative}
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
    .download-btn,.tail-btn{display:inline-flex;align-items:center;gap:.3rem;padding:.28rem .65rem;background:transparent;border-radius:4px;font-family:'JetBrains Mono',monospace;font-size:.67rem;text-decoration:none;transition:all .12s;white-space:nowrap}
    .download-btn{border:1px solid var(--border);color:var(--muted)}
    .download-btn:hover{border-color:var(--accent);color:var(--accent)}
    .tail-btn{border:1px solid rgba(16,185,129,.3);color:var(--accent)}
    .tail-btn:hover{background:rgba(16,185,129,.08)}
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
    <div class="log-table">
      <div class="table-header">
        <div>Name</div><div>Size</div><div>Modified</div><div>Actions</div>
      </div>
      {{if .SubPath}}
      <div class="log-row">
        <div class="entry-name">
          <div class="entry-icon">↑</div>
          <a class="entry-link up-link" href="{{.ParentURL}}">..</a>
        </div>
        <div></div><div></div><div></div>
      </div>
      {{end}}
      {{range .Entries}}
      <div class="log-row{{if .IsDir}} is-dir{{end}}">
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

    <div class="footer">LogVault · Minimal Log Server · by Lokendra Bhat</div>
  </div>
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
    .btn.accent:hover{background:rgba(16,185,129,.07)}
    .log-viewport{flex:1;overflow-y:auto;padding:.75rem 1.2rem}
    .log-viewport::-webkit-scrollbar{width:4px}
    .log-viewport::-webkit-scrollbar-track{background:transparent}
    .log-viewport::-webkit-scrollbar-thumb{background:var(--border);border-radius:2px}
    #log-output{font-size:.76rem;line-height:1.7;white-space:pre-wrap;word-break:break-all}
    .log-line{color:var(--text)}
    .log-line.new{animation:fadeIn .15s ease}
    @keyframes fadeIn{from{opacity:0}to{opacity:1}}
    .empty-msg{color:var(--muted);font-size:.75rem;margin-top:1rem}
    .cursor{display:inline-block;width:7px;height:13px;background:var(--accent);vertical-align:middle;margin-left:2px;animation:blink 1s infinite}
    .statusbar{flex-shrink:0;display:flex;align-items:center;justify-content:space-between;padding:.35rem 1.2rem;background:var(--surface);border-top:1px solid var(--border);font-size:.65rem;color:var(--muted)}
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
    const output = document.getElementById('log-output');
    const emptyMsg = document.getElementById('empty-msg');
    const statusPill = document.getElementById('status-pill');
    const statusText = document.getElementById('status-text');
    const lineCountEl = document.getElementById('line-count');
    const viewport = document.getElementById('viewport');
    let lineCount = 0;
    let es = null;
    let autoScroll = true;

    viewport.addEventListener('scroll', () => {
      const atBottom = viewport.scrollHeight - viewport.scrollTop - viewport.clientHeight < 40;
      autoScroll = atBottom;
    });

    function setStatus(connected) {
      statusPill.className = 'status-pill ' + (connected ? 'connected' : 'disconnected');
      statusText.textContent = connected ? 'live' : 'disconnected';
    }

    function appendLine(text) {
      if (lineCount === 0) emptyMsg.style.display = 'none';
      const div = document.createElement('div');
      div.className = 'log-line new';
      div.textContent = text;
      output.appendChild(div);
      lineCount++;
      lineCountEl.textContent = lineCount + ' line' + (lineCount !== 1 ? 's' : '');
      if (autoScroll) viewport.scrollTop = viewport.scrollHeight;
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
      output.innerHTML = '';
      lineCount = 0;
      lineCountEl.textContent = '0 lines';
      emptyMsg.style.display = '';
      autoScroll = true;
      startStream();
    }

    startStream();
  </script>
</body>
</html>`))
