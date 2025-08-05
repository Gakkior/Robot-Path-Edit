// 鏈哄櫒浜鸿矾寰勭紪杈戝櫒鍓嶇搴旂敤
console.log('鏈哄櫒浜鸿矾寰勭紪杈戝櫒鍓嶇搴旂敤鍔犺浇瀹屾垚');

// API鍩虹閰嶇疆
const API_BASE = '/api/v1';

// 绠€鍗曠殑API瀹㈡埛绔?
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

// 椤甸潰鍔犺浇瀹屾垚鍚庡垵濮嬪寲
document.addEventListener('DOMContentLoaded', function() {
    console.log('椤甸潰鍔犺浇瀹屾垚锛屽垵濮嬪寲搴旂敤...');
    
    // 娴嬭瘯API杩炴帴
    api.get('/nodes')
        .then(data => {
            console.log('鑺傜偣鏁版嵁:', data);
        })
        .catch(error => {
            console.error('鑾峰彇鑺傜偣鏁版嵁澶辫触:', error);
        });
        
    api.get('/paths')
        .then(data => {
            console.log('璺緞鏁版嵁:', data);
        })
        .catch(error => {
            console.error('鑾峰彇璺緞鏁版嵁澶辫触:', error);
        });
});