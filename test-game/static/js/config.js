// static/js/config.js
export default {
    // Servidor
    WS_URL: 'ws://localhost:8080/ws',
    
    // Juego
    WORLD_WIDTH: 800,
    WORLD_HEIGHT: 600,
    BASE_SPEED: 200,
    MAX_SPEED: 250,
    MOVE_SEND_INTERVAL: 50,
    
    // Jugador

 
    
    // FPS
    TARGET_FPS: 60,
    UPDATE_INTERVAL: 1000,
    
    // Items
    ITEMS_COUNT: 20,
    ITEM_COLLECT_RADIUS: 30,
    ITEM_SPIKES: 6,
    ITEM_OUTER_RADIUS: 10,
    ITEM_INNER_RADIUS: 4,
    
    // Debug - Muestra líneas de dirección y ángulos
    DEBUG_MODE: true  // Cambiar a false para desactivar
};