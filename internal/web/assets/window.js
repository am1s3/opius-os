class WindowManager {
    constructor() {
        this.windows = new Map();
        this.zIndex = 100;
        this.container = document.getElementById('windows');
    }

    open(id, title, contentHTML, opts = {}) {
        if (this.windows.has(id)) {
            const existing = this.windows.get(id);
            if (existing.classList.contains('minimized')) {
                existing.classList.remove('minimized');
            }
            this.focus(id);
            return existing;
        }

        const win = document.createElement('div');
        win.className = 'window';
        win.dataset.id = id;
        win.style.width = (opts.width || 720) + 'px';
        win.style.height = (opts.height || 500) + 'px';

        const offset = this.windows.size * 28;
        win.style.left = (Math.min(120 + offset, window.innerWidth - 400)) + 'px';
        win.style.top = (Math.min(60 + offset, window.innerHeight - 300)) + 'px';
        win.style.zIndex = ++this.zIndex;

        win.innerHTML = `
            <div class="window-titlebar">
                <div class="window-buttons">
                    <button class="window-btn close" title="Close"></button>
                    <button class="window-btn minimize" title="Minimize"></button>
                    <button class="window-btn maximize" title="Maximize"></button>
                </div>
                <div class="window-title">${title}</div>
            </div>
            <div class="window-body">${contentHTML}</div>
            <div class="window-resize"></div>
        `;

        this.container.appendChild(win);
        this.windows.set(id, win);

        this.attachDrag(win);
        this.attachResize(win);
        this.attachButtons(win, id);
        this.attachFocus(win, id);

        this.focus(id);
        return win;
    }

    close(id) {
        const win = this.windows.get(id);
        if (!win) return;

        document.dispatchEvent(new CustomEvent('window-closing', { detail: { id } }));

        win.classList.add('closing');
        setTimeout(() => {
            win.remove();
            this.windows.delete(id);
            document.dispatchEvent(new CustomEvent('window-closed', { detail: { id } }));
        }, 180);
    }

    minimize(id) {
        const win = this.windows.get(id);
        if (!win) return;
        win.classList.add('minimized');
        document.dispatchEvent(new CustomEvent('window-minimized', { detail: { id } }));
    }

    restore(id) {
        const win = this.windows.get(id);
        if (!win) return;
        win.classList.remove('minimized');
        this.focus(id);
        document.dispatchEvent(new CustomEvent('window-restored', { detail: { id } }));
    }

    focus(id) {
        const win = this.windows.get(id);
        if (!win) return;

        if (win.classList.contains('minimized')) {
            win.classList.remove('minimized');
        }

        this.windows.forEach(w => w.classList.remove('focused'));
        win.classList.add('focused');
        win.style.zIndex = ++this.zIndex;
        document.dispatchEvent(new CustomEvent('window-focused', { detail: { id } }));
    }

    attachDrag(win) {
        const titlebar = win.querySelector('.window-titlebar');
        let dragging = false, startX, startY, startLeft, startTop;

        titlebar.addEventListener('mousedown', (e) => {
            if (e.target.classList.contains('window-btn')) return;
            dragging = true;
            startX = e.clientX;
            startY = e.clientY;
            startLeft = win.offsetLeft;
            startTop = win.offsetTop;
            this.focus(win.dataset.id);
            e.preventDefault();
        });

        document.addEventListener('mousemove', (e) => {
            if (!dragging) return;
            const newLeft = Math.max(0, Math.min(window.innerWidth - 200, startLeft + e.clientX - startX));
            const newTop = Math.max(32, Math.min(window.innerHeight - 100, startTop + e.clientY - startY));
            win.style.left = newLeft + 'px';
            win.style.top = newTop + 'px';
        });

        document.addEventListener('mouseup', () => {
            dragging = false;
        });
    }

    attachResize(win) {
        const handle = win.querySelector('.window-resize');
        let resizing = false, startX, startY, startW, startH;

        handle.addEventListener('mousedown', (e) => {
            resizing = true;
            startX = e.clientX;
            startY = e.clientY;
            startW = win.offsetWidth;
            startH = win.offsetHeight;
            e.preventDefault();
        });

        document.addEventListener('mousemove', (e) => {
            if (!resizing) return;
            win.style.width = Math.max(360, startW + e.clientX - startX) + 'px';
            win.style.height = Math.max(240, startH + e.clientY - startY) + 'px';
            document.dispatchEvent(new CustomEvent('window-resized', { detail: { id: win.dataset.id } }));
        });

        document.addEventListener('mouseup', () => {
            resizing = false;
        });
    }

    attachButtons(win, id) {
        win.querySelector('.window-btn.close').addEventListener('click', () => this.close(id));
        win.querySelector('.window-btn.minimize').addEventListener('click', () => this.minimize(id));
        win.querySelector('.window-btn.maximize').addEventListener('click', () => {
            if (win.dataset.maximized === 'true') {
                win.style.left = win.dataset.prevLeft;
                win.style.top = win.dataset.prevTop;
                win.style.width = win.dataset.prevW;
                win.style.height = win.dataset.prevH;
                win.dataset.maximized = 'false';
            } else {
                win.dataset.prevLeft = win.style.left;
                win.dataset.prevTop = win.style.top;
                win.dataset.prevW = win.style.width;
                win.dataset.prevH = win.style.height;
                win.style.left = '8px';
                win.style.top = '40px';
                win.style.width = 'calc(100% - 16px)';
                win.style.height = 'calc(100vh - 40px - 76px)';
                win.dataset.maximized = 'true';
            }
            document.dispatchEvent(new CustomEvent('window-resized', { detail: { id } }));
        });
    }

    attachFocus(win, id) {
        win.addEventListener('mousedown', () => this.focus(id));
    }

    getBody(id) {
        const win = this.windows.get(id);
        return win ? win.querySelector('.window-body') : null;
    }
}

window.wm = new WindowManager();