// static/js/input.js
class InputManager {
    constructor() {
        this.keys = new Set();
        this.moveCallbacks = [];
        this.setupListeners();
    }
    
    setupListeners() {
        document.addEventListener('keydown', (e) => {
            const key = e.key;
            if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', ' '].includes(key)) {
                e.preventDefault();
            }
            
            if (!this.keys.has(key)) {
                this.keys.add(key);
                this.notifyMove();
            }
        });
        
        document.addEventListener('keyup', (e) => {
            const key = e.key;
            if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', ' '].includes(key)) {
                e.preventDefault();
            }
            
            this.keys.delete(key);
            this.notifyMove();
        });
        
        // Perder foco - detener movimiento
        document.addEventListener('blur', () => {
            this.keys.clear();
            this.notifyMove();
        });
    }
    
    getDirection() {
        let x = 0, y = 0;
        
        if (this.keys.has('ArrowLeft') || this.keys.has('a') || this.keys.has('A')) x -= 1;
        if (this.keys.has('ArrowRight') || this.keys.has('d') || this.keys.has('D')) x += 1;
        if (this.keys.has('ArrowUp') || this.keys.has('w') || this.keys.has('W')) y -= 1;
        if (this.keys.has('ArrowDown') || this.keys.has('s') || this.keys.has('S')) y += 1;
        
        const normalized = Utils.normalize(x, y);
        return { x: normalized.x, y: normalized.y };
    }
    
    getVelocity(speed) {
        const dir = this.getDirection();
        return {
            vx: dir.x * speed,
            vy: dir.y * speed
        };
    }
    
    isMoving() {
        return this.keys.size > 0;
    }
    
    onMove(callback) {
        this.moveCallbacks.push(callback);
    }
    
    notifyMove() {
        this.moveCallbacks.forEach(callback => callback());
    }
}