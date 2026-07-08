const ITEM_SPIKES = 5;
const ITEM_OUTER_RADIUS = 10;
const ITEM_INNER_RADIUS = 4;
export default class CoinRender {
    constructor() {

    }

    static render(ctx, item) {
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
        const spikes = ITEM_SPIKES;
        const outerRadius = ITEM_OUTER_RADIUS;
        const innerRadius = ITEM_INNER_RADIUS;
        
       
        
        ctx.closePath();
        ctx.fillStyle = '#66BB6A';
        ctx.fill();
        ctx.strokeStyle = '#2E7D32';
        ctx.lineWidth = 1;
        ctx.stroke();
        ctx.shadowBlur = 0;
    }
}