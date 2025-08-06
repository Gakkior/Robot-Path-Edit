/**
 * Canvas工具函数
 * 
 * 提供画布相关的计算和转换功能
 */

import type { Position, Viewport } from '@/types' // Changed ViewportState to Viewport

/**
 * 屏幕坐标转换为画布坐标
 */
export function screenToCanvas(
  screenPos: Position,
  viewport: Viewport
): Position {
  return {
    x: (screenPos.x - viewport.x) / viewport.scale,
    y: (screenPos.y - viewport.y) / viewport.scale,
  }
}

/**
 * 画布坐标转换为屏幕坐标
 */
export function canvasToScreen(
  canvasPos: Position,
  viewport: Viewport
): Position {
  return {
    x: canvasPos.x * viewport.scale + viewport.x,
    y: canvasPos.y * viewport.scale + viewport.y,
  }
}

/**
 * 计算两点之间的距离
 */
export function distance(pos1: Position, pos2: Position): number {
  const dx = pos1.x - pos2.x
  const dy = pos1.y - pos2.y
  return Math.sqrt(dx * dx + dy * dy)
}

/**
 * 计算两点之间的角度（弧度）
 */
export function angle(pos1: Position, pos2: Position): number {
  return Math.atan2(pos2.y - pos1.y, pos2.x - pos1.x)
}

/**
 * 限制数值在指定范围内
 */
export function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max)
}

/**
 * 生成唯一ID
 */
export function generateId(): string {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

/**
 * 计算线段中心点
 * @param points [x1, y1, x2, y2]
 */
export function getLineCenter(points: number[]): Position {
  const x1 = points[0]
  const y1 = points[1]
  const x2 = points[2]
  const y2 = points[3]
  return {
    x: (x1 + x2) / 2,
    y: (y1 + y2) / 2,
  }
}
