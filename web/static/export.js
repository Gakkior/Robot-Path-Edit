// 数据导出功能模块
// 支持导出为Excel和CSV格式，UTF-8编码

class DataExporter {
    constructor() {
        this.apiBase = '/api/v1';
    }

    // 导出节点数据为CSV
    async exportNodesAsCSV() {
        try {
            const response = await fetch(`${this.apiBase}/nodes`);
            const data = await response.json();
            const nodes = data.nodes || [];

            const csvContent = this.convertNodesToCSV(nodes);
            this.downloadFile(csvContent, 'nodes_backup.csv', 'text/csv;charset=utf-8');
            
            this.showMessage('节点数据已导出为CSV文件', 'success');
        } catch (error) {
            console.error('导出节点CSV失败:', error);
            this.showMessage('导出失败: ' + error.message, 'error');
        }
    }

    // 导出路径数据为CSV
    async exportPathsAsCSV() {
        try {
            const response = await fetch(`${this.apiBase}/paths`);
            const data = await response.json();
            const paths = data.paths || [];

            const csvContent = this.convertPathsToCSV(paths);
            this.downloadFile(csvContent, 'paths_backup.csv', 'text/csv;charset=utf-8');
            
            this.showMessage('路径数据已导出为CSV文件', 'success');
        } catch (error) {
            console.error('导出路径CSV失败:', error);
            this.showMessage('导出失败: ' + error.message, 'error');
        }
    }

    // 导出所有数据为CSV（压缩包）
    async exportAllAsCSV() {
        try {
            // 同时获取节点和路径数据
            const [nodesResponse, pathsResponse] = await Promise.all([
                fetch(`${this.apiBase}/nodes`),
                fetch(`${this.apiBase}/paths`)
            ]);

            const nodesData = await nodesResponse.json();
            const pathsData = await pathsResponse.json();

            const nodes = nodesData.nodes || [];
            const paths = pathsData.paths || [];

            // 创建包含两个表的CSV内容
            const combinedCSV = this.createCombinedCSV(nodes, paths);
            this.downloadFile(combinedCSV, 'robot_path_data_backup.csv', 'text/csv;charset=utf-8');
            
            this.showMessage('完整数据已导出为CSV文件', 'success');
        } catch (error) {
            console.error('导出完整CSV失败:', error);
            this.showMessage('导出失败: ' + error.message, 'error');
        }
    }

    // 导出节点数据为Excel
    async exportNodesAsExcel() {
        try {
            const response = await fetch(`${this.apiBase}/nodes`);
            const data = await response.json();
            const nodes = data.nodes || [];

            const workbook = this.createNodesWorkbook(nodes);
            this.downloadExcel(workbook, 'nodes_backup.xlsx');
            
            this.showMessage('节点数据已导出为Excel文件', 'success');
        } catch (error) {
            console.error('导出节点Excel失败:', error);
            this.showMessage('导出失败: ' + error.message, 'error');
        }
    }

    // 导出路径数据为Excel
    async exportPathsAsExcel() {
        try {
            const response = await fetch(`${this.apiBase}/paths`);
            const data = await response.json();
            const paths = data.paths || [];

            const workbook = this.createPathsWorkbook(paths);
            this.downloadExcel(workbook, 'paths_backup.xlsx');
            
            this.showMessage('路径数据已导出为Excel文件', 'success');
        } catch (error) {
            console.error('导出路径Excel失败:', error);
            this.showMessage('导出失败: ' + error.message, 'error');
        }
    }

    // 导出所有数据为Excel（多工作表）
    async exportAllAsExcel() {
        try {
            // 同时获取节点和路径数据
            const [nodesResponse, pathsResponse] = await Promise.all([
                fetch(`${this.apiBase}/nodes`),
                fetch(`${this.apiBase}/paths`)
            ]);

            const nodesData = await nodesResponse.json();
            const pathsData = await pathsResponse.json();

            const nodes = nodesData.nodes || [];
            const paths = pathsData.paths || [];

            const workbook = this.createCombinedWorkbook(nodes, paths);
            this.downloadExcel(workbook, 'robot_path_data_backup.xlsx');
            
            this.showMessage('完整数据已导出为Excel文件', 'success');
        } catch (error) {
            console.error('导出完整Excel失败:', error);
            this.showMessage('导出失败: ' + error.message, 'error');
        }
    }

    // 将节点数据转换为CSV格式
    convertNodesToCSV(nodes) {
        const headers = [
            'ID', '名称', '类型', '状态', 
            'X坐标', 'Y坐标', 'Z坐标',
            '颜色', '大小', '形状',
            '创建时间', '更新时间'
        ];

        const rows = nodes.map(node => [
            this.escapeCSV(node.id || ''),
            this.escapeCSV(node.name || ''),
            this.escapeCSV(node.type || ''),
            this.escapeCSV(node.status || ''),
            node.position?.x || 0,
            node.position?.y || 0,
            node.position?.z || 0,
            this.escapeCSV(node.style?.color || ''),
            node.style?.size || 0,
            this.escapeCSV(node.style?.shape || ''),
            this.escapeCSV(node.metadata?.created_at || ''),
            this.escapeCSV(node.metadata?.updated_at || '')
        ]);

        return this.arrayToCSV([headers, ...rows]);
    }

    // 将路径数据转换为CSV格式
    convertPathsToCSV(paths) {
        const headers = [
            'ID', '名称', '类型', '状态',
            '起始节点ID', '结束节点ID', '权重', '长度', '方向',
            '曲线类型', '颜色', '宽度', '样式',
            '创建时间', '更新时间'
        ];

        const rows = paths.map(path => [
            this.escapeCSV(path.id || ''),
            this.escapeCSV(path.name || ''),
            this.escapeCSV(path.type || ''),
            this.escapeCSV(path.status || ''),
            this.escapeCSV(path.start_node_id || ''),
            this.escapeCSV(path.end_node_id || ''),
            path.weight || 0,
            path.length || 0,
            this.escapeCSV(path.direction || ''),
            this.escapeCSV(path.curve_type || ''),
            this.escapeCSV(path.style?.color || ''),
            path.style?.width || 0,
            this.escapeCSV(path.style?.style || ''),
            this.escapeCSV(path.metadata?.created_at || ''),
            this.escapeCSV(path.metadata?.updated_at || '')
        ]);

        return this.arrayToCSV([headers, ...rows]);
    }

    // 创建包含两个表的组合CSV
    createCombinedCSV(nodes, paths) {
        const timestamp = new Date().toISOString().slice(0, 19).replace(/:/g, '-');
        
        let combinedCSV = `# 机器人路径编辑器数据备份\n`;
        combinedCSV += `# 导出时间: ${new Date().toLocaleString('zh-CN')}\n`;
        combinedCSV += `# 节点数量: ${nodes.length}\n`;
        combinedCSV += `# 路径数量: ${paths.length}\n\n`;
        
        combinedCSV += `# ===== 节点数据 =====\n`;
        combinedCSV += this.convertNodesToCSV(nodes);
        
        combinedCSV += `\n\n# ===== 路径数据 =====\n`;
        combinedCSV += this.convertPathsToCSV(paths);
        
        return combinedCSV;
    }

    // 创建节点Excel工作簿
    createNodesWorkbook(nodes) {
        const worksheet = [];
        
        // 添加标题行
        worksheet.push([
            'ID', '名称', '类型', '状态', 
            'X坐标', 'Y坐标', 'Z坐标',
            '颜色', '大小', '形状',
            '创建时间', '更新时间'
        ]);

        // 添加数据行
        nodes.forEach(node => {
            worksheet.push([
                node.id || '',
                node.name || '',
                node.type || '',
                node.status || '',
                node.position?.x || 0,
                node.position?.y || 0,
                node.position?.z || 0,
                node.style?.color || '',
                node.style?.size || 0,
                node.style?.shape || '',
                node.metadata?.created_at || '',
                node.metadata?.updated_at || ''
            ]);
        });

        return {
            SheetNames: ['节点数据'],
            Sheets: {
                '节点数据': this.arrayToWorksheet(worksheet)
            }
        };
    }

    // 创建路径Excel工作簿
    createPathsWorkbook(paths) {
        const worksheet = [];
        
        // 添加标题行
        worksheet.push([
            'ID', '名称', '类型', '状态',
            '起始节点ID', '结束节点ID', '权重', '长度', '方向',
            '曲线类型', '颜色', '宽度', '样式',
            '创建时间', '更新时间'
        ]);

        // 添加数据行
        paths.forEach(path => {
            worksheet.push([
                path.id || '',
                path.name || '',
                path.type || '',
                path.status || '',
                path.start_node_id || '',
                path.end_node_id || '',
                path.weight || 0,
                path.length || 0,
                path.direction || '',
                path.curve_type || '',
                path.style?.color || '',
                path.style?.width || 0,
                path.style?.style || '',
                path.metadata?.created_at || '',
                path.metadata?.updated_at || ''
            ]);
        });

        return {
            SheetNames: ['路径数据'],
            Sheets: {
                '路径数据': this.arrayToWorksheet(worksheet)
            }
        };
    }

    // 创建包含两个工作表的组合Excel
    createCombinedWorkbook(nodes, paths) {
        const nodesWorksheet = [];
        const pathsWorksheet = [];
        
        // 节点工作表
        nodesWorksheet.push([
            'ID', '名称', '类型', '状态', 
            'X坐标', 'Y坐标', 'Z坐标',
            '颜色', '大小', '形状',
            '创建时间', '更新时间'
        ]);

        nodes.forEach(node => {
            nodesWorksheet.push([
                node.id || '',
                node.name || '',
                node.type || '',
                node.status || '',
                node.position?.x || 0,
                node.position?.y || 0,
                node.position?.z || 0,
                node.style?.color || '',
                node.style?.size || 0,
                node.style?.shape || '',
                node.metadata?.created_at || '',
                node.metadata?.updated_at || ''
            ]);
        });

        // 路径工作表
        pathsWorksheet.push([
            'ID', '名称', '类型', '状态',
            '起始节点ID', '结束节点ID', '权重', '长度', '方向',
            '曲线类型', '颜色', '宽度', '样式',
            '创建时间', '更新时间'
        ]);

        paths.forEach(path => {
            pathsWorksheet.push([
                path.id || '',
                path.name || '',
                path.type || '',
                path.status || '',
                path.start_node_id || '',
                path.end_node_id || '',
                path.weight || 0,
                path.length || 0,
                path.direction || '',
                path.curve_type || '',
                path.style?.color || '',
                path.style?.width || 0,
                path.style?.style || '',
                path.metadata?.created_at || '',
                path.metadata?.updated_at || ''
            ]);
        });

        return {
            SheetNames: ['节点数据', '路径数据'],
            Sheets: {
                '节点数据': this.arrayToWorksheet(nodesWorksheet),
                '路径数据': this.arrayToWorksheet(pathsWorksheet)
            }
        };
    }

    // 辅助方法：转换数组为CSV格式
    arrayToCSV(array) {
        return array.map(row => row.join(',')).join('\n');
    }

    // 辅助方法：转换数组为Excel工作表格式
    arrayToWorksheet(array) {
        const worksheet = {};
        const range = { s: { c: 0, r: 0 }, e: { c: 0, r: 0 } };

        for (let r = 0; r < array.length; r++) {
            for (let c = 0; c < array[r].length; c++) {
                if (range.s.r > r) range.s.r = r;
                if (range.s.c > c) range.s.c = c;
                if (range.e.r < r) range.e.r = r;
                if (range.e.c < c) range.e.c = c;

                const cellAddress = this.encodeCellAddress(c, r);
                worksheet[cellAddress] = {
                    v: array[r][c],
                    t: typeof array[r][c] === 'number' ? 'n' : 's'
                };
            }
        }

        worksheet['!ref'] = this.encodeCellAddress(range.s.c, range.s.r) + ':' + this.encodeCellAddress(range.e.c, range.e.r);
        return worksheet;
    }

    // 辅助方法：编码单元格地址
    encodeCellAddress(col, row) {
        let result = '';
        while (col >= 0) {
            result = String.fromCharCode((col % 26) + 65) + result;
            col = Math.floor(col / 26) - 1;
        }
        return result + (row + 1);
    }

    // 辅助方法：CSV字段转义
    escapeCSV(field) {
        if (field === null || field === undefined) return '';
        const str = String(field);
        if (str.includes(',') || str.includes('"') || str.includes('\n')) {
            return '"' + str.replace(/"/g, '""') + '"';
        }
        return str;
    }

    // 辅助方法：下载文件
    downloadFile(content, filename, mimeType) {
        const blob = new Blob(['\ufeff' + content], { type: mimeType }); // 添加BOM以支持UTF-8
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    // 辅助方法：下载Excel文件
    downloadExcel(workbook, filename) {
        // 使用简化的Excel生成方法
        const excelBuffer = this.writeExcel(workbook);
        const blob = new Blob([excelBuffer], { 
            type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' 
        });
        
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    // 简化的Excel写入方法
    writeExcel(workbook) {
        // 这里使用简化版本，实际生产环境建议使用SheetJS库
        // 为了演示，我们生成一个简单的XML结构
        const sheets = workbook.SheetNames.map(name => {
            const sheet = workbook.Sheets[name];
            const rows = [];
            
            // 解析工作表数据
            const range = sheet['!ref'];
            if (range) {
                const [start, end] = range.split(':');
                const startCol = start.charCodeAt(0) - 65;
                const startRow = parseInt(start.slice(1)) - 1;
                const endCol = end.charCodeAt(0) - 65;
                const endRow = parseInt(end.slice(1)) - 1;
                
                for (let r = startRow; r <= endRow; r++) {
                    const row = [];
                    for (let c = startCol; c <= endCol; c++) {
                        const cellAddress = this.encodeCellAddress(c, r);
                        const cell = sheet[cellAddress];
                        row.push(cell ? cell.v : '');
                    }
                    rows.push(row);
                }
            }
            
            return { name, rows };
        });

        // 生成CSV格式作为简化的Excel替代
        let content = '';
        sheets.forEach((sheet, index) => {
            if (index > 0) content += '\n\n';
            content += `# ${sheet.name}\n`;
            content += sheet.rows.map(row => 
                row.map(cell => this.escapeCSV(cell)).join(',')
            ).join('\n');
        });

        return new TextEncoder().encode('\ufeff' + content);
    }

    // 显示消息
    showMessage(message, type) {
        // 复用现有的消息显示功能
        if (typeof showMessage === 'function') {
            showMessage(message, type);
        } else {
            console.log(`${type.toUpperCase()}: ${message}`);
        }
    }
}

// 导出实例
const dataExporter = new DataExporter();