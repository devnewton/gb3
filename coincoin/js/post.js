class Gb3Norloge extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gb3-norloge', Gb3Norloge);

class Gb3Moule extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gb3-moule', Gb3Moule);

class Gb3Message extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gb3-message', Gb3Message);

class Gb3Tribune extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gb3-tribune', Gb3Tribune);

class Gb3Post extends HTMLElement {

    constructor() {
        super();
    }

    setup(post) {
        this.id = "".concat(post.id, '@', post.tribune);
        this.postId = post.id;
        this.tribune = post.tribune;

        let timeElement = document.createElement('gb3-norloge');
        timeElement.innerText = this.formatTime(post.time);
        this.appendChild(timeElement);

        let tribuneElement = document.createElement('gb3-tribune');
        tribuneElement.innerText = post.tribune;
        this.appendChild(tribuneElement);

        let citeElement = document.createElement('gb3-moule');
        citeElement.innerText = post.login || post.info || 'coward';
        citeElement.title = post.info;
        this.appendChild(citeElement);

        let pElement = document.createElement('gb3-message');
        pElement.innerText = post.message;
        this.appendChild(pElement);
    }

    formatTime(t) {
        return "".concat(t.substr(8, 2), ':',t.substr(10, 2), ':', t.substr(12, 2));
    }

}
customElements.define('gb3-post', Gb3Post);