// Charts
let memoryChart, swapChart, diskChart, cpuCoreChart;
let topCpuChart, topMemChart, topThreadsChart;
let netConnChart, userProcChart;
let memVmsRssChart, memHeapStackChart;
let dockerContainerChart, dockerDiskChart;

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

    // Docker Container Status Chart
    dockerContainerChart = new Chart(document.getElementById('dockerContainerChart').getContext('2d'), {
        type: 'doughnut',
        data: {
            labels: ['Running', 'Paused', 'Stopped'],
            datasets: [{
                data: [0, 0, 0],
                backgroundColor: ['#22c55e', '#f59e0b', '#ef4444'],
                borderWidth: 0
            }]
        },
        options: doughnutOptions
    });

    // Docker Disk Usage Chart
    dockerDiskChart = new Chart(document.getElementById('dockerDiskChart').getContext('2d'), {
        type: 'bar',
        data: {
            labels: ['Containers', 'Images', 'Build Cache'],
            datasets: [{
                label: 'Size (GB)',
                data: [],
                backgroundColor: ['#38bdf8', '#a78bfa', '#fbbf24'],
                borderRadius: 4
            }]
        },
        options: {
            ...commonOptions,
            plugins: { legend: { display: false } },
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

        // Update Docker
        const dockerSection = document.getElementById('docker-section');
        if (data.docker) {
            dockerSection.style.display = 'block';
            const d = data.docker;

            document.getElementById('docker-version').textContent = 'v' + d.APIVersion;
            document.getElementById('docker-platform').textContent = d.OS + ' / ' + d.Arch;
            document.getElementById('docker-images').textContent = d.TotalImages;
            document.getElementById('docker-volumes').textContent = d.TotalVolumes;
            document.getElementById('docker-ncpu').textContent = d.NCpu || '0';
            document.getElementById('docker-mem-total').textContent = d.MemTotal ? (d.MemTotal / 1024 / 1024 / 1024).toFixed(2) + ' GB' : '0 GB';

            dockerContainerChart.data.datasets[0].data = [d.ContainersRunning, d.ContainersPaused, d.ContainersStopped];
            dockerContainerChart.update();

            dockerDiskChart.data.datasets[0].data = [
                (d.ContainersDiskUsage / 1024 / 1024 / 1024).toFixed(2),
                (d.ImagesDiskUsage / 1024 / 1024 / 1024).toFixed(2),
                (d.BuildCacheDiskUsage / 1024 / 1024 / 1024).toFixed(2)
            ];
            dockerDiskChart.update();

            if (d.ContainerStats) {
                let html = '';
                d.ContainerStats.forEach(c => {
                    let ports = '';
                    if (c.ContainerPorts) {
                        ports = c.ContainerPorts.map(p => (p.PublicPort ? p.PublicPort + ':' : '') + p.PrivatePort + '/' + p.Type).join(', ');
                    }

                    html += '<tr>' +
                        '<td>' + (c.ContainerNames ? c.ContainerNames.join(', ') : 'unknown') + '</td>' +
                        '<td>' + c.ImageName + '</td>' +
                        '<td><span class="badge" style="background:' + (c.ContainerState === 'running' ? '#059669' : '#b91c1c') + '">' + c.ContainerState + '</span></td>' +
                        '<td>' + formatBytes(c.ContainerRootSizeInBytes || 0) + '</td>' +
                        '<td>' + ports + '</td>' +
                        '</tr>';
                });
                document.getElementById('docker-container-body').innerHTML = html;
            }

            if (d.ImageStats) {
                let html = '';
                d.ImageStats.forEach(img => {
                    html += '<tr>' +
                        '<td>' + (img.ImageNames ? img.ImageNames.join('<br>') : 'unnamed') + '</td>' +
                        '<td title="' + img.ID + '"><code>' + img.ID.substring(7, 19) + '</code></td>' +
                        '<td>' + formatBytes(img.ImageSize || 0) + '</td>' +
                        '<td>' + img.CreatedDate + '</td>' +
                        '<td>' + img.NumberOfContainersUsingThisImage + '</td>' +
                        '</tr>';
                });
                document.getElementById('docker-images-body').innerHTML = html || '<tr><td colspan="5">No images found</td></tr>';
            }

            if (d.VolumeStats) {
                let html = '';
                d.VolumeStats.forEach(vol => {
                    html += '<tr>' +
                        '<td>' + vol.VolumeName + '</td>' +
                        '<td>' + vol.Driver + '</td>' +
                        '<td>' + vol.Scope + '</td>' +
                        '<td style="font-size: 0.8em; word-break: break-all;">' + vol.MountPoint + '</td>' +
                        '<td>' + formatBytes(vol.VolumeSize || 0) + '</td>' +
                        '</tr>';
                });
                document.getElementById('docker-volumes-body').innerHTML = html || '<tr><td colspan="5">No volumes found</td></tr>';
            }

        } else {
            dockerSection.style.display = 'none';
        }

    } catch (err) {
        console.error("Error fetching metrics:", err);
    }
}

initCharts();
updateData();
setInterval(updateData, 5000);
