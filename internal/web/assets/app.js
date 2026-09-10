class OpiusWebUI {
    constructor() {
        this.devices = [];
        this.packages = [];
        this.init();
    }

    async init() {
        await this.loadDevices();
        await this.loadPackages();
        this.setupEventListeners();
    }

    async loadDevices() {
        try {
            const response = await fetch('/api/devices');
            const data = await response.json();

            if (data.success) {
                this.devices = data.data || [];
                this.renderDevices();
                this.updatePortSelect();
            } else {
                this.showError('Failed to load devices: ' + data.error);
            }
        } catch (err) {
            this.showError('Failed to load devices: ' + err.message);
        }
    }

    async loadPackages() {
        try {
            const response = await fetch('/api/packages');
            const data = await response.json();

            if (data.success) {
                this.packages = data.data || [];
                this.renderPackages();
                this.updatePackageSelect();
            } else {
                this.showError('Failed to load packages: ' + data.error);
            }
        } catch (err) {
            this.showError('Failed to load packages: ' + err.message);
        }
    }

    renderDevices() {
        const container = document.getElementById('devices-list');

        if (this.devices.length === 0) {
            container.innerHTML = '<div class="empty">No devices connected</div>';
            return;
        }

        container.innerHTML = this.devices.map(d => `
            <div class="device-item">
                <div>
                    <span class="device-port">${d.port}</span>
                    <span class="device-status status-${d.status}">${d.status.toUpperCase()}</span>
                </div>
                <div>
                    <span class="device-chip">${d.chip}</span> · ${d.board}
                </div>
                ${d.hint ? `<div style="color: #888; font-size: 0.9em; margin-top: 5px;">${d.hint}</div>` : ''}
            </div>
        `).join('');
    }

    renderPackages() {
        const container = document.getElementById('packages-list');

        if (this.packages.length === 0) {
            container.innerHTML = '<div class="empty">No packages in registry</div>';
            return;
        }

        container.innerHTML = this.packages.map(p => `
            <div class="package-item">
                <div>
                    <span class="package-name">${p.name}</span>
                    <span class="package-version">v${p.version}</span>
                </div>
                <div class="package-description">${p.description || 'No description'}</div>
                <div style="color: #888; font-size: 0.85em; margin-top: 5px;">
                    Chips: ${(p.chips || []).join(', ')}
                </div>
            </div>
        `).join('');
    }

    updatePortSelect() {
        const select = document.getElementById('port-select');
        select.innerHTML = '<option value="">Select a device...</option>';

        this.devices.forEach(d => {
            const option = document.createElement('option');
            option.value = d.port;
            option.textContent = `${d.port} (${d.chip})`;
            select.appendChild(option);
        });
    }

    updatePackageSelect() {
        const select = document.getElementById('package-select');
        select.innerHTML = '<option value="">Select a package...</option>';

        this.packages.forEach(p => {
            const option = document.createElement('option');
            option.value = p.name;
            option.textContent = `${p.name} v${p.version}`;
            select.appendChild(option);
        });
    }

    setupEventListeners() {
        document.getElementById('flash-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            await this.handleFlash();
        });
    }

    async handleFlash() {
        const packageName = document.getElementById('package-select').value;
        const port = document.getElementById('port-select').value;
        const chip = document.getElementById('chip-input').value;
        const erase = document.getElementById('erase-check').checked;

        if (!packageName || !port) {
            this.showError('Please select both package and device');
            return;
        }

        const resultDiv = document.getElementById('flash-result');
        resultDiv.className = 'result';
        resultDiv.style.display = 'block';
        resultDiv.textContent = 'Flashing...';

        try {
            const response = await fetch('/api/flash', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    package_name: packageName,
                    port: port,
                    chip: chip,
                    erase: erase,
                }),
            });

            const data = await response.json();

            if (data.success) {
                resultDiv.className = 'result success';
                resultDiv.textContent = data.data.message;
            } else {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + data.error;
            }
        } catch (err) {
            resultDiv.className = 'result error';
            resultDiv.textContent = 'Error: ' + err.message;
        }
    }

    showError(message) {
        console.error(message);
        // Could show a toast notification here
    }
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new OpiusWebUI();
});