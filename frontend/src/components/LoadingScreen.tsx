/**
 * 加载屏幕组件
 * 
 * 应用启动时的加载界面
 */

import React, { useEffect, useState } from 'react'
import { motion } from 'framer-motion'

export const LoadingScreen: React.FC = () => {
  const [progress, setProgress] = useState(0)
  
  useEffect(() => {
    const timer = setInterval(() => {
      setProgress(prev => {
        if (prev >= 100) {
          clearInterval(timer)
          return 100
        }
        return prev + Math.random() * 15
      })
    }, 200)
    
    return () => clearInterval(timer)
  }, [])
  
  return (
    <div className="fixed inset-0 bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center z-50">
      <motion.div
        initial={{ opacity: 0, scale: 0.9 }}
        animate={{ opacity: 1, scale: 1 }}
        className="text-center"
      >
        {/* Logo动画 */}
        <motion.div
          animate={{ 
            rotate: 360,
            scale: [1, 1.1, 1]
          }}
          transition={{ 
            rotate: { duration: 2, repeat: Infinity, ease: "linear" },
            scale: { duration: 1, repeat: Infinity }
          }}
          className="text-6xl mb-6"
        >
          🤖
        </motion.div>
        
        {/* 标题 */}
        <motion.h1
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.2 }}
          className="text-2xl font-bold text-white mb-2"
        >
          机器人路径编辑器
        </motion.h1>
        
        <motion.p
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="text-gray-300 mb-8"
        >
          正在加载现代化的路径管理工具...
        </motion.p>
        
        {/* 进度条 */}
        <div className="w-64 h-2 bg-gray-700 rounded-full overflow-hidden">
          <motion.div
            className="h-full bg-gradient-to-r from-blue-500 to-purple-500"
            initial={{ width: 0 }}
            animate={{ width: `${progress}%` }}
            transition={{ duration: 0.3 }}
          />
        </div>
        
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="text-gray-400 text-sm mt-4"
        >
          {progress < 30 && "初始化组件..."}
          {progress >= 30 && progress < 60 && "加载数据..."}
          {progress >= 60 && progress < 90 && "准备界面..."}
          {progress >= 90 && "即将完成..."}
        </motion.p>
      </motion.div>
    </div>
  )
}
