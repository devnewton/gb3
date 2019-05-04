class Gc2Controls extends HTMLElement {
    constructor() {
        super();
    }
}
customElements.define('gc2-controls', Gc2Controls);

class Gc2Main extends HTMLElement {

    constructor() {
        super();
        this.tribunes = new Map();
    }

    connectedCallback() {
        this.setupControls();
        this.setupTribuneSelect();
        this.setupMessageInput();
        this.setupBackend2html();
    }

    setupControls() {
        this.controls = document.createElement("gc2-controls");
        this.appendChild(this.controls);
    }

    setupTribuneSelect() {
        this.tribuneSelect = document.createElement("select");
        this.tribuneSelect.onchange = () => {
            let selectedTribune = this.tribuneSelect.value;
            this.tribunes.forEach((tribuneElement, tribune) => {
                tribuneElement.style.display = tribune === selectedTribune ? "" : "none";
            });
        }
        this.controls.appendChild(this.tribuneSelect);
    }

    setupMessageInput() {
        this.messageInput = document.createElement("input");
        this.messageInput.type = "text";
        this.messageInput.spellcheck = true;
        this.controls.appendChild(this.messageInput);
    }

    setupBackend2html() {
        fetch("/peg/backend2html.pegjs")
            .then((response) => {
                return response.text();
            })
            .then((text) => {
                this.backend2html = peg.generate(text);
                this.startPoll();
            });
    }

    startPoll() {
        let postSource = new EventSource("/api/poll");
        postSource.onmessage = (event) => {
            let post = JSON.parse(event.data);
            post.message = this.backend2html.parse(post.message);
            this.getTribuneElement(post.tribune).insertPost(post);
        };
    }

    getTribuneElement(tribune) {
        let tribuneElement = this.tribunes.get(tribune);
        if (!tribuneElement) {
            var option = document.createElement("option");
            option.text = option.value = tribune;
            this.tribuneSelect.add(option);

            tribuneElement = document.createElement('gc2-tribune');
            this.tribunes.set(tribune, tribuneElement);
            tribuneElement.style.display = tribune === this.tribuneSelect.value ? "" : "none";
            this.appendChild(tribuneElement);
        }
        return tribuneElement;
    }

}
customElements.define('gc2-main', Gc2Main);