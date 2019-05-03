class Gb3Coincoin extends HTMLElement {

    constructor() {
        super();
        this.tribunes = new Set();
    }

    connectedCallback() {
        this.setupTribuneSelect();
        this.setupMessageInput();
        this.setupPostsElement();
        this.startPoll();
    }

    setupTribuneSelect() {
        this.tribuneSelect = document.createElement("select");
        let viewAllTribuneOption = document.createElement("option");
        viewAllTribuneOption.value = "*";
        viewAllTribuneOption.text = "All";
        this.tribuneSelect.add(viewAllTribuneOption);
        this.tribuneSelect.onchange = () => {
            let tribune = this.tribuneSelect.value;
            if (tribune === "*") {
                this.querySelectorAll(`gb3-post`).forEach((e) => e.style.display = "initial");
                this.querySelectorAll(`gb3-tribune`).forEach((e) => e.style.display = "initial");
            } else {
                this.querySelectorAll(`gb3-post[id$="${tribune}"]`).forEach((e) => e.style.display = "initial");
                this.querySelectorAll(`gb3-post:not([id$="${tribune}"])`).forEach((e) => e.style.display = "none");
                this.querySelectorAll(`gb3-tribune`).forEach((e) => e.style.display = "none");
            }
        }
        this.appendChild(this.tribuneSelect);
    }

    setupMessageInput() {
        this.messageInput = document.createElement("input");
        this.messageInput.type = "text";
        this.messageInput.spellcheck = true;
        this.appendChild(this.messageInput);
    }

    setupPostsElement() {
        this.postsElement = document.createElement('gb3-posts');
        this.appendChild(this.postsElement);
    }

    startPoll() {
        let postSource = new EventSource("/api/poll");
        postSource.onmessage = (event) => {
            let post = JSON.parse(event.data);
            this.postsElement.insertPost(post);
            this.addTribune(post.tribune);
        };
    }

    addTribune(tribune) {
        if (!this.tribunes.has(tribune)) {
            this.tribunes.add(tribune);
            var option = document.createElement("option");
            option.text = option.value = tribune;
            this.tribuneSelect.add(option);
        }
    }

}
customElements.define('gb3-coincoin', Gb3Coincoin);