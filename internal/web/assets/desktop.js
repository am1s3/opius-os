class DesktopShell {
    constructor() {
        this.initClock();
        this.initDock();
        this.initDesktopIcons();
        this.initAppsMenu();
        this.initSystemIndicator();
        this.initWindowEvents();

        // Автозапуск терминала через небольшую задержку
        setTimeout(() => apps.launch('terminal'), 400);
    }

    initClock() {
        const clock = document.getElementById('clock');
        const update = () => {
            const d = new Date();
            clock.textContent = d.toLocaleTimeString('en-US', {
                hour: '2-digit',
                minute: '2-digit',
                hour12: false,
            });
        };
        update();
        setInterval(update, 1000);
    }

    initDock() {
        document.querySelectorAll('.dock-app').forEach(el => {
            el.onclick = () => apps.launch(el.dataset.app);
        });
        this.updateDock();
    }

    updateDock() {
        document.querySelectorAll('.dock-app').forEach(el => {
            const win = wm.windows.get(el.dataset.app);
            if (win && !win.classList.contains('minimized')) {
                el.classList.add('active');
            } else if (win) {
                // Окно есть, но свёрнуто — всё равно показываем точку
                el.classList.add('active');
            } else {
                el.classList.remove('active');
            }
        });
    }

    initWindowEvents() {
        ['window-closed', 'window-focused', 'window-minimized', 'window-restored'].forEach(evt => {
            document.addEventListener(evt, () => this.updateDock());
        });
    }

    initDesktopIcons() {
        document.querySelectorAll('.desktop-icon').forEach(el => {
            el.ondblclick = () => apps.launch(el.dataset.app);
        });
    }

    initAppsMenu() {
        const btn = document.getElementById('apps-menu-btn');
        const menu = document.getElementById('apps-menu');

        btn.onclick = (e) => {
            e.stopPropagation();
            menu.classList.toggle('visible');
        };

        document.addEventListener('click', () => menu.classList.remove('visible'));

        menu.querySelectorAll('.apps-menu-item').forEach(el => {
            el.onclick = () => {
                apps.launch(el.dataset.app);
                menu.classList.remove('visible');
            };
        });
    }

    initSystemIndicator() {
        const ind = document.getElementById('system-indicator');
        const check = async () => {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();
                if (data.success) {
                    ind.classList.remove('error');
                    ind.title = 'System OK';
                } else {
                    ind.classList.add('error');
                    ind.title = 'System error';
                }
            } catch {
                ind.classList.add('error');
                ind.title = 'Disconnected';
            }
        };
        check();
        setInterval(check, 5000);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    window.shell = new DesktopShell();
});

// Toast notification system
class ToastManager {
    constructor() {
        this.container = document.createElement('div');
        this.container.style.cssText = 'position: fixed; top: 44px; right: 16px; z-index: 10000; display: flex; flex-direction: column; gap: 8px;';
        document.body.appendChild(this.container);
    }

    show(message, type = 'info', duration = 4000) {
        const toast = document.createElement('div');
        toast.style.cssText = `
            background: #1a1d23;
            border: 1px solid rgba(255,255,255,0.12);
            border-radius: 8px;
            padding: 12px 16px;
            color: #e5e7eb;
            font-size: 13px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.5);
            opacity: 0;
            transform: translateX(20px);
            transition: all 0.3s;
            max-width: 320px;
        `;

        const colors = {
            success: '#10b981',
            error: '#ef4444',
            warning: '#f59e0b',
            info: '#3b82f6',
        };

        toast.innerHTML = `<span style="display: inline-block; width: 8px; height: 8px; background: ${colors[type]}; border-radius: 50%; margin-right: 8px;"></span>${message}`;
        this.container.appendChild(toast);

        requestAnimationFrame(() => {
            toast.style.opacity = '1';
            toast.style.transform = 'translateX(0)';
        });

        setTimeout(() => {
            toast.style.opacity = '0';
            toast.style.transform = 'translateX(20px)';
            setTimeout(() => toast.remove(), 300);
        }, duration);
    }
}

window.toastManager = new ToastManager();
window.showToast = (msg, type) => window.toastManager.show(msg, type);