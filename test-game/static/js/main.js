// static/js/main.js
// Punto de entrada principal
document.addEventListener('DOMContentLoaded', () => {
    console.log('🎮 Iniciando juego...');
    console.log(`📱 ID del jugador: ${Utils.generateId()}`);
    console.log(`🌐 Servidor: ${CONFIG.WS_URL}`);
    
    try {
        const game = new Game();
        console.log('✅ Juego inicializado correctamente');
        
        // Exponer game para debugging
        window.game = game;
    } catch (error) {
        console.error('❌ Error al inicializar el juego:', error);
        document.getElementById('connectionStatus').textContent = '❌ Error de inicio';
    }
});