// static/js/renderer.js
class Renderer {
    constructor(canvasId) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext('2d');
        this.width = CONFIG.WORLD_WIDTH;
        this.height = CONFIG.WORLD_HEIGHT;
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
        const x = item.X;
        const y = item.Y;
        
        // Efecto de brillo
        const gradient = ctx.createRadialGradient(x, y, 2, x, y, 18);
        gradient.addColorStop(0, '#4CAF50');
        gradient.addColorStop(0.5, '#66BB6A');
        gradient.addColorStop(1, 'rgba(76, 175, 80, 0)');
        ctx.fillStyle = gradient;
        ctx.beginPath();
        ctx.arc(x, y, 18, 0, Math.PI * 2);
        ctx.fill();
        
        // Estrella
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
    
    drawEye(ctx, x, y, radius, pupilRadius, color) {
        // Fondo del ojo
        ctx.shadowBlur = 0;
        ctx.beginPath();
        ctx.arc(x, y, radius, 0, Math.PI * 2);
        ctx.fillStyle = 'white';
        ctx.fill();
        ctx.strokeStyle = '#333';
        ctx.lineWidth = 0.5;
        ctx.stroke();
        
        // Pupila
        const pupilOffsetX = radius * 0.3;
        const pupilOffsetY = -radius * 0.2;
        ctx.beginPath();
        ctx.arc(x + pupilOffsetX, y + pupilOffsetY, pupilRadius, 0, Math.PI * 2);
        ctx.fillStyle = color || '#2C3E50';
        ctx.fill();
        
        // Brillo
        ctx.beginPath();
        ctx.arc(
            x + pupilOffsetX + pupilRadius * 0.3,
            y + pupilOffsetY - pupilRadius * 0.3,
            pupilRadius * 0.3,
            0,
            Math.PI * 2
        );
        ctx.fillStyle = 'rgba(255,255,255,0.8)';
        ctx.fill();
    }
    
    drawPlayer(player, playerId, isLocal) {
        const ctx = this.ctx;
        const size = isLocal ? CONFIG.PLAYER_SIZE_LOCAL : CONFIG.PLAYER_SIZE_OTHER;
        const x = player.x;
        const y = player.y;
        
        // Calcular ángulo
        let angle = 0;
        if (player.vx !== undefined && player.vy !== undefined) {
            if (Math.abs(player.vx) > 0.5 || Math.abs(player.vy) > 0.5) {
                angle = Math.atan2(player.vy, player.vx);
                player.lastAngle = angle;
            } else {
                angle = player.lastAngle || 0;
            }
        }
        
        ctx.save();
        ctx.translate(x, y);
        ctx.rotate(angle);
        
        // Sombra
        ctx.shadowColor = 'rgba(0,0,0,0.5)';
        ctx.shadowBlur = 15;
        ctx.shadowOffsetX = 3;
        ctx.shadowOffsetY = 3;
        
        // Triángulo
        ctx.beginPath();
        ctx.moveTo(0, -size);
        ctx.lineTo(-size * 0.8, size * 0.7);
        ctx.lineTo(size * 0.8, size * 0.7);
        ctx.closePath();
        
        const gradient = ctx.createLinearGradient(0, -size, 0, size);
        if (isLocal) {
            gradient.addColorStop(0, '#64B5F6');
            gradient.addColorStop(0.5, '#1E88E5');
            gradient.addColorStop(1, '#0D47A1');
        } else {
            gradient.addColorStop(0, '#FF8A65');
            gradient.addColorStop(0.5, '#F4511E');
            gradient.addColorStop(1, '#BF360C');
        }
        
        ctx.fillStyle = gradient;
        ctx.fill();
        
        ctx.shadowBlur = 0;
        ctx.strokeStyle = isLocal ? '#1565C0' : '#4A148C';
        ctx.lineWidth = 2;
        ctx.stroke();
        
        // Ojos
        const eyeOffsetX = size * 0.3;
        const eyeOffsetY = -size * 0.1;
        const eyeRadius = size * CONFIG.EYE_RADIUS_RATIO;
        const pupilRadius = eyeRadius * CONFIG.PUPIL_RADIUS_RATIO;
        const eyeColor = isLocal ? '#1A237E' : '#4A148C';
        
        this.drawEye(ctx, -eyeOffsetX, eyeOffsetY, eyeRadius, pupilRadius, eyeColor);
        this.drawEye(ctx, eyeOffsetX, eyeOffsetY, eyeRadius, pupilRadius, eyeColor);
        
        ctx.restore();
        
        // Nombre y puntuación
        ctx.save();
        ctx.shadowBlur = 0;
        ctx.fillStyle = 'white';
        ctx.font = 'bold 11px "Segoe UI", Arial, sans-serif';
        ctx.textAlign = 'center';
        ctx.fillText('👤 ' + playerId.substring(0, 8), x, y - size - 18);
        ctx.fillStyle = '#FFD54F';
        ctx.font = '11px "Segoe UI", Arial, sans-serif';
        ctx.fillText('⭐ ' + (player.score || 0), x, y + size + 22);
        ctx.restore();
    }
    
    render(gameState, playerId) {
        this.clear();
        this.drawGrid();
        
        // Dibujar items
        gameState.items.forEach(item => this.drawItem(item));
        
        // Dibujar jugadores
        Object.entries(gameState.players).forEach(([id, player]) => {
            const isLocal = id === playerId;
            this.drawPlayer(player, id, isLocal);
        });
    }
}