// static/js/network.js
class NetworkManager {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.isConnected = false;
        this.messageHandlers = [];
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 2000;
    }
    
    connect() {
        return new Promise((resolve, reject) => {
            try {
                this.ws = new WebSocket(this.url);
                this.ws.onopen = () => {
                    this.isConnected = true;
                    this.reconnectAttempts = 0;
                    console.log('✅ Conectado al servidor');
                    this.updateStatus('connected');
                    resolve();
                };
                
                this.ws.onmessage = (event) => {
                    try {
                        const data = JSON.parse(event.data);
                        this.handleMessage(data);
                    } catch (e) {
                        console.error('Error al procesar mensaje:', e);
                    }
                };
                
                this.ws.onclose = () => {
                    this.isConnected = false;
                    console.log('❌ Desconectado del servidor');
                    this.updateStatus('disconnected');
                    this.attemptReconnect();
                };
                
                this.ws.onerror = (error) => {
                    console.error('Error WebSocket:', error);
                    this.updateStatus('error');
                    reject(error);
                };
            } catch (error) {
                console.error('Error al conectar:', error);
                reject(error);
            }
        });
    }
    
    attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Intentando reconectar (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
            this.updateStatus('connecting');
            setTimeout(() => {
                this.connect().catch(() => {});
            }, this.reconnectDelay);
        } else {
            console.error('❌ No se pudo reconectar después de varios intentos');
            this.updateStatus('error');
        }
    }
    
    handleMessage(data) {
        this.messageHandlers.forEach(handler => handler(data));
    }
    
    onMessage(handler) {
        this.messageHandlers.push(handler);
    }
    
    send(data) {
        if (this.isConnected && this.ws) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.warn('No se puede enviar: conexión no disponible');
        }
    }
    
    sendInit(playerId) {
        this.send({
            type: 'init',
            payload: { playerId }
        });
    }
    
    sendMove(velocityX, velocityY) {
        this.send({
            type: 'move',
            payload: { velocityX, velocityY }
        });
    }
    
    sendAction(action) {
        this.send({
            type: 'action',
            payload: { action }
        });
    }
    
    updateStatus(status) {
        const statusEl = document.getElementById('connectionStatus');
        if (!statusEl) return;
        
        statusEl.className = `status-${status}`;
        switch(status) {
            case 'connected':
                statusEl.textContent = '🟢 Conectado';
                break;
            case 'disconnected':
                statusEl.textContent = '🔴 Desconectado';
                break;
            case 'connecting':
                statusEl.textContent = '🟡 Conectando...';
                break;
            case 'error':
                statusEl.textContent = '⚠️ Error de conexión';
                break;
        }
    }
}