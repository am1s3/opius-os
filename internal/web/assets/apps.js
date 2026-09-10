const ICONS = {
    terminal: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>`,
    devices: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="4" width="16" height="16" rx="2"/><rect x="9" y="9" width="6" height="6"/><line x1="9" y1="1" x2="9" y2="4"/><line x1="15" y1="1" x2="15" y2="4"/><line x1="9" y1="20" x2="9" y2="23"/><line x1="15" y1="20" x2="15" y2="23"/><line x1="20" y1="9" x2="23" y2="9"/><line x1="20" y1="14" x2="23" y2="14"/><line x1="1" y1="9" x2="4" y2="9"/><line x1="1" y1="14" x2="4" y2="14"/></svg>`,
    deviceStudio: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="7" width="20" height="14" rx="2" ry="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/><line x1="6" y1="11" x2="6" y2="11"/><line x1="10" y1="11" x2="10" y2="11"/><line x1="14" y1="11" x2="14" y2="11"/><line x1="18" y1="11" x2="18" y2="11"/></svg>`,
    packages: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>`,
    files: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>`,
    system: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>`,
    settings: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>`,
    about: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>`,
    folder: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>`,
    file: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><polyline points="13 2 13 9 20 9"/></svg>`,
    github: `<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/></svg>`,
    refresh: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>`,
    zap: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>`,
    trash: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>`,
    monitor: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>`,
    rotateCw: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>`,
};

window.ICONS = ICONS;

class AppRegistry {
    constructor() {
        this.apps = {
            terminal: { title: 'Terminal', factory: () => new TerminalApp() },
            devices: { title: 'Devices', factory: () => new DevicesApp() },
            deviceStudio: { title: 'Device Studio', factory: () => new DeviceStudioApp() },
            packages: { title: 'Packages', factory: () => new PackagesApp() },
            files: { title: 'Files', factory: () => new FilesApp() },
            system: { title: 'System Monitor', factory: () => new SystemApp() },
            settings: { title: 'Settings', factory: () => new SettingsApp() },
            about: { title: 'About Opius OS', factory: () => new AboutApp() },
        };
        this.instances = new Map();

        document.addEventListener('window-closing', (e) => {
            const instance = this.instances.get(e.detail.id);
            if (instance && instance.destroy) {
                instance.destroy();
            }
        });
    }

    launch(id) {
        const app = this.apps[id];
        if (!app) return;

        if (wm.windows.has(id)) {
            const win = wm.windows.get(id);
            if (win.classList.contains('minimized')) {
                wm.restore(id);
            } else {
                wm.focus(id);
            }
            return;
        }

        const instance = app.factory();
        this.instances.set(id, instance);
        const win = wm.open(id, app.title, instance.render(), instance.opts || {});

        instance.mount(wm.getBody(id));

        document.addEventListener('window-closed', (e) => {
            if (e.detail.id === id) {
                if (instance.destroy) instance.destroy();
                this.instances.delete(id);
            }
        });

        document.addEventListener('window-resized', (e) => {
            if (e.detail.id === id && instance.onResize) {
                instance.onResize();
            }
        });
    }
}

class TerminalApp {
    constructor() {
        this.opts = { width: 820, height: 520 };
    }

    render() {
        return '<div class="app-content terminal-app"><div class="terminal-container" style="height: 100%;"></div></div>';
    }

    mount(container) {
        this.container = container.querySelector('.terminal-container');
        this.term = new Terminal({
            cursorBlink: true,
            fontSize: 13,
            fontFamily: "'SF Mono', Monaco, 'Menlo', 'Consolas', monospace",
            theme: {
                background: '#0a0c10',
                foreground: '#e5e7eb',
                cursor: '#60a5fa',
                selectionBackground: 'rgba(96, 165, 250, 0.3)',
                black: '#0a0c10',
                red: '#ef4444',
                green: '#10b981',
                yellow: '#f59e0b',
                blue: '#3b82f6',
                magenta: '#a855f7',
                cyan: '#06b6d4',
                white: '#e5e7eb',
            },
        });
        this.fitAddon = new FitAddon.FitAddon();
        this.term.loadAddon(this.fitAddon);
        this.term.open(this.container);

        setTimeout(() => {
            this.fitAddon.fit();
            this.connect();
        }, 50);
    }

    connect() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws/terminal?id=${Date.now()}`;
        this.ws = new WebSocket(wsUrl);
        this.ws.binaryType = 'arraybuffer';

        this.ws.onopen = () => {
            this.term.writeln('\x1b[1;34mOPIUS OS Terminal\x1b[0m · connected');
            this.term.writeln('');
            this.sendResize();
        };

        this.ws.onmessage = (ev) => {
            this.term.write(new Uint8Array(ev.data));
        };

        this.ws.onclose = () => {
            this.term.writeln('\r\n\x1b[31m[disconnected]\x1b[0m');
        };

        this.term.onData((data) => {
            if (this.ws.readyState === WebSocket.OPEN) {
                this.ws.send(data);
            }
        });

        this.term.onResize(() => {
            this.sendResize();
        });
    }

    sendResize() {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;
        const cols = this.term.cols;
        const rows = this.term.rows;
        const buf = new Uint8Array(5);
        buf[0] = 0x01;
        buf[1] = (cols >> 8) & 0xff;
        buf[2] = cols & 0xff;
        buf[3] = (rows >> 8) & 0xff;
        buf[4] = rows & 0xff;
        this.ws.send(buf);
    }

    onResize() {
        if (this.fitAddon) {
            setTimeout(() => this.fitAddon.fit(), 50);
        }
    }

    destroy() {
        if (this.ws) this.ws.close();
        if (this.term) this.term.dispose();
    }
}

class DevicesApp {
    render() {
        return `<div class="app-content">
            <h2>Connected Devices</h2>
            <div id="devices-content" class="loading">Loading devices...</div>
            <div style="margin-top: 16px; display: flex; gap: 12px; align-items: center;">
                <button id="dev-refresh" style="display: flex; align-items: center; gap: 6px;">
                    ${ICONS.refresh} Refresh
                </button>
                <label style="display: flex; align-items: center; gap: 6px; color: #9ca3af; font-size: 13px; cursor: pointer;">
                    <input type="checkbox" id="dev-auto"> Auto-refresh (5s)
                </label>
            </div>
        </div>`;
    }

    mount(container) {
        this.container = container;
        this.autoInterval = null;
        this.load();
        container.querySelector('#dev-refresh').onclick = () => this.load();
        container.querySelector('#dev-auto').onchange = (e) => {
            if (e.target.checked) {
                this.autoInterval = setInterval(() => this.load(), 5000);
            } else {
                clearInterval(this.autoInterval);
            }
        };
    }

    async load() {
        try {
            const res = await fetch('/api/devices');
            const data = await res.json();
            const devs = data.data || [];
            const c = this.container.querySelector('#devices-content');

            if (devs.length === 0) {
                c.className = 'empty';
                c.innerHTML = 'No devices connected';
                return;
            }

            c.className = '';
            c.innerHTML = `<table>
                <tr><th>Port</th><th>Chip</th><th>Board</th><th>Status</th></tr>
                ${devs.map(d => `
                    <tr>
                        <td style="font-family: 'SF Mono', monospace; font-size: 12px;">${d.port}</td>
                        <td style="font-weight: 500;">${d.chip}</td>
                        <td style="color: #9ca3af;">${d.board}</td>
                        <td><span class="badge badge-${d.status === 'ready' ? 'ok' : d.status === 'boot' ? 'boot' : 'warn'}">${d.status}</span></td>
                    </tr>
                `).join('')}
            </table>`;
        } catch (e) {
            const c = this.container.querySelector('#devices-content');
            c.className = 'empty';
            c.innerHTML = '<span class="badge badge-err">Error: ' + e.message + '</span>';
        }
    }

    destroy() {
        if (this.autoInterval) clearInterval(this.autoInterval);
    }
}

class DeviceStudioApp {
    constructor() {
        this.opts = { width: 1100, height: 680 };
        this.selectedDevice = null;
        this.devices = [];
        this.refreshInterval = null;
        this.serialWS = null;
        this.monitorActive = false;
    }

    render() {
        return `<div class="app-content" style="padding: 0; height: 100%;">
            <div class="device-studio">
                <div class="device-sidebar">
                    <h3>Devices</h3>
                    <div id="ds-device-list"></div>
                    <button class="secondary" style="width: 100%; margin-top: 12px;" id="ds-refresh">
                        ${ICONS.refresh} Refresh
                    </button>
                </div>
                <div class="device-main" id="ds-main">
                    <div class="empty">Select a device to get started</div>
                </div>
            </div>
        </div>`;
    }

    mount(container) {
        this.container = container;
        this.loadDevices();
        this.refreshInterval = setInterval(() => this.loadDevices(), 5000);
        container.querySelector('#ds-refresh').onclick = () => this.loadDevices();
    }

    async loadDevices() {
        try {
            const res = await fetch('/api/devices');
            const data = await res.json();
            this.devices = data.data || [];
            this.renderDeviceList();
        } catch (e) {
            console.error('Failed to load devices:', e);
        }
    }

    renderDeviceList() {
        const list = this.container.querySelector('#ds-device-list');
        if (this.devices.length === 0) {
            list.innerHTML = '<div class="empty" style="padding: 20px;">No devices</div>';
            return;
        }

        list.innerHTML = this.devices.map(d => `
            <div class="device-list-item ${this.selectedDevice && this.selectedDevice.port === d.port ? 'selected' : ''}" data-port="${d.port}">
                <div class="port">${d.port}</div>
                <div class="info">${d.chip} · ${d.board}</div>
                <div style="margin-top: 4px;"><span class="badge badge-${d.status === 'ready' ? 'ok' : d.status === 'boot' ? 'boot' : 'warn'}">${d.status}</span></div>
            </div>
        `).join('');

        list.querySelectorAll('.device-list-item').forEach(el => {
            el.onclick = () => this.selectDevice(el.dataset.port);
        });
    }

    selectDevice(port) {
        this.stopMonitor();
        this.selectedDevice = this.devices.find(d => d.port === port);
        if (!this.selectedDevice) return;
        this.renderDeviceList();
        this.renderDeviceMain();
    }

    renderDeviceMain() {
        const main = this.container.querySelector('#ds-main');
        const d = this.selectedDevice;

        main.innerHTML = `
            <h2>${d.chip} · ${d.board}</h2>
            
            <div class="device-info">
                <div class="device-info-grid">
                    <div class="info-item">
                        <div class="info-label">Port</div>
                        <div class="info-value">${d.port}</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Chip</div>
                        <div class="info-value">${d.chip}</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Board</div>
                        <div class="info-value">${d.board}</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Status</div>
                        <div class="info-value"><span class="badge badge-${d.status === 'ready' ? 'ok' : d.status === 'boot' ? 'boot' : 'warn'}">${d.status}</span></div>
                    </div>
                    ${d.vid ? `<div class="info-item"><div class="info-label">VID</div><div class="info-value">${d.vid}</div></div>` : ''}
                    ${d.pid ? `<div class="info-item"><div class="info-label">PID</div><div class="info-value">${d.pid}</div></div>` : ''}
                </div>
            </div>

            <h3>Quick Actions</h3>
            <div class="device-actions">
                <button id="ds-flash-btn">${ICONS.zap}<span>Flash Package</span></button>
                <button id="ds-flash-custom-btn" class="secondary">${ICONS.zap}<span>Flash Custom</span></button>
                <button id="ds-erase-btn" class="secondary">${ICONS.trash}<span>Erase</span></button>
                <button id="ds-reset-btn" class="secondary">${ICONS.rotateCw}<span>Reset</span></button>
                <button id="ds-monitor-btn" class="secondary">${ICONS.monitor}<span>Monitor</span></button>
                <button id="ds-release-btn" class="danger">${ICONS.trash}<span>Release Port</span></button>
            </div>

            <div id="ds-flash-section" style="display: none;">
                <div class="flash-form">
                    <h3>Flash Package</h3>
                    <div class="form-row">
                        <label>Package</label>
                        <select id="ds-pkg-select"></select>
                    </div>
                    <div class="form-row">
                        <label></label>
                        <label><input type="checkbox" id="ds-erase-check"> Erase before flash</label>
                    </div>
                    <div style="display: flex; gap: 8px;">
                        <button id="ds-flash-go">Flash Device</button>
                        <button id="ds-flash-cancel" class="secondary">Cancel</button>
                    </div>
                    <div id="ds-flash-result"></div>
                </div>
            </div>

            <div id="ds-flash-custom-section" style="display: none;">
                <div class="flash-form">
                    <h3>Flash Custom Binaries</h3>
                    <p style="color: #9ca3af; font-size: 13px; margin-bottom: 16px;">
                        Upload binary files and specify memory addresses for each.
                    </p>
                    <div class="form-row">
                        <label>Chip</label>
                        <select id="ds-custom-chip">
                            <option value="esp32" selected>ESP32</option>
                            <option value="esp32-s2">ESP32-S2</option>
                            <option value="esp32-s3">ESP32-S3</option>
                            <option value="esp32-c3">ESP32-C3</option>
                            <option value="esp8266">ESP8266</option>
                            <option value="rp2040">RP2040</option>
                            <option value="avr">AVR (Arduino)</option>
                        </select>
                    </div>
                    <div id="ds-custom-bins"></div>
                    <button id="ds-add-bin" class="secondary" style="margin-bottom: 16px;">+ Add Binary</button>
                    <div class="form-row">
                        <label></label>
                        <label><input type="checkbox" id="ds-custom-erase-check"> Erase before flash</label>
                    </div>
                    <div style="display: flex; gap: 8px;">
                        <button id="ds-custom-flash-go">Flash Custom</button>
                        <button id="ds-custom-flash-cancel" class="secondary">Cancel</button>
                    </div>
                    <div id="ds-custom-flash-result"></div>
                </div>
            </div>

            <div id="ds-monitor-section" style="display: none;">
                <h3>Serial Monitor</h3>
                <div class="monitor-toolbar" style="margin-bottom: 12px; display: flex; gap: 8px;">
                    <input type="text" id="ds-monitor-input" placeholder="Send command to device..." style="flex: 1;">
                    <button id="ds-monitor-send">Send</button>
                    <button id="ds-monitor-clear" class="secondary">Clear</button>
                </div>
                <div class="monitor-output" id="ds-monitor-output"></div>
                <button id="ds-monitor-close" class="secondary" style="margin-top: 12px;">Stop Monitor</button>
            </div>
        `;

        this.attachActions();
        this.loadPackages();
    }

    attachActions() {
        const main = this.container.querySelector('#ds-main');

        main.querySelector('#ds-flash-btn').onclick = () => {
            main.querySelector('#ds-flash-section').style.display = 'block';
            main.querySelector('#ds-flash-custom-section').style.display = 'none';
            main.querySelector('#ds-monitor-section').style.display = 'none';
            this.stopMonitor();
        };

        main.querySelector('#ds-flash-custom-btn').onclick = () => {
            main.querySelector('#ds-flash-custom-section').style.display = 'block';
            main.querySelector('#ds-flash-section').style.display = 'none';
            main.querySelector('#ds-monitor-section').style.display = 'none';
            this.stopMonitor();
            if (this.container.querySelectorAll('#ds-custom-bins .form-row').length === 0) {
                this.addCustomBinRow();
            }
        };

        main.querySelector('#ds-erase-btn').onclick = () => {
            if (confirm('Erase device flash? This cannot be undone.')) {
                this.eraseDevice();
            }
        };

        main.querySelector('#ds-reset-btn').onclick = () => this.resetDevice();

        main.querySelector('#ds-monitor-btn').onclick = () => {
            main.querySelector('#ds-monitor-section').style.display = 'block';
            main.querySelector('#ds-flash-section').style.display = 'none';
            main.querySelector('#ds-flash-custom-section').style.display = 'none';
            this.startMonitor();
        };

        main.querySelector('#ds-release-btn').onclick = () => {
            this.stopMonitor();
            this.notify('Port released. Wait 1 second before reconnecting.', 'info');
        };

        main.querySelector('#ds-flash-go').onclick = () => this.flashDevice();
        main.querySelector('#ds-flash-cancel').onclick = () => {
            main.querySelector('#ds-flash-section').style.display = 'none';
        };

        main.querySelector('#ds-add-bin').onclick = () => this.addCustomBinRow();
        main.querySelector('#ds-custom-flash-go').onclick = () => this.flashCustom();
        main.querySelector('#ds-custom-flash-cancel').onclick = () => {
            main.querySelector('#ds-flash-custom-section').style.display = 'none';
        };

        main.querySelector('#ds-monitor-close').onclick = () => {
            main.querySelector('#ds-monitor-section').style.display = 'none';
            this.stopMonitor();
        };

        main.querySelector('#ds-monitor-send').onclick = () => this.sendCommand();
        main.querySelector('#ds-monitor-input').addEventListener('keydown', (e) => {
            if (e.key === 'Enter') this.sendCommand();
        });

        main.querySelector('#ds-monitor-clear').onclick = () => {
            main.querySelector('#ds-monitor-output').innerHTML = '';
        };
    }

    async loadPackages() {
        try {
            const res = await fetch('/api/packages');
            const data = await res.json();
            const pkgs = data.data || [];
            const select = this.container.querySelector('#ds-pkg-select');
            select.innerHTML = pkgs.map(p => `<option value="${p.name}">${p.name} v${p.version}</option>`).join('');
        } catch (e) {
            console.error('Failed to load packages:', e);
        }
    }

    async flashDevice() {
        const pkgName = this.container.querySelector('#ds-pkg-select').value;
        const erase = this.container.querySelector('#ds-erase-check').checked;
        const result = this.container.querySelector('#ds-flash-result');

        result.className = 'flash-result info';
        result.textContent = 'Flashing...';

        try {
            const res = await fetch('/api/flash', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    package_name: pkgName,
                    port: this.selectedDevice.port,
                    erase: erase,
                }),
            });

            const data = await res.json();

            if (data.success) {
                result.className = 'flash-result success';
                result.textContent = data.data.message;
                this.notify('Flash complete', 'success');
            } else {
                result.className = 'flash-result error';
                result.textContent = 'Error: ' + data.error;
                this.notify('Flash failed', 'error');
            }
        } catch (e) {
            result.className = 'flash-result error';
            result.textContent = 'Error: ' + e.message;
            this.notify('Flash failed', 'error');
        }
    }

    addCustomBinRow() {
        const container = this.container.querySelector('#ds-custom-bins');
        const row = document.createElement('div');
        row.className = 'form-row';
        row.innerHTML = `
            <label>Binary</label>
            <input type="file" accept=".bin" class="custom-bin-file">
            <input type="text" placeholder="0x10000" class="custom-bin-address" style="max-width: 120px;">
            <button class="danger" onclick="this.parentElement.remove()">Remove</button>
        `;
        container.appendChild(row);
    }

    async flashCustom() {
        const rows = this.container.querySelectorAll('#ds-custom-bins .form-row');
        const erase = this.container.querySelector('#ds-custom-erase-check').checked;
        const chip = this.container.querySelector('#ds-custom-chip').value;
        const result = this.container.querySelector('#ds-custom-flash-result');

        const binaries = [];
        for (const row of rows) {
            const fileInput = row.querySelector('.custom-bin-file');
            const addressInput = row.querySelector('.custom-bin-address');

            if (!fileInput.files[0] || !addressInput.value) {
                result.className = 'flash-result error';
                result.textContent = 'All binaries must have file and address';
                return;
            }

            const file = fileInput.files[0];
            const address = addressInput.value.trim();

            const formData = new FormData();
            formData.append('file', file);
            formData.append('address', address);

            const uploadRes = await fetch('/api/upload', {
                method: 'POST',
                body: formData,
            });

            const uploadData = await uploadRes.json();
            if (!uploadData.success) {
                result.className = 'flash-result error';
                result.textContent = 'Upload failed: ' + uploadData.error;
                return;
            }

            binaries.push({
                path: uploadData.data.path,
                address: address.startsWith('0x') ? address : '0x' + address,
            });
        }

        if (binaries.length === 0) {
            result.className = 'flash-result error';
            result.textContent = 'Add at least one binary';
            return;
        }

        result.className = 'flash-result info';
        result.textContent = 'Flashing ' + chip + '...';

        try {
            const res = await fetch('/api/flash-custom', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    port: this.selectedDevice.port,
                    chip: chip,
                    binaries: binaries,
                    erase: erase,
                }),
            });

            const data = await res.json();

            if (data.success) {
                result.className = 'flash-result success';
                result.textContent = data.data.message;
                this.notify('Custom flash complete', 'success');
            } else {
                result.className = 'flash-result error';
                result.textContent = 'Error: ' + data.error;
                this.notify('Custom flash failed', 'error');
            }
        } catch (e) {
            result.className = 'flash-result error';
            result.textContent = 'Error: ' + e.message;
            this.notify('Custom flash failed', 'error');
        }
    }

    async eraseDevice() {
        this.notify('Erasing device...', 'info');
        try {
            const res = await fetch('/api/erase', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ port: this.selectedDevice.port }),
            });
            const data = await res.json();
            if (data.success) {
                this.notify('Device erased', 'success');
            } else {
                this.notify('Erase failed: ' + data.error, 'error');
            }
        } catch (e) {
            this.notify('Erase failed: ' + e.message, 'error');
        }
    }

    async resetDevice() {
        try {
            const res = await fetch('/api/reset', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ port: this.selectedDevice.port }),
            });
            const data = await res.json();
            if (data.success) {
                this.notify('Reset signal sent', 'success');
            } else {
                this.notify('Reset failed: ' + data.error, 'error');
            }
        } catch (e) {
            this.notify('Reset failed: ' + e.message, 'error');
        }
    }

    startMonitor() {
        this.stopMonitor();

        const output = this.container.querySelector('#ds-monitor-output');
        output.innerHTML = '';

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws/serial?port=${encodeURIComponent(this.selectedDevice.port)}&baud=115200`;

        this.appendMonitorLine('status', 'Connecting to ' + this.selectedDevice.port + '...');

        this.serialWS = new WebSocket(wsUrl);
        this.monitorActive = true;

        this.serialWS.onopen = () => {
            this.appendMonitorLine('status', 'Monitor connected');
        };

        this.serialWS.onmessage = (ev) => {
            try {
                const msg = JSON.parse(ev.data);
                this.appendMonitorLine(msg.type, msg.content, msg.time);
            } catch (e) {
                this.appendMonitorLine('data', ev.data);
            }
        };

        this.serialWS.onerror = () => {
            this.appendMonitorLine('error', 'Connection error');
        };

        this.serialWS.onclose = () => {
            this.monitorActive = false;
            this.appendMonitorLine('status', 'Monitor disconnected');
        };
    }

    stopMonitor() {
        if (this.serialWS) {
            try {
                if (this.serialWS.readyState === WebSocket.OPEN) {
                    this.serialWS.send(JSON.stringify({ action: 'close' }));
                }
            } catch (e) {}
            this.serialWS.close();
            this.serialWS = null;
        }
        this.monitorActive = false;
    }

    appendMonitorLine(type, content, time) {
        const output = this.container.querySelector('#ds-monitor-output');
        if (!output) return;

        const timeStr = time || new Date().toLocaleTimeString('en-US', { hour12: false });
        const line = document.createElement('div');
        line.className = 'log-line';

        let colorClass = '';
        if (type === 'error') colorClass = ' style="color: #ef4444;"';
        else if (type === 'status') colorClass = ' style="color: #3b82f6;"';

        line.innerHTML = `<span class="log-time">${timeStr}</span><span${colorClass}>${this.escapeHtml(content)}</span>`;
        output.appendChild(line);
        output.scrollTop = output.scrollHeight;
    }

    sendCommand() {
        if (!this.serialWS || !this.monitorActive) return;
        const input = this.container.querySelector('#ds-monitor-input');
        const cmd = input.value.trim();
        if (!cmd) return;

        this.serialWS.send(JSON.stringify({ action: 'send', content: cmd }));
        this.appendMonitorLine('data', '> ' + cmd);
        input.value = '';
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    notify(message, type) {
        if (window.showToast) {
            window.showToast(message, type);
        } else {
            console.log(`[${type}] ${message}`);
        }
    }

    destroy() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
        }
        this.stopMonitor();
        this.selectedDevice = null;
        this.devices = [];
    }
}

class PackagesApp {
    render() {
        return `<div class="app-content">
            <h2>Package Store</h2>
            <div id="pkg-content" class="loading">Loading packages...</div>
            <div style="margin-top: 16px;">
                <button id="pkg-sync-btn" style="display: flex; align-items: center; gap: 6px;">
                    ${ICONS.refresh} Sync Registry
                </button>
            </div>
        </div>`;
    }

    mount(container) {
        this.container = container;
        this.load();
        container.querySelector('#pkg-sync-btn').onclick = () => this.syncRegistry();
    }

    async load() {
        try {
            const res = await fetch('/api/packages');
            const data = await res.json();
            const pkgs = data.data || [];
            const c = this.container.querySelector('#pkg-content');

            if (pkgs.length === 0) {
                c.className = 'empty';
                c.innerHTML = 'No packages available. Click <strong>Sync Registry</strong> to download packages.';
                return;
            }

            c.className = '';
            c.innerHTML = `<table>
                <tr><th>Name</th><th>Version</th><th>Description</th><th>Chips</th></tr>
                ${pkgs.map(p => `
                    <tr>
                        <td style="font-weight: 500;">${p.name}</td>
                        <td><span class="badge badge-ok">v${p.version}</span></td>
                        <td style="color: #9ca3af;">${p.description || '—'}</td>
                        <td style="color: #9ca3af; font-size: 12px;">${(p.chips || []).join(', ') || '—'}</td>
                    </tr>
                `).join('')}
            </table>`;
        } catch (e) {
            const c = this.container.querySelector('#pkg-content');
            c.className = 'empty';
            c.innerHTML = '<span class="badge badge-err">Error: ' + e.message + '</span>';
        }
    }

    async syncRegistry() {
        if (window.showToast) window.showToast('Syncing registry...', 'info');
        try {
            const res = await fetch('/api/sync', { method: 'POST' });
            const data = await res.json();
            if (data.success) {
                if (window.showToast) window.showToast('Registry synced', 'success');
                this.load();
            } else {
                if (window.showToast) window.showToast('Sync failed: ' + data.error, 'error');
            }
        } catch (e) {
            if (window.showToast) window.showToast('Sync failed: ' + e.message, 'error');
        }
    }
}

class FilesApp {
    render() {
        return `<div class="app-content">
            <h2>Files <span id="files-path" style="color: #9ca3af; font-size: 13px; font-weight: 400; margin-left: 8px;"></span></h2>
            <div id="files-list" class="file-list"></div>
        </div>`;
    }

    mount(container) {
        this.container = container;
        this.load('');
    }

    async load(path) {
        try {
            const url = '/api/fs/list' + (path ? '?path=' + encodeURIComponent(path) : '');
            const res = await fetch(url);
            const data = await res.json();
            const c = this.container.querySelector('#files-list');

            this.container.querySelector('#files-path').textContent = data.data.path;

            let html = '';
            if (!data.data.path.endsWith('.opius')) {
                const parent = data.data.path.split('/').slice(0, -1).join('/');
                html += `<div class="file-item dir" data-path="${parent}">
                    <span class="name">${ICONS.folder} ..</span><span class="size">Up</span>
                </div>`;
            }

            for (const e of data.data.entries) {
                const fullPath = data.data.path + '/' + e.name;
                if (e.is_dir) {
                    html += `<div class="file-item dir" data-path="${fullPath}">
                        <span class="name">${ICONS.folder} ${e.name}</span>
                        <span class="size">directory</span>
                    </div>`;
                } else {
                    const size = this.formatSize(e.size);
                    html += `<div class="file-item" data-path="${fullPath}">
                        <span class="name">${ICONS.file} ${e.name}</span>
                        <span class="size">${size}</span>
                    </div>`;
                }
            }

            c.innerHTML = html;
            c.querySelectorAll('.file-item.dir').forEach(el => {
                el.onclick = () => this.load(el.dataset.path);
            });
        } catch (e) {
            this.container.querySelector('#files-list').innerHTML = '<div class="empty"><span class="badge badge-err">Error: ' + e.message + '</span></div>';
        }
    }

    formatSize(bytes) {
        if (bytes >= 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB';
        if (bytes >= 1024) return (bytes / 1024).toFixed(1) + ' KB';
        return bytes + ' B';
    }
}

class SystemApp {
    render() {
        return `<div class="app-content">
            <h2>System Monitor</h2>
            <div id="sys-content" class="stat-grid"><div class="loading">Loading...</div></div>
        </div>`;
    }

    mount(container) {
        this.container = container;
        this.load();
        this.interval = setInterval(() => this.load(), 2000);
    }

    async load() {
        try {
            const res = await fetch('/api/system');
            const data = await res.json();
            const s = data.data;
            const c = this.container.querySelector('#sys-content');
            c.innerHTML = `
                <div class="stat-box"><div class="stat-label">Operating System</div><div class="stat-value">${s.os} / ${s.arch}</div></div>
                <div class="stat-box"><div class="stat-label">Hostname</div><div class="stat-value">${s.hostname}</div></div>
                <div class="stat-box"><div class="stat-label">CPU Cores</div><div class="stat-value">${s.cpus}</div></div>
                <div class="stat-box"><div class="stat-label">Go Runtime</div><div class="stat-value">${s.go_version}</div></div>
                <div class="stat-box"><div class="stat-label">Goroutines</div><div class="stat-value">${s.goroutines}</div></div>
                <div class="stat-box"><div class="stat-label">Opius Storage</div><div class="stat-value">${this.fmt(s.opius_size)}</div></div>
                <div class="stat-box" style="grid-column: 1/-1;"><div class="stat-label">Last Update</div><div class="stat-value">${new Date(s.time).toLocaleTimeString()}</div></div>
            `;
        } catch (e) {
            this.container.querySelector('#sys-content').innerHTML = '<div class="empty"><span class="badge badge-err">Error: ' + e.message + '</span></div>';
        }
    }

    fmt(b) {
        if (b >= 1024 * 1024) return (b / 1024 / 1024).toFixed(1) + ' MB';
        if (b >= 1024) return (b / 1024).toFixed(1) + ' KB';
        return b + ' B';
    }

    destroy() {
        if (this.interval) clearInterval(this.interval);
    }
}

class SettingsApp {
    render() {
        return `<div class="app-content">
            <h2>Settings</h2>
            <p style="color: #9ca3af; margin-bottom: 16px;">Settings are managed via the CLI. The configuration file is located at:</p>
            <pre>~/.config/opius/opius.toml</pre>
            <h2 style="margin-top: 32px;">Common Commands</h2>
            <pre>opius config set package_backend macports
opius config set registry_url https://github.com/am1s3/opius-os-pkg
opius pkg update
opius doctor</pre>
        </div>`;
    }
    mount() {}
}

class AboutApp {
    render() {
        return `<div class="app-content">
            <div class="about-hero">
                <div style="width: 80px; height: 80px; margin: 0 auto 24px; color: #3b82f6;">
                    ${ICONS.terminal}
                </div>
                <h1>Opius OS</h1>
                <p class="tagline">Terminal-first embedded workspace for microcontrollers</p>
                <div class="version">Version 1.0.0</div>
                <div class="about-links">
                    <a class="about-link" href="https://github.com/am1s3/opius-os" target="_blank">
                        ${ICONS.github} Main Repository
                    </a>
                    <a class="about-link" href="https://github.com/am1s3/opius-os-pkg" target="_blank">
                        ${ICONS.packages} Package Registry
                    </a>
                </div>
                <p style="color: #9ca3af; margin-top: 40px; font-size: 13px;">Built with Go, xterm.js, and ❤️ for embedded hackers</p>
            </div>
        </div>`;
    }
    mount() {}
}

window.apps = new AppRegistry();