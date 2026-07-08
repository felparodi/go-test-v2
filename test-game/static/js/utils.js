// static/js/utils.js
export default class Utils {
    // Generar ID único
    static generateId() {
        return 'player_' + Math.random().toString(36).substr(2, 9);
    }
    
    // Calcular distancia entre dos puntos
    static distance(x1, y1, x2, y2) {
        const dx = x2 - x1;
        const dy = y2 - y1;
        return Math.sqrt(dx * dx + dy * dy);
    }
    
    // Normalizar vector
    static normalize(x, y) {
        const length = Math.sqrt(x * x + y * y);
        if (length === 0) return { x: 0, y: 0 };
        return { x: x / length, y: y / length };
    }
    
    // Limitar valor
    static clamp(value, min, max) {
        return Math.min(Math.max(value, min), max);
    }
    
    // Interpolación lineal
    static lerp(a, b, t) {
        return a + (b - a) * t;
    }
    
    // Convertir grados a radianes
    static toRadians(degrees) {
        return degrees * Math.PI / 180;
    }
    
    // Convertir radianes a grados
    static toDegrees(radians) {
        return radians * 180 / Math.PI;
    }

    static getIdToColor(str, saturación = 80, luminosidad = 50) {
        // Función hash para strings
        const hashString = (s) => {
            let hash = 0;
            for (let i = 0; i < s.length; i++) {
                const char = s.charCodeAt(i);
                hash = ((hash << 5) - hash) + char;
                hash = hash & hash; // Convertir a 32-bit integer
            }
            return Math.abs(hash);
        };
        
        // Generar hue basado en el hash
        const hash = hashString(str);
        const hue = (hash % 360);
        
        // Convertir HSL a Hex
        const hslToHex = (h, s, l) => {
            h = h / 360;
            s = s / 100;
            l = l / 100;
            
            let r, g, b;
            if (s === 0) {
                r = g = b = l;
            } else {
                const hue2rgb = (p, q, t) => {
                    if (t < 0) t += 1;
                    if (t > 1) t -= 1;
                    if (t < 1/6) return p + (q - p) * 6 * t;
                    if (t < 1/2) return q;
                    if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
                    return p;
                };
                const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
                const p = 2 * l - q;
                r = hue2rgb(p, q, h + 1/3);
                g = hue2rgb(p, q, h);
                b = hue2rgb(p, q, h - 1/3);
            }
            
            const toHex = (x) => {
                const hex = Math.round(x * 255).toString(16);
                return hex.length === 1 ? '0' + hex : hex;
            };
            
            return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
        };
    
        return hslToHex(hue, saturación, luminosidad);
    }
}