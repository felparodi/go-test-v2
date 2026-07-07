class Game {
    constructor() {
        this.playerId = Utils.generateId();
        this.gameState = {
            players: {},
            items: []
        };
        this.fps = 0;
        this.frameCount = 0;
        this.lastFpsUpdate = Date.now();
        this.lastMoveSend = 0;
        
        // Inicializar componentes
        this.network = new NetworkManager(CONFIG.WS_URL);
        this.input = new InputManager();
        this.renderer = new Renderer('gameCanvas');
        
        // Configurar callbacks
        this.setupNetworkHandlers();
        this.setupInputHandlers();
        
        // Iniciar juego
        this.init();
    }
    
    async init() {
        try {
            await this.network.connect();
            this.network.sendInit(this.playerId);
            
            // Iniciar bucle de juego
            this.gameLoop();
            
            // Actualizar estado periódicamente
            setInterval(() => this.updateStats(), CONFIG.UPDATE_INTERVAL);
            
        } catch (error) {
            console.error('Error al iniciar el juego:', error);
        }
    }
    
    setupNetworkHandlers() {
        this.network.onMessage((data) => {
            this.gameState = data;
            this.renderer.render(this.gameState, this.playerId);
            this.updateStats();
            this.updateFPS();
        });
    }
    
    setupInputHandlers() {
        this.input.onMove(() => {
            this.sendMovement();
        });
    }
    
    sendMovement() {
        const now = Date.now();
        if (now - this.lastMoveSend < CONFIG.MOVE_SEND_INTERVAL) return;
        
        const velocity = this.input.getVelocity(CONFIG.BASE_SPEED);
        
        // Asegurar que enviamos valores numéricos
        const vx = velocity.vx || 0;
        const vy = velocity.vy || 0;
        
        this.network.sendMove(vx, vy);
        this.lastMoveSend = now;
    }
    
    gameLoop() {
        // Envío periódico de movimiento
        setInterval(() => {
            if (this.input.isMoving()) {
                this.sendMovement();
            }
        }, 100);
    }
    
    updateStats() {
        const player = this.gameState.players[this.playerId];
        if (player) {
            document.getElementById('score').textContent = player.score;
        }
        
        const playerCount = Object.keys(this.gameState.players).length;
        document.getElementById('players').textContent = playerCount;
    }
    
    updateFPS() {
        this.frameCount++;
        const now = Date.now();
        if (now - this.lastFpsUpdate >= CONFIG.UPDATE_INTERVAL) {
            this.fps = this.frameCount;
            this.frameCount = 0;
            this.lastFpsUpdate = now;
            document.getElementById('fps').textContent = this.fps;
        }
    }
}