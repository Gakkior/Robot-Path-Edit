// 机器人路径编辑器前端应用
console.log('机器人路径编辑器前端应用加载完成');

// API基础配置
const API_BASE = '/api/v1';

// 简单的API客户�?
class ApiClient {
    async get(endpoint) {
        const response = await fetch(`${API_BASE}${endpoint}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    }
    
    async post(endpoint, data) {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    }
}

const api = new ApiClient();

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    console.log('页面加载完成，初始化应用...');
    
    // 测试API连接
    api.get('/nodes')
        .then(data => {
            console.log('节点数据:', data);
        })
        .catch(error => {
            console.error('获取节点数据失败:', error);
        });
        
    api.get('/paths')
        .then(data => {
            console.log('路径数据:', data);
        })
        .catch(error => {
            console.error('获取路径数据失败:', error);
        });
});