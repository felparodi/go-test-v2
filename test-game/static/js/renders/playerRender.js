export default class PlayerRender {
    debugLineV1(ctx, {x, y, vx, vy, angle}) {
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

    debugLineV2(ctx, {vx, vy, angle}) {
        ctx.save();
        ctx.strokeStyle = 'rgba(255, 255, 0, 0.6)';
        ctx.lineWidth = 3;
        ctx.setLineDash([8, 4]);
        ctx.beginPath();
        const dirLength = 60;
        const endX = x + Math.cos(0) * dirLength;
        const endY = y + Math.sin(0) * dirLength;
        ctx.lineTo(endX, endY);
        ctx.stroke();
        ctx.setLineDash([]);
        
        // Flecha en la punta
        const arrowSize = 10;
        const arrowAngle = 0.5;
        const endX2 = x + Math.cos(0) * (dirLength - 5);
        const endY2 = y + Math.sin(0) * (dirLength - 5);
        ctx.beginPath();
        ctx.moveTo(endX, endY);
        ctx.lineTo(
            endX2 - Math.cos(0 - arrowAngle) * arrowSize,
            endY2 - Math.sin(0 - arrowAngle) * arrowSize
        );
        ctx.moveTo(endX, endY);
        ctx.lineTo(
            endX2 - Math.cos(0 + arrowAngle) * arrowSize,
            endY2 - Math.sin(0 + arrowAngle) * arrowSize
        );
        ctx.stroke();
        
        // Mostrar ángulo
        ctx.setLineDash([]);
        ctx.fillStyle = 'rgba(255, 255, 0, 0.8)';
        ctx.font = '12px monospace';
        ctx.textAlign = 'left';
        const angleDeg = (angle * 180 / Math.PI).toFixed(0);
        ctx.fillText(`Ángulo: ${angleDeg}°`, 10, 20);
        ctx.fillText(`vx: ${vx.toFixed(0)}`, 10, 5);
        ctx.fillText(`vy: ${vy.toFixed(0)}`, 10, 10);
        
        ctx.restore();
    }

}