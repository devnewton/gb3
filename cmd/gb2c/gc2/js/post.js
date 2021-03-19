function highlightNorloges(bouchot, norlogeText) {
    let tribunes = document.getElementsByTagName("gc2-tribune");
    for (let t of tribunes) {
        let tribuneName = t.getAttribute("name");
        let norloges = t.querySelectorAll(`gc2-norloge[title$="${norlogeText.substr(-8)}"],gc2-post-norloge[title*=" ${norlogeText.substr(-8)}"]`);
        for (let n of norloges) {
            if (tribuneName === bouchot || n.getAttribute("bouchot") === bouchot) {
                n.classList.toggle("gc2-highlighted", true);
            }
        }
    }
}

function unhighlightNorloges() {
    let norloges = document.querySelectorAll(".gc2-highlighted");
    for (let n of norloges) {
        n.classList.toggle("gc2-highlighted", false);
    }
}

class Gc2PostNorloge extends HTMLElement {
    constructor() {
        super();

        this.onmouseenter = (e) => {
            highlightNorloges(this.closest('gc2-tribune').getAttribute('name'), this.title);
        }

        this.onmouseleave = (e) => {
            unhighlightNorloges();
        }

        this.onclick = (e) => {
            let message = document.getElementById("gc2-message");
            message.value += `${message.value && ' '}${this.formatNorloge()} `;
            message.focus();
        }
    }

    formatNorloge() {
        let style = localStorage.norlogeStyle || "longlong";
        switch (style) {
            case "auto":
                let now = new Date();
                let dateLocal = new Date(now.getTime() - now.getTimezoneOffset() * 60 * 1000);
                let dateLocalStr = dateLocal.toISOString().slice(0, 10);
                if (this.title.substr(0, 10) === dateLocalStr) {
                    style = "normal";
                } else {
                    style = "longlong"
                }
                break;
            case "rand":
                let styles = ["longlong", "iso", "long", "normal", "short", "id"];
                style = styles[Math.round(Math.random() * styles.length)];
        }
        switch (style) {
            case "iso":
                return `${this.title.substr(0, 10)}T${this.title.substr(11, 8)}`;
            case "long":
                return `${this.title.substr(5, 2)}/${this.title.substr(8, 2)}#${this.title.substr(11, 8)}`;
            case "normal":
                return this.title.substr(11, 8);
            case "short":
                return this.title.substr(11, 5);
            case "id":
                return `#${this.parentElement.id.split("@", 1)[0]}`;
            case "longlong":
            default:
                return this.title;
        }
    }
}
customElements.define('gc2-post-norloge', Gc2PostNorloge);

class Gc2Norloge extends HTMLElement {
    constructor() {
        super();

        this.onmouseenter = (e) => {
            highlightNorloges(this.findBouchot(), this.title);
        }

        this.onmouseleave = (e) => {
            unhighlightNorloges();
        }

        this.onclick = (e) => {
            let n = this.findQuotedNorloge();
            if (n) {
                n.scrollIntoView();
            }
        }
    }

    findBouchot() {
        return this.getAttribute("bouchot") || this.closest('gc2-tribune').getAttribute('name');
    }

    findQuotedNorloge() {
        let bouchot = this.findBouchot();
        return document.querySelector(`gc2-tribune[name="${bouchot}"] gc2-post-norloge[title*=" ${this.title.substr(-8)}"]`);
    }
}
customElements.define('gc2-norloge', Gc2Norloge);

class Gc2Moule extends HTMLElement {
    constructor() {
        super();
        this.onclick = (e) => {
            let message = document.getElementById("gc2-message");
            message.value += `${message.value && ' '}${this.innerText}< `;
            message.focus();
        }
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

    setupWithTribune(post) {
        let tribuneElement = document.createElement('gc2-post-tribune');
        tribuneElement.innerText = post.tribune;
        this.appendChild(tribuneElement);
        this.setup(post);
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