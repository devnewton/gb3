class Gc2Main extends HTMLElement {

    constructor() {
        super();
        this.tribunes = new Map();
    }

    connectedCallback() {
        this.setupOrder();
        this.setupControls();
        this.setupBackend2html();
        this.setupGesture();
    }

    setupOrder() {
        this.classList.toggle(`gb3-postorder-${localStorage.postOrder || "reverse-chronological"}`, true);
    }

    setupControls() {
        this.controls = document.createElement("form");
        this.controls.classList.add("gc2-controls");
        this.appendChild(this.controls);
        this.setupTribuneSelect();
        this.setupMessageInput();
        this.setupMenuButton();

        this.controls.onsubmit = (e) => {
            if (this.messageInput.value && this.tribuneSelect.value) {
                let data = new URLSearchParams();
                data.set('message', this.messageInput.value);
                data.set('tribune', this.tribuneSelect.value);
                this.messageInput.value = "";
                this.messageInput.classList.toggle("gc2-loading", true);
                let headers = new Headers();
                if (localStorage.nickname) {
                    headers.set('User-agent', localStorage.nickname);
                }
                fetch("/api/post", {
                    body: data,
                    method: "POST",
                    headers: headers
                }).then((data) => {
                    this.messageInput.classList.toggle("gc2-loading", false);
                }).catch((error) => {
                    this.messageInput.classList.toggle("gc2-loading", false);
                    console.log(`Cannot post message '${formData.get('message')}'. Error: `, error);
                });
            }
            e.preventDefault()
        };
    }

    setupMenuButton() {
        let menuButton = document.createElement("button");
        menuButton.innerText = "⋯";
        this.controls.appendChild(menuButton);
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
        let hammertime = new Hammer(document.getElementsByTagName("gc2-main")[0], {
            inputClass: Hammer.TouchInput
        });
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