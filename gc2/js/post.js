class Gc2Norloge extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gc2-norloge', Gc2Norloge);

class Gc2Moule extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gc2-moule', Gc2Moule);

class Gc2Message extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gc2-message', Gc2Message);

class Gc2Post extends HTMLElement {

    constructor() {
        super();
    }

    setup(post) {
        this.id = "".concat(post.id, '@', post.tribune);
        this.postId = post.id;
        this.tribune = post.tribune;

        let timeElement = document.createElement('gc2-norloge');
        timeElement.innerText = this.formatTime(post.time);
        this.appendChild(timeElement);

        let citeElement = document.createElement('gc2-moule');
        citeElement.innerText = post.login || post.info || 'coward';
        citeElement.title = post.info;
        this.appendChild(citeElement);

        let pElement = document.createElement('gc2-message');
        pElement.innerHTML = post.message;
        this.appendChild(pElement);
    }

    formatTime(t) {
        return "".concat(t.substr(8, 2), ':',t.substr(10, 2), ':', t.substr(12, 2));
    }

}
customElements.define('gc2-post', Gc2Post);