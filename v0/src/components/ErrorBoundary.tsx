/**
 * 错误边界组件
 * 
 * 捕获和处理React组件树中的错误
 */

import React from 'react'
import { motion } from 'framer-motion'
import { AlertTriangle, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/Button'

interface ErrorBoundaryState {
  hasError: boolean
  error?: Error
}

export class ErrorBoundary extends React.Component<
  React.PropsWithChildren<{}>,
  ErrorBoundaryState
> {
  constructor(props: React.PropsWithChildren<{}>) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo)
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-center p-8 max-w-md"
          >
            <AlertTriangle className="h-16 w-16 text-red-500 mx-auto mb-4" />
            <h1 className="text-2xl font-bold text-gray-900 mb-2">
              应用出现错误
            </h1>
            <p className="text-gray-600 mb-6">
              很抱歉，应用遇到了一个意外错误。请尝试刷新页面。
            </p>
            <div className="space-y-3">
              <Button
                onClick={() => window.location.reload()}
                className="w-full"
              >
                <RefreshCw className="h-4 w-4 mr-2" />
                刷新页面
              </Button>
              <details className="text-left">
                <summary className="cursor-pointer text-sm text-gray-500 hover:text-gray-700">
                  查看错误详情
                </summary>
                <pre className="mt-2 p-3 bg-gray-100 rounded text-xs overflow-auto">
                  {this.state.error?.stack}
                </pre>
              </details>
            </div>
          </motion.div>
        </div>
      )
    }

    return this.props.children
  }
}
