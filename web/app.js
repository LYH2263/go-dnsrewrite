async function refresh() {
  const rules = await (await fetch('/api/rules')).json();
  const box = document.getElementById('rules');
  box.innerHTML = rules.map(r => `
    <div class="rule">
      <span><b>${r.ID}</b> ${r.Kind}:${r.Pattern} → ${r.Action}
      [${(r.Targets||[]).join(',')}] ${r.Upstream||''} ttl=${r.TTL||0}</span>
      <button data-id="${r.ID}">删除</button>
    </div>`).join('') || '<div class="rule">暂无规则</div>';
  box.querySelectorAll('button[data-id]').forEach(btn => {
    btn.onclick = async () => {
      await fetch('/api/rules/' + btn.dataset.id, { method: 'DELETE' });
      refresh();
    };
  });
  const stats = await (await fetch('/api/stats')).json();
  document.getElementById('stats').textContent = JSON.stringify(stats, null, 2);
}

document.getElementById('add-form').onsubmit = async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  const targets = String(fd.get('targets')||'').split(',').map(s=>s.trim()).filter(Boolean);
  const body = {
    ID: fd.get('id'), Pattern: fd.get('pattern'), Kind: fd.get('kind'),
    Action: fd.get('action'), Targets: targets, Upstream: fd.get('upstream'),
    TTL: Number(fd.get('ttl')||300), Priority: 0, Enabled: true,
  };
  await fetch('/api/rules', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(body) });
  e.target.reset();
  refresh();
};

document.getElementById('try-form').onsubmit = async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  const body = { name: fd.get('name'), type: Number(fd.get('type')) };
  const res = await fetch('/api/try', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(body) });
  document.getElementById('try-out').textContent = await res.text();
};

refresh();
setInterval(refresh, 5000);
