'use strict'

class SummaryView {
  constructor(el) {
    if (!el) {
      throw new Error('Root element is null')
    }
    this.root = el;
    this.initialized = false;
  }

  load(path) {
    let url = `ws://${window.location.host}${path}`
    let socket = new WebSocket(url)
    socket.addEventListener('message', this.onmessage.bind(this))
  }

  onmessage(event) {
    if (!this.initialized) {
      this.root.innerHTML = '';
      this.initialized = true;
    }

    let div = document.createElement('div');
    div.innerText = event.data
    this.root.appendChild(div)
  }
}
