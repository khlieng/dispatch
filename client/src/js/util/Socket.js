import Backoff from 'backo';

export default class Socket {
  constructor(host) {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    this.url = `${protocol}://${host}/ws?path=${window.location.pathname}`;

    this.connectTimeout = 20000;
    this.pingTimeout = 30000;
    this.backoff = new Backoff({
      min: 1000,
      max: 5000,
      jitter: 0.25
    });
    this.handlers = [];

    this.connect();
  }

  connect() {
    this.ws = new WebSocket(this.url);

    this.timeoutConnect = setTimeout(() => {
      this.ws.close();
      this.retry();
    }, this.connectTimeout);

    this.ws.onopen = () => {
      clearTimeout(this.timeoutConnect);
      this.backoff.reset();
      this.setTimeoutPing();
    };

    this.ws.onclose = () => {
      clearTimeout(this.timeoutConnect);
      clearTimeout(this.timeoutPing);
      if (!this.closing) {
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
      this.closing = true;
      this.ws.close();
      this.connect();
    }, this.pingTimeout);
  }

  onMessage(handler) {
    this.handlers.push(handler);
  }

  emit(type, data) {
    for (let i = 0; i < this.handlers.length; i++) {
      this.handlers[i](type, data);
    }
  }
}
