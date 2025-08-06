/**
 * 类名合并工具
 * 
 * 基于clsx和tailwind-merge的组合工具
 * 用于智能合并Tailwind CSS类名，避免冲突
 * 
 * 参考：shadcn/ui的设计模式
 */

import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
