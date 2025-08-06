/**
* 应用入口文件
* 
* 负责：
* 1. React应用的挂载
* 2. 全局样式的引入
* 3. 开发环境的配置
*/

import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'

// 移除加载屏幕
const loadingScreen = document.getElementById('loading-screen')
if (loadingScreen) {
  loadingScreen.remove()
}

// 确保DOM完全加载后再渲染React应用
document.addEventListener('DOMContentLoaded', () => {
  const rootElement = document.getElementById('root')
  if (rootElement) {
    ReactDOM.createRoot(rootElement).render(
      <React.StrictMode>
        <App />
      </React.StrictMode>,
    )
  } else {
    console.error('Root element not found. Make sure an element with id="root" exists in your index.html.')
  }
})

// 将App组件作为main.tsx的默认导出，以满足某些预览环境的要求
export default App;
