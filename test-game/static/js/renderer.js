import CONFIG from './config.js';
import BulletRender from './renders/bulletRender.js';
import CoinRender from './renders/coinRender.js';
import PlayerRender from './renders/playerRender.js';
export default class Renderer {
    constructor(canvasId) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext('2d');
        this.width = CONFIG.WORLD_WIDTH;
        this.height = CONFIG.WORLD_HEIGHT;
        this.debugMode = CONFIG.DEBUG_MODE;
    }
    
    clear() {
        this.ctx.clearRect(0, 0, this.width, this.height);
    }
    
    drawGrid() {
        const ctx = this.ctx;
        ctx.strokeStyle = 'rgba(255,255,255,0.03)';
        ctx.lineWidth = 1;
        
        for (let x = 0; x < this.width; x += 50) {
            ctx.beginPath();
            ctx.moveTo(x, 0);
            ctx.lineTo(x, this.height);
            ctx.stroke();
        }
        for (let y = 0; y < this.height; y += 50) {
            ctx.beginPath();
            ctx.moveTo(0, y);
            ctx.lineTo(this.width, y);
            ctx.stroke();
        }
    }
    
    drawItem(item) {
        const ctx = this.ctx;
        const { x, y } = item;
        
        const gradient = ctx.createRadialGradient(x, y, 2, x, y, 18);
        gradient.addColorStop(0, '#4CAF50');
        gradient.addColorStop(0.5, '#66BB6A');
        gradient.addColorStop(1, 'rgba(76, 175, 80, 0)');
        ctx.fillStyle = gradient;
        ctx.beginPath();
        ctx.arc(x, y, 18, 0, Math.PI * 2);
        ctx.fill();
        
        ctx.shadowBlur = 15;
        ctx.shadowColor = 'rgba(76, 175, 80, 0.5)';
        ctx.beginPath();
        const spikes = CONFIG.ITEM_SPIKES;
        const outerRadius = CONFIG.ITEM_OUTER_RADIUS;
        const innerRadius = CONFIG.ITEM_INNER_RADIUS;
        
        for (let i = 0; i < spikes * 2; i++) {
            const radius = i % 2 === 0 ? outerRadius : innerRadius;
            const angle = (i / (spikes * 2)) * Math.PI * 2 - Math.PI / 2;
            const px = x + Math.cos(angle) * radius;
            const py = y + Math.sin(angle) * radius;
            if (i === 0) ctx.moveTo(px, py);
            else ctx.lineTo(px, py);
        }
        ctx.closePath();
        ctx.fillStyle = '#66BB6A';
        ctx.fill();
        ctx.strokeStyle = '#2E7D32';
        ctx.lineWidth = 1;
        ctx.stroke();
        ctx.shadowBlur = 0;
    }
    
    drawPlayer(player, isLocal) {
        const ctx = this.ctx;
        PlayerRender.draw(ctx, player, isLocal)
    }
    
    render(gameState, playerId) {
        this.clear();
        this.drawGrid();
        
        // Dibujar items
        if (gameState.items) {
            gameState.items.forEach(item => {
                if (item) {
                    const { type } = item
                    switch(type) {
                        case "coin": CoinRender.render(this.ctx, item); break;
                        case "bullet": BulletRender.render(this.ctx, item); break;
                    }
                }  
            });
        }
        
        // Dibujar jugadores
        if (gameState.players) {
            Object.entries(gameState.players).forEach(([id, player]) => {
                const isLocal = id === playerId;
                this.drawPlayer(player, isLocal);
            });
        }
    }
}