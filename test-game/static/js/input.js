// static/js/input.js
import Utils from "./utils.js";

const MOVE_UP = ['ArrowUp', 'W', 'w']
const MOVE_DOWN = ['ArrowDown', 'S', 's']
const MOVE_RIGTH = ['ArrowRight', 'D', 'd']
const MOVE_LEFT = ['ArrowLeft', 'A', 'a']
const MOVE_BUTTON = [...MOVE_UP,...MOVE_DOWN,...MOVE_LEFT,...MOVE_RIGTH]
const ACTION_BUTTON = ['b']
export default class InputManager {
    constructor() {
        this.keys = new Set();
        this.moveCallbacks = [];
        this.setupListeners();
    }
    
    setupListeners() {
        document.addEventListener('keydown', (e) => {
            const key = e.key;
            if (!this.keys.has(key)) {
                this.keys.add(key);
            }
            if (MOVE_BUTTON.includes(key)) {
                e.preventDefault();
                this.notifyMove();
            } else if(ACTION_BUTTON.includes(key)) {
                this.notifyAction(key)
            }
        });
        
        document.addEventListener('keyup', (e) => {
            const key = e.key;
            this.keys.delete(key);
            if (MOVE_BUTTON.includes(key)) {
                e.preventDefault();
                this.notifyMove();
            }
        });
        
        // Perder foco - detener movimiento
        document.addEventListener('blur', () => {
            this.keys.clear();
            this.notifyMove();
        });
    }
    
    getDirection() {
        let x = 0, y = 0;
        
        if (MOVE_LEFT.filter(x => this.keys.has(x)).length > 0)  x -= 1;
        if (MOVE_RIGTH.filter(x => this.keys.has(x)).length > 0) x += 1;
        if (MOVE_UP.filter(x => this.keys.has(x)).length > 0) y -= 1;
        if (MOVE_DOWN.filter(x => this.keys.has(x)).length > 0) y += 1;
        
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
        return MOVE_BUTTON.filter(x => this.keys.has(x)).length > 0;
    }
    
    onInput(callback) {
        this.moveCallbacks.push(callback);
    }
    
    notifyMove() {
        this.moveCallbacks.forEach(callback => callback('move'));
    }

    notifyAction(key) {
        this.moveCallbacks.forEach(callback => callback('action', key));
    }
}