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
        this.setupBouchotSuffixor();
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
        menuButton.type = "button";
        menuButton.innerText = "⋯";
        menuButton.onclick = () => {
            this.style.display = "none";
            let menu = document.querySelector('gc2-menu');
            menu.style.display = "flex";
            menu.showSelector();
        }
        this.controls.appendChild(menuButton);
    }

    setupTribuneSelect() {
        this.tribuneSelect = document.createElement("select");
        this.tribuneSelect.onchange = () => {
            this.setActiveTribune(this.tribuneSelect.value);
        }
        this.controls.appendChild(this.tribuneSelect);
    }

    setupMessageInput() {
        this.messageInput = document.createElement("input");
        this.messageInput.id = "gc2-message";
        this.messageInput.type = "text";
        this.messageInput.spellcheck = true;
        this.messageInput.onkeydown = (e) => {
            if (event.altKey) {
                if (this.handleAltShortcut(event.key)) {
                    event.stopPropagation();
                    event.preventDefault();
                }
            }
        }
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
            tribuneElement.setAttribute("name", tribune);
            tribuneElement.style.display = "none";
            this.appendChild(tribuneElement);

            this.tribunes.set(tribune, tribuneElement);
            this.setActiveTribune(this.tribunes.size === 1 ? tribune : this.tribuneSelect.value);
        }
        return tribuneElement;
    }

    setActiveTribune(selectedTribune) {
        this.addBouchotSuffixInMessageInput(this.activeTribune);
        this.activeTribune = selectedTribune;
        this.messageInput.placeholder = selectedTribune;
        this.tribunes.forEach((tribuneElement, tribune) => {
            tribuneElement.style.display = tribune === selectedTribune ? "" : "none";
        });
    }

    setupGesture() {
        delete Hammer.defaults.cssProps.userSelect;
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

    handleAltShortcut(keychar) {
        switch (keychar) {
            case 'o':
                this.insertTextInMessage('_o/* <b>BLAM</b>! ');
                return true;
            case 'm':
                this.insertTextInMessage(`====> <b>Moment ${this.getSelectedText()}</b> <====`, 16);
                return true;
            case 'f':
                this.insertTextInMessage('\u03C6');
                return true;
            case 'b':
                this.insertTextInMessage(`<b>${this.getSelectedText()}</b>`, 3);
                return true;
            case 'i':
                this.insertTextInMessage(`<i>${this.getSelectedText()}</i>`, 3);
                return true;
            case 'u':
                this.insertTextInMessage(`<u>${this.getSelectedText()}</u>`, 3);
                return true;
            case 's':
                this.insertTextInMessage(`<s>${this.getSelectedText()}</s>`, 3);
                return true;
            case 't':
                this.insertTextInMessage(`<tt>${this.getSelectedText()}</tt>`, 4);
                return true;
            case 'c':
                this.insertTextInMessage(`<code>${this.getSelectedText()}</code>`, 6);
                return true;
            case 'd':
                this.insertTextInMessage(`<spoiler>${this.getSelectedText()}</spoiler>`, 9);
                return true;
            case 'p':
                this.insertTextInMessage('_o/* <b>paf!</b> ');
                return true;
            case 'a':
                this.insertTextInMessage(`\u266A <i>${this.getSelectedText()}</i> \u266A`, 5);
                return true;
        }
        return false;
    }

    getSelectedText() {
        return this.messageInput.value.substring(this.messageInput.selectionStart, this.messageInput.selectionEnd);
    }

    insertTextInMessage(text, pos) {
        if (!pos) {
            pos = text.length;
        }
        var selectionEnd = this.messageInput.selectionStart + pos;
        this.messageInput.value = this.messageInput.value.substring(0, this.messageInput.selectionStart) + text + this.messageInput.value.substr(this.messageInput.selectionEnd);
        this.messageInput.focus();
        this.messageInput.setSelectionRange(selectionEnd, selectionEnd);
    }

    setupBouchotSuffixor() {
        fetch("/peg/bouchotsuffixor.pegjs")
            .then((response) => {
                return response.text();
            })
            .then((text) => {
                this.bouchotSuffixor = peg.generate(text);
            });
    }

    addBouchotSuffixInMessageInput(tribune) {
        if (tribune && this.bouchotSuffixor) {
            this.messageInput.value = this.bouchotSuffixor.parse(this.messageInput.value, { bouchot: tribune });
        }
    }

}
customElements.define('gc2-main', Gc2Main);