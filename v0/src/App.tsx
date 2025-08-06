/**
 * 主应用组件
 * 
 * 应用的根组件，负责：
 * 1. 提供全局上下文（React Query、主题等）
 * 2. 错误边界处理
 * 3. 加载状态管理
 * 4. 全局样式和动画配置
 */

import React, { Suspense } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { motion, AnimatePresence } from 'framer-motion'
import { ToastProvider, ToastViewport } from '@/components/ui/Toast'
import AppLayout from '@/components/Layout/AppLayout' // Changed to default import
import { ErrorBoundary } from '@/components/ErrorBoundary'
import { LoadingScreen } from '@/components/LoadingScreen'

// 创建React Query客户端
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 1000, // 30秒
      gcTime: 5 * 60 * 1000, // 5分钟
      retry: 3,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 1,
    },
  },
})

function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <ToastProvider>
          <div className="min-h-screen bg-gray-50">
            <Suspense fallback={<LoadingScreen />}>
              <AnimatePresence mode="wait">
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.3 }}
                >
                  <AppLayout />
                </motion.div>
              </AnimatePresence>
            </Suspense>
          </div>
          
          {/* Toast容器 */}
          <ToastViewport />
        </ToastProvider>
        
        {/* 开发工具 */}
        {process.env.NODE_ENV === 'development' && (
          <ReactQueryDevtools initialIsOpen={false} />
        )}
      </QueryClientProvider>
    </ErrorBoundary>
  )
}

export default App
