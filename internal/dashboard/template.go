package dashboard

var tmpl string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Monitor</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #0f172a;
            --card-bg: #1e293b;
            --text-primary: #f8fafc;
            --text-secondary: #94a3b8;
            --accent-color: #38bdf8;
            --border-color: #334155;
            --success-color: #22c55e;
            --danger-color: #ef4444;
        }

        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg-color);
            color: var(--text-primary);
            margin: 0;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
        }

        header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--border-color);
        }

        h1 {
            font-size: 1.5rem;
            font-weight: 700;
            margin: 0;
            color: var(--accent-color);
        }

        .status-badge {
            background: #059669;
            color: white;
            padding: 4px 12px;
            border-radius: 9999px;
            font-size: 0.875rem;
            font-weight: 500;
        }

        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }

        .section-header {
            font-size: 1.25rem;
            font-weight: 600;
            margin: 30px 0 15px 0;
            color: var(--text-primary);
            border-left: 4px solid var(--accent-color);
            padding-left: 10px;
        }

        .card {
            background: var(--card-bg);
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            display: flex;
            flex-direction: column;
        }

        .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
        }

        .card-title {
            font-size: 1.1rem;
            font-weight: 600;
            color: var(--text-primary);
        }

        .stat-value {
            font-size: 1.5rem;
            font-weight: 700;
            color: var(--accent-color);
        }

        .info-row {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid var(--border-color);
            font-size: 0.9rem;
        }
        
        .info-row:last-child {
            border-bottom: none;
        }

        .info-label {
            color: var(--text-secondary);
        }

        canvas {
            max-height: 250px;
            width: 100% !important;
        }

        /* Improved Table Styles */
        .table-container {
            width: 100%;
            overflow-x: auto;
            max-height: 400px;
            overflow-y: auto;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 0.9rem;
            text-align: left;
        }

        th {
            background-color: #334155;
            color: var(--text-primary);
            font-weight: 600;
            position: sticky;
            top: 0;
            z-index: 10;
        }

        th, td {
            padding: 12px 15px;
            border-bottom: 1px solid var(--border-color);
        }

        tr:hover {
            background-color: #334155;
        }

        .text-right {
            text-align: right;
        }
        
        .text-center {
            text-align: center;
        }

        .badge {
            display: inline-block;
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 0.8rem;
            background: #334155;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>System Monitor</h1>
            <div id="connection-status" class="status-badge">Live</div>
        </header>

        <div class="section-header">Overview</div>
        <div class="grid">
            <div class="card">
                <div class="card-header">
                    <span class="card-title">System Info</span>
                    <span id="uptime-badge" class="status-badge" style="background: #eab308; display: none;">Uptime: 0h</span>
                </div>
                <div id="system-info">Loading...</div>
                <div id="host-info"></div>
            </div>
            <div class="card">
                <div class="card-header">
                    <span class="card-title">CPU Load</span>
                    <span id="cpu-total" class="stat-value">0%</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Load Avg (1m | 5m | 15m)</span>
                    <span id="load-avg" style="font-weight: 600;">0.00 | 0.00 | 0.00</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Physical Cores</span>
                    <span id="cpu-physical">0</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Logical Cores</span>
                    <span id="cpu-logical">0</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Processes</span>
                    <span id="proc-count">0</span>
                </div>
            </div>
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Memory Usage</span>
                    <span id="mem-total-percent" class="stat-value">0%</span>
                </div>
                <canvas id="memoryChart"></canvas>
            </div>
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Swap Memory</span>
                    <span id="swap-total-percent" class="stat-value">0%</span>
                </div>
                <canvas id="swapChart"></canvas>
            </div>
        </div>

        <div class="section-header">Hardware Metrics</div>
        <div class="grid">
            <div class="card" style="grid-column: span 2;">
                <div class="card-header">
                    <span class="card-title">CPU Usage Per Core</span>
                </div>
                <canvas id="cpuCoreChart"></canvas>
            </div>
            <div class="card" style="grid-column: span 1;">
                <div class="card-header">
                    <span class="card-title">Disk Usage</span>
                </div>
                <canvas id="diskChart"></canvas>
            </div>
        </div>
        
        <div class="section-header">Network Deep Dive</div>
        <div class="grid">
            <div class="card">
                 <div class="card-header">
                    <span class="card-title">Total Connections</span>
                    <span id="total-conn-count" class="stat-value">0</span>
                </div>
                 <div class="card-header">
                    <span class="card-title">Connection States</span>
                </div>
                <canvas id="netConnChart"></canvas>
            </div>
             <div class="card" style="grid-column: span 2;">
                <div class="card-header">
                    <span class="card-title">Network Interface Details</span>
                </div>
                <div class="table-container">
                    <table>
                        <thead>
                            <tr>
                                <th>Interface</th>
                                <th class="text-right">Bytes Sent</th>
                                <th class="text-right">Bytes Recv</th>
                                <th class="text-right">Sent Pkts</th>
                                <th class="text-right">Recv Pkts</th>
                                <th class="text-right">Errors In</th>
                                <th class="text-right">Drops In</th>
                            </tr>
                        </thead>
                        <tbody id="io-table-body">
                            <tr><td colspan="7">Loading...</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
        <div class="grid">
             <div class="card" style="grid-column: span 3;">
                <div class="card-header">
                    <span class="card-title">Established Connections</span>
                </div>
                 <div class="table-container" style="max-height: 300px;">
                    <table>
                        <thead>
                            <tr>
                                <th>Local Address</th>
                                <th>Remote Address</th>
                                <th>PID</th>
                                <th>Status</th>
                            </tr>
                        </thead>
                        <tbody id="conn-table-body">
                            <tr><td colspan="4">Loading...</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div class="section-header">Memory Deep Dive</div>
        <div class="grid">
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Virtual vs Physical (Top Procs)</span>
                </div>
                <canvas id="memVmsRssChart"></canvas>
            </div>
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Heap vs Stack (Top Procs)</span>
                </div>
                <canvas id="memHeapStackChart"></canvas>
            </div>
        </div>

        <div class="section-header">Top Processes</div>
        <div class="grid">
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Top 5 Processes (CPU)</span>
                </div>
                <canvas id="topCpuChart"></canvas>
            </div>
            <div class="card">
                <div class="card-header">
                    <span class="card-title">Top 5 Processes (Memory)</span>
                </div>
                <canvas id="topMemChart"></canvas>
            </div>
             <div class="card">
                <div class="card-header">
                    <span class="card-title">Top 5 Processes (Threads)</span>
                </div>
                <canvas id="topThreadsChart"></canvas>
            </div>
             <div class="card">
                <div class="card-header">
                    <span class="card-title">Processes by User</span>
                </div>
                <canvas id="userProcChart"></canvas>
            </div>
        </div>
    </div>

    <script>
        // Charts
        let memoryChart, swapChart, diskChart, cpuCoreChart;
        let topCpuChart, topMemChart, topThreadsChart;
        let netConnChart, userProcChart;
        let memVmsRssChart, memHeapStackChart;

        function initCharts() {
            const commonOptions = {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { labels: { color: '#94a3b8' } } },
                scales: {
                    x: { grid: { color: '#334155' }, ticks: { color: '#94a3b8' } },
                    y: { grid: { color: '#334155' }, ticks: { color: '#94a3b8' } }
                }
            };
            
            const doughnutOptions = {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { position: 'right', labels: { color: '#94a3b8' } }
                }
            };

            // Memory Chart
            memoryChart = new Chart(document.getElementById('memoryChart').getContext('2d'), {
                type: 'doughnut',
                data: {
                    labels: ['Used', 'Free', 'Other'],
                    datasets: [{
                        data: [0, 0, 0],
                        backgroundColor: ['#ef4444', '#22c55e', '#3b82f6'],
                        borderWidth: 0
                    }]
                },
                options: doughnutOptions
            });

            // Swap Chart
            swapChart = new Chart(document.getElementById('swapChart').getContext('2d'), {
                type: 'doughnut',
                data: {
                    labels: ['Used', 'Free'],
                    datasets: [{
                        data: [0, 0],
                        backgroundColor: ['#f59e0b', '#10b981'],
                        borderWidth: 0
                    }]
                },
                options: doughnutOptions
            });

            // CPU Core Chart
            cpuCoreChart = new Chart(document.getElementById('cpuCoreChart').getContext('2d'), {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Load %',
                        data: [],
                        backgroundColor: '#0ea5e9',
                        borderRadius: 4
                    }]
                },
                options: {
                    ...commonOptions,
                    plugins: { legend: { display: false } },
                    scales: {
                        y: { beginAtZero: true, max: 100, grid: { color: '#334155' }, ticks: { color: '#94a3b8' } },
                        x: { grid: { color: '#334155' }, ticks: { color: '#94a3b8' } }
                    }
                }
            });

            // Disk Chart
            diskChart = new Chart(document.getElementById('diskChart').getContext('2d'), {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Used (GB)',
                        data: [],
                        backgroundColor: '#38bdf8',
                        borderRadius: 4
                    }, {
                        label: 'Free (GB)',
                        data: [],
                        backgroundColor: '#334155',
                        borderRadius: 4
                    }]
                },
                options: {
                    indexAxis: 'y',
                    ...commonOptions,
                    scales: {
                        x: { stacked: true, grid: { color: '#334155' }, ticks: { color: '#94a3b8' } },
                        y: { stacked: true, grid: { display: false }, ticks: { color: '#94a3b8' } }
                    }
                }
            });

            // Top CPU Chart
            topCpuChart = new Chart(document.getElementById('topCpuChart').getContext('2d'), {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'CPU %',
                        data: [],
                        backgroundColor: '#f472b6',
                        borderRadius: 4
                    }]
                },
                options: commonOptions
            });

            // Top Mem Chart
            topMemChart = new Chart(document.getElementById('topMemChart').getContext('2d'), {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Memory %',
                        data: [],
                        backgroundColor: '#a78bfa',
                        borderRadius: 4
                    }]
                },
                options: commonOptions
            });
            
             // Top Threads Chart
            topThreadsChart = new Chart(document.getElementById('topThreadsChart').getContext('2d'), {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Threads',
                        data: [],
                        backgroundColor: '#ec4899',
                        borderRadius: 4
                    }]
                },
                options: commonOptions
            });

            // Network Connections Chart
            netConnChart = new Chart(document.getElementById('netConnChart').getContext('2d'), {
                type: 'doughnut',
                data: {
                    labels: ['Established', 'Other'],
                    datasets: [{
                        data: [0, 0],
                        backgroundColor: ['#22c55e', '#64748b'],
                        borderWidth: 0
                    }]
                },
                options: doughnutOptions
            });

            // User Proc Chart
            userProcChart = new Chart(document.getElementById('userProcChart').getContext('2d'), {
                type: 'doughnut',
                data: {
                    labels: [],
                    datasets: [{
                        data: [],
                        backgroundColor: ['#38bdf8', '#f472b6', '#a78bfa', '#facc15', '#4ade80', '#fbbf24', '#94a3b8'],
                        borderWidth: 0
                    }]
                },
                options: doughnutOptions
            });

            // Mem Deep Dive: VMS vs RSS
            memVmsRssChart = new Chart(document.getElementById('memVmsRssChart').getContext('2d'), {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Virtual (MB)',
                        data: [],
                        backgroundColor: '#6366f1',
                        borderRadius: 4
                    }, {
                        label: 'Physical (MB)',
                        data: [],
                        backgroundColor: '#14b8a6',
                        borderRadius: 4
                    }]
                },
                options: commonOptions
            });

            // Mem Deep Dive: Heap vs Stack
            memHeapStackChart = new Chart(document.getElementById('memHeapStackChart').getContext('2d'), {
                type: 'bar',
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Heap (MB)',
                        data: [],
                        backgroundColor: '#f59e0b',
                        borderRadius: 4
                    }, {
                        label: 'Stack (MB)',
                        data: [],
                        backgroundColor: '#8b5cf6',
                        borderRadius: 4
                    }]
                },
                options: {
                    ...commonOptions,
                    scales: {
                        x: { stacked: true, grid: { color: '#334155' }, ticks: { color: '#94a3b8' } },
                        y: { stacked: true, grid: { color: '#334155' }, ticks: { color: '#94a3b8' } }
                    }
                }
            });
        }

        function formatBytes(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        async function updateData() {
            try {
                const response = await fetch('/api/metrics');
                const data = await response.json();
                
                // Update System Info
                if (data.user) {
                    document.getElementById('system-info').innerHTML = 
                        '<div class="info-row"><span class="info-label">User</span><span>' + data.user.Username + '</span></div>' +
                        '<div class="info-row"><span class="info-label">Hostname</span><span>' + (data.user.FullName || 'N/A') + '</span></div>';
                }
                
                if (data.host) {
                    document.getElementById('host-info').innerHTML = 
                        '<div class="info-row"><span class="info-label">OS</span><span>' + data.host.OS + ' ' + (data.host.PlatformVer || '') + '</span></div>' +
                        '<div class="info-row"><span class="info-label">Kernel</span><span>' + data.host.KernelVersion + '</span></div>' + 
                        '<div class="info-row"><span class="info-label">Platform</span><span>' + data.host.Platform + '</span></div>';

                    // Uptime
                    const uptimeSec = data.host.Uptime;
                    const days = Math.floor(uptimeSec / 86400);
                    const hours = Math.floor((uptimeSec % 86400) / 3600);
                    const minutes = Math.floor((uptimeSec % 3600) / 60);
                    
                    let uptimeStr = '';
                    if (days > 0) uptimeStr += days + 'd ';
                    uptimeStr += hours + 'h ' + minutes + 'm';
                    
                    const uptimeBadge = document.getElementById('uptime-badge');
                    uptimeBadge.style.display = 'inline-block';
                    uptimeBadge.textContent = 'Up: ' + uptimeStr;

                    // Load Avg
                    document.getElementById('load-avg').textContent = 
                        data.host.LoadAvg1.toFixed(2) + ' | ' + 
                        data.host.LoadAvg5.toFixed(2) + ' | ' + 
                        data.host.LoadAvg15.toFixed(2);
                }

                // Update CPU & Cores
                if (data.cpu) {
                    document.getElementById('cpu-total').textContent = data.cpu.AveragePercentages.toFixed(1) + '%';
                    document.getElementById('cpu-physical').textContent = data.cpu.PhysicalCores;
                    document.getElementById('cpu-logical').textContent = data.cpu.LogicalCores;
                    document.getElementById('proc-count').textContent = data.cpu.Processes ? data.cpu.Processes.length : 0;

                    if (data.cpu.Percentages) {
                        cpuCoreChart.data.labels = data.cpu.Percentages.map((_, i) => 'Core ' + i);
                        cpuCoreChart.data.datasets[0].data = data.cpu.Percentages;
                        cpuCoreChart.update();
                    }

                    if (data.cpu.Processes && data.cpu.Processes.length > 0) {
                        // Top CPU
                        const sortedByCpu = [...data.cpu.Processes].sort((a, b) => b.CpuPercent - a.CpuPercent).slice(0, 5);
                        topCpuChart.data.labels = sortedByCpu.map(p => p.ProcessName || p.Pid);
                        topCpuChart.data.datasets[0].data = sortedByCpu.map(p => p.CpuPercent);
                        topCpuChart.update();
                        
                        // Top Threads
                        const sortedByThreads = [...data.cpu.Processes].sort((a, b) => b.NumThreads - a.NumThreads).slice(0, 5);
                        topThreadsChart.data.labels = sortedByThreads.map(p => p.ProcessName || p.Pid);
                        topThreadsChart.data.datasets[0].data = sortedByThreads.map(p => p.NumThreads);
                        topThreadsChart.update();

                        // User Processes Distribution
                        const userCounts = {};
                        data.cpu.Processes.forEach(p => {
                            const u = p.Username || 'unknown';
                            userCounts[u] = (userCounts[u] || 0) + 1;
                        });
                        userProcChart.data.labels = Object.keys(userCounts);
                        userProcChart.data.datasets[0].data = Object.values(userCounts);
                        userProcChart.update();
                    }
                }

                // Update Memory, Swap & Details
                if (data.memory) {
                    if (data.memory.Vmemory) {
                        const vmem = data.memory.Vmemory;
                        document.getElementById('mem-total-percent').textContent = vmem.UsedPercentage.toFixed(1) + '%';
                        
                        const used = (vmem.Used / 1024 / 1024 / 1024).toFixed(2);
                        const free = (vmem.Free / 1024 / 1024 / 1024).toFixed(2);
                        const other = ((vmem.Total - vmem.Used - vmem.Free) / 1024 / 1024 / 1024).toFixed(2);

                        memoryChart.data.datasets[0].data = [used, free, other];
                        memoryChart.update();
                    }
                    
                    const swapUsed = (data.memory.SwapMemoryUsed / 1024 / 1024 / 1024).toFixed(2);
                    const swapFree = (data.memory.SwapMemoryFree / 1024 / 1024 / 1024).toFixed(2);
                    document.getElementById('swap-total-percent').textContent = data.memory.SwapMemoryUsedPercent.toFixed(1) + '%';
                    swapChart.data.datasets[0].data = [swapUsed, swapFree];
                    swapChart.update();

                    if (data.memory.ProcessInfo && data.memory.ProcessInfo.length > 0) {
                        const sortedProcs = [...data.memory.ProcessInfo].sort((a, b) => b.MemPercent - a.MemPercent).slice(0, 5);
                        topMemChart.data.labels = sortedProcs.map(p => p.ProcessName || p.Pid);
                        topMemChart.data.datasets[0].data = sortedProcs.map(p => p.MemPercent);
                        topMemChart.update();

                        // Update Memory Deep Dive Charts
                        memVmsRssChart.data.labels = sortedProcs.map(p => p.ProcessName || p.Pid);
                        memVmsRssChart.data.datasets[0].data = sortedProcs.map(p => (p.VirtualMemorySize / 1024 / 1024).toFixed(2));
                        memVmsRssChart.data.datasets[1].data = sortedProcs.map(p => (p.PhysicalMemorySize / 1024 / 1024).toFixed(2));
                        memVmsRssChart.update();

                        memHeapStackChart.data.labels = sortedProcs.map(p => p.ProcessName || p.Pid);
                        memHeapStackChart.data.datasets[0].data = sortedProcs.map(p => (p.MemoryUsedByHeap / 1024 / 1024).toFixed(2));
                        memHeapStackChart.data.datasets[1].data = sortedProcs.map(p => (p.MemoryUsedByStack / 1024 / 1024).toFixed(2));
                        memHeapStackChart.update();
                    }
                }

                // Update Disk
                if (data.disk && data.disk.UsageStat) {
                    const paths = [];
                    const used = [];
                    const free = [];
                    for (const [path, stat] of Object.entries(data.disk.UsageStat)) {
                        paths.push(path);
                        used.push((stat.UsedDisk / 1024 / 1024 / 1024).toFixed(2));
                        free.push((stat.FreeDisk / 1024 / 1024 / 1024).toFixed(2));
                    }
                    diskChart.data.labels = paths;
                    diskChart.data.datasets[0].data = used;
                    diskChart.data.datasets[1].data = free;
                    diskChart.update();
                }
                
                // Update Network & Tables
                if (data.network) {
                     // IO Stats Table
                     if (data.network.IOStats) {
                         let html = '';
                         data.network.IOStats.forEach(stat => {
                             html += '<tr>' +
                                    '<td><span class="badge">' + (stat.name || 'unknown') + '</span></td>' +
                                    '<td class="text-right">' + formatBytes(stat.bytesSent || 0) + '</td>' +
                                    '<td class="text-right">' + formatBytes(stat.bytesRecv || 0) + '</td>' +
                                    '<td class="text-right">' + (stat.packetsSent || 0).toLocaleString() + '</td>' +
                                    '<td class="text-right">' + (stat.packetsRecv || 0).toLocaleString() + '</td>' +
                                    '<td class="text-right">' + (stat.errin || 0).toLocaleString() + '</td>' +
                                    '<td class="text-right">' + (stat.dropin || 0).toLocaleString() + '</td>' +
                                '</tr>';
                         });
                         document.getElementById('io-table-body').innerHTML = html;
                     }

                     // Connections
                     const totalConn = data.network.NumTotalConnections || 0;
                     const estConn = data.network.NumEstablishedConnections || 0;
                     const otherConn = totalConn - estConn;
                     netConnChart.data.datasets[0].data = [estConn, otherConn];
                     netConnChart.update();
                     document.getElementById('total-conn-count').innerText = totalConn;

                     // Connections Table
                     if (data.network.Connections) {
                        let connHtml = '';
                         // Limit to top 50 for performance
                         data.network.Connections.slice(0, 50).forEach(c => {
                             if (!c.LocalAddr || !c.RemoteAddr) return;
                             connHtml += '<tr>' +
                                    '<td>' + c.LocalAddr.ip + ':' + c.LocalAddr.port + '</td>' +
                                    '<td>' + c.RemoteAddr.ip + ':' + c.RemoteAddr.port + '</td>' +
                                    '<td>' + c.Pid + '</td>' +
                                    '<td><span class="badge">' + (c.Status || 'UNKNOWN') + '</span></td>' +
                                '</tr>';
                         });
                         document.getElementById('conn-table-body').innerHTML = connHtml || '<tr><td colspan="4">No visible established connections</td></tr>';
                     }
                }

            } catch (err) {
                console.error("Error fetching metrics:", err);
            }
        }

        initCharts();
        updateData();
        setInterval(updateData, 5000); 
    </script>
</body>
</html>

`
