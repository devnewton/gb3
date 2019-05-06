class Gc2Norloge extends HTMLElement {
    constructor() {
        super();

        this.onmouseenter = (e) => {
            let norloges = document.querySelectorAll(`gc2-norloge[title$="${this.title.substr(-8)}"]`);
            for (let n of norloges) {
                n.classList.toggle('gc2-highlighted', true);
            }
        }

        this.onmouseleave = (e) => {
            let norloges = document.querySelectorAll(`.gc2-highlighted`);
            for (let n of norloges) {
                n.classList.toggle('gc2-highlighted', false);
            }
        }
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
        let t = post.time;
        let timeText = "".concat(t.substr(8, 2), ':', t.substr(10, 2), ':', t.substr(12, 2));
        let dateText = "".concat(t.substr(0, 4), '-', t.substr(4, 2), '-', t.substr(6, 2));
        timeElement.innerText = timeText;
        timeElement.title = dateText.concat("T", timeText);
        this.appendChild(timeElement);

        let citeElement = document.createElement('gc2-moule');
        citeElement.innerText = post.login || post.info || 'coward';
        citeElement.title = post.info;
        this.appendChild(citeElement);

        let pElement = document.createElement('gc2-message');
        pElement.innerHTML = post.message;
        this.appendChild(pElement);
    }
}
customElements.define('gc2-post', Gc2Post);