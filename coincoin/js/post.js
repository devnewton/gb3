class Gb3Post extends HTMLElement {

    constructor() {
        super();
    }

    setup(post) {
        this.id = "".concat(post.id, '@', post.tribune);
        this.postId = post.id;
        this.tribune = post.tribune;

        let timeElement = document.createElement('time');
        timeElement.innerText = this.formatTime(post.time);
        this.appendChild(timeElement);

        let citeElement = document.createElement('cite');
        citeElement.innerText = post.login || post.info || 'coward';
        citeElement.title = post.info;
        this.appendChild(citeElement);

        let pElement = document.createElement('p');
        pElement.innerText = post.message;
        this.appendChild(pElement);
    }

    formatTime(t) {
        return "".concat(t.substr(8, 2), ':',t.substr(10, 2), ':', t.substr(12, 2));
    }

}
customElements.define('gb3-post', Gb3Post);