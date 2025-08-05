// 琛ㄦ牸瑙嗗浘鑴氭湰
console.log('鍔犺浇琛ㄦ牸瑙嗗浘鑴氭湰');

const API_BASE = '/';

class TableView {
  constructor() {
    this.currentView = 'nodes'; // 'nodes' 鎴?'paths'
    this.data = { nodes: {}, paths: {} };
    this.init();
  }

  async init() {
    await this.loadData();
    this.setupEventListeners();
    this.renderTable();
  }

  async loadData() {
    try {
      const response = await fetch(API_BASE + 'canvas-data');
      if (!response.ok) throw new Error('鑾峰彇鏁版嵁澶辫触');
      this.data = await response.json();
    } catch (error) {
      console.error('鍔犺浇鏁版嵁澶辫触:', error);
    }
  }

  setupEventListeners() {
    // 瑙嗗浘鍒囨崲鎸夐挳
    const nodeViewBtn = document.getElementById('nodeViewBtn');
    const pathViewBtn = document.getElementById('pathViewBtn');
    const refreshBtn = document.getElementById('refreshBtn');
    const addBtn = document.getElementById('addBtn');

    if (nodeViewBtn) {
      nodeViewBtn.addEventListener('click', () => {
        this.switchView('nodes');
      });
    }

    if (pathViewBtn) {
      pathViewBtn.addEventListener('click', () => {
        this.switchView('paths');
      });
    }

    if (refreshBtn) {
      refreshBtn.addEventListener('click', async () => {
        await this.loadData();
        this.renderTable();
      });
    }

    if (addBtn) {
      addBtn.addEventListener('click', () => {
        this.addNewItem();
      });
    }
  }

  switchView(viewType) {
    this.currentView = viewType;
    
    // 鏇存柊鎸夐挳鐘舵€?
    document.getElementById('nodeViewBtn').classList.toggle('active', viewType === 'nodes');
    document.getElementById('pathViewBtn').classList.toggle('active', viewType === 'paths');
    
    this.renderTable();
  }

  renderTable() {
    const container = document.getElementById('tableContainer');
    if (!container) return;

    const items = this.currentView === 'nodes' ? 
      Object.values(this.data.nodes || {}) : 
      Object.values(this.data.paths || {});

    let html = '';

    if (this.currentView === 'nodes') {
      html = this.renderNodeTable(items);
    } else {
      html = this.renderPathTable(items);
    }

    container.innerHTML = html;
    this.attachTableEventListeners();
  }

  renderNodeTable(nodes) {
    let html = `
      <div class="table-header">
        <h3>鑺傜偣绠＄悊 (${nodes.length} 涓妭鐐?</h3>
      </div>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>鍚嶇О</th>
              <th>绫诲瀷</th>
              <th>鐘舵€?/th>
              <th>X鍧愭爣</th>
              <th>Y鍧愭爣</th>
              <th>Z鍧愭爣</th>
              <th>鍒涘缓鏃堕棿</th>
              <th>鎿嶄綔</th>
            </tr>
          </thead>
          <tbody>
    `;

    nodes.forEach(node => {
      html += `
        <tr data-id="${node.id}">
          <td>${node.id}</td>
          <td><input type="text" value="${node.name}" data-field="name" class="editable-input"></td>
          <td><input type="text" value="${node.type || 'default'}" data-field="type" class="editable-input"></td>
          <td><input type="text" value="${node.status || 'active'}" data-field="status" class="editable-input"></td>
          <td><input type="number" value="${node.position.x}" data-field="position.x" class="editable-input" step="0.01"></td>
          <td><input type="number" value="${node.position.y}" data-field="position.y" class="editable-input" step="0.01"></td>
          <td><input type="number" value="${node.position.z || 0}" data-field="position.z" class="editable-input" step="0.01"></td>
          <td>${new Date(node.created_at).toLocaleString()}</td>
          <td>
            <button class="btn-small btn-save" data-id="${node.id}">淇濆瓨</button>
            <button class="btn-small btn-delete" data-id="${node.id}">鍒犻櫎</button>
          </td>
        </tr>
      `;
    });

    html += `
          </tbody>
        </table>
      </div>
    `;

    return html;
  }

  renderPathTable(paths) {
    let html = `
      <div class="table-header">
        <h3>璺緞绠＄悊 (${paths.length} 涓矾寰?</h3>
      </div>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>鍚嶇О</th>
              <th>绫诲瀷</th>
              <th>鐘舵€?/th>
              <th>璧峰鑺傜偣</th>
              <th>缁撴潫鑺傜偣</th>
              <th>鍒涘缓鏃堕棿</th>
              <th>鎿嶄綔</th>
            </tr>
          </thead>
          <tbody>
    `;

    paths.forEach(path => {
      html += `
        <tr data-id="${path.id}">
          <td>${path.id}</td>
          <td><input type="text" value="${path.name}" data-field="name" class="editable-input"></td>
          <td><input type="text" value="${path.type || 'default'}" data-field="type" class="editable-input"></td>
          <td><input type="text" value="${path.status || 'active'}" data-field="status" class="editable-input"></td>
          <td><select data-field="start_node_id" class="editable-input">${this.renderNodeOptions(path.start_node_id)}</select></td>
          <td><select data-field="end_node_id" class="editable-input">${this.renderNodeOptions(path.end_node_id)}</select></td>
          <td>${new Date(path.created_at).toLocaleString()}</td>
          <td>
            <button class="btn-small btn-save" data-id="${path.id}">淇濆瓨</button>
            <button class="btn-small btn-delete" data-id="${path.id}">鍒犻櫎</button>
          </td>
        </tr>
      `;
    });

    html += `
          </tbody>
        </table>
      </div>
    `;

    return html;
  }

  renderNodeOptions(selectedId) {
    const nodes = Object.values(this.data.nodes || {});
    return nodes.map(node => 
      `<option value="${node.id}" ${node.id === selectedId ? 'selected' : ''}>${node.name} (${node.id})</option>`
    ).join('');
  }

  attachTableEventListeners() {
    // 淇濆瓨鎸夐挳
    document.querySelectorAll('.btn-save').forEach(btn => {
      btn.addEventListener('click', async (e) => {
        const id = e.target.dataset.id;
        await this.saveItem(id);
      });
    });

    // 鍒犻櫎鎸夐挳
    document.querySelectorAll('.btn-delete').forEach(btn => {
      btn.addEventListener('click', async (e) => {
        const id = e.target.dataset.id;
        if (confirm(`纭畾鍒犻櫎杩欎釜${this.currentView === 'nodes' ? '鑺傜偣' : '璺緞'}鍚楋紵`)) {
          await this.deleteItem(id);
        }
      });
    });
  }

  async saveItem(id) {
    const row = document.querySelector(`tr[data-id="${id}"]`);
    if (!row) return;

    const data = {};
    const inputs = row.querySelectorAll('.editable-input');
    
    inputs.forEach(input => {
      const field = input.dataset.field;
      let value = input.value;
      
      // 澶勭悊鏁板瓧绫诲瀷
      if (input.type === 'number') {
        value = parseFloat(value) || 0;
      }
      
      // 澶勭悊宓屽瀛楁 (濡?position.x)
      if (field.includes('.')) {
        const parts = field.split('.');
        if (!data[parts[0]]) data[parts[0]] = {};
        data[parts[0]][parts[1]] = value;
      } else {
        data[field] = value;
      }
    });

    try {
      const endpoint = this.currentView === 'nodes' ? 'nodes' : 'paths';
      const response = await fetch(`${API_BASE}api/v1/${endpoint}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      });

      if (!response.ok) throw new Error('淇濆瓨澶辫触');
      
      // 鍒锋柊鏁版嵁
      await this.loadData();
      this.renderTable();
      
      // 鏄剧ず鎴愬姛娑堟伅
      this.showMessage('淇濆瓨鎴愬姛', 'success');
      
    } catch (error) {
      console.error('淇濆瓨澶辫触:', error);
      this.showMessage('淇濆瓨澶辫触: ' + error.message, 'error');
    }
  }

  async deleteItem(id) {
    try {
      const endpoint = this.currentView === 'nodes' ? 'nodes' : 'paths';
      const response = await fetch(`${API_BASE}api/v1/${endpoint}/${id}`, {
        method: 'DELETE'
      });

      if (!response.ok) throw new Error('鍒犻櫎澶辫触');
      
      // 鍒锋柊鏁版嵁
      await this.loadData();
      this.renderTable();
      
      // 鏄剧ず鎴愬姛娑堟伅
      this.showMessage('鍒犻櫎鎴愬姛', 'success');
      
    } catch (error) {
      console.error('鍒犻櫎澶辫触:', error);
      this.showMessage('鍒犻櫎澶辫触: ' + error.message, 'error');
    }
  }

  async addNewItem() {
    const data = this.currentView === 'nodes' ? 
      {
        name: '鏂拌妭鐐?,
        type: 'default',
        status: 'active',
        position: { x: 100, y: 100, z: 0 }
      } : 
      {
        name: '鏂拌矾寰?,
        type: 'default',
        status: 'active',
        start_node_id: Object.keys(this.data.nodes)[0] || '',
        end_node_id: Object.keys(this.data.nodes)[1] || ''
      };

    try {
      const endpoint = this.currentView === 'nodes' ? 'nodes' : 'paths';
      const response = await fetch(`${API_BASE}api/v1/${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      });

      if (!response.ok) throw new Error('鍒涘缓澶辫触');
      
      // 鍒锋柊鏁版嵁
      await this.loadData();
      this.renderTable();
      
      // 鏄剧ず鎴愬姛娑堟伅
      this.showMessage('鍒涘缓鎴愬姛', 'success');
      
    } catch (error) {
      console.error('鍒涘缓澶辫触:', error);
      this.showMessage('鍒涘缓澶辫触: ' + error.message, 'error');
    }
  }

  showMessage(message, type) {
    // 鍒涘缓娑堟伅鍏冪礌
    const msgEl = document.createElement('div');
    msgEl.className = `message message-${type}`;
    msgEl.textContent = message;
    msgEl.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      padding: 12px 24px;
      border-radius: 4px;
      color: white;
      font-weight: 500;
      z-index: 1000;
      background: ${type === 'success' ? '#27ae60' : '#e74c3c'};
    `;
    
    document.body.appendChild(msgEl);
    
    // 3绉掑悗鑷姩绉婚櫎
    setTimeout(() => {
      if (msgEl.parentNode) {
        msgEl.parentNode.removeChild(msgEl);
      }
    }, 3000);
  }
}

// 椤甸潰鍔犺浇瀹屾垚鍚庡垵濮嬪寲琛ㄦ牸瑙嗗浘
document.addEventListener('DOMContentLoaded', () => {
  new TableView();
});