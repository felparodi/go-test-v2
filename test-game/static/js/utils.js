// static/js/utils.js
class Utils {
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
}