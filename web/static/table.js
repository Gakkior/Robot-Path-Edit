// 表格视图脚本
console.log('加载表格视图脚本');

const API_BASE = '/';

class TableView {
  constructor() {
    this.currentView = 'nodes'; // 'nodes' �?'paths'
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
      if (!response.ok) throw new Error('获取数据失败');
      this.data = await response.json();
    } catch (error) {
      console.error('加载数据失败:', error);
    }
  }

  setupEventListeners() {
    // 视图切换按钮
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
    
    // 更新按钮状�?
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
        <h3>节点管理 (${nodes.length} 个节�?</h3>
      </div>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>类型</th>
              <th>状�?/th>
              <th>X坐标</th>
              <th>Y坐标</th>
              <th>Z坐标</th>
              <th>创建时间</th>
              <th>操作</th>
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
            <button class="btn-small btn-save" data-id="${node.id}">保存</button>
            <button class="btn-small btn-delete" data-id="${node.id}">删除</button>
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
        <h3>路径管理 (${paths.length} 个路�?</h3>
      </div>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>类型</th>
              <th>状�?/th>
              <th>起始节点</th>
              <th>结束节点</th>
              <th>创建时间</th>
              <th>操作</th>
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
            <button class="btn-small btn-save" data-id="${path.id}">保存</button>
            <button class="btn-small btn-delete" data-id="${path.id}">删除</button>
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
    // 保存按钮
    document.querySelectorAll('.btn-save').forEach(btn => {
      btn.addEventListener('click', async (e) => {
        const id = e.target.dataset.id;
        await this.saveItem(id);
      });
    });

    // 删除按钮
    document.querySelectorAll('.btn-delete').forEach(btn => {
      btn.addEventListener('click', async (e) => {
        const id = e.target.dataset.id;
        if (confirm(`确定删除这个${this.currentView === 'nodes' ? '节点' : '路径'}吗？`)) {
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
      
      // 处理数字类型
      if (input.type === 'number') {
        value = parseFloat(value) || 0;
      }
      
      // 处理嵌套字段 (�?position.x)
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

      if (!response.ok) throw new Error('保存失败');
      
      // 刷新数据
      await this.loadData();
      this.renderTable();
      
      // 显示成功消息
      this.showMessage('保存成功', 'success');
      
    } catch (error) {
      console.error('保存失败:', error);
      this.showMessage('保存失败: ' + error.message, 'error');
    }
  }

  async deleteItem(id) {
    try {
      const endpoint = this.currentView === 'nodes' ? 'nodes' : 'paths';
      const response = await fetch(`${API_BASE}api/v1/${endpoint}/${id}`, {
        method: 'DELETE'
      });

      if (!response.ok) throw new Error('删除失败');
      
      // 刷新数据
      await this.loadData();
      this.renderTable();
      
      // 显示成功消息
      this.showMessage('删除成功', 'success');
      
    } catch (error) {
      console.error('删除失败:', error);
      this.showMessage('删除失败: ' + error.message, 'error');
    }
  }

  async addNewItem() {
    const data = this.currentView === 'nodes' ? 
      {
        name: '新节�?,
        type: 'default',
        status: 'active',
        position: { x: 100, y: 100, z: 0 }
      } : 
      {
        name: '新路�?,
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

      if (!response.ok) throw new Error('创建失败');
      
      // 刷新数据
      await this.loadData();
      this.renderTable();
      
      // 显示成功消息
      this.showMessage('创建成功', 'success');
      
    } catch (error) {
      console.error('创建失败:', error);
      this.showMessage('创建失败: ' + error.message, 'error');
    }
  }

  showMessage(message, type) {
    // 创建消息元素
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
    
    // 3秒后自动移除
    setTimeout(() => {
      if (msgEl.parentNode) {
        msgEl.parentNode.removeChild(msgEl);
      }
    }, 3000);
  }
}

// 页面加载完成后初始化表格视图
document.addEventListener('DOMContentLoaded', () => {
  new TableView();
});