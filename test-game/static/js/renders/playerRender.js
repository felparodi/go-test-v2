const EYE_RADIUS_RATIO = 0.2;
const PUPIL_RADIUS_RATIO = 0.5;
const PLAYER_SIZE_LOCAL = 25;
const PLAYER_SIZE_OTHER = 22;
export default class PlayerRender {
    static drawDebugLine(ctx, { x, y, vx=0, vy=0, angle }) {
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

    static drawBody(ctx, {x , y, angle}, isLocal ) {
        const size = isLocal ? PLAYER_SIZE_LOCAL : PLAYER_SIZE_OTHER;
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
        const eyeRadius = size * EYE_RADIUS_RATIO;
        const pupilRadius = eyeRadius * PUPIL_RADIUS_RATIO;
        const eyeColor = isLocal ? '#1A237E' : '#4A148C';
        
        // Ojo superior
        this.drawEyes(ctx, { x: eyeOffsetX, y:-eyeOffsetY, radius:eyeRadius, pupilRadius, color:eyeColor });
        
        // Ojo inferior
        this.drawEyes(ctx, { x: eyeOffsetX, y: eyeOffsetY, radius:eyeRadius, pupilRadius, color:eyeColor });
        
        ctx.restore();
    }

    static drawEyes(ctx, { x, y, radius, pupilRadius, color }) {
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

    static drawPlayerInfo(ctx, { x, y, size, player, playerId }) {
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
        ctx.fillText('⭐ NO' + (player.score || 0), x, y + size + 28);
        
        ctx.restore();
    }
}