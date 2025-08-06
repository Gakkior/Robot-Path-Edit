// Konva 画布初始化脚本
console.log('加载 Konva 画布脚本');

const API_BASE = '/';

let stage, layer;
const sidebar = document.getElementById('sidebar');
let selectedPathId = null;

// 命令模式 - 撤销/重做系统
class CommandManager {
  constructor() {
    this.history = [];
    this.currentIndex = -1;
    this.maxHistory = 50;
  }

  async executeCommand(command) {
    // 清除当前位置之后的历史
    this.history = this.history.slice(0, this.currentIndex + 1);
    
    // 执行命令
    await command.execute();
    
    // 添加到历史
    this.history.push(command);
    this.currentIndex++;
    
    // 限制历史长度
    if (this.history.length > this.maxHistory) {
      this.history.shift();
      this.currentIndex--;
    }
    
    this.updateUI();
  }

  async undo() {
    if (this.canUndo()) {
      const command = this.history[this.currentIndex];
      await command.undo();
      this.currentIndex--;
      this.updateUI();
    }
  }

  async redo() {
    if (this.canRedo()) {
      this.currentIndex++;
      const command = this.history[this.currentIndex];
      await command.execute();
      this.updateUI();
    }
  }

  canUndo() {
    return this.currentIndex >= 0;
  }

  canRedo() {
    return this.currentIndex < this.history.length - 1;
  }

  updateUI() {
    // 更新撤销/重做按钮状态
    const undoBtn = document.getElementById('undoBtn');
    const redoBtn = document.getElementById('redoBtn');
    if (undoBtn) undoBtn.disabled = !this.canUndo();
    if (redoBtn) redoBtn.disabled = !this.canRedo();
  }
}

const commandManager = new CommandManager();

// 移动节点命令
class MoveNodeCommand {
  constructor(nodeId, oldPosition, newPosition) {
    this.nodeId = nodeId;
    this.oldPosition = oldPosition;
    this.newPosition = newPosition;
  }

  async execute() {
    await this.updateNodePosition(this.newPosition);
  }

  async undo() {
    await this.updateNodePosition(this.oldPosition);
  }

  async updateNodePosition(position) {
    await fetch(API_BASE + 'api/v1/nodes/' + this.nodeId + '/position', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(position)
    });
    // 重新加载画布数据
    await loadCanvasData();
  }
}

// 创建路径命令
class CreatePathCommand {
  constructor(pathData) {
    this.pathData = pathData;
    this.pathId = null;
  }

  async execute() {
    const response = await fetch(API_BASE + 'api/v1/paths', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(this.pathData)
    });
    const result = await response.json();
    this.pathId = result.path.id;
    await loadCanvasData();
  }

  async undo() {
    if (this.pathId) {
      await fetch(API_BASE + 'api/v1/paths/' + this.pathId, { method: 'DELETE' });
      await loadCanvasData();
    }
  }
}

// 删除路径命令
class DeletePathCommand {
  constructor(pathId, pathData) {
    this.pathId = pathId;
    this.pathData = pathData;
  }

  async execute() {
    await fetch(API_BASE + 'api/v1/paths/' + this.pathId, { method: 'DELETE' });
    await loadCanvasData();
  }

  async undo() {
    const response = await fetch(API_BASE + 'api/v1/paths', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(this.pathData)
    });
    const result = await response.json();
    this.pathId = result.path.id;
    await loadCanvasData();
  }
}
const nodeForm = document.getElementById('nodeForm');

async function fetchCanvasData() {
  const res = await fetch(API_BASE + 'canvas-data');
  if (!res.ok) throw new Error('获取画布数据失败');
  return await res.json();
}

function openSidebar(node, circle) {
  sidebar.style.display = 'block';
  document.getElementById('nodeId').value = node.id;
  document.getElementById('nodeName').value = node.name;
  document.getElementById('posX').value = circle.x().toFixed(2);
  document.getElementById('posY').value = circle.y().toFixed(2);
  document.getElementById('posZ').value = node.position.z || 0;
}

nodeForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const id = document.getElementById('nodeId').value;
  const name = document.getElementById('nodeName').value;
  const x = parseFloat(document.getElementById('posX').value);
  const y = parseFloat(document.getElementById('posY').value);
  const z = parseFloat(document.getElementById('posZ').value);
  await fetch(API_BASE + 'api/v1/nodes/' + id, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, position: { x, y, z } }),
  });
  await fetch(API_BASE + 'api/v1/nodes/' + id + '/position', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ x, y, z }),
  });
  location.reload();
});

async function createPath(startId, endId) {
  const pathData = { name: '新路径', start_node_id: startId, end_node_id: endId };
  const command = new CreatePathCommand(pathData);
  await commandManager.executeCommand(command);
}

function drawCanvas(data) {
  const width = window.innerWidth;
  const height = window.innerHeight;

  stage = new Konva.Stage({
    container: 'canvas-container',
    width,
    height,
    draggable: true,
  });

  layer = new Konva.Layer();
  stage.add(layer);

  const nodeMap = {};
let selectedNodeId = null;

  // 绘制节点
  Object.values(data.nodes).forEach((n) => {
    const circle = new Konva.Circle({
      x: n.position.x,
      y: n.position.y,
      radius: 20,
      fill: '#3498db',
      stroke: '#2980b9',
      strokeWidth: 2,
      draggable: true,
    });

    const text = new Konva.Text({
      x: n.position.x - 20,
      y: n.position.y - 6,
      text: n.name,
      fontSize: 12,
      fill: '#fff',
      width: 40,
      align: 'center',
    });

    circle.on('click', (e) => {
      e.cancelBubble = true;
      openSidebar(n, circle);
      if (e.evt.shiftKey) {
        if (!selectedNodeId) {
          selectedNodeId = n.id;
        } else if (selectedNodeId && selectedNodeId !== n.id) {
          createPath(selectedNodeId, n.id);
          selectedNodeId = null;
        }
      }
    });

    let originalPosition = { x: n.position.x, y: n.position.y, z: n.position.z || 0 };
    
    circle.on('dragstart', () => {
      originalPosition = { x: circle.x(), y: circle.y(), z: n.position.z || 0 };
    });

    circle.on('dragmove', () => {
      text.position({ x: circle.x() - 20, y: circle.y() - 6 });
      redrawPaths();
    });

    circle.on('dragend', async () => {
      const newPosition = { x: circle.x(), y: circle.y(), z: n.position.z || 0 };
      const command = new MoveNodeCommand(n.id, originalPosition, newPosition);
      await commandManager.executeCommand(command);
    });

    layer.add(circle);
    layer.add(text);
    nodeMap[n.id] = { circle, text };
  });

  // 绘制路径
  function redrawPaths() {
    // 清理旧路径
    layer.find('Line').forEach((l) => l.destroy());

    data.paths && Object.values(data.paths).forEach((p) => {
      const startNode = nodeMap[p.start_node_id];
      const endNode = nodeMap[p.end_node_id];
      if (!startNode || !endNode) return;
      const line = new Konva.Line({
        points: [startNode.circle.x(), startNode.circle.y(), endNode.circle.x(), endNode.circle.y()],
        stroke: '#34495e',
        strokeWidth: 2,
        lineCap: 'round',
        id: p.id,
      });
      line.on('click', (e) => {
        e.cancelBubble = true;
        if (selectedPathId === p.id) {
          selectedPathId = null;
          line.stroke('#34495e');
        } else {
          // 取消其他选中
          layer.find('Line').forEach((l) => l.stroke('#34495e'));
          selectedPathId = p.id;
          line.stroke('#e74c3c');
        }
        layer.batchDraw();
      });
      layer.add(line);
    });
    layer.batchDraw();
  }

  redrawPaths();
  layer.draw();
}

async function updateNodePosition(id, x, y, z) {
  await fetch(API_BASE + 'api/v1/nodes/' + id + '/position', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ x, y, z }),
  });
}

let currentCanvasData = null;

async function loadCanvasData() {
  try {
    currentCanvasData = await fetchCanvasData();
    drawCanvas(currentCanvasData);
  } catch (err) {
    console.error('加载画布数据失败:', err);
  }
}

// 应用布局算法
async function applyLayout(algorithm) {
  try {
    const response = await fetch(API_BASE + 'api/v1/layout/apply', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ algorithm })
    });
    
    if (!response.ok) {
      throw new Error('布局应用失败');
    }
    
    const result = await response.json();
    console.log('布局应用成功:', result);
    
    // 重新加载画布数据
    await loadCanvasData();
    
    // 显示成功消息
    showMessage(`${algorithm}布局应用成功，影响了${result.affected_nodes}个节点`, 'success');
    
  } catch (error) {
    console.error('应用布局失败:', error);
    showMessage('布局应用失败: ' + error.message, 'error');
  }
}

// 生成路径
async function generatePaths(algorithm, params) {
  try {
    const response = await fetch(API_BASE + `api/v1/path-generation/${algorithm}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(params)
    });
    
    if (!response.ok) {
      throw new Error('路径生成失败');
    }
    
    const result = await response.json();
    console.log('路径生成成功:', result);
    
    // 重新加载画布数据
    await loadCanvasData();
    
    // 显示成功消息
    const algorithmNames = {
      'nearest-neighbor': '最近邻',
      'full-connectivity': '完全连通',
      'grid': '网格'
    };
    showMessage(`${algorithmNames[algorithm]}路径生成成功，创建了${result.created_paths}条路径`, 'success');
    
  } catch (error) {
    console.error('生成路径失败:', error);
    showMessage('路径生成失败: ' + error.message, 'error');
  }
}

// 显示消息
function showMessage(message, type) {
  const msgEl = document.createElement('div');
  msgEl.style.cssText = `
    position: fixed;
    top: 60px;
    right: 20px;
    padding: 12px 24px;
    border-radius: 4px;
    color: white;
    font-weight: 500;
    z-index: 1000;
    background: ${type === 'success' ? '#27ae60' : '#e74c3c'};
    animation: slideIn 0.3s ease-out;
  `;
  msgEl.textContent = message;
  
  document.body.appendChild(msgEl);
  
  setTimeout(() => {
    if (msgEl.parentNode) {
      msgEl.parentNode.removeChild(msgEl);
    }
  }, 3000);
}

(async () => {
  try {
    await loadCanvasData();
    
    // 初始化撤销/重做按钮事件
    const undoBtn = document.getElementById('undoBtn');
    const redoBtn = document.getElementById('redoBtn');
    
    if (undoBtn) {
      undoBtn.addEventListener('click', async () => {
        await commandManager.undo();
      });
    }
    
    if (redoBtn) {
      redoBtn.addEventListener('click', async () => {
        await commandManager.redo();
      });
    }
    
    // 初始化布局按钮事件
    const gridLayoutBtn = document.getElementById('gridLayoutBtn');
    const forceLayoutBtn = document.getElementById('forceLayoutBtn');
    const circularLayoutBtn = document.getElementById('circularLayoutBtn');
    
    if (gridLayoutBtn) {
      gridLayoutBtn.addEventListener('click', () => applyLayout('grid'));
    }
    
    if (forceLayoutBtn) {
      forceLayoutBtn.addEventListener('click', () => applyLayout('force-directed'));
    }
    
    if (circularLayoutBtn) {
      circularLayoutBtn.addEventListener('click', () => applyLayout('circular'));
    }

    // 初始化路径生成按钮事件
    const nearestPathBtn = document.getElementById('nearestPathBtn');
    const fullConnectBtn = document.getElementById('fullConnectBtn');
    const gridPathBtn = document.getElementById('gridPathBtn');
    
    if (nearestPathBtn) {
      nearestPathBtn.addEventListener('click', () => generatePaths('nearest-neighbor', { max_distance: 200 }));
    }
    
    if (fullConnectBtn) {
      fullConnectBtn.addEventListener('click', () => generatePaths('full-connectivity', {}));
    }
    
    if (gridPathBtn) {
      gridPathBtn.addEventListener('click', () => generatePaths('grid', { connect_diagonal: false }));
    }
    
    // 初始化按钮状态
    commandManager.updateUI();
    
  } catch (err) {
    console.error(err);
  }
})();

// 快捷键处理
window.addEventListener('keydown', async (e) => {
  // 删除选中路径
  if (e.key === 'Delete' && selectedPathId) {
    if (confirm('确定删除所选路径?')) {
      // 获取路径数据用于撤销
      const pathData = currentCanvasData.paths[selectedPathId];
      if (pathData) {
        const command = new DeletePathCommand(selectedPathId, pathData);
        await commandManager.executeCommand(command);
        selectedPathId = null;
      }
    }
  }
  
  // 撤销 (Ctrl+Z)
  if (e.ctrlKey && e.key === 'z' && !e.shiftKey) {
    e.preventDefault();
    await commandManager.undo();
  }
  
  // 重做 (Ctrl+Shift+Z 或 Ctrl+Y)
  if ((e.ctrlKey && e.shiftKey && e.key === 'Z') || (e.ctrlKey && e.key === 'y')) {
    e.preventDefault();
    await commandManager.redo();
  }
});