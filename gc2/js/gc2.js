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
        this.setupGesture();
    }

    setupControls() {
        this.controls = document.createElement("gc2-controls");
        this.appendChild(this.controls);
    }

    setupTribuneSelect() {
        this.tribuneSelect = document.createElement("select");
        this.tribuneSelect.onchange = () => {
            this.setActiveTribune(this.tribuneSelect.value);
            let selectedTribune = this.tribuneSelect.value;
        }
        this.controls.appendChild(this.tribuneSelect);
    }

    setupMessageInput() {
        this.messageInput = document.createElement("input");
        this.messageInput.id = "gc2-message";
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
            tribuneElement.style.display = "none";
            this.appendChild(tribuneElement);

            this.tribunes.set(tribune, tribuneElement);
            this.setActiveTribune(this.tribunes.size === 1 ? tribune : this.tribuneSelect.value);
        }
        return tribuneElement;
    }

    setActiveTribune(selectedTribune) {
        this.messageInput.placeholder = selectedTribune;
        this.tribunes.forEach((tribuneElement, tribune) => {
            tribuneElement.style.display = tribune === selectedTribune ? "" : "none";
        });
    }

    setupGesture() {
        let hammertime = new Hammer(document.getElementsByTagName("gc2-main")[0]);
        hammertime.on('swipeleft', (e) => {
            if (this.tribuneSelect.selectedIndex === 0) {
                this.tribuneSelect.selectedIndex = this.tribuneSelect.options.length - 1;
            } else {
                this.tribuneSelect.selectedIndex = this.tribuneSelect.selectedIndex - 1;
            }
            this.setActiveTribune(this.tribuneSelect.value);
        }
        );
        hammertime.on('swiperight', (e) => {
            if (this.tribuneSelect.selectedIndex >= (this.tribuneSelect.options.length - 1)) {
                this.tribuneSelect.selectedIndex = 0;
            } else {
                this.tribuneSelect.selectedIndex = this.tribuneSelect.selectedIndex + 1;
            }
            this.setActiveTribune(this.tribuneSelect.value);
        });

    }

}
customElements.define('gc2-main', Gc2Main);