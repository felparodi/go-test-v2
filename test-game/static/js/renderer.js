// static/js/renderer.js
class Renderer {
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
        const x = item.X;
        const y = item.Y;
        
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
        
        // Pupila - mirando hacia adelante
        const pupilOffsetX = radius * 0.4;
        const pupilOffsetY = -radius * 0.3;
        ctx.beginPath();
        ctx.arc(x + pupilOffsetX, y + pupilOffsetY, pupilRadius, 0, Math.PI * 2);
        ctx.fillStyle = color || '#2C3E50';
        ctx.fill();
        
        // Brillo en la pupila
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
        
        // --- CALCULAR ÁNGULO DE DIRECCIÓN ---
        let angle = 0;
        let isMoving = false;
        let {vx, vy} = player;
        vx = vx ? vx : 0;
        vy = vy ? vy : 0;
        if (vx !== undefined && vy !== undefined) {
            const speed = Math.sqrt(vx * vx + vy * vy);
            
            if (speed > 5) { // Umbral de movimiento
                isMoving = true;
                // Calcular ángulo en radianes
                angle = Math.atan2(vy, vx);
                // Guardar la última dirección
                player.lastAngle = angle;
            } else {
                // Usar la última dirección guardada
                angle = player.lastAngle || 0;
            }
        }
        
        // --- DIBUJAR LÍNEA DE DIRECCIÓN (DEBUG) ---
        if (this.debugMode && isLocal) {
            ctx.save();
            ctx.strokeStyle = 'rgba(255, 255, 0, 0.6)';
            ctx.lineWidth = 3;
            ctx.setLineDash([8, 4]);
            ctx.beginPath();
            ctx.moveTo(x, y);
            const dirLength = 60;
            const endX = x + Math.cos(angle) * dirLength;
            const endY = y + Math.sin(angle) * dirLength;
            ctx.lineTo(endX, endY);
            ctx.stroke();
            ctx.setLineDash([]);
            
            // Flecha en la punta
            const arrowSize = 10;
            const arrowAngle = 0.5;
            const endX2 = x + Math.cos(angle) * (dirLength - 5);
            const endY2 = y + Math.sin(angle) * (dirLength - 5);
            ctx.beginPath();
            ctx.moveTo(endX, endY);
            ctx.lineTo(
                endX2 - Math.cos(angle - arrowAngle) * arrowSize,
                endY2 - Math.sin(angle - arrowAngle) * arrowSize
            );
            ctx.moveTo(endX, endY);
            ctx.lineTo(
                endX2 - Math.cos(angle + arrowAngle) * arrowSize,
                endY2 - Math.sin(angle + arrowAngle) * arrowSize
            );
            ctx.stroke();
            
            // Mostrar ángulo
            ctx.setLineDash([]);
            ctx.fillStyle = 'rgba(255, 255, 0, 0.8)';
            ctx.font = '12px monospace';
            ctx.textAlign = 'left';
            const angleDeg = (angle * 180 / Math.PI).toFixed(0);
            ctx.fillText(`Ángulo: ${angleDeg}°`, x + 10, y - 20);
            ctx.fillText(`vx: ${vx.toFixed(0)}`, x + 10, y - 5);
            ctx.fillText(`vy: ${vy.toFixed(0)}`, x + 10, y + 10);
            
            ctx.restore();
        }
        
        // --- DIBUJAR EL TRIÁNGULO ---
        ctx.save();
        ctx.translate(x, y);
        
        // Aplicar rotación
        // El triángulo base apunta hacia ARRIBA (ángulo 0 = arriba)
        // En canvas, el ángulo 0 apunta a la DERECHA
        // Rotamos para que la punta del triángulo apunte en la dirección del movimiento
        ctx.rotate(angle);
        
        // Sombra
        ctx.shadowColor = 'rgba(0,0,0,0.5)';
        ctx.shadowBlur = 15;
        ctx.shadowOffsetX = 3;
        ctx.shadowOffsetY = 3;
        
        // --- DIBUJAR TRIÁNGULO (apuntando hacia la DERECHA) ---
        // Cambiamos el triángulo para que apunte a la DERECHA por defecto
        // Esto hace que la rotación sea más intuitiva con atan2
        ctx.beginPath();
        ctx.moveTo(size, 0);              // Punta (derecha)
        ctx.lineTo(-size * 0.6, -size * 0.7); // Superior izquierda
        ctx.lineTo(-size * 0.6, size * 0.7);  // Inferior izquierda
        ctx.closePath();
        
        // Gradiente de color
        const gradient = ctx.createLinearGradient(-size, 0, size, 0);
        if (isLocal) {
            gradient.addColorStop(0, '#0D47A1');
            gradient.addColorStop(0.5, '#1E88E5');
            gradient.addColorStop(1, '#64B5F6');
        } else {
            gradient.addColorStop(0, '#BF360C');
            gradient.addColorStop(0.5, '#F4511E');
            gradient.addColorStop(1, '#FF8A65');
        }
        
        ctx.fillStyle = gradient;
        ctx.fill();
        
        // Borde
        ctx.shadowBlur = 0;
        ctx.strokeStyle = isLocal ? '#1565C0' : '#4A148C';
        ctx.lineWidth = 2;
        ctx.stroke();
        
        // --- DIBUJAR OJOS EN LA PUNTA ---
        // Los ojos se colocan cerca de la punta (lado derecho)
        const eyeOffsetX = size * 0.3;
        const eyeOffsetY = size * 0.2;
        const eyeRadius = size * CONFIG.EYE_RADIUS_RATIO;
        const pupilRadius = eyeRadius * CONFIG.PUPIL_RADIUS_RATIO;
        const eyeColor = isLocal ? '#1A237E' : '#4A148C';
        
        // Ojo superior
        this.drawEye(ctx, eyeOffsetX, -eyeOffsetY, eyeRadius, pupilRadius, eyeColor);
        
        // Ojo inferior
        this.drawEye(ctx, eyeOffsetX, eyeOffsetY, eyeRadius, pupilRadius, eyeColor);
        
        ctx.restore();
        
        // --- NOMBRE Y PUNTUACIÓN ---
        ctx.save();
        ctx.shadowBlur = 0;
        
        // Nombre del jugador (arriba)
        ctx.fillStyle = 'white';
        ctx.font = 'bold 11px "Segoe UI", Arial, sans-serif';
        ctx.textAlign = 'center';
        const displayName = playerId.substring(0, 8);
        ctx.fillText('👤 ' + displayName, x, y - size - 18);
        
        // Puntuación (abajo)
        ctx.fillStyle = '#FFD54F';
        ctx.font = '11px "Segoe UI", Arial, sans-serif';
        ctx.fillText('⭐ ' + (player.score || 0), x, y + size + 28);
        
        ctx.restore();
    }
    
    render(gameState, playerId) {
        this.clear();
        this.drawGrid();
        
        // Dibujar items
        if (gameState.items) {
            gameState.items.forEach(item => this.drawItem(item));
        }
        
        // Dibujar jugadores
        if (gameState.players) {
            Object.entries(gameState.players).forEach(([id, player]) => {
                const isLocal = id === playerId;
                this.drawPlayer(player, id, isLocal);
            });
        }
    }
}