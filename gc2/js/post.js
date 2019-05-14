class Gc2PostNorloge extends HTMLElement {
    constructor() {
        super();

        this.onmouseenter = (e) => {
            let norloges = document.querySelectorAll(`gc2-norloge[title$="${this.title.substr(-8)}"],gc2-post-norloge[title$="${this.title.substr(-8)}"]`);
            for (let n of norloges) {
                n.classList.toggle("gc2-highlighted", true);
            }
        }

        this.onmouseleave = (e) => {
            let norloges = document.querySelectorAll(".gc2-highlighted");
            for (let n of norloges) {
                n.classList.toggle("gc2-highlighted", false);
            }
        }

        this.onclick = (e) => {
            let message = document.getElementById("gc2-message");
            message.value += `${message.value && ' '}${this.title} `;
            message.focus();
        }
    }
}
customElements.define('gc2-post-norloge', Gc2PostNorloge);

class Gc2Norloge extends HTMLElement {
    constructor() {
        super();

        this.onmouseenter = (e) => {
            let norloges = document.querySelectorAll(`gc2-norloge[title$="${this.title.substr(-8)}"],gc2-post-norloge[title$="${this.title.substr(-8)}"]`);
            for (let n of norloges) {
                n.classList.toggle("gc2-highlighted", true);
            }
        }

        this.onmouseleave = (e) => {
            let norloges = document.querySelectorAll(".gc2-highlighted");
            for (let n of norloges) {
                n.classList.toggle("gc2-highlighted", false);
            }
        }

        this.onclick = (e) => {
            let n = this.findQuotedNorloge();
            if(n) {
                n.scrollIntoView();
            }
        }
    }

    findQuotedNorloge() {
        let norloges = document.querySelectorAll("gc2-norloge");
        let time = this.title.substr(-8);
        for (let n of norloges) {
            if(!this.isSameNode(n) && time === n.title.substr(-8)) {
                return n
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

        let timeElement = document.createElement('gc2-post-norloge');
        let t = post.time;
        let dateText = `${t.substr(0, 4)}-${t.substr(4, 2)}-${t.substr(6, 2)}`;
        let timeText = `${t.substr(8, 2)}:${t.substr(10, 2)}:${t.substr(12, 2)}`;
        timeElement.innerText = timeText;
        timeElement.title = `${dateText} ${timeText}`;
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