class Gc2Main extends HTMLElement {

    constructor() {
        super();
        this.tribunes = new Map();
    }

    async connectedCallback() {
        this.setupTribunesContainer();
        this.setupControls();
        this.setupTribuneNavigator();
        this.setupGesture();
        await this.setupBackend2html();
        await this.setupBouchotSuffixor();
        await this.listTribunes();
        this.selectTribune();
        this.appendChild(this.controls);
        this.startPoll();
    }

    setupTribuneNavigator() {
        this.tribuneNavigator = document.querySelector("gc2-tribune-navigator");
        this.tribuneNavigator.onnavigate = (tribuneName) => {
            this.tribuneSelect.value = tribuneName;
            this.setActiveTribune(tribuneName);
            this.scrollToBottom();
        }
    }

    setupTribunesContainer() {
        this.tribunesContainer = document.createElement("gc2-tribunes");
        this.appendChild(this.tribunesContainer);
    }

    async setupControls() {
        this.controls = document.createElement("form");
        this.controls.classList.add("gc2-controls");
        this.setupTribuneSelect();
        this.setupMessageInput();
        this.setupMenuButton();

        this.controls.onsubmit = (e) => {
            if (this.messageInput.value && this.tribuneSelect.value) {
                let data = new URLSearchParams();//NEVER use FormData, it forces multipart and backend dont support multipart
                data.set('message', this.messageInput.value);
                data.set('tribune', this.tribuneSelect.value);
                this.messageInput.value = "";
                this.messageInput.classList.toggle("gc2-loading", true);
                let headers = new Headers();
                if (localStorage.nickname) {
                    headers.set('User-agent', localStorage.nickname);
                }
                if (this.tribuneSelect.value === "dlfp") {
                    let accessToken = localStorage.getItem("linuxfr_access_token");
                    let expiresAt = parseInt(localStorage.getItem("linuxfr_expires_at"), 10) || 0;
                    if (!accessToken || expiresAt < Date.now()) {
                        localStorage.setItem("linuxfr_unposted_message", data.get('message'));
                        window.location.href = "/gb2c/linuxfr/authorize";
                    } else {
                        headers.set('Authorization', `Bearer ${accessToken}`);
                    }
                }
                fetch("/gb2c/post", {
                    body: data,
                    method: "POST",
                    headers: headers
                }).then((data) => {
                    this.messageInput.classList.toggle("gc2-loading", false);
                }).catch((error) => {
                    this.messageInput.classList.toggle("gc2-loading", false);
                    console.log(`Cannot post message '${data.get('message')}'. Error: `, error);
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
            menu.style.display = "";
            menu.showSelector();
        }
        this.controls.appendChild(menuButton);
    }

    setupTribuneSelect() {
        this.tribuneSelect = document.createElement("select");
        this.tribuneSelect.onchange = () => {
            this.setActiveTribune(this.tribuneSelect.value);
            this.scrollToBottom();
        }
        this.controls.appendChild(this.tribuneSelect);
    }

    setupMessageInput() {
        this.messageInput = document.createElement("input");
        this.messageInput.id = "gc2-message";
        this.messageInput.type = "text";
        this.messageInput.spellcheck = true;
        this.messageInput.onkeydown = (e) => {
            if (e.altKey) {
                if (this.handleAltShortcut(e.key)) {
                    e.stopPropagation();
                    e.preventDefault();
                }
            }
        }
        this.controls.appendChild(this.messageInput);
    }

    async setupBackend2html() {
        let response = await fetch("/peg/backend2html.pegjs");
        let text = await response.text();
        this.backend2html = peg.generate(text);
    }

    selectTribune() {
        let unposted = localStorage.getItem("linuxfr_unposted_message");
        if (unposted) {
            this.messageInput.value = unposted;
            localStorage.removeItem("linuxfr_unposted_message");
            this.setActiveTribune("dlfp", false);
        } else {
            let tribuneName = window.location.hash.substr(1);
            if(this.isTribune(tribuneName)) {
                this.setActiveTribune(tribuneName, false);
            }
        }        
    }

    startPoll() {
        let postSource = new EventSource("/gb2c/poll");
        postSource.onmessage = (event) => {
            let wasAtbottom = this.isScrollAtBottom();
            let post = JSON.parse(event.data);
            post.message = this.backend2html.parse(post.message);
            this.getTribuneElement(post.tribune).insertPost(post);
            this.updateNotifications();
            if (wasAtbottom) {
                this.scrollToBottom();
            }
        };
        postSource.onerror = (err) => {
            console.log(`Lost connection, retry in 10 seconds`);
            postSource.close();
            setTimeout(() => this.startPoll(), 10000);
        };
    }

    isScrollAtBottom() {
        return (this.tribunesContainer.scrollTop + this.tribunesContainer.clientHeight) >= this.tribunesContainer.scrollHeight;
    }

    scrollToBottom() {
        this.tribunesContainer.scrollTop = this.tribunesContainer.scrollHeight;
    }

    isTribune(tribuneName) {
        return this.tribunes.get(tribuneName);
    }

    getTribuneElement(tribuneName) {
        let tribuneElement = this.tribunes.get(tribuneName);
        if (!tribuneElement) {
            var option = document.createElement("option");
            option.text = option.value = tribuneName;
            this.tribuneSelect.add(option);

            if(this.tribuneNavigator) {
                this.tribuneNavigator.addTribune(tribuneName);
            }

            tribuneElement = document.createElement('gc2-tribune');
            tribuneElement.setAttribute("name", tribuneName);
            tribuneElement.style.display = "none";
            this.tribunesContainer.appendChild(tribuneElement);

            this.tribunes.set(tribuneName, tribuneElement);
            this.setActiveTribune(this.tribunes.size === 1 ? tribuneName : this.tribuneSelect.value);
        }
        return tribuneElement;
    }

    setActiveTribune(selectedTribuneName, markPreviousAsRead = true) {
        this.addBouchotSuffixInMessageInput(this.activeTribune);
        let previousTribuneElement = this.tribunes.get(this.activeTribune);
        if (previousTribuneElement && markPreviousAsRead) {
            previousTribuneElement.markAsRead();
        }
        this.activeTribune = selectedTribuneName;
        this.tribuneSelect.value = selectedTribuneName;
        this.messageInput.placeholder = selectedTribuneName;
        if(this.tribuneNavigator) {
            this.tribuneNavigator.setActiveTribune(selectedTribuneName);
        }
        this.tribunes.forEach((tribuneElement, tribuneName) => {
            tribuneElement.style.display = tribuneName === selectedTribuneName ? "" : "none";
        });
        this.updateNotifications();
    }

    setupGesture() {
        delete Hammer.defaults.cssProps.userSelect;
        let hammertime = new Hammer(document.querySelector("gc2-main"), {
            inputClass: Hammer.TouchInput
        });
        hammertime.on('swipeleft', (e) => {
            if (this.tribuneSelect.selectedIndex === 0) {
                this.tribuneSelect.selectedIndex = this.tribuneSelect.options.length - 1;
            } else {
                this.tribuneSelect.selectedIndex = this.tribuneSelect.selectedIndex - 1;
            }
            this.setActiveTribune(this.tribuneSelect.value);
            this.scrollToBottom();
        }
        );
        hammertime.on('swiperight', (e) => {
            if (this.tribuneSelect.selectedIndex >= (this.tribuneSelect.options.length - 1)) {
                this.tribuneSelect.selectedIndex = 0;
            } else {
                this.tribuneSelect.selectedIndex = this.tribuneSelect.selectedIndex + 1;
            }
            this.setActiveTribune(this.tribuneSelect.value);
            this.scrollToBottom();
        });

    }

    async listTribunes() {
        let response = await fetch("/gb2c/list");
        let tribunes = await response.json();
        tribunes.forEach(t => this.getTribuneElement(t));
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

    async setupBouchotSuffixor() {
        let response = await fetch("/peg/bouchotsuffixor.pegjs");
        let text = await response.text();
        this.bouchotSuffixor = peg.generate(text);
    }

    addBouchotSuffixInMessageInput(tribune) {
        if (tribune && this.bouchotSuffixor) {
            this.messageInput.value = this.bouchotSuffixor.parse(this.messageInput.value, { bouchot: tribune });
        }
    }

    updateNotifications() {
        let bigorno = false;
        let reply = false;
        this.tribunes.forEach((tribuneElement, tribuneName) => {
            let option = this.tribuneSelect.querySelector(`option[value="${tribuneName}"`);
            if (option) {
                if (tribuneElement.hasBigorno) {
                    bigorno = true;
                    if (option.innerText.indexOf("📢") < 0) {
                        option.innerText += "📢";
                    }
                } else {
                    option.innerText = option.innerText.replace("📢", "");
                }
                if (tribuneElement.hasReply) {
                    reply = true;
                    if (option.innerText.indexOf("↩") < 0) {
                        option.innerText += "↩";
                    }
                } else {
                    option.innerText = option.innerText.replace("↩", "");
                }
                if(this.tribuneNavigator) {
                    this.tribuneNavigator.setTribuneNotifications(tribuneName, tribuneElement.hasBigorno, tribuneElement.hasReply);
                }
            }
        });
        if (bigorno) {
            if (document.title.indexOf("📢") < 0) {
                document.title = `📢${document.title}`;
            }
        } else {
            document.title = document.title.replace("📢", "");
        }
        if (reply) {
            if(document.title.indexOf("↩") < 0) {
                document.title = `↩${document.title}`;
            }
        } else {
            document.title = document.title.replace("↩", "");
        }
    }

}
customElements.define('gc2-main', Gc2Main);