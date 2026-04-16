/* ─── State ─── */
let allData = [];
let filteredData = [];
let secret = '';
let currentSort = { col: 'alias', asc: true };
const PAGE_SIZE = 15;
let currentPage = 0;

/* ─── Init ─── */
document.addEventListener('DOMContentLoaded', () => {
    fetchData();
    // close modals on overlay click
    document.querySelectorAll('.modal-overlay').forEach(el => {
        el.addEventListener('click', e => {
            if (e.target === el) closeModal(el.id);
        });
    });
    // close modals on Escape
    document.addEventListener('keydown', e => {
        if (e.key === 'Escape') {
            document.querySelectorAll('.modal-overlay.active').forEach(el => closeModal(el.id));
        }
    });
});

/* ─── Data Fetching ─── */
function fetchData() {
    const url = secret ? `actions/view?secret=${encodeURIComponent(secret)}` : 'actions/view';
    fetch(url)
        .then(r => r.json())
        .then(resp => {
            allData = (resp.data || []).map(d => ({
                alias: d.values.alias,
                fullurl: d.values.fullurl,
                shorturl: d.values.shorturl
            }));
            applySort();
            filterTable();
        })
        .catch(() => toast('Failed to load bookmarks', 'error'));
}

function refreshData() {
    fetchData();
    toast('Data refreshed', 'info');
}

/* ─── Filtering ─── */
function filterTable() {
    const q = document.getElementById('searchInput').value.toLowerCase();
    filteredData = allData.filter(d =>
        d.alias.toLowerCase().includes(q) ||
        d.fullurl.toLowerCase().includes(q) ||
        d.shorturl.toLowerCase().includes(q)
    );
    currentPage = 0;
    renderTable();
}

/* ─── Sorting ─── */
function sortBy(col) {
    if (currentSort.col === col) currentSort.asc = !currentSort.asc;
    else { currentSort.col = col; currentSort.asc = true; }
    applySort();
    filterTable();
}

function applySort() {
    const { col, asc } = currentSort;
    allData.sort((a, b) => {
        const va = (a[col] || '').toLowerCase();
        const vb = (b[col] || '').toLowerCase();
        return asc ? va.localeCompare(vb) : vb.localeCompare(va);
    });
}

/* ─── Rendering ─── */
function renderTable() {
    const tbody = document.getElementById('tableBody');
    const totalPages = Math.ceil(filteredData.length / PAGE_SIZE);
    const start = currentPage * PAGE_SIZE;
    const pageData = filteredData.slice(start, start + PAGE_SIZE);

    document.getElementById('totalCount').textContent = `${filteredData.length} bookmark${filteredData.length !== 1 ? 's' : ''}`;
    document.getElementById('searchCount').textContent = document.getElementById('searchInput').value
        ? `${filteredData.length} result${filteredData.length !== 1 ? 's' : ''}` : '';

    if (pageData.length === 0) {
        tbody.innerHTML = `
            <tr><td colspan="4">
                <div class="empty-state">
                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M19 21l-7-5-7 5V5a2 2 0 012-2h10a2 2 0 012 2z"/></svg>
                    <h3>No bookmarks found</h3>
                    <p>Add your first bookmark to get started.</p>
                </div>
            </td></tr>`;
        document.getElementById('paginator').innerHTML = '';
        return;
    }

    tbody.innerHTML = pageData.map(d => `
        <tr id="row-${escHtml(d.alias)}">
            <td class="col-url">
                <div class="editable" onclick="startEdit(this, '${escAttr(d.alias)}', 'orig', '${escAttr(d.fullurl)}')">
                    <span class="cell-url" title="${escHtml(d.fullurl)}">${escHtml(d.fullurl)}</span>
                </div>
            </td>
            <td class="col-alias">
                <div class="editable" onclick="startEdit(this, '${escAttr(d.alias)}', 'alias', '${escAttr(d.alias)}')">
                    <span class="cell-alias">${escHtml(d.alias)}</span>
                </div>
            </td>
            <td class="col-short">
                <div class="copy-wrap">
                    <a class="cell-short" href="${escHtml(d.shorturl)}" target="_blank">${escHtml(d.shorturl)}</a>
                    <button class="btn-icon copy-btn" onclick="event.stopPropagation(); copyText('${escAttr(d.shorturl)}')" title="Copy short URL">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
                    </button>
                </div>
            </td>
            <td class="col-actions">
                <button class="btn-icon danger" onclick="confirmDelete('${escAttr(d.alias)}')" title="Delete bookmark">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
                </button>
            </td>
        </tr>`).join('');

    renderPaginator(totalPages);
}

/* ─── Pagination ─── */
function renderPaginator(totalPages) {
    const pg = document.getElementById('paginator');
    if (totalPages <= 1) { pg.innerHTML = ''; return; }
    let html = `<button ${currentPage === 0 ? 'disabled' : ''} onclick="goPage(${currentPage - 1})">‹</button>`;
    for (let i = 0; i < totalPages; i++) {
        html += `<button class="${i === currentPage ? 'active' : ''}" onclick="goPage(${i})">${i + 1}</button>`;
    }
    html += `<button ${currentPage === totalPages - 1 ? 'disabled' : ''} onclick="goPage(${currentPage + 1})">›</button>`;
    pg.innerHTML = html;
}

function goPage(p) {
    currentPage = p;
    renderTable();
    document.querySelector('.table-wrapper').scrollIntoView({ behavior: 'smooth', block: 'start' });
}

/* ─── Inline Editing ─── */
function startEdit(el, alias, colname, oldValue) {
    if (el.classList.contains('editing')) return;
    el.classList.add('editing');
    const origText = oldValue;
    el.innerHTML = `<input class="edit-input" type="text" value="${escAttr(origText)}" />`;
    const input = el.querySelector('input');
    input.focus();
    input.select();

    function commit() {
        const newVal = input.value.trim();
        el.classList.remove('editing');
        if (!newVal || newVal === origText) {
            el.innerHTML = `<span class="${colname === 'alias' ? 'cell-alias' : 'cell-url'}" title="${escHtml(origText)}">${escHtml(origText)}</span>`;
            return;
        }
        const params = new URLSearchParams({
            id: alias, newvalue: newVal, oldvalue: origText, colname: colname, secret: secret
        });
        fetch('actions/update', { method: 'POST', body: params })
            .then(r => r.text())
            .then(resp => {
                if (resp === 'ok') {
                    toast(`Updated ${colname}`, 'success');
                    fetchData();
                } else {
                    toast(resp || 'Update failed', 'error');
                    el.innerHTML = `<span class="${colname === 'alias' ? 'cell-alias' : 'cell-url'}" title="${escHtml(origText)}">${escHtml(origText)}</span>`;
                }
            })
            .catch(() => {
                toast('Network error', 'error');
                el.innerHTML = `<span class="${colname === 'alias' ? 'cell-alias' : 'cell-url'}" title="${escHtml(origText)}">${escHtml(origText)}</span>`;
            });
    }

    input.addEventListener('blur', commit);
    input.addEventListener('keydown', e => {
        if (e.key === 'Enter') { e.preventDefault(); input.blur(); }
        if (e.key === 'Escape') {
            input.removeEventListener('blur', commit);
            el.classList.remove('editing');
            el.innerHTML = `<span class="${colname === 'alias' ? 'cell-alias' : 'cell-url'}" title="${escHtml(origText)}">${escHtml(origText)}</span>`;
        }
    });
}

/* ─── Add Bookmark ─── */
function addBookmark() {
    const alias = document.getElementById('addAlias').value.trim();
    const url = document.getElementById('addUrl').value.trim();
    if (!alias || !url) { toast('Please fill both fields', 'error'); return; }

    const params = new URLSearchParams({ short: alias, url: url, secret: secret });
    fetch('actions/add', { method: 'POST', body: params })
        .then(r => r.text())
        .then(resp => {
            if (resp === 'ok') {
                toast(`Bookmark "${alias}" added`, 'success');
                closeModal('addModal');
                document.getElementById('addAlias').value = '';
                document.getElementById('addUrl').value = '';
                fetchData();
            } else {
                toast(resp || 'Failed to add', 'error');
            }
        })
        .catch(() => toast('Network error', 'error'));
}

/* ─── Delete Bookmark ─── */
function confirmDelete(alias) {
    document.getElementById('confirmMessage').textContent = `This will permanently delete the alias "${alias}". This action cannot be undone.`;
    const actionBtn = document.getElementById('confirmAction');
    actionBtn.onclick = () => deleteBookmark(alias);
    openModal('confirmModal');
}

function deleteBookmark(alias) {
    const params = new URLSearchParams({ id: alias, secret: secret });
    fetch('actions/delete', { method: 'POST', body: params })
        .then(r => r.text())
        .then(resp => {
            if (resp === 'ok') {
                toast(`Deleted "${alias}"`, 'success');
                closeModal('confirmModal');
                fetchData();
            } else {
                toast(resp || 'Delete failed', 'error');
            }
        })
        .catch(() => toast('Network error', 'error'));
}

/* ─── Secret Management ─── */
function saveSecret() {
    secret = document.getElementById('secretInput').value;
    closeModal('secretModal');
    updateSecretBadge();
    fetchData();
    toast('Secret saved — data refreshed', 'success');
}

function clearSecret() {
    secret = '';
    document.getElementById('secretInput').value = '';
    closeModal('secretModal');
    updateSecretBadge();
    fetchData();
    toast('Secret cleared', 'info');
}

function updateSecretBadge() {
    const badge = document.getElementById('secretStatus');
    if (secret) {
        badge.className = 'secret-badge';
        badge.innerHTML = '<span class="secret-dot"></span> Authenticated';
    } else {
        badge.className = 'secret-badge none';
        badge.innerHTML = '<span class="secret-dot"></span> No secret';
    }
}

/* ─── Modals ─── */
function openModal(id) { document.getElementById(id).classList.add('active'); }
function closeModal(id) { document.getElementById(id).classList.remove('active'); }
function openAddModal() { openModal('addModal'); setTimeout(() => document.getElementById('addAlias').focus(), 200); }
function openSecretModal() { openModal('secretModal'); setTimeout(() => document.getElementById('secretInput').focus(), 200); }

/* ─── Toast Notifications ─── */
function toast(message, type) {
    type = type || 'info';
    const container = document.getElementById('toastContainer');
    const icons = {
        success: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>',
        error: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ef4444" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>',
        info: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#6366f1" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>',
    };
    const el = document.createElement('div');
    el.className = `toast ${type}`;
    el.innerHTML = `<span class="toast-icon">${icons[type] || icons.info}</span>${escHtml(message)}`;
    container.appendChild(el);
    setTimeout(() => {
        el.classList.add('removing');
        setTimeout(() => el.remove(), 250);
    }, 3000);
}

/* ─── Copy to Clipboard ─── */
function copyText(text) {
    navigator.clipboard.writeText(text)
        .then(() => toast('Copied to clipboard', 'success'))
        .catch(() => toast('Copy failed', 'error'));
}

/* ─── HTML Escape Helpers ─── */
function escHtml(s) {
    const d = document.createElement('div');
    d.appendChild(document.createTextNode(s));
    return d.innerHTML;
}

function escAttr(s) {
    return s.replace(/&/g,'&amp;').replace(/'/g,'&#39;').replace(/"/g,'&quot;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}
