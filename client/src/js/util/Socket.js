import EventEmitter2 from 'eventemitter2';
import Backoff from 'backo';

export default class Socket extends EventEmitter2 {
  constructor(address) {
    super();

    this.address = address;
    this.connectTimeout = 20000;
    this.pingTimeout = 30000;
    this.backoff = new Backoff({
      min: 1000,
      max: 5000,
      jitter: 0.25
    });

    this.connect();
  }

  connect() {
    this.ws = new WebSocket(`ws://${this.address}/ws`);

    this.timeoutConnect = setTimeout(() => {
      this.ws.close();
      this.retry();
    }, this.connectTimeout);

    this.ws.onopen = () => {
      clearTimeout(this.timeoutConnect);
      this.backoff.reset();
      this.emit('connect');
      this.setTimeoutPing();
    };

    this.ws.onclose = () => {
      clearTimeout(this.timeoutConnect);
      clearTimeout(this.timeoutPing);
      if (!this.closing) {
        this.emit('disconnect');
        this.retry();
      }
      this.closing = false;
    };

    this.ws.onerror = () => {
      clearTimeout(this.timeoutConnect);
      clearTimeout(this.timeoutPing);
      this.closing = true;
      this.ws.close();
      this.retry();
    };

    this.ws.onmessage = (e) => {
      this.setTimeoutPing();

      const msg = JSON.parse(e.data);

      if (msg.type === 'ping') {
        this.send('pong');
      }

      this.emit(msg.type, msg.data);
    };
  }

  retry() {
    setTimeout(() => this.connect(), this.backoff.duration());
  }

  send(type, data) {
    this.ws.send(JSON.stringify({ type, data }));
  }

  setTimeoutPing() {
    clearTimeout(this.timeoutPing);
    this.timeoutPing = setTimeout(() => {
      this.emit('disconnect');
      this.closing = true;
      this.ws.close();
      this.connect();
    }, this.pingTimeout);
  }
}
